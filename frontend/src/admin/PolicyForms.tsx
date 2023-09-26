import { useMemoizedFn } from "ahooks";
import { useContext, useId, useMemo, useState } from "react";
import { Button } from "../components/Button";
import { ErrorAlert } from "../components/ErrorAlert";
import { WellknownId } from "../constants";
import {
  CertificateIdentifierType,
  CertificateUsage,
  PolicyParameters,
  PolicyType,
} from "../generated";
import {
  PolicyAadAppCredForm,
  usePolicyAppCredFormState,
} from "./PolicyAadAppCredForm";
import { PolicyContext } from "./PolicyContext";
import {
  IsIntCaNamespace,
  isRootCaNamespace,
  policyTypeNames,
} from "./displayConstants";
import {
  PolicyCertificateRequestForm,
  useCertificateRequestFormState,
} from "./PolicyCertificateRequestForm";

export type PolicyFormsProps = {};

const policyTypes: ReadonlyArray<{
  policyType: PolicyType;
  title: string;
  defaultPolicyId: string;
}> = [
  {
    policyType: "certRequest",
    title: policyTypeNames[PolicyType.PolicyType_CertRequest],
    defaultPolicyId: WellknownId.defaultPolicyIdCertRequest,
  },
  {
    policyType: "certAadAppCred",
    title: policyTypeNames[PolicyType.PolicyType_CertAadAppClientCredential],
    defaultPolicyId: WellknownId.defaultPolicyIdAadAppCred,
  },
];

export function PolicyForms(_props: PolicyFormsProps) {
  const { policyId, putPolicy, putPolicyError, namespaceId } =
    useContext(PolicyContext);
  const inptuIdPrefix = useId();
  const policyTypeOverride = useMemo(() => {
    switch (policyId) {
      case WellknownId.defaultPolicyIdCertRequest:
        return PolicyType.PolicyType_CertRequest;
      case WellknownId.defaultPolicyIdAadAppCred:
        return PolicyType.PolicyType_CertAadAppClientCredential;
    }
    return undefined;
  }, [policyId]);
  const [policyTypeState, setPolicyType] = useState<PolicyType>();
  const policyType = policyTypeOverride ?? policyTypeState;
  const checkboxOnChange = useMemoizedFn<
    React.ChangeEventHandler<HTMLInputElement>
  >((e) => {
    if (e.target.checked) {
      setPolicyType(e.target.value as PolicyType);
    }
  });

  const aadAppCredState = usePolicyAppCredFormState();
  const certReqState = useCertificateRequestFormState();

  const onSubmit = useMemoizedFn<React.FormEventHandler<HTMLFormElement>>(
    (e) => {
      e.preventDefault();
      if (!policyType) {
        return;
      }
      const policyParameters: PolicyParameters = {
        policyType: policyType,
      };
      switch (policyType) {
        case PolicyType.PolicyType_CertRequest:
          if (
            !certReqState.issuerNamespaceId ||
            !certReqState.issuerPolicyId ||
            !certReqState.subjectCN ||
            !certReqState.keyStorePath
          ) {
            return;
          }
          policyParameters.certRequest = {
            issuerNamespaceId: certReqState.issuerNamespaceId,
            issuerPolicyIdentifier: certReqState.issuerPolicyId,
            subject: {
              cn: certReqState.subjectCN,
              ou: certReqState.subjectOU,
              o: certReqState.subjectO,
              c: certReqState.subjectC,
            },
            usage: isRootCaNamespace(namespaceId)
              ? CertificateUsage.Usage_RootCA
              : IsIntCaNamespace(namespaceId)
              ? CertificateUsage.Usage_IntCA
              : certReqState.certUsage,
            keyStorePath: certReqState.keyStorePath,
          };
          break;
        case PolicyType.PolicyType_CertAadAppClientCredential:
          if (!aadAppCredState.certPolicyId) {
            return;
          }
          policyParameters.certAadAppCred = {
            certificateIdentifier: {
              id: aadAppCredState.certPolicyId,
              type: CertificateIdentifierType.CertIdTypePolicyId,
            },
          };
          break;
        default:
          return;
      }
      return putPolicy(policyParameters);
    }
  );

  return (
    <form
      className="divide-y-2 divide-neutral-500 overflow-hidden rounded-lg bg-white shadow p-6 space-y-6"
      onSubmit={onSubmit}
    >
      {policyTypeOverride ? (
        <div>Policy type: {policyTypeNames[policyTypeOverride]}</div>
      ) : (
        <div>
          <label className="text-base font-semibold text-gray-900">
            Select policy type
          </label>
          <fieldset className="mt-4">
            <legend className="sr-only">Notification method</legend>
            <div className="space-y-4">
              {policyTypes.map((p) => (
                <div key={p.policyType} className="flex items-center">
                  <input
                    id={`${inptuIdPrefix}:${p.policyType}`}
                    name="notification-method"
                    type="radio"
                    checked={policyTypeState === p.policyType}
                    className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
                    onChange={checkboxOnChange}
                    value={p.policyType}
                  />
                  <label
                    htmlFor={`${inptuIdPrefix}:${p.policyType}`}
                    className="ml-3 block text-sm font-medium leading-6 text-gray-900"
                  >
                    {p.title}
                  </label>
                </div>
              ))}
            </div>
          </fieldset>
        </div>
      )}
      {putPolicyError && <ErrorAlert error={putPolicyError} />}
      <div className="pt-6">
        {policyType === PolicyType.PolicyType_CertRequest && (
          <PolicyCertificateRequestForm {...certReqState} />
        )}
        {policyType === PolicyType.PolicyType_CertAadAppClientCredential && (
          <PolicyAadAppCredForm {...aadAppCredState} />
        )}
      </div>
      <div className="pt-6 flex flex-row items-center gap-x-6 justify-end">
        <Button variant="primary" type="submit">
          Create or Update
        </Button>
      </div>
    </form>
  );
}
