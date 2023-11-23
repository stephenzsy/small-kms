package cert

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// EnrollCertificate implements admin.ServerInterface.
func (s *CertServer) EnrollCertificate(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, policyID string, params admin.EnrollCertificateParams) (err error) {
	c := ec.(ctx.RequestContext)
	reqIdentity := auth.GetAuthIdentity(c)
	requesterID := reqIdentity.ClientPrincipalID()
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	var requesterProfile *profile.ProfileDoc
	if params.OnBehalfOfApplication != nil && *params.OnBehalfOfApplication && reqIdentity.HasAdminRole() {
		c, requesterProfile, _, err = profile.SyncServicePrincipalProfileByAppID(c, reqIdentity.AppID(), nil)
		if err != nil {
			return err
		}
		namespaceId = requesterProfile.ID
		if requesterID, err = uuid.Parse(requesterProfile.ID); err != nil {
			return err
		}
		namespaceProvider = models.NamespaceProviderServicePrincipal
	}
	if !reqIdentity.HasAdminRole() && !reqIdentity.HasRole(auth.RoleValueAgentActiveHost) && !reqIdentity.HasRole(auth.RoleValueCertificateEnroll) {
		return base.ErrResponseStatusForbidden
	}
	namespaceUUID, err := uuid.Parse(namespaceId)
	if err != nil {
		namespaceUUID = uuid.UUID{}
	}

	canEnroll := false
	if requesterID == namespaceUUID {
		// authroize self
		if requesterProfile == nil {
			gclient := graph.GetServiceMsGraphClient(c)

			requesterProfile, err = profile.SyncProfileInternal(c, requesterID.String(), gclient)
			if err != nil {
				return err
			}
		}
		requesterTemplateVarData := &ResourceTemplateGraphVarData{
			ID: requesterProfile.ID,
		}
		c = c.WithValue(templateContextKeyRequesterGraph, requesterTemplateVarData)
		c = c.WithValue(templateContextKeyNamespaceGraph, requesterTemplateVarData)
		canEnroll = true
	} else if namespaceProvider == models.NamespaceProviderGroup {
		// authorize group member
		var nsProfile *profile.ProfileDoc
		if _, requesterProfile, nsProfile, err = profile.SyncMemberOfInternal(c, requesterID.String(), namespaceId); err == nil {
			c = c.WithValue(templateContextKeyRequesterGraph, &ResourceTemplateGraphVarData{
				ID: requesterProfile.ID,
			})
			c = c.WithValue(templateContextKeyNamespaceGraph, &ResourceTemplateGraphVarData{
				ID: nsProfile.ID,
			})
			canEnroll = true
		}
	}
	if !canEnroll {
		return base.ErrResponseStatusForbidden
	}

	req := &certmodels.EnrollCertificateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	policy, err := GetCertificatePolicyInternal(c, namespaceProvider, namespaceId, policyID)
	if err != nil {
		return err
	} else if !policy.AllowEnroll {
		return fmt.Errorf("%w: policy %s does not allow enroll", base.ErrResponseStatusBadRequest, policyID)
	}

	return s.enrollInternal(c, requesterProfile.TargetNamespaceProvider(), requesterProfile.ID, policy, req)

}

func (s *CertServer) enrollInternal(c ctx.RequestContext, nsProvider models.NamespaceProvider, nsID string, policy *CertPolicyDoc,
	req *certmodels.EnrollCertificateRequest) (err error) {
	certDoc := &certDocEnrollPending{
		certDocPending: certDocPending{
			CertDoc: CertDoc{
				ResourceDoc: resdoc.ResourceDoc{
					PartitionKey: resdoc.PartitionKey{
						NamespaceProvider: nsProvider,
						NamespaceID:       nsID,
						ResourceProvider:  models.ResourceProviderCert,
					},
				},
			},
		},
	}

	if err = certDoc.init(c, nsProvider, nsID, policy, &req.PublicKey); err != nil {
		return
	}

	if req.WithOneTimePkcs12Key != nil && *req.WithOneTimePkcs12Key {
		oneTimeKey, err := s.cryptoStore.GenerateECDSAKeyPair(elliptic.P384())
		if err != nil {
			return err
		}
		certDoc.OneTimePkcs12Key = &cloudkey.JsonWebKey{}
		certDoc.OneTimePkcs12Key.X = oneTimeKey.X.Bytes()
		certDoc.OneTimePkcs12Key.Y = oneTimeKey.Y.Bytes()
		certDoc.OneTimePkcs12Key.D = oneTimeKey.D.Bytes()
		certDoc.OneTimePkcs12Key.Curve = cloudkey.CurveNameP384
		certDoc.OneTimePkcs12Key.KeyType = cloudkey.KeyTypeEC
		certDoc.OneTimePkcs12Key.Alg = "ECDH-ES"
		certDoc.OneTimePkcs12Key.KeyOperations = []cloudkey.JsonWebKeyOperation{cloudkey.JsonWebKeyOperationDeriveKey}
	}

	var signed []byte
	signed, err = x509.CreateCertificate(rand.Reader,
		certDoc.templateX509Cert,
		certDoc.issuerX509Cert,
		certDoc.publicKey,
		certDoc.signer)
	if err != nil {
		return err
	}
	certDoc.collectSignedCert(signed)
	certDoc.Checksum = certDoc.calculateChecksum()

	_, docCreateErr := resdoc.GetDocService(c).Create(c, certDoc, nil)
	if docCreateErr != nil {
		return docCreateErr
	}

	return c.JSON(http.StatusCreated, certDoc.ToModel())

}

type certDocEnrollPending struct {
	certDocPending

	templateX509Cert *x509.Certificate
	issuerX509Cert   *x509.Certificate
	publicKey        crypto.PublicKey
	signer           crypto.Signer
	OneTimePkcs12Key *cloudkey.JsonWebKey `json:"oneTimePkcs12Key"`
}

func (d *certDocEnrollPending) init(
	c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc,
	publicKey *cloudkey.JsonWebKey) error {
	if err := d.certDocPending.commonInitPending(c, nsProvider, nsID, pDoc); err != nil {
		return err
	}

	if d.JsonWebKey.KeyType != publicKey.KeyType {
		return fmt.Errorf("%w: public key type does not match", base.ErrResponseStatusBadRequest)
	}
	if d.JsonWebKey.KeyType == cloudkey.KeyTypeRSA {
		if publicKey.N.BitLen() != *pDoc.KeySpec.KeySize {
			return fmt.Errorf("%w: public key size does not match", base.ErrResponseStatusBadRequest)
		}
	} else if d.JsonWebKey.KeyType == cloudkey.KeyTypeEC {
		if d.JsonWebKey.Curve != publicKey.Curve {
			return fmt.Errorf("%w: public key curve does not match", base.ErrResponseStatusBadRequest)
		}
	}

	d.JsonWebKey.N = publicKey.N
	d.JsonWebKey.E = publicKey.E
	d.JsonWebKey.X = publicKey.X
	d.JsonWebKey.Y = publicKey.Y
	pubKey := publicKey.PublicKey()
	d.publicKey = pubKey

	d.templateX509Cert = d.generateCertificateTemplate()

	issuerPolicy, err := GetCertificatePolicyInternal(c, pDoc.IssuerPolicy.NamespaceProvider, pDoc.IssuerPolicy.NamespaceID, pDoc.IssuerPolicy.ID)
	if err != nil {
		return err
	}
	signerCert, err := issuerPolicy.getIssuerCert(c)
	d.Issuer = signerCert.Identifier()
	if err != nil {
		return err
	} else if signerCert.Status != certmodels.CertificateStatusIssued {
		return fmt.Errorf("issuer certificate is not issued")
	} else if time.Until(signerCert.NotAfter.Time) < 24*time.Hour {
		return fmt.Errorf("issuer certificate is expiring soon or has expired")
	}
	d.issuerCertChain = signerCert.JsonWebKey.CertificateChain
	d.issuerX509Cert, err = x509.ParseCertificate(signerCert.JsonWebKey.CertificateChain[0])

	if err != nil {
		return err
	}
	azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
	d.signer = cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, azKeysClient, signerCert.JsonWebKey.KeyID, cloudkey.JsonWebSignatureAlgorithm(signerCert.JsonWebKey.Alg), true)
	d.templateX509Cert.SignatureAlgorithm = cloudkey.JsonWebSignatureAlgorithm(signerCert.JsonWebKey.Alg).X509SignatureAlgorithm()

	return nil
}

func (d *certDocEnrollPending) ToModel() (m certmodels.Certificate) {
	m = d.certDocPending.ToModel(true)
	if d.OneTimePkcs12Key != nil {
		m.OneTimePkcs12Key = d.OneTimePkcs12Key.PublicJWK()
	}
	return m
}
