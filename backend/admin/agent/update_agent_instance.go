package agentadmin

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func getInstanceID(endpoint string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(endpoint))
}

// CreateAgentInstance implements admin.ServerInterface.
func (*AgentAdminServer) UpdateAgentInstance(ec echo.Context, nsID string) error {
	c := ec.(ctx.RequestContext)
	nsID = ns.ResolveMeNamespace(c, nsID)
	if _, authOk := authz.Authorize(c, authz.AllowHasRole(auth.RoleValueAgentActiveHost)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	req := new(agentmodels.AgentInstanceParameters)
	if err := c.Bind(req); err != nil {
		return err
	}
	if req.Endpoint == "" {
		return base.ErrResponseStatusBadRequest
	}

	docID := getInstanceID(req.Endpoint)

	doc := &AgentInstanceDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderServicePrincipal,
				NamespaceID:       nsID,
				ResourceProvider:  models.ResourceProviderAgentInstance,
			},
			ID: docID.String(),
		},
		Endpoint:         req.Endpoint,
		State:            req.State,
		ConfigVersion:    req.ConfigVersion,
		BuildID:          req.BuildId,
		TlsCertificateID: req.TlsCertificateId,
		JwtVerfyKeyID:    req.JwtVerifyKeyId,
	}

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
