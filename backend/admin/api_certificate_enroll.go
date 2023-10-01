package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// extract owner name

// const (
// 	deviceOwnershipTypeCompany  = "Company"
// 	deviceOwnershipTypePersonal = "Personal"
// )
/*
func (r *CertificateEnrollRequest) createX509Certificate(c *gin.Context) (*x509.Certificate, error) {
	cert := x509.Certificate{}
	if r.IssueToUser != nil && *r.IssueToUser {
		auth.CallerPrincipalName(c)
	}
	return &cert, nil
}*/

func (s *adminServer) BeginEnrollCertificateV2(c *gin.Context, nsID uuid.UUID, templateId TemplateIdParameter) {

	req := new(CertificateEnrollmentRequest)
	if err := c.Bind(req); err != nil {
		respondPublicErrorMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	//c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("user %s does not have permission to enroll certificate for target %s", callerId.String(), targetId.String())})

	p := new(CertificateEnrollmentRequest)
	if err := c.Bind(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	/*
		if p.ValidityInMonths < 1 || p.ValidityInMonths > 120 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "validity in months must be between 1 and 120"})
			return
		}

		// check enrollment policy
		policyDoc, err := s.GetPolicyDoc(c, p.Issuer.IssuerNamespaceID, p.PolicyID)
		if err != nil {
			if common.IsAzNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"message": "no policy found for given issuer and policy id"})
				return
			}
			log.Error().Err(err).Msg("Failed to get enrollment policy")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		if policyDoc.PolicyType != PolicyTypeCertEnroll {
			c.JSON(http.StatusBadRequest, gin.H{"message": "policy is not for certificate enrollment"})
			return
		}

		// validate request against policy
		if policyDoc.CertEnroll.MaxValidityInMonths < p.ValidityInMonths {
			c.JSON(http.StatusBadRequest, gin.H{"message": "validity in months exceeds policy limit"})
			return
		}
		usageValidated := false
		for _, usage := range policyDoc.CertEnroll.AllowedUsages {
			if usage == p.Usage {
				usageValidated = true
				break
			}
		}
		if !usageValidated {
			c.JSON(http.StatusBadRequest, gin.H{"message": "usage is not allowed by policy"})
		}
	*/
}

func (s *adminServer) CompleteCertificateEnrollmentV2(c *gin.Context, namespaceType NamespaceTypeParameter, namespaceId NamespaceIdParameter, certId CertIdParameter, params CompleteCertificateEnrollmentV2Params) {
	// Your code here
}
