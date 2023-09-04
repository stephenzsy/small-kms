package admin

import "crypto/x509/pkix"

// Ptr returns a pointer to the provided value.
func ToPtr[T any](v T) *T {
	return &v
}

func (s *CertificateSubject) ToPkixName() pkix.Name {
	caSubjectOU := []string{}
	caSubjectO := []string{}
	caSubjectC := []string{}
	if s.OU != nil && len(*s.OU) > 0 {
		caSubjectOU = append(caSubjectOU, *s.OU)
	}
	if s.O != nil && len(*s.O) > 0 {
		caSubjectO = append(caSubjectO, *s.O)
	}
	if s.C != nil && len(*s.C) > 0 {
		caSubjectC = append(caSubjectC, *s.C)
	}
	return pkix.Name{
		CommonName:         s.CN,
		OrganizationalUnit: caSubjectOU,
		Organization:       caSubjectO,
		Country:            caSubjectC,
	}
}
