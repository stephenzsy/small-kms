package certda

import certmodels "github.com/stephenzsy/small-kms/backend/models/cert"

type CertDTO interface {
	Status() certmodels.CertificateStatus
}

type CertDTOPendingAuthorization interface {
	CertDTO
}

type certDTOImpl struct {
}
