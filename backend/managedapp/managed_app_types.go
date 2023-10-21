package managedapp

import (
	"github.com/stephenzsy/small-kms/backend/base"
)

type ResourceReference = base.ResourceReference

type managedAppComposed struct {
	ResourceReference
	ManagedAppFields
}
