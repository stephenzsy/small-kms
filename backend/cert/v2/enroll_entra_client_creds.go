package cert

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/key"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func enrollMsEntraClientCredCert(c ctx.RequestContext, policyRID base.Identifier, params *EnrollMsEntraClientCredentialRequest) error {

	// verify jwt is 2048
	if params.PublicKey.Kty != key.JsonWebKeyTypeRSA {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid public key type"})
	}

	pKey, err := params.PublicKey.AsRsaPubicKey()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid public key"})
	}

	nsCtx := ns.GetNSContext(c)

	// verify proof of jwt, so make sure client has possession of the private key
	if token, err := jwt.Parse(params.MsEntraProof, func(token *jwt.Token) (interface{}, error) {
		return pKey, nil
	}); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof"})
	} else if aud, err := token.Claims.GetAudience(); err != nil || !slices.Contains(aud, "00000003-0000-0000-c000-000000000000") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has audience of '00000003-0000-0000-c000-000000000000'"})
	} else if iss, err := token.Claims.GetIssuer(); err != nil || base.StringIdentifier(iss) != nsCtx.Identifier() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid proof, must has issuer of '%s'", nsCtx.Identifier().String())})
	} else if nbf, err := token.Claims.GetNotBefore(); err != nil || time.Until(nbf.Time) > 5*time.Minute || time.Until(nbf.Time) < -5*time.Minute {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has not before within 5 minutes"})
	} else if exp, err := token.Claims.GetExpirationTime(); err != nil || exp.Time != nbf.Time.Add(10*time.Minute) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proof, must has expiration time of 10 minutes"})
	}

	// issue certificate
	certDoc, err := createCertFromPolicy(c, policyRID, pKey)
	if err != nil {
		return err
	}

	certLocator := base.SLocator{
		NID: base.GetDefaultStorageNamespaceID(nsCtx.Kind(), nsCtx.Identifier()),
		RID: certDoc.GetStorageID(c),
	}

	// check existing certificates
	linkDoc := &CertLinkRelDoc{}
	linkDoc.initNamespaceMsEntraClientCredentials(nsCtx.Kind(), nsCtx.Identifier())

	linkDoc, err = getNamespaceLinkRelDoc(c, nsCtx.Kind(), nsCtx.Identifier(), RelNameMsEntraClientCredentials)
	if err != nil {
		if !errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return err
		}
	}

	if linkDoc.Relations == nil {
		linkDoc.Relations = new(base.DocRelations)
	}
	if linkDoc.Relations.NamedToList == nil {
		linkDoc.Relations.NamedToList = make(map[base.RelName][]base.SLocator, 1)
	}
	if l, hasValue := linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials]; !hasValue || len(l) == 0 {
		linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials] = []base.SLocator{certLocator}
	} else {
		linkDoc.Relations.NamedToList[RelNameMsEntraClientCredentials] = []base.SLocator{certLocator, l[0]}
	}

	// materialize certificate
	return c.JSON(http.StatusOK, certDoc)

}
