package agentconfig

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentProfileDocInstalledMsEntraClientCreds struct {
	CertificateID  shared.Identifier             `json:"certificateId"`
	ThumbprintSHA1 shared.CertificateFingerprint `json:"thumbprint"`
	GraphKeyID     uuid.UUID                     `json:"graphKeyId"`
}

type AgentProfileDoc struct {
	kmsdoc.BaseDoc

	Status                                  shared.AgentProfileStatus                    `json:"status"`
	MsEntraClientCredsTemplateID            shared.Identifier                            `json:"msEntraClientCredsTemplateId"`
	MsEntraClientCredsInstalledCertificates []AgentProfileDocInstalledMsEntraClientCreds `json:"msEntraClientCredsInstalledCertificates"`
	AppRoles                                []shared.AgentProfileAppRole                 `json:"appRoles"`
}

func (d *AgentProfileDoc) toModel() *shared.AgentProfile {
	if d == nil {
		return nil
	}
	r := new(shared.AgentProfile)

	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.Status = d.Status
	r.MsEntraClientCredentialCertificateTemplateId = d.MsEntraClientCredsTemplateID
	r.MsEntraClientCredentialInstalledCertificateIds = utils.MapSlice(d.MsEntraClientCredsInstalledCertificates, func(item AgentProfileDocInstalledMsEntraClientCreds) shared.Identifier {
		return item.CertificateID
	})
	r.AppRoles = d.AppRoles
	return r
}

const (
	patchKeyMsEntraClientCredsInstalledCertificates = "/msEntraClientCredsInstalledCertificates"
	patchKeyAppRoles                                = "/appRoles"
)

var agentProfileIdentifier = shared.NewResourceIdentifier(shared.ResourceKindReserved, shared.StringIdentifier("agent-profile"))

type ThumbprintSHA1 = [sha1.Size]byte

func getKeyCredentialThumbprintSHA1(kc gmodels.KeyCredentialable) (ThumbprintSHA1, string, error) {
	fp := utils.CertificateFingerprintSHA1{}
	encodedCustomKeyIdentifier := base64.StdEncoding.EncodeToString(kc.GetCustomKeyIdentifier())
	err := fp.UnmarshalText([]byte(encodedCustomKeyIdentifier))
	return fp, encodedCustomKeyIdentifier, err
}

func getToBeProvisionedClientCreds(c context.Context, templateID shared.Identifier, keyCredentials []gmodels.KeyCredentialable) ([]gmodels.KeyCredentialable, map[ThumbprintSHA1]*cert.CertDoc, error) {
	bad := func(e error) ([]gmodels.KeyCredentialable, map[ThumbprintSHA1]*cert.CertDoc, error) {
		return nil, nil, e
	}

	// latestCert, err := cert.GetAuthorizedLatestCertByTemplateID(c, templateID)
	// if err != nil {
	// 	return bad(err)
	// }

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
	// if _, hasKey := installedKeys[[sha1.Size]byte(latestCert.Thumbprint)]; !hasKey {
	// 	log.Info().Msgf("latest cert not installed, thumbprint: %s", latestCert.Thumbprint.HexString())
	// 	// not installed, install in graph
	// 	kc := gmodels.NewKeyCredential()
	// 	pemBlob, err := latestCert.FetchCertificatePEMBlob(c)
	// 	if err != nil {
	// 		return bad(err)
	// 	}
	// 	block, _ := pem.Decode(pemBlob)
	// 	kc.SetKey(block.Bytes)
	// 	kc.SetUsage(to.Ptr("Verify"))
	// 	kc.SetTypeEscaped(to.Ptr("AsymmetricX509Cert"))
	// 	kc.SetStartDateTime(latestCert.NotBefore.TimePtr())
	// 	kc.SetEndDateTime(latestCert.NotAfter.TimePtr())
	// 	patchKeyCredentials = append(patchKeyCredentials, kc)
	// 	hasChange = true
	// }
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

const (
	AppRoleAgentPushConfig = "Agent.PushConfig"
)

func getToBeProvisionedAppRoles(c context.Context, appRoles []gmodels.AppRoleable) []gmodels.AppRoleable {
	rolesMap := utils.ToMapFunc(appRoles, func(item gmodels.AppRoleable) string {
		return *item.GetValue()
	})
	if _, ok := rolesMap[AppRoleAgentPushConfig]; ok {
		return nil
	}

	pushConfigRole := gmodels.NewAppRole()
	pushConfigRole.SetAllowedMemberTypes([]string{"User", "Application"})
	pushConfigRole.SetDescription(to.Ptr("Agent push config"))
	pushConfigRole.SetDisplayName(to.Ptr(AppRoleAgentPushConfig))
	pushConfigRole.SetId(to.Ptr(uuid.New()))
	pushConfigRole.SetIsEnabled(to.Ptr(true))
	pushConfigRole.SetValue(to.Ptr(AppRoleAgentPushConfig))
	return []gmodels.AppRoleable{pushConfigRole}
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
		c := ctx.Elevate(c)
		graph := common.GetAdminServerClientProvider(c).MsGraphClient()
		applicationIdReqBuilder := graph.Applications().ByApplicationId(nsID.Identifier().UUID().String())
		queryApplicationParameters := &applications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "appId", "appRoles", "identifierUris", "keyCredentials", "oauth2Permissions"},
		}
		applicationable, err := applicationIdReqBuilder.Get(c, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
			QueryParameters: queryApplicationParameters,
		})
		if err != nil {
			return doc, err
		}
		// match key credentials
		patchApplication := gmodels.NewApplication()
		patchKeyCredentials, allowedCerts, err := getToBeProvisionedClientCreds(c, doc.MsEntraClientCredsTemplateID, applicationable.GetKeyCredentials())
		if err != nil {
			return doc, err
		}
		if patchKeyCredentials != nil {
			patchApplication.SetKeyCredentials(patchKeyCredentials)
		}

		patchAppRoles := getToBeProvisionedAppRoles(c, applicationable.GetAppRoles())
		if patchAppRoles != nil {
			patchApplication.SetAppRoles(patchAppRoles)
		}

		identifierUris := applicationable.GetIdentifierUris()
		if len(identifierUris) == 0 {
			patchApplication.SetIdentifierUris([]string{fmt.Sprintf("api://%s", *applicationable.GetAppId())})
		}

		if patchKeyCredentials != nil || patchAppRoles != nil || len(identifierUris) == 0 {
			_, err = applicationIdReqBuilder.Patch(c, patchApplication, nil)
			if err != nil {
				return doc, err
			}
			applicationable, err = applicationIdReqBuilder.Get(c, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
				QueryParameters: queryApplicationParameters,
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
		patchedAppRoles := applicationable.GetAppRoles()
		doc.AppRoles = utils.MapSlice(patchedAppRoles, func(item gmodels.AppRoleable) shared.AgentProfileAppRole {
			return shared.AgentProfileAppRole{
				ID:    *item.GetId(),
				Value: *item.GetValue(),
			}
		})
		patchOps := azcosmos.PatchOperations{}
		patchOps.AppendSet(patchKeyMsEntraClientCredsInstalledCertificates, docCc)
		patchOps.AppendSet(patchKeyAppRoles, doc.AppRoles)
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
