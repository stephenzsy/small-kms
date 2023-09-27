package admin

type PolicyCertEnrollDocSection struct {
	MaxValidityInMonths int32              `json:"maxValidityInMonths"`
	AllowedUsages       []CertificateUsage `json:"allowedUsages"`
}

func (d *PolicyCertEnrollDocSection) toCertificateEnrollPolicyParameters() *CertificateEnrollPolicyParameters {
	if d == nil {
		return nil
	}
	return &CertificateEnrollPolicyParameters{
		MaxValidityInMonths: d.MaxValidityInMonths,
		AllowedUsages:       d.AllowedUsages,
	}
}
