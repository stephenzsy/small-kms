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
	keyParameters := KeyParameters{
		Kty:  KtyRSA,
		Size: ToPtr(KeySize2048),
	}
	if IsRootCANamespace(p.NamespaceID) {
		// override issuer namespace for root ca
		r.CertRequest.IssuerNamespaceID = p.NamespaceID
		if p.CertRequest.Validity.Years == 0 && p.CertRequest.Validity.Months == 0 && p.CertRequest.Validity.Days == 0 {
			r.CertRequest.Validity.Years = 10
		}
		keyParameters.Size = ToPtr(KeySize4096)
		r.CertRequest.Usage = UsageRootCA
	}

	// validate key spec
	switch p.CertRequest.KeyParameters.Kty {
	case KtyEC:
		keyParameters.Kty = KtyEC
		keyParameters.Size = nil
		curve := *p.CertRequest.KeyParameters.Curve
		switch curve {
		case EcCurveP256, EcCurveP384:
			keyParameters.Curve = &curve
		default:
			keyParameters.Curve = ToPtr(EcCurveP384)
		}
	case KtyRSA:
		size := *p.CertRequest.KeyParameters.Size
		switch size {
		case KeySize3072, KeySize4096:
			keyParameters.Size = &size
		}
	}
	r.CertRequest.KeyParameters = &keyParameters

	return &r, nil
}
