package profile

import "github.com/stephenzsy/small-kms/backend/base"

type profileRefComposed struct {
	base.ResourceReference
	ProfileRefFields
}
