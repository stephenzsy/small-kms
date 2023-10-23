package cert

import (
	"github.com/stephenzsy/small-kms/backend/base"
)

type (
	certPolicyRefComposed struct {
		base.ResourceReference
		CertPolicyRefFields
	}

	certPolicyComposed struct {
		CertPolicyRef
		CertPolicyFields
	}
)
