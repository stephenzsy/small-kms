package cert

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// AddMsEntraKeyCredential implements admin.ServerInterface.
func (*CertServer) AddMsEntraKeyCredential(ec echo.Context, namespaceProvider models.NamespaceProvider, nsID string, id string, params admin.AddMsEntraKeyCredentialParams) error {
	c := ec.(ctx.RequestContext)
	nsID = ns.ResolveMeNamespace(c, nsID)
	reqIdentity := auth.GetAuthIdentity(c)
	requesterID := reqIdentity.ClientPrincipalID().String()
	var requesterProfile *profile.ProfileDoc
	var err error
	useUpdate := false
	var sp gmodels.ServicePrincipalable
	if params.OnBehalfOfApplication != nil && *params.OnBehalfOfApplication && reqIdentity.HasAdminRole() {
		c, requesterProfile, sp, err = profile.SyncServicePrincipalProfile(c, nsID, []string{"keyCredentials"})
		if err != nil {
			return err
		}
		nsID = requesterProfile.ID
		namespaceProvider = models.NamespaceProviderServicePrincipal
		requesterID = nsID
		useUpdate = true
	} else {
		c, requesterProfile, sp, err = profile.SyncServicePrincipalProfile(c, nsID, []string{"keyCredentials"})
		if err != nil {
			return err
		}
	}

	if requesterID != nsID {
		return base.ErrResponseStatusForbidden
	}

	req := certmodels.AddMsEntraKeyCredentialRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	cert, err := getCertificateInternal(c, namespaceProvider, nsID, id)
	if err != nil {
		return err
	}
	if cert.Status != certmodels.CertificateStatusIssued {
		return base.ErrResponseStatusBadRequest
	}
	if cert.NotAfter.Before(time.Now()) {
		return base.ErrResponseStatusBadRequest
	}

	if useUpdate {
		return updateApplicationWithCert(c, sp, cert, req)
	}

	return nil
}

func updateApplicationWithCert(c ctx.RequestContext, sp gmodels.ServicePrincipalable, cert *CertDoc, req certmodels.AddMsEntraKeyCredentialRequest) error {
	c, gclient, err := graph.WithDelegatedMsGraphClient(c)
	if err != nil {
		return err
	}
	app, err := gclient.ApplicationsWithAppId(sp.GetAppId()).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId", "keyCredentials"},
		},
	})
	if err != nil {
		return err
	}

	certTpStr := strings.ToLower(cert.JsonWebKey.ThumbprintSHA1.HexString())
	nextKeyCredentials := make([]gmodels.KeyCredentialable, 0, len(app.GetKeyCredentials())+1)

	for _, installedKey := range app.GetKeyCredentials() {
		tp := strings.ToLower(base64.StdEncoding.EncodeToString(installedKey.GetCustomKeyIdentifier()))
		if tp == certTpStr {
			// no action needed as certificate is already installed
			return c.NoContent(http.StatusNoContent)
		}
		installedKey.SetKey(nil)
		nextKeyCredentials = append(nextKeyCredentials, installedKey)
	}

	toAdd := gmodels.NewKeyCredential()
	toAdd.SetTypeEscaped(utils.ToPtr("AsymmetricX509Cert"))
	toAdd.SetUsage(utils.ToPtr("Verify"))
	toAdd.SetStartDateTime(&cert.NotBefore.Time)
	toAdd.SetEndDateTime(&cert.NotAfter.Time)
	toAdd.SetKey(cert.JsonWebKey.CertificateChain[0])
	nextKeyCredentials = append(nextKeyCredentials, toAdd)

	patchApplication := gmodels.NewApplication()
	patchApplication.SetKeyCredentials(nextKeyCredentials)

	_, err = gclient.Applications().ByApplicationId(*app.GetId()).Patch(c, patchApplication, nil)
	if err != nil {
		return graph.HandleMsGraphError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
