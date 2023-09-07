import { WellknownId } from "../constants";
import { PolicyType } from "../generated";

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

export const certRequestPolicyNames: Record<PolicyType, string> = {
  [PolicyType.PolicyType_CertRequest]: "Certificate Request Policy",
};
