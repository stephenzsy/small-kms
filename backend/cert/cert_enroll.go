package cert

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// EnrollMsEntraClientCredential implements ServerInterface.
func (s *server) EnrollCertificate(ec echo.Context, nsKind base.NamespaceKind, nsID ID, policyID ID, params EnrollCertificateParams) error {
	c := ec.(ctx.RequestContext)

	c, nsCtx := ns.WithResovingMeNSContext(c, nsKind, nsID)
	c, authzOk := authz.Authorize(c,
		nsCtx.AllowSelf(),
		s.allowAdminDelegatedServicePrincipal(nsID),
		s.allowGraphGroupMemberOf(nsKind, nsCtx.ID()))
	if !authzOk {
		return base.ErrResponseStatusForbidden
	}

	req := new(EnrollCertificateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	return enrollMsEntraClientCredCert(c, base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindCertPolicy, policyID), req)
}

const graphAudVerify = "00000003-0000-0000-c000-000000000000"

func enrollMsEntraClientCredCert(c ctx.RequestContext, policyLocator base.DocFullIdentifier, params *EnrollCertificateRequest) error {
	matchAud := string(policyLocator.NamespaceID())

	if params.EnrollmentType == EnrollmentTypeMsEntraClientCredential {
		matchAud = graphAudVerify
		// jwt must be RSA
		if params.PublicKey.KeyType != cloudkey.KeyTypeRSA {
			return fmt.Errorf("%w: invalid public key type: %s", base.ErrResponseStatusBadRequest, params.PublicKey.KeyType)
		}
	}

	pKey := params.PublicKey.PublicKey()
	nsCtx := ns.GetNSContext(c)

	// verify proof of jwt, so make sure client has possession of the private key
	if token, err := jwt.Parse(params.Proof, func(token *jwt.Token) (interface{}, error) {
		return pKey, nil
	}); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof"})
	} else if aud, err := token.Claims.GetAudience(); err != nil || !slices.Contains(aud, matchAud) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid proof, must has audience of '%s'", matchAud)})
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
