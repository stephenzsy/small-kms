package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	cloudauthzaz "github.com/stephenzsy/small-kms/backend/cloud/authz/az"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/cloudutils"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/secret"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentConfigServerDoc struct {
	AgentConfigDoc
	GlobalKeyVaultEndpoint string            `json:"globalKeyVaultEndpoint"`
	GlobalACRImageRef      string            `json:"globalAcrImageRef"`
	TLSCertificatePolicyID base.ID           `json:"tlsCertificatePolicyId"`
	TLSCertificateID       base.ID           `json:"tlsCertificateId"`
	JWTKeyCertPolicyID     base.DocLocator   `json:"jwtKeyCertPolicyId"`
	JWTKeyCertIDs          []base.DocLocator `json:"jwtKeyCertIds"`
}

func (d *AgentConfigServerDoc) init(nsKind base.NamespaceKind, nsIdentifier base.ID) {
	d.AgentConfigDoc.init(nsKind, nsIdentifier, base.AgentConfigNameServer)
}

const envVarMessage = "For security reasons, these variables should be set manualy after review in the agent config server. Server must not auto update these variables."

func (d *AgentConfigServerDoc) populateModel(m *AgentConfigServer) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.Env = AgentConfigServerEnv{
		Message:                             envVarMessage,
		EnvVarAzureKeyVaultResourceEndpoint: d.GlobalKeyVaultEndpoint,
		EnvVarAzureContainerRegistryImageRepository: strings.Split(d.GlobalACRImageRef, ":")[0],
	}
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	m.TlsCertificateId = d.TLSCertificateID
	m.JwtKeyCertPolicyId = d.JWTKeyCertPolicyID
	m.JWTKeyCertIDs = d.JWTKeyCertIDs
	m.RefreshAfter = time.Now().Add(24 * time.Hour).UTC()
	m.AzureACRImageRef = d.GlobalACRImageRef
}

// will wrap 404 if doc is not found
func ApiReadAgentConfigDoc(c ctx.RequestContext) (*AgentConfigServerDoc, error) {
	nsCtx := ns.GetNSContext(c)
	doc := &AgentConfigServerDoc{}
	if err := base.GetAzCosmosCRUDService(c).Read(c, base.NewDocLocator(nsCtx.Kind(),
		nsCtx.ID(), base.ResourceKindNamespaceConfig, base.ID(base.AgentConfigNameServer)), doc, nil); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameServer)
		}
		return nil, err
	}
	return doc, nil
}

func apiGetAgentConfigServer(c ctx.RequestContext) error {
	doc, err := ApiReadAgentConfigDoc(c)
	if err != nil {
		return err
	}

	m := &AgentConfigServer{}
	doc.populateModel(m)
	return c.JSON(200, m)
}

func (s *server) apiPutAgentConfigServer(c ctx.RequestContext, param *AgentConfigServerFields) error {

	nsCtx := ns.GetNSContext(c)

	doc := &AgentConfigServerDoc{}
	doc.init(nsCtx.Kind(), nsCtx.ID())

	docSvc := base.GetAzCosmosCRUDService(c)

	digest := md5.New()
	switch nsCtx.Kind() {
	case base.NamespaceKindSystem:
		if nsCtx.ID() != base.ID("default") {
			return fmt.Errorf("%w: only default system namespace is supported", base.ErrResponseStatusBadRequest)
		}
		doc.GlobalKeyVaultEndpoint = s.GetAzKeyVaultEndpoint()
		digest.Write([]byte(doc.GlobalKeyVaultEndpoint))

		doc.GlobalACRImageRef = param.AzureACRImageRef
		// do not write digest for image ref
	case base.NamespaceKindServicePrincipal:

		globalDoc := &AgentConfigServerDoc{}
		if err := docSvc.Read(c,
			base.NewDocLocator(base.NamespaceKindSystem,
				base.ID("default"), base.ResourceKindNamespaceConfig, base.ID(base.AgentConfigNameServer)), globalDoc, nil); err != nil {
			if errors.Is(err, base.ErrAzCosmosDocNotFound) {
				return fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameServer)
			}
			return err
		}

		doc.GlobalKeyVaultEndpoint = globalDoc.GlobalKeyVaultEndpoint
		doc.GlobalACRImageRef = globalDoc.GlobalACRImageRef
		globalDocVersionBytes, err := hex.DecodeString(globalDoc.Version)
		if err != nil {
			return err
		}
		digest.Write(globalDocVersionBytes)

		doc.TLSCertificatePolicyID = param.TlsCertificatePolicyId
		certIds, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
			base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, param.TlsCertificatePolicyId), 1)
		if err != nil {
			return err
		}
		if len(certIds) == 0 {
			return fmt.Errorf("%w: no certificate issued by policy %s", base.ErrResponseStatusBadRequest, param.TlsCertificatePolicyId)
		}
		doc.TLSCertificateID = certIds[0]
		digest.Write([]byte(doc.TLSCertificateID))

		doc.JWTKeyCertPolicyID = param.JwtKeyCertPolicyId
		jwtKeyCertIdentifiers, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
			doc.JWTKeyCertPolicyID, 2)
		if err != nil {
			return err
		}
		doc.JWTKeyCertIDs = utils.MapSlice(jwtKeyCertIdentifiers, func(id base.ID) base.DocLocator {
			fullIdentifier := base.NewDocLocator(doc.JWTKeyCertPolicyID.NamespaceKind(), doc.JWTKeyCertPolicyID.NamespaceID(), base.ResourceKindCert, id)
			digest.Write([]byte(fullIdentifier.String()))
			return fullIdentifier
		})
	default:
		return fmt.Errorf("%w: unsupported namespace kind: %s", base.ErrResponseStatusBadRequest, nsCtx.Kind())
	}
	doc.Version = hex.EncodeToString(digest.Sum(nil))

	err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	s.assignAgentServerRoles(c, nsCtx.ID().UUID(), doc)

	m := &AgentConfigServer{}
	doc.populateModel(m)
	return c.JSON(200, m)
}

var (
	roleDefIDAcrPull             = uuid.MustParse("7f951dda-4ed3-4680-a7ca-43fe172d538d")
	roleDefIDKeyVaultSecretsUser = uuid.MustParse("4633458b-17de-408a-b874-0445c86b69e6")
)

func (s *server) assignAgentServerRoles(c ctx.RequestContext, assignedTo uuid.UUID, doc *AgentConfigServerDoc) error {
	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return err
	}
	subscriptionIDBuilder := &cloudutils.AzureSubscriptionResourceIDBuilder{
		SubscriptionID: s.GetAzSubscriptionID(),
	}
	nsCtx := ns.GetNSContext(c)
	// AcrPull
	{
		acrName, err := acr.ExtractACRName(doc.GlobalACRImageRef)
		if err != nil {
			return err
		}
		if acrResourceGroupName := s.EnvService().Default("AZURE_RESOURCE_GROUP_NAME", "", "ACR_"); acrResourceGroupName != "" {
			p := cloudauthzaz.RoleAssignmentProvisioner{
				RoleDefinitionID: roleDefIDAcrPull,
				Scope:            subscriptionIDBuilder.WithResourceGroup(acrResourceGroupName).WithContainerRegistry(acrName).Build(),
				AssignedTo:       assignedTo,
			}

			roleDefID := subscriptionIDBuilder.WithRoleDefinitionID(roleDefIDAcrPull).Build()
			if isAssigned, err := p.IsRoleAssigned(c, armRAClient, roleDefID); err != nil {
				return err
			} else if !isAssigned {
				if err := p.AssignRole(c, armRAClient, roleDefID); err != nil {
					return err
				}
			}
		}

	}
	// tls cert secret
	{
		certPolicyDoc, err := cert.ReadCertPolicyDoc(c, doc.TLSCertificatePolicyID)
		if err != nil {
			return err
		}

		p := cloudauthzaz.RoleAssignmentProvisioner{
			RoleDefinitionID: roleDefIDKeyVaultSecretsUser,
			Scope: subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), "secrets",
				cert.GetKeyStoreName(nsCtx.Kind(), nsCtx.ID(), certPolicyDoc.ID)).Build(),
			AssignedTo: assignedTo,
		}

		roleDefID := subscriptionIDBuilder.WithRoleDefinitionID(roleDefIDKeyVaultSecretsUser).Build()

		if isAssigned, err := p.IsRoleAssigned(c, armRAClient, roleDefID); err != nil {
			return err
		} else if !isAssigned {
			if err := p.AssignRole(c, armRAClient, roleDefID); err != nil {
				return err
			}
		}
	}
	return nil

}

func (s *server) assignAgentRadiusRoles(c ctx.RequestContext, assignedTo uuid.UUID, doc *AgentConfigRadiusDoc) error {
	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return err
	}
	subscriptionIDBuilder := &cloudutils.AzureSubscriptionResourceIDBuilder{
		SubscriptionID: s.GetAzSubscriptionID(),
	}
	nsCtx := ns.GetNSContext(c)
	// AcrPull
	{
		acrName, err := acr.ExtractACRName(doc.GlobalRadiusServerACRImageRef)
		if err != nil {
			return err
		}
		if acrResourceGroupName := s.EnvService().Default("AZURE_RESOURCE_GROUP_NAME", "", "ACR_"); acrResourceGroupName != "" {
			p := cloudauthzaz.RoleAssignmentProvisioner{
				RoleDefinitionID: roleDefIDAcrPull,
				Scope:            subscriptionIDBuilder.WithResourceGroup(acrResourceGroupName).WithContainerRegistry(acrName).Build(),
				AssignedTo:       assignedTo,
			}

			roleDefID := subscriptionIDBuilder.WithRoleDefinitionID(roleDefIDAcrPull).Build()
			if isAssigned, err := p.IsRoleAssigned(c, armRAClient, roleDefID); err != nil {
				return err
			} else if !isAssigned {
				if err := p.AssignRole(c, armRAClient, roleDefID); err != nil {
					return err
				}
			}
		}

	}
	// client secret
	{

		for _, client := range doc.Clients {
			if client.SecretPolicyId == "" {
				continue
			}
			p := cloudauthzaz.RoleAssignmentProvisioner{
				RoleDefinitionID: roleDefIDKeyVaultSecretsUser,
				Scope: subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), "secrets",
					secret.GetKeyStoreName(nsCtx.Kind(), nsCtx.ID(), client.SecretPolicyId)).Build(),
				AssignedTo: assignedTo,
			}

			roleDefID := subscriptionIDBuilder.WithRoleDefinitionID(roleDefIDKeyVaultSecretsUser).Build()

			if isAssigned, err := p.IsRoleAssigned(c, armRAClient, roleDefID); err != nil {
				return err
			} else if !isAssigned {
				if err := p.AssignRole(c, armRAClient, roleDefID); err != nil {
					return err
				}
			}
		}
	}
	return nil

}
