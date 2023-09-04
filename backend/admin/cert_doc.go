package admin

import (
	"time"

	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertDoc struct {
	kmsdoc.BaseDoc

	// associated policy id
	PolicyID uuid.UUID `json:"policyId"`
	Expires  time.Time `json:"expires"`
	// alias for certs with L prefix
	AliasID *kmsdoc.KmsDocID `json:"aliasId,omitempty"`
	Usage   CertificateUsage `json:"usage"`

	// keyvault certificate id
	CID string `json:"cid"`
	// keyvault certificate key id
	KID string `json:"kid"`
	// keyvault certificate secret id
	SID string `json:"sid"`
	// certificate storage path in blob storage
	CertStorePath string `json:"certStorePath"`
}
