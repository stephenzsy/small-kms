package agentconfig

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentActiveServerDoc struct {
	AgentConfigDoc

	AuthorizedCertificateTemplate shared.ResourceLocator `json:"authorizedCertificateTemplate"`
	ServerCertificateTemplate     shared.ResourceLocator `json:"serverCertificateTemplate"`
	ServerCertificate             shared.ResourceLocator `json:"serverCertificate"`
	AuthorizedCertificate         shared.ResourceLocator `json:"authorizedCertificate"`
}

// toModel implements AgentConfigDocument.
func (d *AgentActiveServerDoc) toModel(isAdmin bool) *shared.AgentConfiguration {
	if d == nil {
		return nil
	}
	refreshTime := d.Updated.Add(24 * time.Hour)
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
		params.AuthorizedCertificateTemplate = &d.AuthorizedCertificateTemplate
		params.ServerCertificateTemplate = &d.ServerCertificateTemplate
	}
	params.AuthorizedCertificate = &d.AuthorizedCertificate
	params.ServerCertificate = &d.ServerCertificate
	m.Config.FromAgentConfigurationAgentActiveServer(params)
	return &m
}

var _ AgentConfigDocument = (*AgentActiveServerDoc)(nil)

var systemNamespaceId shared.NamespaceIdentifier = shared.NewNamespaceIdentifier(shared.NamespaceKindSystem,
	shared.StringIdentifier(ns.SystemServiceNameAgentPush))

func newAgentActiveServerConfigurator() *docConfigurator[AgentConfigDocument] {
	return &docConfigurator[AgentConfigDocument]{
		preparePut: func(
			c context.Context,
			nsID shared.NamespaceIdentifier, params shared.AgentConfigurationParameters) (AgentConfigDocument, error) {
			p, err := params.AsAgentConfigurationAgentActiveServer()
			if err != nil {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}

			if p.ServerCertificateTemplate == nil || p.AuthorizedCertificateTemplate == nil {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}

			d := AgentActiveServerDoc{
				ServerCertificateTemplate:     *p.ServerCertificateTemplate,
				AuthorizedCertificateTemplate: *p.AuthorizedCertificateTemplate,
			}
			d.initLocator(nsID, shared.AgentConfigNameActiveServer)
			digester := md5.New()
			digester.Write([]byte(d.ServerCertificateTemplate.String()))
			digester.Write([]byte(d.AuthorizedCertificateTemplate.String()))
			d.BaseVersion = digester.Sum(nil)

			return &d, nil
		},

		eval: func(c context.Context, doc AgentConfigDocument) (*azcosmos.PatchOperations, error) {
			d, ok := doc.(*AgentActiveServerDoc)
			if !ok {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}
			nsID := d.GetNamespaceID()
			if d.ServerCertificateTemplate.GetNamespaceID() != nsID {
				return nil, fmt.Errorf("%w:invalid input, server cert must be frome the same namespace", common.ErrStatusBadRequest)
			}
			if d.AuthorizedCertificateTemplate.GetNamespaceID() != systemNamespaceId {
				return nil, fmt.Errorf("%w:invalid input, authorized cert must be from system namespace", common.ErrStatusBadRequest)
			}

			// load the last certs
			serverCertLocator := shared.NewResourceLocator(nsID, cert.NewLatestCertificateForTemplateID(d.ServerCertificateTemplate.GetID().Identifier()))
			oldServerCertLocator := d.ServerCertificate
			serverCertDoc, err := cert.ReadCertDocByLocator(c, serverCertLocator)
			if err != nil {
				if !errors.Is(err, common.ErrStatusNotFound) {
					return nil, err
				}
				// should configure to drop the certificate
				d.ServerCertificate = shared.ResourceLocator{}
			} else {
				d.ServerCertificate = serverCertDoc.GetLocator()
			}

			patchOps := azcosmos.PatchOperations{}
			hasChanges := false
			if d.ServerCertificate != oldServerCertLocator {
				patchOps.AppendSet("/serverCertificate", d.ServerCertificate)
				hasChanges = true
			}

			// load latest authorized certs
			authorizedCertLocator := shared.NewResourceLocator(systemNamespaceId,
				cert.NewLatestCertificateForTemplateID(d.AuthorizedCertificateTemplate.GetID().Identifier()))
			authedDoc, err := cert.ReadCertDocByLocator(c, authorizedCertLocator)
			oldAuthorizedCertificateLocator := d.AuthorizedCertificate
			if err != nil {
				if !errors.Is(err, common.ErrStatusNotFound) {
					return nil, err
				}
				// should configure to drop the certificate
				d.AuthorizedCertificate = shared.ResourceLocator{}
			} else {
				d.AuthorizedCertificate = authedDoc.GetLocator()
			}
			if d.AuthorizedCertificate != oldAuthorizedCertificateLocator {
				patchOps.AppendSet("/authorizedCertificate", d.AuthorizedCertificate)
				hasChanges = true
			}

			if hasChanges {
				digester := md5.New()
				digester.Write(d.BaseVersion)
				digester.Write([]byte(d.ServerCertificate.String()))
				digester.Write([]byte(d.AuthorizedCertificate.String()))
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
