package admin

import (
	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyCertIssueDocSection struct {
	IssuerID            kmsdoc.KmsDocID    `json:"issuerID"`
	AllowedRequesters   []uuid.UUID        `json:"allowedRequesters"`
	AllowedUsages       []CertificateUsage `json:"allowedUsages"`
	MaxValidityInMonths int32              `json:"max_validity_months"`
}
