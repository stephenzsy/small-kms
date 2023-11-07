package secret

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type SecretPolicyDoc struct {
	base.BaseDoc

	DisplayName          string                      `json:"displayName"`
	Mode                 SecretGenerateMode          `json:"mode"`
	ExpiryTime           *base.Period                `json:"expiryTime,omitempty"`
	RandomCharacterClass *SecretRandomCharacterClass `json:"randomCharacterClass,omitempty"`
	RandomLength         *int                        `json:"randomLength,omitempty"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

func (d *SecretPolicyDoc) init(nsKind base.NamespaceKind, nsIdentifier base.ID, identifier base.ID) {
	d.BaseDoc.Init(nsKind, nsIdentifier, base.ResourceKindSecretPolicy, identifier)
}

func (d *SecretPolicyDoc) GetID() base.ID {
	return d.ID
}

func (d *SecretPolicyDoc) PopulateModelRef(r *SecretPolicyRef) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&r.ResourceReference)
	r.DisplayName = d.DisplayName
}

func (d *SecretPolicyDoc) PopulateModel(r *SecretPolicy) {
	if d == nil || r == nil {
		return
	}
	d.PopulateModelRef(&r.SecretPolicyRef)
	r.Mode = d.Mode
	r.ExpiryTime = d.ExpiryTime
	r.RandomCharacterClass = d.RandomCharacterClass
	r.RandomLength = d.RandomLength
}

func (d *SecretPolicyDoc) populateByRequestParams(p SecretPolicyParameters) error {

	if p.DisplayName != "" {
		d.DisplayName = p.DisplayName
	} else {
		d.DisplayName = string(d.ID)
	}

	switch p.Mode {
	case SecretGenerateModeManual:
		d.Mode = SecretGenerateModeManual
		return nil
	case SecretGenerateModeServerGeneratedRandom:
		d.Mode = SecretGenerateModeServerGeneratedRandom
	default:
		return fmt.Errorf("%w, invalid mode: %s", base.ErrResponseStatusBadRequest, p.Mode)
	}

	if p.ExpiryTime != nil {
		now := time.Now()
		if base.AddPeriod(now, *p.ExpiryTime).Before(now.AddDate(0, 0, 28)) {
			return fmt.Errorf("%w, expiry time must be at least 28 days", base.ErrResponseStatusBadRequest)
		}
		d.ExpiryTime = p.ExpiryTime
	}

	if p.RandomCharacterClass == nil {
		d.RandomCharacterClass = utils.ToPtr(SecretRandomCharClassBase64RawURL)
	} else {
		switch *p.RandomCharacterClass {
		case SecretRandomCharClassBase64RawURL:
			// ok
		default:
			return fmt.Errorf("%w, invalid random character class: %s", base.ErrResponseStatusBadRequest, *p.RandomCharacterClass)
		}
	}

	if p.RandomLength == nil {
		return fmt.Errorf("%w, missing random length", base.ErrResponseStatusBadRequest)
	}
	if *p.RandomLength < 8 || *p.RandomLength > 1024 {
		return fmt.Errorf("%w, random length must be between 8 and 1024", base.ErrResponseStatusBadRequest)
	}
	d.RandomLength = p.RandomLength
	return nil
}

func apiPutSecretPolicy(c ctx.RequestContext, policyIdentifier base.ID, p SecretPolicyParameters) error {
	nsCtx := ns.GetNSContext(c)
	doc := &SecretPolicyDoc{}
	doc.init(nsCtx.Kind(), nsCtx.ID(), policyIdentifier)
	if err := doc.populateByRequestParams(p); err != nil {
		return err
	}
	if err := base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil); err != nil {
		return err
	}
	model := &SecretPolicy{}
	doc.PopulateModel(model)
	return c.JSON(http.StatusOK, model)
}

func apiGetSecretPolicy(c ctx.RequestContext, policyIdentifier base.ID) error {
	nsCtx := ns.GetNSContext(c)
	doc := &SecretPolicyDoc{}
	if err := base.GetAzCosmosCRUDService(c).Read(c,
		base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindSecretPolicy, policyIdentifier),
		doc,
		nil); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, secret policy not found: %s", base.ErrResponseStatusNotFound, policyIdentifier)
		}
		return err
	}
	model := &SecretPolicy{}
	doc.PopulateModel(model)
	return c.JSON(http.StatusOK, model)
}

func apiListSecretPolicies(c ctx.RequestContext) error {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(queryColumnDisplayName)
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindSecretPolicy)
	pager := base.NewQueryDocPager[*SecretPolicyDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *SecretPolicyDoc) *SecretPolicyRef {
		r := &SecretPolicyRef{}
		d.PopulateModelRef(r)
		return r
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
