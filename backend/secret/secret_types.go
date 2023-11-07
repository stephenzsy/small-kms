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

	secretRefComposed struct {
		base.ResourceReference
		SecretRefFields
	}

	secretComposed struct {
		SecretRef
		SecretFields
	}
)
