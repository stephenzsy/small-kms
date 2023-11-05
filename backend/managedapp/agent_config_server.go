package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	cloudauthzaz "github.com/stephenzsy/small-kms/backend/cloud/authz/az"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/cloudutils"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentConfigServerDoc struct {
	AgentConfigDoc
	GlobalKeyVaultEndpoint string                   `json:"globalKeyVaultEndpoint"`
	GlobalACRImageRef      string                   `json:"globalAcrImageRef"`
	TLSCertificatePolicyID base.Identifier          `json:"tlsCertificatePolicyId"`
	TLSCertificateID       base.Identifier          `json:"tlsCertificateId"`
	JWTKeyCertPolicyID     base.DocFullIdentifier   `json:"jwtKeyCertPolicyId"`
	JWTKeyCertIDs          []base.DocFullIdentifier `json:"jwtKeyCertIds"`
}

func (d *AgentConfigServerDoc) init(nsKind base.NamespaceKind, nsIdentifier base.Identifier) {
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
	if err := base.GetAzCosmosCRUDService(c).Read(c, base.NewDocFullIdentifier(nsCtx.Kind(),
		nsCtx.Identifier(), base.ResourceKindNamespaceConfig, base.StringIdentifier(string(base.AgentConfigNameServer))), doc, nil); err != nil {
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
	doc.init(nsCtx.Kind(), nsCtx.Identifier())

	docSvc := base.GetAzCosmosCRUDService(c)

	digest := md5.New()
	switch nsCtx.Kind() {
	case base.NamespaceKindSystem:
		if nsCtx.Identifier() != base.StringIdentifier("default") {
			return fmt.Errorf("%w: only default system namespace is supported", base.ErrResponseStatusBadRequest)
		}
		doc.GlobalKeyVaultEndpoint = s.GetAzKeyVaultEndpoint()
		digest.Write([]byte(doc.GlobalKeyVaultEndpoint))

		doc.GlobalACRImageRef = param.AzureACRImageRef
		// do not write digest for image ref
	case base.NamespaceKindServicePrincipal:

		globalDoc := &AgentConfigServerDoc{}
		if err := docSvc.Read(c,
			base.NewDocFullIdentifier(base.NamespaceKindSystem,
				base.StringIdentifier("default"), base.ResourceKindNamespaceConfig, base.StringIdentifier(string(base.AgentConfigNameServer))), globalDoc, nil); err != nil {
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
			base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, param.TlsCertificatePolicyId), 1)
		if err != nil {
			return err
		}
		if len(certIds) == 0 {
			return fmt.Errorf("%w: no certificate issued by policy %s", base.ErrResponseStatusBadRequest, param.TlsCertificatePolicyId.String())
		}
		doc.TLSCertificateID = certIds[0]
		digest.Write([]byte(doc.TLSCertificateID.String()))

		doc.JWTKeyCertPolicyID = param.JwtKeyCertPolicyId
		jwtKeyCertIdentifiers, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
			doc.JWTKeyCertPolicyID, 2)
		if err != nil {
			return err
		}
		doc.JWTKeyCertIDs = utils.MapSlice(jwtKeyCertIdentifiers, func(id base.Identifier) base.DocFullIdentifier {
			fullIdentifier := base.NewDocFullIdentifier(doc.JWTKeyCertPolicyID.NamespaceKind(), doc.JWTKeyCertPolicyID.NamespaceIdentifier(), base.ResourceKindCert, id)
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

	s.assignAgentServerRoles(c, nsCtx.Identifier().UUID(), doc)

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
				cert.GetKeyStoreName(nsCtx.Kind(), nsCtx.Identifier(), certPolicyDoc.ID)).Build(),
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

func (s *server) apiListAgentConfigServerRoleAssignments(c ctx.RequestContext) error {
	doc, err := ApiReadAgentConfigDoc(c)
	if err != nil {
		return err
	}

	if doc.GlobalACRImageRef == "" {
		return fmt.Errorf("%w: image ref is not specified", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	if !nsCtx.Identifier().IsUUID() {
		return fmt.Errorf("%w: invalid namespace identifier", base.ErrResponseStatusBadRequest)
	}
	assignedTo := nsCtx.Identifier().UUID()
	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return err
	}

	subscriptionIDBuilder := &cloudutils.AzureSubscriptionResourceIDBuilder{
		SubscriptionID: s.GetAzSubscriptionID(),
	}
	pagers := make([]utils.ItemsPager[*armauthorization.RoleAssignment], 0, 2)
	// ACR Pull
	{
		acrName, err := acr.ExtractACRName(doc.GlobalACRImageRef)
		if err != nil {
			return err
		}
		if acrResourceGroupName := s.EnvService().Default("AZURE_RESOURCE_GROUP_NAME", "", "ACR_"); acrResourceGroupName != "" {
			scope := subscriptionIDBuilder.WithResourceGroup(acrResourceGroupName).WithContainerRegistry(acrName).Build()
			pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
		}
	}
	// Key Vault Secrets User
	{
		scope := subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), "secrets",
			cert.GetKeyStoreName(nsCtx.Kind(), nsCtx.Identifier(), doc.TLSCertificatePolicyID)).Build()
		pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
	}

	resultPager := utils.NewSerializableItemsPager(
		utils.NewMappedItemsPager[*base.AzureRoleAssignment, *armauthorization.RoleAssignment](utils.NewChainedItemPagers(pagers...), func(ra *armauthorization.RoleAssignment) *base.AzureRoleAssignment {
			if ra == nil {
				return nil
			}
			return &base.AzureRoleAssignment{
				ID:               ra.ID,
				Name:             ra.Name,
				RoleDefinitionId: ra.Properties.RoleDefinitionID,
				PrincipalId:      ra.Properties.PrincipalID,
			}
		}))

	return api.RespondPagerList(c, resultPager)
}
