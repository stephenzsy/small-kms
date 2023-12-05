package cert

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// EnrollCertificate implements admin.ServerInterface.
func (s *CertServer) EnrollCertificate(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, policyID string, params admin.EnrollCertificateParams) (err error) {
	c := ec.(ctx.RequestContext)
	reqIdentity := auth.GetAuthIdentity(c)
	requesterID := reqIdentity.ClientPrincipalID()
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	var requesterProfile *profile.ProfileDoc
	if params.OnBehalfOfApplication != nil && *params.OnBehalfOfApplication && reqIdentity.HasAdminRole() {
		c, requesterProfile, _, err = profile.SyncServicePrincipalProfileByAppID(c, reqIdentity.AppID(), nil)
		if err != nil {
			return err
		}
		namespaceId = requesterProfile.ID
		if requesterID, err = uuid.Parse(requesterProfile.ID); err != nil {
			return err
		}
		namespaceProvider = models.NamespaceProviderServicePrincipal
	}
	if !reqIdentity.HasAdminRole() && !reqIdentity.HasRole(auth.RoleValueAgentActiveHost) && !reqIdentity.HasRole(auth.RoleValueCertificateEnroll) {
		return base.ErrResponseStatusForbidden
	}
	namespaceUUID, err := uuid.Parse(namespaceId)
	if err != nil {
		namespaceUUID = uuid.UUID{}
	}

	canEnroll := false
	if requesterID == namespaceUUID {
		// authroize self
		if requesterProfile == nil {
			gclient := graph.GetServiceMsGraphClient(c)

			requesterProfile, err = profile.SyncProfileInternal(c, requesterID.String(), gclient)
			if err != nil {
				return err
			}
		}
		requesterTemplateVarData := &ResourceTemplateGraphVarData{
			ID: requesterProfile.ID,
		}
		c = c.WithValue(templateContextKeyRequesterGraph, requesterTemplateVarData)
		c = c.WithValue(templateContextKeyNamespaceGraph, requesterTemplateVarData)
		canEnroll = true
	} else if namespaceProvider == models.NamespaceProviderGroup {
		// authorize group member
		var nsProfile *profile.ProfileDoc
		if _, requesterProfile, nsProfile, err = profile.SyncMemberOfInternal(c, requesterID.String(), namespaceId); err == nil {
			c = c.WithValue(templateContextKeyRequesterGraph, &ResourceTemplateGraphVarData{
				ID: requesterProfile.ID,
			})
			c = c.WithValue(templateContextKeyNamespaceGraph, &ResourceTemplateGraphVarData{
				ID: nsProfile.ID,
			})
			canEnroll = true
		}
	}
	if !canEnroll {
		return base.ErrResponseStatusForbidden
	}

	req := &certmodels.EnrollCertificateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	policy, err := GetCertificatePolicyInternal(c, namespaceProvider, namespaceId, policyID)
	if err != nil {
		return err
	} else if !policy.AllowEnroll {
		return fmt.Errorf("%w: policy %s does not allow enroll", base.ErrResponseStatusBadRequest, policyID)
	}

	return s.enrollInternal(c, requesterProfile.TargetNamespaceProvider(), requesterProfile.ID, policy, req)

}

func (s *CertServer) enrollInternal(c ctx.RequestContext, nsProvider models.NamespaceProvider, nsID string, policy *CertPolicyDoc,
	req *certmodels.EnrollCertificateRequest) (err error) {
	certDoc := &certDocInternal{}

	if err = certDoc.init(c, nsProvider, nsID, policy, &req.PublicKey); err != nil {
		return
	}

	var csr CertCSR
	if nsProvider != models.NamespaceProviderRootCA {
		csr, err = certDoc.GetCertificateRequest(c, true)
		if err != nil {
			return err
		}
	}

	docSvc := resdoc.GetDocService(c)
	der, err := certDoc.CreateCertificate(c, csr)
	if err != nil {
		return err
	}
	if err := certDoc.CollectSignedCertificate(c, der); err != nil {
		return err
	}

	resp, err := docSvc.Create(c, certDoc, nil)
	if err != nil {
		return err
	}

	m := certDoc.ToModel(true)
	return c.JSON(resp.RawResponse.StatusCode, m)
}
