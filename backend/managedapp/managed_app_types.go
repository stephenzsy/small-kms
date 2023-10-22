package managedapp

import (
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/profile/v2"
)

type ResourceReference = base.ResourceReference

type managedAppRefComposed struct {
	profile.ProfileRef
	ManagedAppRefFields
}
