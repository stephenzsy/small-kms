package agentconfig

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentActiveServerDoc struct {
	AgentConfigDoc

	EndpointURLs                    shared.AgentConfigurationAgentActiveServerEndpointUrls `json:"endpointUrls"`
	AuthorizedCertificateTemplateID shared.Identifier                                      `json:"authorizedCertificateTemplateId"`
	ServerCertificateTemplateID     shared.Identifier                                      `json:"serverCertificateTemplateId"`
	ServerCertificateID             shared.Identifier                                      `json:"serverCertificateId"`
	AuthorizedCertificateIDs        []shared.Identifier                                    `json:"authorizedCertificateIds"`
}

// toModel implements AgentConfigDocument.
func (d *AgentActiveServerDoc) toModel(isAdmin bool) *shared.AgentConfiguration {
	if d == nil {
		return nil
	}
	refreshTime := time.Now().Add(24 * time.Hour)
	refreshToken := getTimeRefreshToken(refreshTime)
	m := shared.AgentConfiguration{
		Version:          d.Version.HexString(),
		NextRefreshAfter: &refreshTime,
		NextRefreshToken: &refreshToken,
	}
	params := shared.AgentConfigurationAgentActiveServer{
		Name: shared.AgentConfigNameActiveServer,
	}
	if isAdmin {
		params.EndpointUrls = &d.EndpointURLs
		params.AuthorizedCertificateTemplateId = &d.AuthorizedCertificateTemplateID
		params.ServerCertificateTemplateId = &d.ServerCertificateTemplateID
	}
	params.AuthorizedCertificateIds = d.AuthorizedCertificateIDs
	params.ServerCertificateId = &d.ServerCertificateID
	m.Config.FromAgentConfigurationAgentActiveServer(params)
	return &m
}

var _ AgentConfigDocument = (*AgentActiveServerDoc)(nil)

func newAgentActiveServerConfigurator() *docConfigurator[AgentConfigDocument] {
	return &docConfigurator[AgentConfigDocument]{
		preparePut: func(
			c context.Context,
			nsID shared.NamespaceIdentifier, params shared.AgentConfigurationParameters) (AgentConfigDocument, error) {
			p, err := params.AsAgentConfigurationAgentActiveServer()
			if err != nil {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}

			if p.ServerCertificateTemplateId == nil || p.AuthorizedCertificateTemplateId == nil {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}

			d := AgentActiveServerDoc{
				ServerCertificateTemplateID:     *p.ServerCertificateTemplateId,
				AuthorizedCertificateTemplateID: *p.AuthorizedCertificateTemplateId,
			}
			if p.EndpointUrls != nil {
				d.EndpointURLs.Primary = p.EndpointUrls.Primary
				d.EndpointURLs.Secondary = p.EndpointUrls.Secondary
			}
			d.initLocator(nsID, shared.AgentConfigNameActiveServer)
			digester := md5.New()
			digester.Write([]byte(d.ServerCertificateTemplateID.String()))
			digester.Write([]byte(d.AuthorizedCertificateTemplateID.String()))
			d.BaseVersion = digester.Sum(nil)

			callbackDoc := AgentActiveServerCallbackDoc{
				AgentCallbackDoc: AgentCallbackDoc{
					BaseDoc: kmsdoc.BaseDoc{
						NamespaceID: nsID,
						ID:          shared.NewResourceIdentifier(shared.ResourceKindAgentCallback, shared.StringIdentifier(shared.AgentConfigNameActiveServer)),
					},
					Name: shared.AgentConfigNameActiveServer,
				},
			}

			err = kmsdoc.Create(c, &callbackDoc)
			log.Ctx(c).Error().Err(err).Msg("failed to create callback doc")

			return &d, nil
		},

		eval: func(c context.Context, doc AgentConfigDocument) (*azcosmos.PatchOperations, error) {
			d, ok := doc.(*AgentActiveServerDoc)
			if !ok {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}
			nsID := d.GetNamespaceID()

			// load the last certs
			serverCertLocator := shared.NewResourceLocator(nsID, cert.NewLatestCertificateForTemplateID(d.ServerCertificateTemplateID))
			prevServerCertId := d.ServerCertificateID
			serverCertDoc, err := cert.ReadCertDocByLocator(c, serverCertLocator)
			if err != nil {
				if !errors.Is(err, common.ErrStatusNotFound) {
					return nil, err
				}
				// should configure to drop the certificate
				d.ServerCertificateID = shared.Identifier{}
			} else {
				d.ServerCertificateID = serverCertDoc.GetLocator().GetID().Identifier()
			}

			patchOps := azcosmos.PatchOperations{}
			hasChanges := false
			if d.ServerCertificateID != prevServerCertId {
				patchOps.AppendSet("/serverCertificateId", d.ServerCertificateID)
				hasChanges = true
			}

			// load latest authorized certs
			_, err = cert.GetAuthorizedLatestCertByTemplateID(c, d.AuthorizedCertificateTemplateID)
			if err != nil {
				return nil, err
			}
			certItems, err := cert.ListActiveCertDocsByTemplateID(c, d.AuthorizedCertificateTemplateID)
			if err != nil {
				return nil, err
			}
			certLocators := utils.MapSlices(certItems, func(item *cert.CertDoc) shared.Identifier {
				return item.ID.Identifier()
			})
			if certLocators == nil {
				certLocators = []shared.Identifier{}
			}

			prevAuthorizedCertIDs := d.AuthorizedCertificateIDs
			d.AuthorizedCertificateIDs = certLocators
			if !slices.Equal(certLocators, prevAuthorizedCertIDs) {
				patchOps.AppendSet("/authorizedCertificateIds", d.AuthorizedCertificateIDs)
				hasChanges = true
			}

			if hasChanges {
				digester := md5.New()
				digester.Write(d.BaseVersion)
				digester.Write([]byte(d.ServerCertificateID.String()))
				for _, id := range d.AuthorizedCertificateIDs {
					digester.Write([]byte(id.String()))
				}
				d.Version = digester.Sum(nil)
				patchOps.AppendSet("/version", d.Version.HexString())
			}

			return &patchOps, nil
		},

		readDoc: func(c context.Context,
			nsID shared.NamespaceIdentifier) (AgentConfigDocument, error) {
			d := AgentActiveServerDoc{}
			err := kmsdoc.Read(c, NewConfigDocLocator(nsID, shared.AgentConfigNameActiveServer), &d)
			return &d, err
		},
	}
}

type AgentActiveServerCallbackDocEndpointState struct {
	Endpoint string                                               `json:"endpoint"`
	State    shared.AgentConfigurationAgentActiveServerReplyState `json:"state"`
	Version  string                                               `json:"version"` // version of the config after evaluation

}

type AgentActiveServerCallbackDoc struct {
	AgentCallbackDoc

	Primary   AgentActiveServerCallbackDocEndpointState `json:"primary"`
	Secondary AgentActiveServerCallbackDocEndpointState `json:"secondary"`
}

func NewAgentCallbackDocLocator(nsID shared.NamespaceIdentifier, configName shared.AgentConfigName) shared.ResourceLocator {
	return shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindAgentCallback, shared.StringIdentifier(configName)))
}

func ApiRecordAgentActiveServerCallback(c RequestContext, req *shared.AgentConfiguration) error {

	if reqConfig, err := req.Config.AsAgentConfigurationAgentActiveServer(); err != nil {
		return fmt.Errorf("%w:invalid input:%s", common.ErrStatusBadRequest, err)
	} else if reqConfig.Reply == nil {
		return fmt.Errorf("%w:invalid input, nil reply", common.ErrStatusBadRequest)
	} else {
		nsID := ns.GetNamespaceContext(c).GetID()
		docLocator := NewAgentCallbackDocLocator(nsID, shared.AgentConfigNameActiveServer)
		patchOps := azcosmos.PatchOperations{}
		var prefix string
		if reqConfig.Reply.SlotId == 0 {
			prefix = "/primary"
		} else {
			prefix = "/secondary"
		}
		endpoint := fmt.Sprintf("https://%s%s", c.RealIP(), reqConfig.Reply.Listener)
		patchOps.AppendSet(prefix, &AgentActiveServerCallbackDocEndpointState{
			Endpoint: endpoint,
			State:    reqConfig.Reply.State,
			Version:  req.Version,
		})

		err = kmsdoc.PatchWithLocator(c, docLocator, patchOps)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusNoContent)
}
