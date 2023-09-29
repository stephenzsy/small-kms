import { WellknownId } from "../constants";
import { CertificateUsage, PolicyType } from "../generated";

export const nsDisplayNames: Record<string, string> = {
  [WellknownId.nsRootCa]: "Root CA",
  [WellknownId.nsTestRootCa]: "Test Root CA",
  [WellknownId.nsIntCaIntranet]: "Intranet CA",
  [WellknownId.nsTestIntCa]: "Test Intermediate CA",
};

export function isRootCaNamespace(namespaceId: string) {
  switch (namespaceId) {
    case WellknownId.nsRootCa:
    case WellknownId.nsTestRootCa:
      return true;
  }
  return false;
}

export function IsIntCaNamespace(namespaceId: string) {
  switch (namespaceId) {
    case WellknownId.nsIntCaIntranet:
    case WellknownId.nsTestIntCa:
      return true;
  }
  return false;
}

export const policyTypeNames: Record<PolicyType, string> = {
  [PolicyType.PolicyType_CertRequest]: "Certificate Request Policy",
  [PolicyType.PolicyType_CertEnroll]: "Certificate Enrollment Policy",
  [PolicyType.PolicyType_CertAadAppClientCredential]:
    "AAD Application Client Credential Certificate Policy",
};

export const certUsageNames: Record<CertificateUsage, string> = {
  [CertificateUsage.Usage_ServerAndClient]: "Server and client",
  [CertificateUsage.Usage_ServerOnly]: "Server only",
  [CertificateUsage.Usage_ClientOnly]: "Client only",
  [CertificateUsage.Usage_RootCA]: "Root CA",
  [CertificateUsage.Usage_IntCA]: "Intermediate CA",
  [CertificateUsage.Usage_AADClientCredential]: "AAD Client Credential",
};
