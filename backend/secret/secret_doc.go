package secret

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
)

type SecretKeyVaultStore struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type SecretDoc struct {
	base.BaseDoc

	Policy        base.DocLocator     `json:"policy"`
	Created       base.NumericDate    `json:"iat"`
	NotBefore     *base.NumericDate   `json:"nbf,omitempty"`
	NotAfter      *base.NumericDate   `json:"exp,omitempty"`
	KeyVaultStore SecretKeyVaultStore `json:"keyVaultStore"`
}

const (
	secretDocQueryColumnVersion = "c.version"
	secretDocQueryColumnCreated = "c.created"
)

func (d *SecretDoc) PopulateModelRef(r *SecretRef) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&r.ResourceReference)
	r.Iat = d.Created
	r.Exp = d.NotAfter
}

// PopulateModel implements base.ModelPopulater.
func (d *SecretDoc) PopulateModel(r *Secret) {
	if d == nil || r == nil {
		return
	}
	d.PopulateModelRef(&r.SecretRef)
	r.ContentType = "text/plain"
	r.Sid = string(d.KeyVaultStore.ID)
}

func GetKeyStoreName(nsKind base.NamespaceKind, nsID base.ID, policyID base.ID) string {
	return fmt.Sprintf("s-%s-%s-%s", nsKind, nsID, policyID)
}

func (d *SecretDoc) init(
	nsKind base.NamespaceKind,
	nsID base.ID,
	pDoc *SecretPolicyDoc) error {
	secretUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindSecret, base.IDFromUUID(secretUUID))
	d.Policy = pDoc.GetStorageFullIdentifier()
	d.KeyVaultStore.Name = GetKeyStoreName(nsKind, nsID, pDoc.ID)
	return nil
}

var _ base.ModelPopulater[Secret] = (*SecretDoc)(nil)
