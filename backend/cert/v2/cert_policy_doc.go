package cert

import (
	"crypto/md5"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/key"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type CertPolicyDoc struct {
	resdoc.ResourceDoc
	DisplayName string `json:"displayName"`

	KeySpec       keymodels.JsonWebKeySpec            `json:"keySpec"`
	KeyExportable bool                                `json:"keyExportable"`
	AllowGenerate bool                                `json:"allowGenerate"`
	AllowEnroll   bool                                `json:"allowEnroll"`
	ExpiryTime    caldur.CalendarDuration             `json:"expiryTime"`
	Subject       certmodels.CertificateSubject       `json:"subject"`
	SANs          *certmodels.SubjectAlternativeNames `json:"sans,omitempty"`
	Flags         []certmodels.CertificateFlag        `json:"flags,omitempty"`
	IssuerPolicy  resdoc.DocIdentifier                `json:"issuerPolicy"`

	Version []byte `json:"version"`
}

const (
	queryColumnDisplayName = "c.displayName"
)

func (d *CertPolicyDoc) init(
	p *certmodels.CreateCertificatePolicyRequest) error {
	if d == nil {
		return nil
	}

	d.DisplayName = d.ID
	if p.DisplayName != "" {
		d.DisplayName = p.DisplayName
	}

	nsProvider := d.PartitionKey.NamespaceProvider
	keySignVerifyOnly := false

	switch nsProvider {
	case models.NamespaceProviderRootCA:
		// TODO: verify key policy with the same ID exists
		keySignVerifyOnly = true

		d.ExpiryTime = caldur.CalendarDuration{
			Year: 10,
		}
		d.KeySpec = keymodels.JsonWebKeySpec{
			Kty:     cloudkey.KeyTypeRSA,
			KeySize: utils.ToPtr(4096),
		}
		d.KeyExportable = false
		d.AllowGenerate = true
		d.AllowEnroll = false
	case models.NamespaceProviderIntermediateCA:
		keySignVerifyOnly = true

		d.IssuerPolicy = resdoc.DocIdentifier{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderRootCA,
				NamespaceID:       "default",
				ResourceProvider:  models.ResourceProviderCertPolicy,
			},
			ID: "default",
		}
		if p.IssuerPolicyIdentifier != "" {
			parsed, err := resdoc.ParseIdentifier(p.IssuerPolicyIdentifier)
			if err != nil {
				return fmt.Errorf("%w: invalid issuer policy identifier: %s", base.ErrResponseStatusBadRequest, p.IssuerPolicyIdentifier)
			}
			if parsed.NamespaceProvider != models.NamespaceProviderRootCA {
				return fmt.Errorf("%w: issuer policy must be root ca", base.ErrResponseStatusBadRequest)
			}
			d.IssuerPolicy.NamespaceID = parsed.NamespaceID
			d.IssuerPolicy.ID = parsed.ID
		}
		// TODO: verify issuer policy

		d.ExpiryTime = caldur.CalendarDuration{
			Year: 3,
		}
		d.KeySpec = keymodels.JsonWebKeySpec{
			Kty:     cloudkey.KeyTypeRSA,
			KeySize: utils.ToPtr(4096),
		}
		d.KeyExportable = false
		d.AllowGenerate = true
		d.AllowEnroll = false
	case models.NamespaceProviderServicePrincipal,
		models.NamespaceProviderGroup:
		d.IssuerPolicy = resdoc.DocIdentifier{
			PartitionKey: resdoc.PartitionKey{
				NamespaceProvider: models.NamespaceProviderIntermediateCA,
				NamespaceID:       "default",
				ResourceProvider:  models.ResourceProviderCertPolicy,
			},
			ID: "default",
		}
		if p.IssuerPolicyIdentifier != "" {
			parsed, err := resdoc.ParseIdentifier(p.IssuerPolicyIdentifier)
			if err != nil {
				return fmt.Errorf("%w: invalid issuer policy identifier: %s", base.ErrResponseStatusBadRequest, p.IssuerPolicyIdentifier)
			}
			if parsed.NamespaceProvider != models.NamespaceProviderIntermediateCA {
				return fmt.Errorf("%w: issuer policy must be intermediate ca", base.ErrResponseStatusBadRequest)
			}
			d.IssuerPolicy.NamespaceID = parsed.NamespaceID
			d.IssuerPolicy.ID = parsed.ID
		}
		// TODO verify issuer policy

		d.KeyExportable = true
		d.AllowGenerate = true
		d.AllowEnroll = true

		if p.KeyExportable != nil {
			d.KeyExportable = *p.KeyExportable
		}
		if p.AllowGenerate != nil {
			d.AllowGenerate = *p.AllowGenerate
		}
		if p.AllowEnroll != nil {
			d.AllowEnroll = *p.AllowEnroll
		}
		if !d.AllowGenerate && !d.AllowEnroll {
			return fmt.Errorf("%w: certificate policy must allow generate or enroll", base.ErrResponseStatusBadRequest)
		}

		d.ExpiryTime = caldur.CalendarDuration{
			Year: 1,
		}
		d.KeySpec = keymodels.JsonWebKeySpec{
			Kty:     cloudkey.KeyTypeRSA,
			KeySize: utils.ToPtr(2048),
		}
		if len(p.Flags) == 0 {
			d.Flags = []certmodels.CertificateFlag{certmodels.CertificateFlagServerAuth, certmodels.CertificateFlagClientAuth}
		} else {
			if slices.Contains(p.Flags, certmodels.CertificateFlagServerAuth) {
				d.Flags = append(d.Flags, certmodels.CertificateFlagServerAuth)
			}
			if slices.Contains(p.Flags, certmodels.CertificateFlagClientAuth) {
				d.Flags = append(d.Flags, certmodels.CertificateFlagClientAuth)
			}
		}
		if len(d.Flags) == 0 {
			return fmt.Errorf("%w: certificate must have at least one usage flag", base.ErrResponseStatusBadRequest)
		}
	default:
		return fmt.Errorf("%w: unsupported namespace provider: %s", base.ErrResponseStatusBadRequest, nsProvider)
	}

	var pAlg cloudkey.JsonWebSignatureAlgorithm

	if p.KeySpec != nil {
		ks := *p.KeySpec
		pAlg = cloudkey.JsonWebSignatureAlgorithm(ks.Alg)
		switch ks.Kty {
		case cloudkey.KeyTypeRSA:
			d.KeySpec.Kty = cloudkey.KeyTypeRSA
			d.KeySpec.KeySize = utils.ToPtr(2048)
			d.KeySpec.Crv = ""
			if ks.KeySize != nil {
				switch *ks.KeySize {
				case 2048, 3072, 4096:
					d.KeySpec.KeySize = ks.KeySize
				}
				// any other value will be using default
			}
		case cloudkey.KeyTypeEC:
			d.KeySpec.Kty = cloudkey.KeyTypeEC
			d.KeySpec.Crv = cloudkey.CurveNameP384
			switch ks.Crv {
			case cloudkey.CurveNameP256,
				cloudkey.CurveNameP384,
				cloudkey.CurveNameP521:
				d.KeySpec.Crv = ks.Crv
				// other values will use default
			}
		default:
			// other values use default
		}
		if !keySignVerifyOnly && len(ks.KeyOperations) > 0 {
			d.KeySpec.KeyOperations = ks.KeyOperations
			if !slices.Contains(ks.KeyOperations, cloudkey.JsonWebKeyOperationSign) ||
				!slices.Contains(ks.KeyOperations, cloudkey.JsonWebKeyOperationVerify) {
				return fmt.Errorf("%w: key operations must include sign and verify", base.ErrResponseStatusBadRequest)
			}
			d.KeySpec.KeyOperations = cloudkey.SanitizeKeyOperations(ks.KeyOperations)

		}
	}

	if len(d.KeySpec.KeyOperations) == 0 {
		d.KeySpec.KeyOperations = []cloudkey.JsonWebKeyOperation{
			cloudkey.JsonWebKeyOperationSign,
			cloudkey.JsonWebKeyOperationVerify,
		}
		if !keySignVerifyOnly {
			switch d.KeySpec.Kty {
			case cloudkey.KeyTypeRSA:
				d.KeySpec.KeyOperations = append(d.KeySpec.KeyOperations,
					cloudkey.JsonWebKeyOperationWrapKey, cloudkey.JsonWebKeyOperationUnwrapKey)
			case cloudkey.KeyTypeEC:
				if d.KeyExportable {
					d.KeySpec.KeyOperations = append(d.KeySpec.KeyOperations,
						cloudkey.JsonWebKeyOperationDeriveKey, cloudkey.JsonWebKeyOperationDeriveBits)
				}
			}
		}
	}

	switch d.KeySpec.Kty {
	case cloudkey.KeyTypeRSA:
		d.KeySpec.Alg = string(cloudkey.SignatureAlgorithmRS256)
		if (*d.KeySpec.KeySize) >= 3072 {
			d.KeySpec.Alg = string(cloudkey.SignatureAlgorithmRS384)
		}
		switch pAlg {
		case cloudkey.SignatureAlgorithmRS256,
			cloudkey.SignatureAlgorithmRS384,
			cloudkey.SignatureAlgorithmRS512,
			cloudkey.SignatureAlgorithmPS256,
			cloudkey.SignatureAlgorithmPS384,
			cloudkey.SignatureAlgorithmPS512:
			d.KeySpec.Alg = string(pAlg)
		}
	case cloudkey.KeyTypeEC:
		switch d.KeySpec.Crv {
		case cloudkey.CurveNameP256:
			d.KeySpec.Alg = string(cloudkey.SignatureAlgorithmES256)
		case cloudkey.CurveNameP384:
			d.KeySpec.Alg = string(cloudkey.SignatureAlgorithmES384)
		case cloudkey.CurveNameP521:
			d.KeySpec.Alg = string(cloudkey.SignatureAlgorithmES512)
		}
	}

	baseTime := time.Now().UTC()
	expMax := baseTime.AddDate(10, 0, 0)
	expMin := baseTime.AddDate(0, 0, 28)
	if p.ExpiryTime != "" {
		expTime, err := caldur.Parse(p.ExpiryTime)
		if err != nil {
			return fmt.Errorf("%w: invalid expiry time format", base.ErrResponseStatusBadRequest)
		}
		if caldur.Shift(baseTime, expTime).After(expMax) {
			return fmt.Errorf("%w: expiry time cannot be more than 10 years", base.ErrResponseStatusBadRequest)
		} else if caldur.Shift(baseTime, expTime).Before(expMin) {
			return fmt.Errorf("%w: expiry time cannot be less than 28 days", base.ErrResponseStatusBadRequest)
		} else {
			d.ExpiryTime = expTime
		}
	}

	d.Subject = p.Subject
	if d.Subject.CommonName == "" {
		return fmt.Errorf("%w: subject common name cannot be empty", base.ErrResponseStatusBadRequest)
	}

	d.SANs = p.SubjectAlternativeNames.Sanitize()

	// get checksum of key fields
	dw := md5.New()
	d.KeySpec.Digest(dw)
	if d.KeyExportable {
		dw.Write([]byte("keyExportable"))
	}
	if d.AllowGenerate {
		dw.Write([]byte("allowGenerate"))
	}
	if d.AllowEnroll {
		dw.Write([]byte("allowEnroll"))
	}
	io.WriteString(dw, d.IssuerPolicy.String())
	dw.Write(d.ExpiryTime.Bytes())
	io.WriteString(dw, d.Subject.String())
	d.SANs.Digest(dw)
	for _, flag := range d.Flags {
		dw.Write([]byte(flag))
	}
	d.Version = dw.Sum(nil)

	return nil
}

// populate ref
func (d *CertPolicyDoc) ToRef() (m models.Ref) {
	m = d.ResourceDoc.ToRef()
	m.DisplayName = &d.DisplayName
	return m
}

func (d *CertPolicyDoc) ToModel() (m certmodels.CertificatePolicy) {
	m.Ref = d.ToRef()
	m.KeySpec = d.KeySpec
	if m.KeySpec.KeyOperations == nil {
		m.KeySpec.KeyOperations = []key.JsonWebKeyOperation{}
	}
	m.KeyExportable = d.KeyExportable
	m.AllowGenerate = d.AllowGenerate
	m.AllowEnroll = d.AllowEnroll
	m.ExpiryTime = d.ExpiryTime.String()
	m.Subject = d.Subject
	m.SubjectAlternativeNames = d.SANs
	m.Flags = d.Flags
	m.IssuerPolicyIdentifier = d.IssuerPolicy.String()
	if m.IssuerPolicyIdentifier == "" {
		m.IssuerPolicyIdentifier = "self"
	}
	return m
}

func (d *CertPolicyDoc) getIssuerCert(c ctx.RequestContext) (*CertDoc, error) {
	linkDoc, err := getPolicyIssuerCertInternal(c, d.PartitionKey.NamespaceProvider, d.PartitionKey.NamespaceID, d.ID)
	if err != nil {
		return nil, err
	}

	return getCertificateInternal(c, linkDoc.LinkTo.NamespaceProvider, linkDoc.LinkTo.NamespaceID, linkDoc.LinkTo.ID)
}
