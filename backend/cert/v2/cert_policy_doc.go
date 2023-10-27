package cert

import (
	"context"
	"crypto/md5"
	"fmt"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/key"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertPolicyDoc struct {
	base.BaseDoc

	DisplayName     string                   `json:"displayName"`
	KeySpec         key.SigningKeySpec       `json:"keySpec"`
	KeyExportable   bool                     `json:"keyExportable"`
	ExpiryTime      base.Period              `json:"expiryTime"`
	LifetimeAction  *key.LifetimeAction      `json:"lifetimeActions,omitempty"`
	Subject         CertificateSubject       `json:"subject"`
	SANs            *SubjectAlternativeNames `json:"sans,omitempty"`
	Flags           []CertificateFlag        `json:"flags"`
	Version         HexDigest                `json:"version"`
	IssuerNamespace base.NamespaceIdentifier `json:"issuerNamespace"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

func (d *CertPolicyDoc) Init(
	c context.Context,
	nsKind base.NamespaceKind,
	nsID base.Identifier,
	policyID base.Identifier,
	p *CertPolicyParameters) error {
	if d == nil {
		return nil
	}
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindCertPolicy, policyID)

	d.DisplayName = policyID.String()
	if p.DisplayName != nil && *p.DisplayName != "" {
		d.DisplayName = *p.DisplayName
	}

	digest := md5.New()

	isSelfSigning := nsKind == base.NamespaceKindRootCA

	switch nsKind {
	case base.NamespaceKindRootCA:
		if !isSelfSigning {
			return fmt.Errorf("%w: issuer namespace must be Self", base.ErrResponseStatusBadRequest)
		}
		d.IssuerNamespace = base.NewNamespaceIdentifier(nsKind, nsID)
	case base.NamespaceKindIntermediateCA:
		if p.IssuerNamespaceKind == nil || p.IssuerNamespaceIdentifier == nil {
			d.IssuerNamespace = base.NewNamespaceIdentifier(base.NamespaceKindRootCA, base.StringIdentifier("default"))
		} else if *p.IssuerNamespaceKind != base.NamespaceKindRootCA {
			return fmt.Errorf("%w: issuer namespace must be root ca", base.ErrResponseStatusBadRequest)
		} else {
			d.IssuerNamespace = base.NewNamespaceIdentifier(*p.IssuerNamespaceKind, *p.IssuerNamespaceIdentifier)
		}
	case base.NamespaceKindServicePrincipal:
		if isSelfSigning {
			d.IssuerNamespace = base.NewNamespaceIdentifier(nsKind, nsID)
		} else if p.IssuerNamespaceKind == nil || p.IssuerNamespaceIdentifier == nil {
			d.IssuerNamespace = base.NewNamespaceIdentifier(base.NamespaceKindIntermediateCA, base.StringIdentifier("default"))
		} else if *p.IssuerNamespaceKind != base.NamespaceKindIntermediateCA {
			return fmt.Errorf("%w: issuer namespace must be intermediate ca", base.ErrResponseStatusBadRequest)
		} else {
			d.IssuerNamespace = base.NewNamespaceIdentifier(*p.IssuerNamespaceKind, *p.IssuerNamespaceIdentifier)
		}
	default:
		return fmt.Errorf("%w: unsupported namespace kind: %s", base.ErrResponseStatusBadRequest, nsKind)
	}
	digest.Write(([]byte)(d.IssuerNamespace.String()))

	if p.KeySpec == nil {
		d.KeySpec = key.SigningKeySpec{
			Alg:     to.Ptr(key.JsonWebKeySignatureAlgorithmRS256),
			Kty:     key.JsonWebKeyTypeRSA,
			KeySize: utils.ToPtr(int32(2048)),
			KeyOperations: []key.JsonWebKeyOperation{
				key.JsonWebKeyOperationSign,
				key.JsonWebKeyOperationVerify,
			},
		}
	} else {
		ks := *p.KeySpec
		switch ks.Kty {
		case key.JsonWebKeyTypeEC:
			d.KeySpec.Kty = key.JsonWebKeyTypeEC
			d.KeySpec.Crv = utils.ToPtr(key.JsonWebKeyCurveNameP384)
			d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmES384)
			if ks.Crv != nil {
				switch *ks.Crv {
				case key.JsonWebKeyCurveNameP256:
					d.KeySpec.Crv = ks.Crv
					d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmES256)
				case key.JsonWebKeyCurveNameP384:
					d.KeySpec.Crv = ks.Crv
					d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmES384)
				case key.JsonWebKeyCurveNameP521:
					d.KeySpec.Crv = ks.Crv
					d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmES512)
				}
			}
		case key.JsonWebKeyTypeRSA:
			d.KeySpec.Kty = key.JsonWebKeyTypeRSA
			d.KeySpec.KeySize = utils.ToPtr(int32(2048))
			if ks.KeySize != nil {
				switch *ks.KeySize {
				case 2048, 3072, 4096:
					d.KeySpec.KeySize = ks.KeySize
				}
				// any other value will be using default
			}
			if d.KeySpec.Alg == nil {
				d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmRS256)
				if *ks.KeySize >= 3072 {
					d.KeySpec.Alg = to.Ptr(key.JsonWebKeySignatureAlgorithmRS384)
				}
			} else {
				switch *d.KeySpec.Alg {
				case key.JsonWebKeySignatureAlgorithmRS256,
					key.JsonWebKeySignatureAlgorithmRS384,
					key.JsonWebKeySignatureAlgorithmRS512,
					key.JsonWebKeySignatureAlgorithmPS256,
					key.JsonWebKeySignatureAlgorithmPS384,
					key.JsonWebKeySignatureAlgorithmPS512:
					// ok
				default:
					return fmt.Errorf("%w: unsupported key signature algorithm: %s", base.ErrResponseStatusBadRequest, *d.KeySpec.Alg)
				}
			}

		default:
			return fmt.Errorf("%w: unsupported key type: %s", base.ErrResponseStatusBadRequest, ks.Kty)
		}
		if len(ks.KeyOperations) == 0 {
			d.KeySpec.KeyOperations = []key.JsonWebKeyOperation{
				key.JsonWebKeyOperationSign,
				key.JsonWebKeyOperationVerify,
			}
		} else {
			d.KeySpec.KeyOperations = ks.KeyOperations
		}
	}
	d.KeySpec.WriteToDigest(digest)

	if p.KeyExportable == nil {
		switch nsKind {
		default:
			d.KeyExportable = false
		}
	} else if *p.KeyExportable {
		d.KeyExportable = true
	}
	digest.Write([]byte(fmt.Sprintf("%t", d.KeyExportable)))

	baseTime := time.Now().UTC()
	expMaxCutoff := baseTime.AddDate(10, 0, 0)
	expMinCutoff := baseTime.AddDate(0, 0, 28)
	if base.AddPeriod(baseTime, p.ExpiryTime).After(expMaxCutoff) {
		return fmt.Errorf("%w: expiry time cannot be more than 10 years", base.ErrResponseStatusBadRequest)
	} else if base.AddPeriod(baseTime, p.ExpiryTime).Before(expMinCutoff) {
		return fmt.Errorf("%w: expiry time cannot be less than 28 days", base.ErrResponseStatusBadRequest)
	} else {
		d.ExpiryTime = p.ExpiryTime
	}
	digest.Write(d.ExpiryTime.Bytes())

	if p.LifetimeAction != nil {
		d.LifetimeAction = p.LifetimeAction
		if d.LifetimeAction.Trigger.TimeAfterCreate != nil {
			timeAfterCutoff := baseTime.AddDate(0, 0, 14)
			if base.AddPeriod(baseTime, *d.LifetimeAction.Trigger.TimeAfterCreate).Before(timeAfterCutoff) {
				return fmt.Errorf("%w: lifetime action trigger after creation cannot be less than 14 days", base.ErrResponseStatusBadRequest)
			}
			if base.AddPeriod(baseTime, *d.LifetimeAction.Trigger.TimeAfterCreate).After(base.AddPeriod(baseTime, d.ExpiryTime)) {
				return fmt.Errorf("%w: lifetime action trigger after creation cannot be after expiry time", base.ErrResponseStatusBadRequest)
			}
		}
		if d.LifetimeAction.Trigger.TimeBeforeExpiry != nil {
			exp := base.AddPeriod(baseTime, d.ExpiryTime)
			timeBeforeMinCutoff := baseTime.AddDate(0, 0, 14)
			if base.AddPeriod(timeBeforeMinCutoff, *d.LifetimeAction.Trigger.TimeBeforeExpiry).After(exp) {
				return fmt.Errorf("%w: lifetime action trigger before expiry cannot be less than 14 days after creation", base.ErrResponseStatusBadRequest)
			}
		}
		if d.LifetimeAction.Trigger.PercentageAfterCreate != nil {
			if *d.LifetimeAction.Trigger.PercentageAfterCreate < 1 || *d.LifetimeAction.Trigger.PercentageAfterCreate > 99 {
				return fmt.Errorf("%w: lifetime action trigger percentage after creation must be between 1 and 99", base.ErrResponseStatusBadRequest)
			}

			duration := base.AddPeriod(baseTime, d.ExpiryTime).Sub(baseTime)
			if baseTime.Add(duration * time.Duration(*d.LifetimeAction.Trigger.PercentageAfterCreate) / 100).Before(baseTime.AddDate(0, 0, 14)) {
				return fmt.Errorf("%w: lifetime action trigger percentage after creation cannot be less than 14 days", base.ErrResponseStatusBadRequest)
			}
		}
		d.LifetimeAction.WriteToDigest(digest)
	}

	d.Subject = p.Subject
	if d.Subject.CommonName == "" {
		return fmt.Errorf("%w: subject common name cannot be empty", base.ErrResponseStatusBadRequest)
	}
	digest.Write([]byte(d.Subject.CommonName))

	d.SANs = p.SubjectAlternativeNames.Sanitize()
	d.SANs.WriteToDigest(digest)

	switch nsKind {
	case base.NamespaceKindRootCA:
		d.Flags = []CertificateFlag{CertificateFlagCA, CertificateFlagRootCA}
	case base.NamespaceKindIntermediateCA:
		d.Flags = []CertificateFlag{CertificateFlagCA}
	default:
		d.Flags = make([]CertificateFlag, 0, 2)
		if len(p.Flags) == 0 {
			d.Flags = append(d.Flags, CertificateFlagServerAuth, CertificateFlagClientAuth)
		} else {
			if slices.Contains(p.Flags, CertificateFlagServerAuth) {
				d.Flags = append(d.Flags, CertificateFlagServerAuth)
			}
			if slices.Contains(p.Flags, CertificateFlagClientAuth) {
				d.Flags = append(d.Flags, CertificateFlagClientAuth)
			}
		}
	}
	if len(d.Flags) == 0 {
		return fmt.Errorf("%w: certificate must have at least one usage flag", base.ErrResponseStatusBadRequest)
	}
	for _, flag := range d.Flags {
		digest.Write([]byte(flag))
	}
	d.Version = digest.Sum(nil)

	return nil
}

// populate ref
func (d *CertPolicyDoc) PopulateModelRef(m *CertPolicyRef) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.DisplayName = d.DisplayName
}

func (d *CertPolicyDoc) PopulateModel(m *CertPolicy) {
	if d == nil || m == nil {
		return
	}
	d.PopulateModelRef(&m.CertPolicyRef)
	m.KeySpec = d.KeySpec
	if m.KeySpec.KeyOperations == nil {
		m.KeySpec.KeyOperations = []key.JsonWebKeyOperation{}
	}
	m.KeyExportable = d.KeyExportable
	m.ExpiryTime = d.ExpiryTime
	m.LifetimeAction = d.LifetimeAction
	m.Subject = d.Subject
	m.SubjectAlternativeNames = d.SANs
	m.Flags = d.Flags
	m.Version = d.Version
	m.IssuerNamespaceKind = d.IssuerNamespace.Kind()
	m.IssuerNamespaceIdentifier = d.IssuerNamespace.Identifier()
}

var _ base.ModelRefPopulater[CertPolicyRef] = (*CertPolicyDoc)(nil)
var _ base.ModelPopulater[CertPolicy] = (*CertPolicyDoc)(nil)
