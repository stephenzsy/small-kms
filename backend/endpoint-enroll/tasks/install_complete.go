package tasks

import (
	"io"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/client"
)

func InstallComplete(receiptIn io.Reader) error {
	receipt, err := readReceipt(receiptIn)
	if err != nil {
		return err
	}

	tenantID := common.MustGetenv(common.DefaultEnvVarAzureTenantId)
	clientID := uuid.MustParse(common.MustGetenv(common.DefaultEnvVarAzureClientId))
	endpointClientID := common.MustGetenv(common.DefaultEnvVarAppAzureClientId)

	templateGroupID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_GROUP_ID"))
	templateID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_TEMPLATE_ID"))
	deviceObjectID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_OBJECT_ID"))
	deviceLinkID := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_DEVICE_LINK_ID"))
	servicePrincipalId := uuid.MustParse(common.MustGetenv("SMALLKMS_ENROLL_SERVICE_PRINCIPAL_ID"))

	serviceClient, err := newServiceClientForInstall(clientID.String(), tenantID, endpointClientID)
	if err != nil {
		return err
	}

	body := client.CertificateEnrollmentReplyFinalize{}
	/*
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
	*/
	return nil
}
