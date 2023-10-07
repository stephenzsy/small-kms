import { WellknownId } from "../constants";
import { PolicyType } from "../generated";
import { CertificateUsage } from "../generated3";

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
  [PolicyType.PolicyType_CertEnroll]: "Certificate Enrollment Policy",
};

export const certUsageNames: Record<CertificateUsage, string> = {
  [CertificateUsage.CertUsageCA]: "Certificate Authority",
  [CertificateUsage.CertUsageCARoot]: "Root Certificate Authority",
  [CertificateUsage.CertUsageClientAuth]: "Client Authentication",
  [CertificateUsage.CertUsageServerAuth]: "Server Authentication",
};
