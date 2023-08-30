package admin

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
)

type createCertificateInternalParameters struct {
	usage              CertificateUsage
	kty                KeyParametersKty
	size               KeyParametersSize
	namespaceID        uuid.UUID
	keyVaultKeyName    string
	keyVaultKeyVersion string
	subject            CertificateSubject
	createdBy          string
	issuer             CertDBItem
}

func (s *adminServer) validateCreateCertificateOptions(c context.Context, out *createCertificateInternalParameters, namespaceID NamespaceID, p *CreateCertificateParameters) (err error) {
	// allow root ca
	out.usage = p.Usage
	out.subject = p.Subject
	out.namespaceID = namespaceID
	switch p.Usage {
	case UsageRootCA:
		if namespaceID != wellKnownNamespaceID_RootCA {
			return fmt.Errorf("invalid namespace for root ca: %s", namespaceID.String())
		}
		// fill in default parameters
		out.kty = KtyRSA
		out.size = KeySize4096
	case UsageIntCA:
		switch namespaceID {
		case
			wellKnownNamespaceID_IntCAService,
			wellKnownNamespaceID_IntCaSCEPIntranet:
			if p.IssuerNamespace != wellKnownNamespaceID_RootCA {
				return fmt.Errorf("invalid issuer namespace for intermediate ca: %s", p.IssuerNamespace.String())
			}
			out.issuer, err = s.readCertDBItem(c, p.IssuerNamespace, p.Issuer)
			if err != nil || out.issuer.ID == uuid.Nil {
				return fmt.Errorf("invalid issuer: %s/%s", p.IssuerNamespace.String(), p.Issuer.String())
			}
			out.kty = KtyRSA
			out.size = KeySize2048
		default:
			return fmt.Errorf("invalid namespace for intermediate ca: %s", namespaceID.String())
		}
	default:
		return fmt.Errorf("unsupported usage: %s", p.Usage)
	}
	return
}

func (s *adminServer) CreateCertificateV1(c *gin.Context, namespaceID NamespaceID) {
	body := CreateCertificateParameters{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}
	if !auth.CallerHasAdminAppRole(c) {
		c.JSON(403, gin.H{"message": "User must have admin role"})
		return
	}
	p := createCertificateInternalParameters{
		createdBy: auth.GetCallerID(c),
	}
	if err := s.validateCreateCertificateOptions(c, &p, namespaceID, &body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	lastCertificate, err := s.findLatestCertificate(c.Request.Context(), p.namespaceID, p.subject.CN)
	if err != nil {
		log.Printf("Error find latest certificate: %s", err.Error())
		c.JSON(500, gin.H{"message": "internal error"})
		return
	}
	if lastCertificate.ID != uuid.Nil {
		if len(lastCertificate.KeyStore) > 0 {
			keyId := azkeys.ID(lastCertificate.KeyStore)
			if body.Options != nil && body.Options.KeepKeyVersion != nil && *body.Options.KeepKeyVersion {
				p.keyVaultKeyName = keyId.Name()
				p.keyVaultKeyVersion = keyId.Version()
			} else if body.Options == nil || body.Options.NewKeyName == nil || !*body.Options.NewKeyName {
				p.keyVaultKeyName = keyId.Name()
			}
		}
	}
	if len(p.keyVaultKeyName) == 0 {
		p.keyVaultKeyName = generateRandomHexSuffix(namespacePrefixMapping[p.namespaceID])
	}

	switch p.usage {
	case UsageRootCA:
	case UsageIntCA:
		certCreated, err := s.createCACertificate(c, p)
		if err != nil {
			c.JSON(400, gin.H{"message": "Failed to create certificate", "error": err.Error()})
			log.Printf("Failed to create cert: %s", err.Error())
			return
		}
		c.JSON(201, &certCreated.CertificateRef)
	default:
		c.JSON(400, gin.H{"message": "Usage not supported", "usage": body.Usage})
		return
	}
}
