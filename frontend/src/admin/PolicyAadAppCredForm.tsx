import { useUpdateEffect } from "ahooks";
import React from "react";
import { CertificateIdentifierType, PolicyType } from "../generated";
import { PolicyContext } from "./PolicyContext";
import { PolicySelector } from "./Selectors";

export interface PolicyAadAppCredFormProps {
  certPolicyId: string | undefined;
  setCertPolicyId: (policyId: string | undefined) => void;
}

export function usePolicyAppCredFormState(): PolicyAadAppCredFormProps {
  const { policy } = React.useContext(PolicyContext);

  const [certPolicyId, setCertPolicyId] = React.useState<string | undefined>();
  React.useEffect(() => {
    if (policy?.certAadAppCred?.certificateIdentifier) {
      if (
        policy.certAadAppCred.certificateIdentifier.type ===
        CertificateIdentifierType.CertIdTypePolicyId
      ) {
        setCertPolicyId(policy.certAadAppCred.certificateIdentifier.id);
      } else {
        setCertPolicyId(undefined);
      }
    }
  }, [policy]);

  return {
    certPolicyId,
    setCertPolicyId,
  };
}

export function PolicyAadAppCredForm({
  certPolicyId,
  setCertPolicyId,
}: PolicyAadAppCredFormProps) {
  const { namespaceId } = React.useContext(PolicyContext);
  return (
    <div className="divide-y divide-neutral-200 space-y-6">
      <h2 className="text-2xl font-semibold">AAD Client Credential Policy</h2>
      <div className="pt-6">
        <PolicySelector
          selectedPolicyId={certPolicyId ?? ""}
          onChange={setCertPolicyId}
          namespaceId={namespaceId}
          policyType={PolicyType.PolicyType_CertRequest}
          label="Select certificate policy"
        />
      </div>
    </div>
  );
}
