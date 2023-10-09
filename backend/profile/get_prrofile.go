package profile

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type RequestContext = common.RequestContext

func getProfileDoc(c RequestContext, locator models.ResourceLocator) (doc *ProfileDoc, err error) {
	if locator.GetNamespaceID() == docNsIDProfileBuiltIn {
		docID := locator.GetID()
		if docID.Kind() == models.ResourceKindCaRoot {
			if a, ok := rootCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		} else if docID.Kind() == models.ResourceKindCaInt {
			if a, ok := intCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		}
	}
	doc = &ProfileDoc{}
	err = kmsdoc.Read(c, locator, doc)
	return
}

func resolveTenantProfileLocatorFromNamespaceID(nsID models.NamespaceID) models.ResourceLocator {
	return models.NewResourceLocator(docNsIDProfileTenant, common.NewIdentifierWithKind(models.ResourceKindMsGraph, nsID.Identifier()))
}

func resolveProfileLocatorFromNamespaceID(nsID models.NamespaceID) models.ResourceLocator {
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot:
		return models.NewResourceLocator(docNsIDProfileBuiltIn, common.NewIdentifierWithKind(models.ResourceKindCaRoot, nsID.Identifier()))
	case models.NamespaceKindCaInt:
		return models.NewResourceLocator(docNsIDProfileBuiltIn, common.NewIdentifierWithKind(models.ResourceKindCaInt, nsID.Identifier()))
	default:
		return resolveTenantProfileLocatorFromNamespaceID(nsID)
	}
}

// GetProfile implements ProfileService.
func GetProfile(c RequestContext) (*models.ProfileComposed, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	var profileNsID models.NamespaceID
	var resourceKind models.ResourceKind
	switch nsID.Kind() {
	case models.NamespaceKindCaRoot:
		resourceKind = models.ResourceKindCaRoot
		profileNsID = docNsIDProfileBuiltIn
	case models.NamespaceKindCaInt:
		resourceKind = models.ResourceKindCaInt
		profileNsID = docNsIDProfileBuiltIn
	case models.NamespaceKindProfile:
		return nil, fmt.Errorf("profile.GetProfile: invalid namespace kind: %s", nsID.Kind())
	default:
		resourceKind = models.ResourceKindMsGraph
		profileNsID = docNsIDProfileTenant
	}
	doc, err := getProfileDoc(c, models.NewResourceLocator(profileNsID, common.NewIdentifierWithKind(resourceKind, nsID.Identifier())))
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}
