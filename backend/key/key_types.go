package key

import "github.com/stephenzsy/small-kms/backend/base"

type keyPolicyRefComposed struct {
	base.ResourceReference
	KeyPolicyRefFields
}

type keyPolicyComposed struct {
	KeyPolicyRef
	KeyPolicyFields
}
