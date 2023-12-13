package cert

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/admin/profile"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// AddMsEntraKeyCredential implements admin.ServerInterface.
func (*CertServer) AddMsEntraKeyCredential(ec echo.Context, namespaceProvider models.NamespaceProvider, nsID string, id string) error {
	c := ec.(ctx.RequestContext)

	if namespaceProvider != models.NamespaceProviderServicePrincipal {
		return base.ErrResponseStatusNotFound
	}

	nsID = ns.ResolveMeNamespace(c, nsID)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(nsID)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	reqIdentity := auth.GetAuthIdentity(c)
	requesterID := reqIdentity.ClientPrincipalID().String()
	var namespaceProfile *profile.ProfileDoc
	var err error
	var gclient *msgraphsdkgo.GraphServiceClient

	if reqIdentity.ClientPrincipalID().String() == nsID {
		// we need client ID
		gclient = graph.GetServiceMsGraphClient(c)

		namespaceProfile, err = profile.SyncProfileInternal(c, nsID, gclient)
		if err != nil {
			return err
		}
	} else {
		c, gclient, err = graph.WithDelegatedMsGraphClient(c)
		if err != nil {
			return err
		}
		c, namespaceProfile, _, err = profile.SyncServicePrincipalProfile(c, nsID, []string{"keyCredentials"})
		if err != nil {
			return err
		}
	}

	if requesterID != nsID {
		return base.ErrResponseStatusForbidden
	}

	cert, err := GetCertificateInternal(c, namespaceProvider, nsID, id)
	if err != nil {
		return err
	}
	if cert.GetStatus() != certmodels.CertificateStatusIssued {
		return base.ErrResponseStatusBadRequest
	}
	if cert.IsExpired() {
		return base.ErrResponseStatusBadRequest
	}

	return updateApplicationWithCert(c, namespaceProfile, cert, gclient)
}

func updateApplicationWithCert(c ctx.RequestContext, profile *profile.ProfileDoc, cert CertDocument, gclient *msgraphsdkgo.GraphServiceClient) error {
	app, err := gclient.ApplicationsWithAppId(profile.AppId).Get(c, &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
		QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
			Select: []string{"id", "appId", "keyCredentials"},
		},
	})
	if err != nil {
		return err
	}

	certTpStr := strings.ToLower(cert.GetJsonWebKey().ThumbprintSHA1.HexString())
	nextKeyCredentials := make([]gmodels.KeyCredentialable, 0, len(app.GetKeyCredentials())+1)

	for _, installedKey := range app.GetKeyCredentials() {
		tp := strings.ToLower(base64.StdEncoding.EncodeToString(installedKey.GetCustomKeyIdentifier()))
		if tp == certTpStr {
			// no action needed as certificate is already installed
			return c.NoContent(http.StatusNoContent)
		}
		installedKey.SetKey(nil)
		installedKeyEnd := installedKey.GetEndDateTime()
		if installedKeyEnd == nil || installedKeyEnd.After(time.Now()) {
			// only add non-expired keys
			nextKeyCredentials = append(nextKeyCredentials, installedKey)
		}
	}

	toAdd := gmodels.NewKeyCredential()
	toAdd.SetTypeEscaped(utils.ToPtr("AsymmetricX509Cert"))
	toAdd.SetUsage(utils.ToPtr("Verify"))
	toAdd.SetStartDateTime(utils.ToPtr(cert.GetNotBefore()))
	toAdd.SetEndDateTime(utils.ToPtr(cert.GetNotAfter()))
	toAdd.SetKey(cert.GetCertificateBytes())
	nextKeyCredentials = append(nextKeyCredentials, toAdd)

	patchApplication := gmodels.NewApplication()
	patchApplication.SetKeyCredentials(nextKeyCredentials)

	_, err = gclient.Applications().ByApplicationId(*app.GetId()).Patch(c, patchApplication, nil)
	if err != nil {
		return graph.HandleMsGraphError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
