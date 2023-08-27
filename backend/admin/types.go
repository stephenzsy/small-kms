package admin

type CertDBItem struct {
	CertificateRef
	KeyStore  string `json:"keyStore,omitempty"`
	CertStore string `json:"certStore,omitempty"`
}
