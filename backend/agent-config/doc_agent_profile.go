package agentconfig

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
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

func ApiGetAgentProfile(c RequestContext) error {
	nsID := ns.GetNamespaceContext(c).GetID()
	docLocator := shared.NewResourceLocator(nsID, agentProfileIdentifier)
	doc, err := readAgentProfile(c, docLocator)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, doc.toModel())
}
