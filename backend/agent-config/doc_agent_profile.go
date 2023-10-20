package agentconfig

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentProfileDocInstalledMsEntraClientCreds struct {
	CertificateID  shared.Identifier             `json:"certificateId"`
	ThumbprintSHA1 shared.CertificateFingerprint `json:"x5tHex"`
	GraphKeyID     uuid.UUID                     `json:"graphKeyId"`
}

type AgentProfileDoc struct {
	kmsdoc.BaseDoc

	Status                                  shared.AgentProfileStatus                    `json:"status"`
	MsEntraClientCredsTemplateID            shared.Identifier                            `json:"msEntraClientCredsTemplateId"`
	MsEntraClientCredsInstalledCertificates []AgentProfileDocInstalledMsEntraClientCreds `json:"msEntraClientCredsInstalledCertificates"`
}

func (d *AgentProfileDoc) toModel() *shared.AgentProfile {
	if d == nil {
		return nil
	}
	r := new(shared.AgentProfile)

	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.Status = d.Status
	r.MsEntraClientCredentialCertificateTemplateId = d.MsEntraClientCredsTemplateID
	r.MsEntraClientCredentialInstalledCertificateIds = utils.MapSlices(d.MsEntraClientCredsInstalledCertificates, func(item AgentProfileDocInstalledMsEntraClientCreds) shared.Identifier {
		return item.CertificateID
	})
	return r
}

const (
	patchKeyMsEntraClientCredsInstalledCertificates = "/msEntraClientCredsInstalledCertificates"
)

var agentProfileIdentifier = shared.NewResourceIdentifier(shared.ResourceKindReserved, shared.StringIdentifier("agent-profile"))

type ThumbprintSHA1 = [sha1.Size]byte

func getKeyCredentialThumbprintSHA1(kc gmodels.KeyCredentialable) (ThumbprintSHA1, string, error) {
	fp := utils.CertificateFingerprintSHA1{}
	encodedCustomKeyIdentifier := base64.StdEncoding.EncodeToString(kc.GetCustomKeyIdentifier())
	err := fp.UnmarshalText([]byte(encodedCustomKeyIdentifier))
	return fp, encodedCustomKeyIdentifier, err
}

func getToBeProvisionClientCreds(c context.Context, templateID shared.Identifier, keyCredentials []gmodels.KeyCredentialable) ([]gmodels.KeyCredentialable, map[ThumbprintSHA1]*cert.CertDoc, error) {
	bad := func(e error) ([]gmodels.KeyCredentialable, map[ThumbprintSHA1]*cert.CertDoc, error) {
		return nil, nil, e
	}

	latestCert, err := cert.GetAuthorizedLatestCertByTemplateID(c, templateID)
	if err != nil {
		return bad(err)
	}

	certDocs, err := cert.ListActiveCertDocsByTemplateID(c, templateID)
	if err != nil {
		return bad(err)
	}

	installedKeys := utils.ToMapFunc(keyCredentials, func(item gmodels.KeyCredentialable) ThumbprintSHA1 {
		fp, str, err := getKeyCredentialThumbprintSHA1(item)
		if err != nil {
			log.Error().Err(err).Msgf("failed to parse key credential thumbprint: %s", str)
		}
		return fp
	})
	allowedCerts := utils.ToMapFunc(certDocs, func(item *cert.CertDoc) ThumbprintSHA1 {
		return ThumbprintSHA1(item.Thumbprint)
	})

	hasChange := false
	patchKeyCredentials := make([]gmodels.KeyCredentialable, 0, len(certDocs))
	if _, hasKey := installedKeys[[sha1.Size]byte(latestCert.Thumbprint)]; !hasKey {
		log.Info().Msgf("latest cert not installed, thumbprint: %s", latestCert.Thumbprint.HexString())
		// not installed, install in graph
		kc := gmodels.NewKeyCredential()
		pemBlob, err := latestCert.FetchCertificatePEMBlob(c)
		if err != nil {
			return bad(err)
		}
		block, _ := pem.Decode(pemBlob)
		kc.SetKey(block.Bytes)
		kc.SetUsage(to.Ptr("Verify"))
		kc.SetTypeEscaped(to.Ptr("AsymmetricX509Cert"))
		kc.SetStartDateTime(latestCert.NotBefore.TimePtr())
		kc.SetEndDateTime(latestCert.NotAfter.TimePtr())
		patchKeyCredentials = append(patchKeyCredentials, kc)
		hasChange = true
	}
	for tp, key := range installedKeys {
		if _, ok := allowedCerts[tp]; ok {
			// key still allowed, keep it
			patchKeyCredentials = append(patchKeyCredentials, key)
		} else {
			// notify has changes
			hasChange = true
		}

	}

	if !hasChange {
		return nil, allowedCerts, nil
	}
	return patchKeyCredentials, allowedCerts, nil
}

func readAgentProfile(c context.Context, docLocator shared.ResourceLocator) (*AgentProfileDoc, error) {
	doc := AgentProfileDoc{}
	err := kmsdoc.Read(c, docLocator, &doc)
	return &doc, err
}

func provisionAgentProfile(c RequestContext, params *shared.AgentProfileParameters) (*AgentProfileDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	docLocator := shared.NewResourceLocator(nsID, agentProfileIdentifier)
	doc, err := readAgentProfile(c, docLocator)
	if err != nil {
		if errors.Is(err, common.ErrStatusNotFound) {
			doc = &AgentProfileDoc{
				BaseDoc: kmsdoc.BaseDoc{
					NamespaceID: nsID,
					ID:          agentProfileIdentifier,
				},
				MsEntraClientCredsTemplateID: params.MsEntraClientCredentialCertificateTemplateId,
				Status:                       shared.AgentProfileStatusPending,
			}

			err = kmsdoc.Create(c, doc)
		}
		if err != nil {
			return doc, err
		}
	}

	{
		c := c.Elevate()
		graph := common.GetAdminServerClientProvider(c).MsGraphClient()
		applicationIdReqBuilder := graph.Applications().ByApplicationId(nsID.Identifier().UUID().String())
		applicationable, err := applicationIdReqBuilder.Get(c, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "appId", "keyCredentials"},
			},
		})
		if err != nil {
			return doc, err
		}
		// match key credentials
		patchApplication := gmodels.NewApplication()
		patchKeyCredentials, allowedCerts, err := getToBeProvisionClientCreds(c, doc.MsEntraClientCredsTemplateID, applicationable.GetKeyCredentials())
		if err != nil {
			return doc, err
		}
		if patchKeyCredentials != nil {
			patchApplication.SetKeyCredentials(patchKeyCredentials)
		}

		if patchKeyCredentials != nil {
			_, err = applicationIdReqBuilder.Patch(c, patchApplication, nil)
			if err != nil {
				return doc, err
			}
			applicationable, err = applicationIdReqBuilder.Get(c, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "displayName", "appId", "keyCredentials"},
				},
			})
			if err != nil {
				return doc, err
			}
		}
		// apply patch to doc
		patchedKeyCredentials := applicationable.GetKeyCredentials()
		docCc := make([]AgentProfileDocInstalledMsEntraClientCreds, 0, len(patchedKeyCredentials))
		for _, kc := range patchedKeyCredentials {
			if thumbprint, _, err := getKeyCredentialThumbprintSHA1(kc); err == nil {
				if matchedCert, ok := allowedCerts[thumbprint]; ok {
					docCc = append(docCc, AgentProfileDocInstalledMsEntraClientCreds{
						CertificateID:  matchedCert.ID.Identifier(),
						ThumbprintSHA1: thumbprint[:],
						GraphKeyID:     *kc.GetKeyId(),
					})
				}
			}
		}
		patchOps := azcosmos.PatchOperations{}
		patchOps.AppendSet(patchKeyMsEntraClientCredsInstalledCertificates, docCc)
		doc.MsEntraClientCredsInstalledCertificates = docCc
		err = kmsdoc.Patch(c, doc, patchOps, nil)
		if err != nil {
			return doc, err
		}

	}
	return doc, nil
}

func ApiProvisionAgentProfile(c RequestContext, params *shared.AgentProfileParameters) error {
	doc, err := provisionAgentProfile(c, params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.toModel())
}

func ApiGetAgentProfile(c RequestContext) error {
	nsID := ns.GetNamespaceContext(c).GetID()
	docLocator := shared.NewResourceLocator(nsID, agentProfileIdentifier)
	doc, err := readAgentProfile(c, docLocator)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.toModel())
}
