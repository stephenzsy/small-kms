package key

import "github.com/stephenzsy/small-kms/backend/base"

type keySpecRefComposed struct {
	base.ResourceReference
	KeySpecRefFields
}

type keySpecComposed struct {
	KeySpecRef
	KeySpecFields
}
