package profile

import (
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetSelfProfileDoc implements ProfileContextService.
func GetResourceProfileDoc(c RequestContext) (*ProfileDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	return getProfileDoc(c, resolveProfileLocatorFromNamespaceID(nsID))
}
