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
	ServerCertificateID           shared.Identifier      `json:"serverCertificateId"`
	AuthorizedCertificateID       shared.Identifier      `json:"authorizedCertificateId"`
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
	params.AuthorizedCertificateId = &d.AuthorizedCertificateID
	params.ServerCertificateId = &d.ServerCertificateID
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
			authorizedCertLocator := shared.NewResourceLocator(systemNamespaceId,
				cert.NewLatestCertificateForTemplateID(d.AuthorizedCertificateTemplate.GetID().Identifier()))
			authedDoc, err := cert.ReadCertDocByLocator(c, authorizedCertLocator)
			prevAuthorizedCertID := d.AuthorizedCertificateID
			if err != nil {
				if !errors.Is(err, common.ErrStatusNotFound) {
					return nil, err
				}
				// should configure to drop the certificate
				d.AuthorizedCertificateID = shared.Identifier{}
			} else {
				linkedDoc, err := cert.CreateLinkedCertificate(c, authedDoc)
				if err != nil {
					return nil, err
				}
				d.AuthorizedCertificateID = linkedDoc.GetLocator().GetID().Identifier()
			}
			if d.AuthorizedCertificateID != prevAuthorizedCertID {
				patchOps.AppendSet("/authorizedCertificateId", d.AuthorizedCertificateID)
				hasChanges = true
			}

			if hasChanges {
				digester := md5.New()
				digester.Write(d.BaseVersion)
				digester.Write([]byte(d.ServerCertificateID.String()))
				digester.Write([]byte(d.AuthorizedCertificateID.String()))
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
