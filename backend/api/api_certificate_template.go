package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

// GetCertificateTemplate implements models.ServerInterface.
func (s *server) GetCertificateTemplate(ec echo.Context, namespaceKind shared.NamespaceKind, namespaceId shared.Identifier, templateID shared.Identifier) error {
	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	respData, respErr := (func() (*models.CertificateTemplateComposed, error) {

		c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
		if err != nil {
			return nil, err
		}
		c, err = ct.WithCertificateTemplateContext(c, templateID)
		if err != nil {
			return nil, err
		}
		return ct.GetCertificateTemplate(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// PutCertificateTemplate implements models.ServerInterface.
func (s *server) PutCertificateTemplate(ec echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	templateID shared.Identifier) error {
	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	respData, respErr := (func() (*models.CertificateTemplateComposed, error) {
		req := models.CertificateTemplateParameters{}
		err := ec.Bind(&req)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid input body", common.ErrStatusBadRequest)
		}

		c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceId)
		if err != nil {
			return nil, err
		}
		c, err = ct.WithCertificateTemplateContext(c, templateID)
		if err != nil {
			return nil, err
		}
		return ct.PutCertificateTemplate(c, req)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// ListProfiles implements models.ServerInterface.
func (s *server) ListCertificateTemplates(ec echo.Context, namespaceKind shared.NamespaceKind, namespaceId shared.Identifier) error {
	c := ec.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	respData, respErr := (func() ([]*models.CertificateTemplateRefComposed, error) {
		c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
		if err != nil {
			return nil, err
		}
		return ct.ListCertificateTemplates(c)
	})()
	return wrapResponse(ec, http.StatusOK, respData, respErr)
}

// CreateLinkedCertificateTemplate implements models.ServerInterface.
func (*server) CreateLinkedCertificateTemplate(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier) error {
	params := models.CreateLinkedCertificateTemplateParameters{}
	err := ctx.Bind(&params)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	return wrapEchoResponse(c, ct.ApiCreateLinkedCertificateTemplate(c, params))
}

// ListKeyVaultRoleAssignments implements models.ServerInterface.
func (*server) ListKeyVaultRoleAssignments(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	templateID shared.Identifier) error {
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	bad := func(e error) error {
		return wrapResponse[[]*models.AzureRoleAssignment](ctx, http.StatusOK, nil, e)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return bad(err)
	}
	result, err := ct.ListKeyVaultRoleAssignments(c)
	if err == nil && result == nil {
		result = []*models.AzureRoleAssignment{}
	}
	return wrapResponse[[]*models.AzureRoleAssignment](c, http.StatusOK, result, err)
}

func (*server) RemoveKeyVaultRoleAssignment(ctx echo.Context,
	namespaceKind shared.NamespaceKind,
	namespaceId shared.Identifier,
	templateID shared.Identifier,
	roleAssignmentID string) error {

	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	bad := func(e error) error {
		return wrapResponse[any](ctx, http.StatusNoContent, nil, e)
	}

	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return bad(err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return bad(err)
	}
	err = ct.DeleteKeyVaultRoleAssignment(c, roleAssignmentID)
	return wrapResponse[any](c, http.StatusNoContent, nil, err)
}

// AddKeyVaultRoleAssignment implements models.ServerInterface.
func (*server) AddKeyVaultRoleAssignment(ctx echo.Context, namespaceKind shared.NamespaceKind, namespaceID shared.Identifier, templateID shared.Identifier,
	params models.AddKeyVaultRoleAssignmentParams) error {
	bad := func(e error) error {
		return wrapResponse[*models.AzureRoleAssignment](ctx, http.StatusOK, nil, e)
	}

	roleDefID, err := ct.ValidateRoleDefnitionIDForAdd(params.RoleDefinitionId)
	if err != nil {
		return bad(err)
	}
	c := ctx.(RequestContext)

	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err = ns.WithNamespaceContext(c, namespaceKind, namespaceID)
	if err != nil {
		return bad(err)
	}
	c, err = ct.WithCertificateTemplateContext(c, templateID)
	if err != nil {
		return bad(err)
	}
	result, err := ct.AddKeyVaultRoleAssignment(c, roleDefID)
	return wrapResponse[*models.AzureRoleAssignment](c, http.StatusOK, result, err)
}

// DeleteCertificateTemplate implements models.ServerInterface.
func (*server) DeleteCertificateTemplate(ctx echo.Context, namespaceKind shared.NamespaceKind, namespaceId shared.Identifier, templateId shared.Identifier) error {
	c := ctx.(RequestContext)
	if ok := auth.AuthorizeAdminOnly(c); !ok {
		return respondRequireAdmin(c)
	}
	c, err := ns.WithNamespaceContext(c, namespaceKind, namespaceId)
	if err != nil {
		return wrapEchoResponse(c, err)
	}
	return wrapEchoResponse(c, errors.New("not implemented"))
}
