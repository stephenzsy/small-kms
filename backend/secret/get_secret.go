package secret

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
)

// GetSecret implements ServerInterface.
func (s *server) GetSecret(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, secretID base.ID, params GetSecretParams) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	doc := &SecretDoc{}
	docSvc := base.GetAzCosmosCRUDService(c)
	if err := docSvc.Read(c, base.NewDocLocator(nsKind, nsID, base.ResourceKindSecret, secretID), doc, nil); err != nil {
		return wrapSecretNotFoundError(err, secretID)
	}

	model := &Secret{}
	doc.PopulateModel(model)

	if params.WithValue != nil && *params.WithValue {
		// this operation is senstive
		c, client, err := kv.WithDelegatedAzSecretsClient(c, s.GetAzKeyVaultEndpoint())
		if err != nil {
			return err
		}
		sid := azsecrets.ID(doc.KeyVaultStore.ID)
		resp, err := client.GetSecret(c, sid.Name(), sid.Version(), nil)
		if err != nil {
			return err
		}
		model.Value = *resp.Value
		model.ContentType = *resp.ContentType
	}

	return c.JSON(http.StatusOK, model)
}

func wrapSecretNotFoundError(err error, secretID base.ID) error {
	if errors.Is(err, base.ErrAzCosmosDocNotFound) {
		return fmt.Errorf("%w, secret not found: %s", base.ErrResponseStatusNotFound, secretID)
	}
	return err
}
