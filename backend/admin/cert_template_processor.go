package admin

import (
	"crypto/x509"
	"net/url"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	certtemplate "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type certTemplateProcessor struct {
	tmplDoc *CertificateTemplateDoc
	data    *certtemplate.TemplateVarData
}

// func (p *certTemplateProcessor) processSubject() (name pkix.Name, err error) {
// 	s := p.tmplDoc.Subject
// 	name.CommonName = processTemplate(s.CN, p.data)
// 	if s.C != nil && len(*s.C) > 0 {
// 		a := processTemplate(*s.C, p.data)
// 		if len(a) > 0 {
// 			name.Country = []string{a}
// 		}
// 	}
// 	if s.O != nil && len(*s.O) > 0 {
// 		a := processTemplate(*s.O, p.data)
// 		if len(a) > 0 {
// 			name.Organization = []string{a}
// 		}
// 	}
// 	if s.OU != nil && len(*s.OU) > 0 {
// 		a := processTemplate(*s.OU, p.data)
// 		if len(a) > 0 {
// 			name.OrganizationalUnit = []string{a}
// 		}
// 	}
// 	return
// }

func processTemplate(tmplStr string, data *certtemplate.TemplateVarData) string {
	tmplStr = strings.TrimSpace(tmplStr)
	tmpl, err := parseCertificateRequestTemplate(tmplStr)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to parse template: %s", tmplStr)
		return ""
	}
	if tmpl == nil {
		return tmplStr
	}
	transformed, err := executeTemplate(tmpl, data)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to execute template: %s", tmplStr)
		return ""
	}
	return transformed
}

func parseCertificateRequestTemplate(s string) (*template.Template, error) {
	return nil, nil
	// transformed, hasTemplate, err := preprocessTemplate(s)
	// if err != nil {
	// 	return nil, err
	// }
	// if hasTemplate {
	// 	return template.New(s).Parse(transformed)
	// }
	// return nil, nil
}

func executeTemplate(t *template.Template, data *certtemplate.TemplateVarData) (string, error) {
	if t == nil {
		return "", nil
	}
	sb := strings.Builder{}
	err := t.Execute(&sb, data)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (p *certTemplateProcessor) processSubjectAltNames(cert *x509.Certificate, data *certtemplate.TemplateVarData) {
	sans := p.tmplDoc.SubjectAlternativeNames
	if sans == nil {
		return
	}

	cert.EmailAddresses = utils.NilIfZeroLen(
		utils.FilterSlice(
			utils.MapSlices(sans.EmailAddresses, func(emailAddrStr string) string {
				return processTemplate(emailAddrStr, data)
			}),
			func(emailAddrStr string) bool {
				return len(emailAddrStr) > 0
			}))

	cert.URIs = utils.NilIfZeroLen[*url.URL](
		utils.FilterSlice[*url.URL](
			utils.MapSlices(sans.URIs, func(uriStr string) *url.URL {
				uriStr = processTemplate(uriStr, data)
				if len(uriStr) > 0 {
					uri, _ := url.Parse(uriStr)
					return uri
				}
				return nil
			}),
			func(uri *url.URL) bool {
				return uri != nil
			}))
}

func (p *certTemplateProcessor) processTemplate(c *x509.Certificate) (err error) {
	// c.Subject, err = p.processSubject()
	// if err != nil {
	// 	return err
	// }
	p.processSubjectAltNames(c, p.data)
	c.NotAfter = c.NotBefore.AddDate(0, int(p.tmplDoc.ValidityInMonths), 0)
	t := p.tmplDoc

	if t.Usage == UsageServerAndClient || t.Usage == UsageServerOnly {
		c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if t.Usage == UsageServerAndClient || t.Usage == UsageClientOnly {
		c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}
	return nil
}
