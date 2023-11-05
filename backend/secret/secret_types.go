package secret

import "github.com/stephenzsy/small-kms/backend/base"

type (
	secretPolicyRefComposed struct {
		base.ResourceReference
		SecretPolicyRefFields
	}

	secretPolicyComposed struct {
		SecretPolicyRef
		SecretPolicyFields
	}
)
