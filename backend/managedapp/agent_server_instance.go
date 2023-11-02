package managedapp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	agentauth "github.com/stephenzsy/small-kms/backend/agent/auth"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	cloudkeyx "github.com/stephenzsy/small-kms/backend/cloud/key/x"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentInstanceDoc struct {
	base.BaseDoc
	AgentInstanceFields
}

func (d *AgentInstanceDoc) init(nsKind base.NamespaceKind, nsID base.Identifier, rID base.Identifier, req AgentInstanceFields) {
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindAgentInstance, rID)
	d.AgentInstanceFields = req
}

func (d *AgentInstanceDoc) PopulateModel(r *AgentInstance) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&r.ResourceReference)
	r.AgentInstanceFields = d.AgentInstanceFields
}

func apiPutAgentInstance(c ctx.RequestContext, instanceID base.Identifier, req AgentInstanceFields) error {
	nsCtx := ns.GetNSContext(c)
	doc := &AgentInstanceDoc{}
	doc.init(nsCtx.Kind(), nsCtx.Identifier(), instanceID, req)

	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Upsert(c, doc, nil); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

type AgentInstanceQueryDoc struct {
	base.QueryBaseDoc
	AgentInstanceFields
}

func apiListAgentInstances(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)

	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns("c.endpoint", "c.version", "c.buildId").
		WithOrderBy("c.ts DESC")
	docSvc := base.GetAzCosmosCRUDService(c)
	pager := base.NewQueryDocPager[*AgentInstanceQueryDoc](docSvc, qb, base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindAgentInstance))
	sPager := utils.NewSerializableItemsPager(c, pager)
	return c.JSON(http.StatusOK, sPager)
}

// wraps not found with 404
func apiReadAgentInstanceDoc(c ctx.RequestContext, instanceID base.Identifier) (*AgentInstanceDoc, error) {
	nsCtx := ns.GetNSContext(c)
	doc := &AgentInstanceDoc{}
	docSvc := base.GetAzCosmosCRUDService(c)
	err := docSvc.Read(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindAgentInstance, instanceID), doc, nil)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: agent instance with id %s not found", base.ErrResponseStatusNotFound, instanceID.String())
		}
		return nil, err
	}
	return doc, err
}

func apiGetAgentInstance(c ctx.RequestContext, instanceID base.Identifier) error {
	doc, err := apiReadAgentInstanceDoc(c, instanceID)
	if err != nil {
		return err
	}
	m := &AgentInstance{}
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}

func apiCreateAgentInstanceProxyAuthToken(c ctx.RequestContext, resourceIdentifier base.Identifier) error {

	instanceDoc, err := apiReadAgentInstanceDoc(c, resourceIdentifier)
	if err != nil {
		return err
	}
	if instanceDoc.Endpoint == "" {
		return fmt.Errorf("%w: no endpoint found", base.ErrResponseStatusBadRequest)
	}

	configDoc, err := apiReadAgentConfigDoc(c)
	if err != nil {
		return err
	}
	if len(configDoc.JWTKeyCertIDs) == 0 {
		return fmt.Errorf("%w: no JWT key cert IDs configured", base.ErrResponseStatusNotFound)
	}

	certDoc, err := cert.ReadCertDocByFullIdentifier(c, configDoc.JWTKeyCertIDs[0])
	if err != nil {
		return err
	}

	azKeyVaultService := c.Value(kv.AzKeyVaultServiceContextKey).(kv.AzKeyVaultService)

	if certDoc.KeySpec.KeyID == nil {
		return fmt.Errorf("%w: no key ID found", base.ErrResponseStatusBadRequest)
	}
	ck := cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, azKeyVaultService.AzKeysClient(), *certDoc.KeySpec.KeyID)
	var jwtSigningMethod jwt.SigningMethod
	switch certDoc.KeySpec.Kty {

	case cloudkey.KeyTypeEC:
		switch *certDoc.KeySpec.Crv {
		case cloudkey.CurveNameP256:
			jwtSigningMethod = cloudkeyx.NewJWTSigningMethod(cloudkey.SignatureAlgorithmES256)
		case cloudkey.CurveNameP384:
			jwtSigningMethod = cloudkeyx.NewJWTSigningMethod(cloudkey.SignatureAlgorithmES384)
		case cloudkey.CurveNameP521:
			jwtSigningMethod = cloudkeyx.NewJWTSigningMethod(cloudkey.SignatureAlgorithmES512)
		}
	}
	if jwtSigningMethod == nil {
		return fmt.Errorf("%w: unsupported key type", base.ErrResponseStatusBadRequest)
	}
	identity := auth.GetAuthIdentity(c)
	accessToken, err := agentauth.NewSignedAgentAuthJWT(jwtSigningMethod, identity.ClientPrincipalID().String(), instanceDoc.Endpoint, ck)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &AuthResult{
		AccessToken: accessToken,
	})
}
