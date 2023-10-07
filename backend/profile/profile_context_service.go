package profile

import (
	"github.com/stephenzsy/small-kms/backend/common"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetSelfProfileDoc implements ProfileContextService.
func GetResourceProfileDoc(c common.ServiceContext) (*ProfileDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	return getProfileDoc(c, resolveProfileLocatorFromNamespaceID(nsID))
}
