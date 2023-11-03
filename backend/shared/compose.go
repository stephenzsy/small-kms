package shared

type certificateRefComposed struct {
	ResourceRef
	CertificateRefFields
}

type certificateInfoComposed struct {
	CertificateRef
	CertificateInfoFields
}
