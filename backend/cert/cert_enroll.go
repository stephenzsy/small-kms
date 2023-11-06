package cert

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// EnrollMsEntraClientCredential implements ServerInterface.
func (s *server) EnrollCertificate(ec echo.Context, nsKind base.NamespaceKind, nsID ID, policyID ID, params EnrollCertificateParams) error {
	c := ec.(ctx.RequestContext)

	certNsKind, certNsID, err := authEnrollCertificate(c, nsKind, nsID)
	if err != nil {
		return err
	}

	req := new(EnrollCertificateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	c = ns.WithDefaultNSContext(c, certNsKind, certNsID)

	return enrollMsEntraClientCredCert(c, base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindCertPolicy, policyID), req)
}

func authEnrollCertificate(c ctx.RequestContext, policyNsKind base.NamespaceKind, policyNsID ID) (base.NamespaceKind, ID, error) {
	bad := func(e error) (base.NamespaceKind, ID, error) {
		return policyNsKind, policyNsID, e
	}

	identity := auth.GetAuthIdentity(c)
	requesterID := base.IDFromUUID(identity.ClientPrincipalID())
	if requesterID == policyNsID {
		return policyNsKind, policyNsID, nil
	}
	var err error
	var gclient *msgraphsdkgo.GraphServiceClient
	if c, gclient, err = graph.WithDelegatedMsGraphClient(c); err != nil {
		return bad(err)
	}

	// authorize for admin role
	if identity.HasAdminRole() {
		requestAppID := identity.AppID()
		if requestAppID == "" {
			if sp, err := gclient.ServicePrincipalsWithAppId(&requestAppID).Get(c, &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetRequestConfiguration{
				QueryParameters: &serviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderGetQueryParameters{
					Select: []string{"id"},
				},
			}); err != nil {
				return bad(err)
			} else if *sp.GetId() == string(policyNsID) {
				return policyNsKind, policyNsID, nil
			}
		}
	}

	// authorize for requester
	dirObjBuilder := gclient.DirectoryObjects().ByDirectoryObjectId(string(requesterID))
	if dirObj, err := dirObjBuilder.Get(c, &directoryobjects.DirectoryObjectItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &directoryobjects.DirectoryObjectItemRequestBuilderGetQueryParameters{
			Select: []string{"id"},
		},
	}); err != nil {
		return bad(err)
	} else if policyNsKind != base.NamespaceKindGroup {
		return bad(fmt.Errorf("%w: only group policies can be enrolled by member", base.ErrResponseStatusBadRequest))
	} else {
		requestBody := directoryobjects.NewItemCheckMemberGroupsPostRequestBody()
		requestBody.SetGroupIds([]string{string(policyNsID)})
		resp, err := dirObjBuilder.CheckMemberGroups().Post(c, requestBody, nil)
		if err != nil {
			return bad(err)
		}
		if !slices.Contains(resp.GetValue(), string(policyNsID)) {
			return bad(fmt.Errorf("%w: requester %s is not a member of the group %s", base.ErrResponseStatusBadRequest, requesterID, policyNsID))
		}
		switch dirObj.(type) {
		case gmodels.Userable:
			return base.NamespaceKindUser, requesterID, nil
		case gmodels.ServicePrincipalable:
			return base.NamespaceKindServicePrincipal, requesterID, nil
		default:
			return bad(fmt.Errorf("%w: unsupported requester type: %s", base.ErrResponseStatusBadRequest, *dirObj.GetOdataType()))
		}
	}

}

const graphAudVerify = "00000003-0000-0000-c000-000000000000"

func enrollMsEntraClientCredCert(c ctx.RequestContext, policyLocator base.DocFullIdentifier, params *EnrollCertificateRequest) error {

	// verify jwt is 2048
	if params.PublicKey.Kty != cloudkey.KeyTypeRSA {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid public key type"})
	}

	pKey, err := params.PublicKey.AsRsaPubicKey()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid public key"})
	}

	nsCtx := ns.GetNSContext(c)

	matchAud := string(policyLocator.NamespaceID())
	if params.EnrollmentType == EnrollmentTypeMsEntraClientCredential {
		matchAud = graphAudVerify
	}

	// verify proof of jwt, so make sure client has possession of the private key
	if token, err := jwt.Parse(params.Proof, func(token *jwt.Token) (interface{}, error) {
		return pKey, nil
	}); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof"})
	} else if aud, err := token.Claims.GetAudience(); err != nil || !slices.Contains(aud, matchAud) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has audience of '00000003-0000-0000-c000-000000000000'"})
	} else if iss, err := token.Claims.GetIssuer(); err != nil || base.ParseID(iss) != nsCtx.ID() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid proof, must has issuer of '%s'", nsCtx.ID())})
	} else if nbf, err := token.Claims.GetNotBefore(); err != nil || time.Until(nbf.Time) > 5*time.Minute || time.Until(nbf.Time) < -5*time.Minute {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has not before within 5 minutes"})
	} else if exp, err := token.Claims.GetExpirationTime(); err != nil || exp.Time != nbf.Time.Add(10*time.Minute) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has expiration time of 10 minutes"})
	}

	// issue certificate
	certDoc, err := createCertFromPolicy(c, policyLocator, pKey)
	if err != nil {
		return err
	}

	// check existing certificates
	// linkDoc := &CertLinkRelDoc{}
	// linkDoc.initNamespaceMsEntraClientCredentials(nsCtx.Kind(), nsCtx.Identifier())

	// linkDoc, err = getNamespaceLinkRelDoc(c, nsCtx.Kind(), nsCtx.Identifier(), RelNameMsEntraClientCredentials)
	// if err != nil {
	// 	if !errors.Is(err, base.ErrAzCosmosDocNotFound) {
	// 		return err
	// 	}
	// }

	// if linkDoc.Relations == nil {
	// 	linkDoc.Relations = new(base.DocRelations)
	// }
	// if linkDoc.Relations.NamedToList == nil {
	// 	linkDoc.Relations.NamedToList = make(map[base.RelName][]base.SLocator, 1)
	// }
	// if l, hasValue := linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials]; !hasValue || len(l) == 0 {
	// 	linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials] = []base.SLocator{certLocator}
	// } else {
	// 	linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials] = []base.SLocator{certLocator, l[0]}
	// }

	m := new(Certificate)
	certDoc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)

}
