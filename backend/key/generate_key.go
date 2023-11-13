package key

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func GetKeyStoreName(nsKind base.NamespaceKind, nsID base.ID, policyID base.ID) string {
	return fmt.Sprintf("k-%s-%s-%s", nsKind, nsID, policyID)
}

// GenerateKey implements ServerInterface.
func (*server) GenerateKey(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, policyID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceId)

	policy, err := apiGetKeyPolicyDoc(c, policyID)
	if err != nil {
		return err
	}

	doc := &KeyDoc{}
	err = doc.init(c, policy)
	if err != nil {
		return err
	}

	kc := kv.GetAzKeyVaultService(c).AzKeysClient()
	c = c.Elevate()
	resp, err := kc.CreateKey(c, GetKeyStoreName(namespaceKind, namespaceId, policyID), doc.toAzCreateKeyParameters(), nil)
	if err != nil {
		return err
	}
	doc.KeyID = string(*resp.Key.KID)
	doc.Y = resp.Key.Y
	doc.X = resp.Key.X
	doc.N = resp.Key.N
	doc.E = resp.Key.E
	err = base.GetAzCosmosCRUDService(c).Create(c, doc, nil)
	if err != nil {
		return err
	}

	panic("unimplemented")
}
