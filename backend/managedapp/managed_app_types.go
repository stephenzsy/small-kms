package managedapp

import (
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile"
)

type ResourceReference = base.ResourceReference

type managedAppRefComposed struct {
	profile.ProfileRef
	ManagedAppRefFields
}
