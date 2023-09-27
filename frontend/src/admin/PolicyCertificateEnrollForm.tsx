import { useMemoizedFn, useSetState } from "ahooks";
import React from "react";
import { CertificateUsage } from "../generated";
import { CertificateUsageSelector } from "./CertificateUsageSelector";
import { InputField } from "./InputField";
import { PolicyContext } from "./PolicyContext";

export interface PolicyAadAppCredFormProps {
  validityInMonths: number | undefined;
  setValidityInMonths: (value: number | undefined) => void;
  certificateUsages: Partial<Record<CertificateUsage, boolean>>;
  setCertificateUsages: (
    value: Partial<Record<CertificateUsage, boolean>>
  ) => void;
}

export function usePolicyCertEnrollFormState(): PolicyAadAppCredFormProps {
  const { policy } = React.useContext(PolicyContext);

  const [validityInMonths, setValidityInMonths] = React.useState<number>();

  const [certPolicyId, setCertPolicyId] = React.useState<string | undefined>();
  const [certificateUsages, setCertificateUsages] = useSetState<
    Partial<Record<CertificateUsage, boolean>>
  >({});
  React.useEffect(() => {
    const certEnroll = policy?.certEnroll;
    if (!certEnroll) {
      return;
    }
    if (certEnroll.maxValidityInMonths) {
      setValidityInMonths(certEnroll.maxValidityInMonths);
    }
    setCertificateUsages(
      certEnroll.allowedUsages.reduce<
        Partial<Record<CertificateUsage, boolean>>
      >((acc, usage) => {
        acc[usage] = true;
        return acc;
      }, {})
    );
  }, [policy]);

  return {
    validityInMonths,
    setValidityInMonths,
    certificateUsages,
    setCertificateUsages,
  };
}

export function PolicyCertificateEnrollForm(props: PolicyAadAppCredFormProps) {
  const { namespaceId } = React.useContext(PolicyContext);
  const certUsageInputOnChange = useMemoizedFn<
    (e: React.ChangeEvent<HTMLInputElement>) => void
  >((e) => {
    props.setCertificateUsages({ [e.target.value]: e.target.checked });
  });

  const certUsageIsChecked = useMemoizedFn((usage: CertificateUsage) => {
    return !!props.certificateUsages[usage];
  });
  return (
    <div className="divide-y divide-neutral-200 space-y-4">
      <h2 className="text-2xl font-semibold">Certificate Enrollment Policy</h2>
      <div className="pt-4 space-y-4">
        <InputField
          labelContent="Max validity in months"
          type="number"
          inputMode="numeric"
          placeholder="12"
          value={props.validityInMonths}
          onChange={props.setValidityInMonths as any}
        />
        <CertificateUsageSelector
          inputType="checkbox"
          onChange={certUsageInputOnChange}
          isChecked={certUsageIsChecked}
          label="Allowed certificate usages"
          availableUsages={[
            CertificateUsage.Usage_ServerAndClient,
            CertificateUsage.Usage_ClientOnly,
          ]}
        />
      </div>
    </div>
  );
}
