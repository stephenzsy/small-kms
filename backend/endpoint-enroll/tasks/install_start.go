package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/auth"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/client"
)

func newServiceClientForInstall(
	clientID string,
	tenantID string,
	endpointClientID string,
) (*client.Client, error) {

	deviceCodeCredential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientID,
		TenantID: tenantID,
		UserPrompt: func(_ context.Context, msg azidentity.DeviceCodeMessage) error {
			_, err := fmt.Println(msg.Message)
			return err
		},
	})
	if err != nil {
		return nil, err
	}
	creds := auth.GetCachedTokenCredential(deviceCodeCredential)
	scopes := []string{fmt.Sprintf("%s/.default", endpointClientID)}
	return client.NewClientWithCreds(
		common.MustGetenv("SMALLKMS_SERVER_BASE_URL"), creds, scopes, tenantID)
}

func InstallStart(out io.Writer, outDescription string) error {
	tenantID := common.MustGetenv(common.IdentityEnvVarNameAzTenantID)
	clientID := uuid.MustParse(common.MustGetenv(common.IdentityEnvVarNameAzClientID))
	endpointClientID := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientID, "")
	templateGroupID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_GROUP_ID"))
	templateID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_ID"))
	deviceObjectID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_OBJECT_ID"))
	deviceLinkID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_LINK_ID"))
	servicePrincipalId := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_SERVICE_PRINCIPAL_ID"))

	serviceClient, err := newServiceClientForInstall(clientID.String(), tenantID, endpointClientID)
	if err != nil {
		return err
	}

	body := client.CertificateEnrollmentRequest{}
	body.FromCertificateEnrollmentRequestDeviceLinkedServicePrincipal(client.CertificateEnrollmentRequestDeviceLinkedServicePrincipal{
		AppID:              clientID,
		DeviceNamespaceID:  deviceObjectID,
		DeviceLinkID:       deviceLinkID,
		ServicePrincipalID: servicePrincipalId,
		Type:               client.CertEnrollTargetTypeDeviceLinkedServicePrincipal,
	})
	resp, err := serviceClient.BeginEnrollCertificateV2(context.Background(), templateGroupID, templateID, body)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reciept := client.CertificateEnrollmentReceipt{}
	if err := json.Unmarshal(bodyBytes, &reciept); err != nil {
		return err
	}
	marshalled, err := json.MarshalIndent(reciept, "", "  ")
	if err != nil {
		return err
	}
	out.Write(marshalled)
	_, err = fmt.Printf("Certificate enrollment receipt received: %s\n", outDescription)
	return err
}
