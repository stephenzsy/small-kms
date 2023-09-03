package admin

import (
	"errors"
	"fmt"
)

type PolicyDBItem struct {
	Policy
}

type CertIssuePolicy interface {
	MaxValidity() DurationSpec
}

type certIssuePolicy struct {
	PolicyDBItem
}

func (p *certIssuePolicy) MaxValidity() DurationSpec {
	return p.PolicyDBItem.CertIssue.MaxValidity
}

func (p *PolicyDBItem) ToCertIssuePolicy() (CertIssuePolicy, error) {
	if p.Type != PolicyTypeCertIssue {
		return nil, fmt.Errorf("mismatched policy type: %s, expected: %s", p.Type, PolicyTypeCertIssue)
	}
	if p.CertIssue == nil {
		return nil, errors.New("missing CertIssue property")
	}
	result := certIssuePolicy{
		PolicyDBItem: *p,
	}
	var err error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type CertRequestPolicy interface {
	DBItem() *PolicyDBItem
}

type certRequestPolicy struct {
	PolicyDBItem
}

func (item *PolicyDBItem) DBItem() *PolicyDBItem {
	return item
}

func (p *PolicyDBItem) ToCertRequestPolicy() (CertRequestPolicy, error) {
	if p.Type != PolicyTypeCertRequest {
		return nil, fmt.Errorf("mismatched policy type: %s, expected: %s", p.Type, PolicyTypeCertRequest)
	}
	if p.CertRequest == nil {
		return nil, errors.New("missing CertRequest property")
	}
	r := certRequestPolicy{
		*p,
	}
	r.CertIssue = nil
	keyProps := KeyProperties{
		KeyType: KtyRSA,
		KeySize: ToPtr(KeySize2048),
	}
	validityInMonths := 12
	if IsRootCANamespace(p.NamespaceID) {
		// override issuer namespace for root ca
		r.CertRequest.IssuerNamespaceID = p.NamespaceID
		if p.CertRequest.ValidityInMonths == nil || *p.CertRequest.ValidityInMonths == 0 {
			validityInMonths = 120
		}
		keyProps.KeySize = ToPtr(KeySize4096)
		r.CertRequest.Usage = UsageRootCA
	}

	// validate key spec
	if p.CertRequest.KeyProperties != nil {
		switch p.CertRequest.KeyProperties.KeyType {
		case KtyEC:
			keyProps.KeyType = KtyEC
			keyProps.KeySize = nil
			curve := *p.CertRequest.KeyProperties.CurveName
			switch curve {
			case EcCurveP256, EcCurveP384:
				keyProps.CurveName = &curve
			default:
				keyProps.CurveName = ToPtr(EcCurveP384)
			}
		case KtyRSA:
			size := *p.CertRequest.KeyProperties.KeySize
			switch size {
			case KeySize3072, KeySize4096:
				keyProps.KeySize = &size
			}
		}

		if p.CertRequest.ValidityInMonths != nil && *p.CertRequest.ValidityInMonths > 0 && *p.CertRequest.ValidityInMonths < 120 {
			validityInMonths = *p.CertRequest.ValidityInMonths
		}
	}
	r.CertRequest.KeyProperties = &keyProps
	r.CertRequest.ValidityInMonths = &validityInMonths

	return &r, nil
}
