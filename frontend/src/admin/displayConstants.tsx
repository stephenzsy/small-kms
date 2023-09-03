import {
  PolicyType,
  TestNamespaceId,
  WellKnownNamespaceId,
} from "../generated";

export const namespaceFriendlierNames: Record<string, string> = {
  [WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA]: "Root CA",
  [TestNamespaceId.TestNamespaceIDStr_RootCA]: "Test Root CA",
};

export function isRootCANamespace(namespaceId: string) {
  switch (namespaceId) {
    case WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA:
    case TestNamespaceId.TestNamespaceIDStr_RootCA:
      return true;
  }
  return false;
}

export const certRequestPolicyNames: Record<PolicyType, string> = {
  [PolicyType.PolicyType_CertRequest]: "Certificate Request Policy",
  [PolicyType.PolicyType_CertIssue]: "Certificate Issurance Policy",
};
