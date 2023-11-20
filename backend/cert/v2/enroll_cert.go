package cert

import (
	"crypto"
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
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// EnrollCertificate implements admin.ServerInterface.
func (*CertServer) EnrollCertificate(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, policyID string, params admin.EnrollCertificateParams) error {
	c := ec.(ctx.RequestContext)
	reqIdentity := auth.GetAuthIdentity(c)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if params.OnBehalfOfApplication != nil && *params.OnBehalfOfApplication && reqIdentity.HasAdminRole() {
		// TODO
		return c.JSON(http.StatusNotImplemented, map[string]string{"message": "not implemented"})
		// appID := reqIdentity.AppID()
		// profile.SyncProfileWithAppID(c, appID)
	}
	if !reqIdentity.HasAdminRole() && !reqIdentity.HasRole(auth.RoleValueAgentActiveHost) && !reqIdentity.HasRole(auth.RoleValueCertificateEnroll) {
		return base.ErrResponseStatusForbidden
	}
	namespaceUUID, err := uuid.Parse(namespaceId)
	if err != nil {
		namespaceUUID = uuid.UUID{}
	}

	canEnroll := false
	var requesterProfile *profile.ProfileDoc
	if reqIdentity.ClientPrincipalID() == namespaceUUID {
		// authroize self
		requesterProfile, err = profile.SyncProfileInternal(c, reqIdentity.ClientPrincipalID().String())
		if err != nil {
			return err
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
		if _, requesterProfile, nsProfile, err = profile.SyncMemberOfInternal(c, reqIdentity.ClientPrincipalID().String(), namespaceId); err == nil {
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

	policy, err := getCertificatePolicyInternal(c, namespaceProvider, namespaceId, policyID)
	if err != nil {
		return err
	} else if !policy.AllowEnroll {
		return fmt.Errorf("%w: policy %s does not allow enroll", base.ErrResponseStatusBadRequest, policyID)
	}

	return enrollInternal(c, requesterProfile.TargetNamespaceProvider(), requesterProfile.ID, policy, req)

}

func enrollInternal(c ctx.RequestContext, nsProvider models.NamespaceProvider, nsID string, policy *CertPolicyDoc,
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

	return c.JSON(http.StatusCreated, certDoc.ToModel(true))

}

type certDocEnrollPending struct {
	certDocPending

	templateX509Cert *x509.Certificate
	issuerX509Cert   *x509.Certificate
	issuerCertChain  []cloudkey.Base64RawURLEncodableBytes
	publicKey        crypto.PublicKey
	signer           crypto.Signer
}

func (d *certDocEnrollPending) init(
	c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc,
	publicKey *cloudkey.JsonWebSignatureKey) error {
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

	issuerPolicy, err := getCertificatePolicyInternal(c, pDoc.IssuerPolicy.NamespaceProvider, pDoc.IssuerPolicy.NamespaceID, pDoc.IssuerPolicy.ID)
	if err != nil {
		return err
	}
	signerCert, err := issuerPolicy.getIssuerCert(c)
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
	d.signer = cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, azKeysClient, signerCert.JsonWebKey.KeyID, signerCert.JsonWebKey.Alg)
	d.templateX509Cert.SignatureAlgorithm = signerCert.JsonWebKey.Alg.X509SignatureAlgorithm()

	return nil
}
