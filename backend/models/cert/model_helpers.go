package certmodels

import (
	"github.com/stephenzsy/small-kms/backend/models"
)

type (
	certificatePolicyComposed struct {
		models.Ref
		CertificatePolicyFields
	}
)
