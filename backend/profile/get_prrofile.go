package profile

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func getProfileDoc(c RequestContext, locator shared.ResourceLocator) (doc *ProfileDoc, err error) {
	if locator.GetNamespaceID() == docNsIDProfileBuiltIn {
		docID := locator.GetID()
		if docID.Kind() == shared.ResourceKindCaRoot {
			if a, ok := rootCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		} else if docID.Kind() == shared.ResourceKindCaInt {
			if a, ok := intCaProfileDocs[docID.Identifier()]; ok {
				return &a, nil
			}
			return nil, common.ErrStatusNotFound
		} else if docID.Kind() == shared.ResourceKindReserved {
			if docID.Identifier().String() == "agent-push" {
				return &ProfileDoc{
					BaseDoc: kmsdoc.BaseDoc{
						NamespaceID: locator.GetNamespaceID(),
						ID:          locator.GetID(),
					},
					DispalyName: "Agent push",
					ProfileType: shared.NamespaceKindSystem,
				}, nil
			}
			return nil, common.ErrStatusNotFound
		}
	}
	doc = &ProfileDoc{}
	err = kmsdoc.Read(c, locator, doc)
	return
}

func resolveTenantProfileLocatorFromNamespaceID(nsID shared.NamespaceIdentifier) shared.ResourceLocator {
	return shared.NewResourceLocator(docNsIDProfileTenant, shared.NewResourceIdentifier(shared.ResourceKindMsGraph, nsID.Identifier()))
}

func resolveProfileLocatorFromNamespaceID(nsID shared.NamespaceIdentifier) shared.ResourceLocator {
	switch nsID.Kind() {
	case shared.NamespaceKindSystem:
		return shared.NewResourceLocator(docNsIDProfileBuiltIn, shared.NewResourceIdentifier(shared.ResourceKindReserved, nsID.Identifier()))
	case shared.NamespaceKindCaRoot:
		return shared.NewResourceLocator(docNsIDProfileBuiltIn, shared.NewResourceIdentifier(shared.ResourceKindCaRoot, nsID.Identifier()))
	case shared.NamespaceKindCaInt:
		return shared.NewResourceLocator(docNsIDProfileBuiltIn, shared.NewResourceIdentifier(shared.ResourceKindCaInt, nsID.Identifier()))
	default:
		return resolveTenantProfileLocatorFromNamespaceID(nsID)
	}
}

// GetProfile implements ProfileService.
func GetProfile(c RequestContext) (*models.ProfileComposed, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	var profileNsID shared.NamespaceIdentifier
	var resourceKind shared.ResourceKind
	switch nsID.Kind() {
	case shared.NamespaceKindSystem:
		resourceKind = shared.ResourceKindReserved
		profileNsID = docNsIDProfileBuiltIn
	case shared.NamespaceKindCaRoot:
		resourceKind = shared.ResourceKindCaRoot
		profileNsID = docNsIDProfileBuiltIn
	case shared.NamespaceKindCaInt:
		resourceKind = shared.ResourceKindCaInt
		profileNsID = docNsIDProfileBuiltIn
	case shared.NamespaceKindProfile:
		return nil, fmt.Errorf("profile.GetProfile: invalid namespace kind: %s", nsID.Kind())
	default:
		resourceKind = shared.ResourceKindMsGraph
		profileNsID = docNsIDProfileTenant
	}
	doc, err := getProfileDoc(c, shared.NewResourceLocator(profileNsID, shared.NewResourceIdentifier(resourceKind, nsID.Identifier())))
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}
