package cert

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/applicationswithappid"
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type CertRuleMsEntraClientCredDoc = CertRuleIssuerLastNCertificateDoc

// PopulateModel implements base.ModelPopulater.
func (d *CertRuleMsEntraClientCredDoc) PopulateModel(r *CertificateRuleMsEntraClientCredential) {
	if d == nil || r == nil {
		return
	}
	r.CertificateIds = d.CertificateIDs
	r.PolicyId = d.PolicyID
}

// var _ base.ModelPopulater[CertificateRuleIssuer] = (*CertRuleIssuerDoc)(nil)

func readCertRuleMsEntraClientDoc(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertRuleIssuerLastNCertificateDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerLastNCertificateDoc)
	err := docSvc.Read(c, getNamespaceCertificateRuleDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.ID(), base.CertRuleNameMsEntraClientCredential), ruleDoc, nil)
	return ruleDoc, err
}

// func getNamespaceIssuerCert(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertDoc, error) {

// 	ruleDoc, err := readCertRuleIssuerDoc(c, nsIdentifier)
// 	if err != nil {
// 		return nil, err
// 	}

// 	docSvc := base.GetAzCosmosCRUDService(c)
// 	certDoc := new(CertDoc)
// 	err = docSvc.Read(c, base.NewDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), base.ResourceKindCert, ruleDoc.CertificateID), certDoc, nil)
// 	return certDoc, err
// }

func apiGetCertRuleMsEntraClientCredential(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)
	ruleDoc, err := readCertRuleMsEntraClientDoc(c, base.NewNamespaceIdentifier(nsCtx.Kind(), nsCtx.ID()))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, ms entra credential configuration not found: %s", base.ErrResponseStatusNotFound, base.CertRuleNameMsEntraClientCredential)
		}
		return err
	}
	m := new(CertificateRuleMsEntraClientCredential)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}

func apiPutCertRuleMsEntraClientCredentrial(c ctx.RequestContext, p *CertificateRuleMsEntraClientCredential) error {
	nsCtx := ns.GetNSContext(c)
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleMsEntraClientCredDoc)
	ruleDoc.init(nsCtx.Kind(), nsCtx.ID(), base.CertRuleNameMsEntraClientCredential)
	ruleDoc.PolicyID = p.PolicyId
	if len(p.CertificateIds) == 0 {
		certIds, err := QueryLatestCertificateIdsIssuedByPolicy(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, p.PolicyId), 2)
		if err != nil {
			return err
		}
		ruleDoc.CertificateIDs = certIds
	}
	{
		c := ctx.Elevate(c)
		err := docSvc.Upsert(c, ruleDoc, nil)
		if err != nil {
			return err
		}
		if err := applyMsEntraClientCredential(c, nsCtx.ID().UUID(), ruleDoc.CertificateIDs); err != nil {
			return err
		}
	}
	m := new(CertificateRuleMsEntraClientCredential)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}

func applyMsEntraClientCredential(c context.Context, servicePrincipalID uuid.UUID, certIDs []base.ID) error {
	provisioningCerts := make(map[string]*CertDoc, len(certIDs))
	for _, certID := range certIDs {
		certDoc, err := ApiReadCertDocByID(c, certID)
		if err != nil {
			continue
		}
		provisioningCerts[certDoc.KeySpec.X5t.HexString()] = certDoc
	}

	gclient := graph.GetServiceMsGraphClient(c)
	sp, err := gclient.ServicePrincipals().ByServicePrincipalId(servicePrincipalID.String()).Get(c,
		&serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "appId"},
			},
		})
	if err != nil {
		return err
	}
	application, err := gclient.ApplicationsWithAppId(sp.GetAppId()).Get(c,
		&applicationswithappid.ApplicationsWithAppIdRequestBuilderGetRequestConfiguration{
			QueryParameters: &applicationswithappid.ApplicationsWithAppIdRequestBuilderGetQueryParameters{
				Select: []string{"id", "keyCredentials"},
			},
		})
	if err != nil {
		return err
	}
	nextKeyCredentials := make([]gmodels.KeyCredentialable, 0, len(certIDs))
	hasPatch := false
	for _, installedKey := range application.GetKeyCredentials() {
		tp := strings.ToLower(base64.StdEncoding.EncodeToString(installedKey.GetCustomKeyIdentifier()))
		if _, hasValue := provisioningCerts[tp]; hasValue {
			nextKeyCredentials = append(nextKeyCredentials, installedKey)
			delete(provisioningCerts, tp)
		} else {
			hasPatch = true
		}
	}
	for _, certDoc := range provisioningCerts {
		kc := gmodels.NewKeyCredential()
		kc.SetKey(certDoc.KeySpec.CertificateChain[0])
		kc.SetUsage(to.Ptr("Verify"))
		kc.SetTypeEscaped(to.Ptr("AsymmetricX509Cert"))
		kc.SetStartDateTime(&certDoc.NotBefore.Time)
		kc.SetEndDateTime(&certDoc.NotAfter.Time)
		nextKeyCredentials = append(nextKeyCredentials, kc)
		hasPatch = true
	}

	if hasPatch {
		patchApplication := gmodels.NewApplication()
		patchApplication.SetKeyCredentials(nextKeyCredentials)
		_, err = gclient.Applications().ByApplicationId(*application.GetId()).Patch(c, patchApplication, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
