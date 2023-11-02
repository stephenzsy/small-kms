package agentconfig

import (
	"context"
	"crypto/sha1"
	"net/http"

	"github.com/google/uuid"
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
	return r
}

var agentProfileIdentifier = shared.NewResourceIdentifier(shared.ResourceKindReserved, shared.StringIdentifier("agent-profile"))

type ThumbprintSHA1 = [sha1.Size]byte

const (
	AppRoleAgentPushConfig = "Agent.PushConfig"
)

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
