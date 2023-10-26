package agentclient

import "github.com/stephenzsy/small-kms/backend/base"

type (
	certificateRefComposed struct {
		base.ResourceReference
		CertificateRefFields
	}

	certificateComposed struct {
		CertificateRef
		CertificateFields
	}
)
