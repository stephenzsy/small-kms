package cert

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
	certtemplate "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func NewCertificateID(certId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindCert, certId)
}

func NewLatestCertificateForTemplateID(templateId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateId)
}

func ReadCertDocByLocator(c context.Context, locator shared.ResourceLocator) (*CertDoc, error) {
	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, locator, certDoc)
	return certDoc, err
}

func ApiGetCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) error {
	cert, err := getCertificate(c, certificateId, params)
	if err != nil {
		return err
	}
	return c.JSON(200, cert)
}

func getLatestCertificateByTemplateDoc(c RequestContext, templateLocator shared.ResourceLocator) (doc *CertDoc, err error) {
	doc = &CertDoc{}
	err = kmsdoc.Read[*CertDoc](c,
		shared.NewResourceLocator(templateLocator.GetNamespaceID(), shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateLocator.GetID().Identifier())), doc)
	return
}

func getCertificate(c RequestContext, certificateId shared.Identifier, params models.GetCertificateParams) (*shared.CertificateInfo, error) {
	logger := log.Ctx(c)
	var certDocLocator shared.ResourceLocator
	if certificateId.IsUUID() {
		nsID := ns.GetNamespaceContext(c).GetID()
		certDocLocator = shared.NewResourceLocator(nsID, NewCertificateID(certificateId))
	} else {
		return nil, fmt.Errorf("%w: invalid certificate ID: %s", common.ErrStatusBadRequest, certificateId)
	}

	certDoc, err := ReadCertDocByLocator(c, certDocLocator)
	if err != nil {
		return nil, err
	}
	m := certDoc.toModel()

	if params.IncludeCertificate != nil && *params.IncludeCertificate {
		// fetch cert from blob
		pemBlob, err := certDoc.FetchCertificatePEMBlob(c)
		if err != nil {
			return m, err
		}
		m.Pem = utils.ToPtr(string(pemBlob))
		var certRaw []byte
		for block, rest := pem.Decode(pemBlob); block != nil; block, rest = pem.Decode(rest) {
			if certRaw == nil && block.Type == "CERTIFICATE" {
				certRaw = block.Bytes
			}
			m.Jwk.CertificateChain = append(m.Jwk.CertificateChain, block.Bytes)
		}
		if x509Cert, err := x509.ParseCertificate(certRaw); err != nil {
			logger.Error().Err(err).Msg("failed to parse certificate")
		} else {
			switch x509Cert.PublicKeyAlgorithm {
			case x509.RSA:
				m.Jwk.Kty = shared.KeyTypeRSA
				if rsaPubKey, ok := x509Cert.PublicKey.(*rsa.PublicKey); ok {
					m.Jwk.N = rsaPubKey.N.Bytes()
					m.Jwk.E = big.NewInt(int64(rsaPubKey.E)).Bytes()
				} else {
					logger.Error().Msg("failed to parse RSA public key")
				}
			case x509.ECDSA:
				m.Jwk.Kty = shared.KeyTypeEC
				if ecdsaPubKey, ok := x509Cert.PublicKey.(*ecdsa.PublicKey); ok {
					m.Jwk.X = ecdsaPubKey.X.Bytes()
					m.Jwk.Y = ecdsaPubKey.Y.Bytes()
				} else {
					logger.Error().Msg("failed to parse ECDSA public key")
				}
			default:
				logger.Error().Msgf("unsupported public key algorithm: %s", x509Cert.PublicKeyAlgorithm.String())
			}
		}
	}

	return m, nil
}

// This call has writes, please do not use for regular query
func GetAuthorizedLatestCertByTemplateID(c context.Context, templateID shared.Identifier) (*CertDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	if !templateID.IsUUID() || templateID.UUID().Version() != 5 {
		return ReadCertDocByLocator(c, shared.NewResourceLocator(nsID, NewLatestCertificateForTemplateID(templateID)))
	}

	// read linked doc
	localTemplateLocator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, templateID))
	localTemplateDoc, err := certtemplate.GetCertificateTemplateDoc(c, localTemplateLocator)
	if err != nil {
		return nil, err
	}
	if localTemplateDoc.Owner == nil {
		return nil, fmt.Errorf("%w: template is not linked", common.ErrStatusBadRequest)
	} else if localTemplateDoc.LinkProperties == nil || localTemplateDoc.LinkProperties.Usage != models.LinkedCertificateTemplateUsageClientAuthorization {
		return nil, fmt.Errorf("%w: template is not linked for client authorization", common.ErrStatusBadRequest)
	}
	remoteTemplateLocator := *localTemplateDoc.Owner
	remoteCertLinkLocator := remoteTemplateLocator.WithIDKind(shared.ResourceKindLatestCertForTemplate)
	remoteCertLinkDoc, err := ReadCertDocByLocator(c, remoteCertLinkLocator)
	if err != nil {
		return nil, err
	}
	// create link
	targetFinalLocator := remoteCertLinkDoc.GetLocator()
	targetCertDoc, err := ReadCertDocByLocator(c, targetFinalLocator)
	if err != nil {
		return nil, err
	}

	linkedCertDoc := *targetCertDoc
	linkedCertDoc.NamespaceID = nsID
	linkedCertDoc.ID = targetFinalLocator.GetID()
	linkedCertDoc.Owner = &targetFinalLocator
	linkedCertDoc.Owns = nil
	linkedCertDoc.Template = localTemplateLocator

	eCtx := common.ElevateContext(c)
	err = kmsdoc.Upsert(eCtx, &linkedCertDoc)
	if err != nil {
		return nil, err
	}
	patchOps := azcosmos.PatchOperations{}
	if targetCertDoc.Owns == nil {
		patchOps.AppendSet(kmsdoc.PatchPathOwns, map[shared.NamespaceIdentifier]shared.ResourceLocator{
			nsID: linkedCertDoc.GetLocator(),
		})
	} else {
		patchOps.AppendSet(fmt.Sprintf("%s/%s", kmsdoc.PatchPathOwns, nsID), linkedCertDoc.GetLocator())
	}
	err = kmsdoc.Patch(eCtx, targetCertDoc, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: &targetCertDoc.ETag,
	})
	return &linkedCertDoc, err
}
