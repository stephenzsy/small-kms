import React from "react";
import { CertificateUsage, NamespaceType, PolicyType } from "../generated";
import { InputField } from "./InputField";
import { PolicyContext } from "./PolicyContext";
import { IssuerSelector, PolicySelector } from "./Selectors";
import { useMemoizedFn } from "ahooks";
import { CertificateUsageSelector } from "./CertificateUsageSelector";

export interface PolicyCertificateRequestFormProps {
  issuerNamespaceId: string | undefined;
  setIssuerNamespaceId: (value: string | undefined) => void;
  issuerPolicyId: string | undefined;
  setIssuerPolicyId: (value: string | undefined) => void;
  subjectCN: string | undefined;
  setSubjectCN: (value: string | undefined) => void;
  subjectOU: string | undefined;
  setSubjectOU: (value: string | undefined) => void;
  subjectO: string | undefined;
  setSubjectO: (value: string | undefined) => void;
  subjectC: string | undefined;
  setSubjectC: (value: string | undefined) => void;
  validityInMonths?: number;
  setValidityInMonths: (value: number | undefined) => void;
  keyStorePath?: string;
  setKeyStorePath: (value: string | undefined) => void;
  certUsage: CertificateUsage;
  setCertUsage: (value: CertificateUsage) => void;
}

export function useCertificateRequestFormState(): PolicyCertificateRequestFormProps {
  const { policy, namespaceProfile } = React.useContext(PolicyContext);

  const [issuerNamespaceId, setIssuerNamespaceId] = React.useState<
    string | undefined
  >();
  const [issuerPolicyId, setIssuerPolicyId] = React.useState<
    string | undefined
  >();
  const [subjectCN, setSubjectCN] = React.useState<string>();
  const [subjectOU, setSubjectOU] = React.useState<string>();
  const [subjectO, setSubjectO] = React.useState<string>();
  const [subjectC, setSubjectC] = React.useState<string>();
  const [validityInMonths, setValidityInMonths] = React.useState<number>();
  const [keyStorePath, setKeyStorePath] = React.useState<string>();
  const [certUsage, setCertUsage] = React.useState<CertificateUsage>(
    CertificateUsage.Usage_ServerAndClient
  );

  React.useEffect(() => {
    const certReq = policy?.certRequest;
    if (certReq) {
      setIssuerNamespaceId(certReq.issuerNamespaceId);
      if (certReq.issuerPolicyIdentifier) {
        setIssuerPolicyId(certReq.issuerPolicyIdentifier);
      } else {
        setIssuerPolicyId(undefined);
      }
      setSubjectCN(certReq.subject.cn);
      setSubjectOU(certReq.subject.ou);
      setSubjectO(certReq.subject.o);
      setSubjectC(certReq.subject.c);
      setValidityInMonths(certReq.validityMonths);
      setCertUsage(certReq.usage);
    }
  }, [policy]);

  return {
    issuerNamespaceId,
    setIssuerNamespaceId,
    issuerPolicyId,
    setIssuerPolicyId,
    subjectCN,
    setSubjectCN,
    subjectOU,
    setSubjectOU,
    subjectO,
    setSubjectO,
    subjectC,
    setSubjectC,
    validityInMonths,
    setValidityInMonths,
    keyStorePath,
    setKeyStorePath,
    certUsage,
    setCertUsage,
  };
}

export function PolicyCertificateRequestForm({
  issuerNamespaceId,
  setIssuerNamespaceId,
  issuerPolicyId,
  setIssuerPolicyId,
  subjectCN,
  setSubjectCN,
  subjectOU,
  setSubjectOU,
  subjectO,
  setSubjectO,
  subjectC,
  setSubjectC,
  validityInMonths,
  setValidityInMonths,
  keyStorePath,
  setKeyStorePath,
  certUsage,
  setCertUsage,
}: PolicyCertificateRequestFormProps) {
  const { policy, namespaceId, namespaceProfile } =
    React.useContext(PolicyContext);

  const certUsageInputOnChange = useMemoizedFn<
    (e: React.ChangeEvent<HTMLInputElement>) => void
  >((e) => {
    if (e.target.checked) {
      setCertUsage(e.target.value as any);
    }
  });

  const certUsageIsChecked = useMemoizedFn((usage: CertificateUsage) => {
    return usage === certUsage;
  });
  return (
    <div className="divide-y divide-neutral-200 space-y-6">
      <h2 className="text-2xl font-semibold">Certificate Request Policy</h2>
      <div className="pt-6 space-y-4">
        <IssuerSelector
          selectedIssuerId={issuerNamespaceId ?? ""}
          onChange={setIssuerNamespaceId}
        />
        {issuerNamespaceId && (
          <PolicySelector
            namespaceId={issuerNamespaceId}
            policyType={PolicyType.PolicyType_CertRequest}
            selectedPolicyId={issuerPolicyId ?? ""}
            onChange={setIssuerPolicyId}
            label="Select certificate policy"
          />
        )}
      </div>
      <div className="pt-6 space-y-4">
        <h3 className="text-base font-semibold leading-7 text-gray-900">
          Subject
        </h3>
        <InputField
          labelContent="Common Name (CN)"
          placeholder="Sample Common Name"
          required
          value={subjectCN ?? ""}
          onChange={setSubjectCN}
        />
        <InputField
          labelContent="Organizational Unit (OU)"
          placeholder="Sample Organizational Unit"
          value={subjectOU ?? ""}
          onChange={setSubjectOU}
        />
        <InputField
          labelContent="Organization (O)"
          placeholder="Sample Organization"
          value={subjectO ?? ""}
          onChange={setSubjectO}
        />
        <InputField
          labelContent="Country or Region (C)"
          placeholder="US"
          value={subjectC ?? ""}
          onChange={setSubjectC}
        />
      </div>
      <div className="pt-6 space-y-6">
        <InputField
          labelContent="Validity in months"
          type="number"
          inputMode="numeric"
          placeholder="12"
          value={validityInMonths}
          onChange={setValidityInMonths as any}
        />
        <InputField
          labelContent="Key Store Path"
          required
          value={keyStorePath}
          onChange={setKeyStorePath}
        />
      </div>
      {namespaceProfile?.objectType !==
        NamespaceType.NamespaceType_BuiltInCaRoot &&
        namespaceProfile?.objectType !==
          NamespaceType.NamespaceType_BuiltInCaInt && (
          <CertificateUsageSelector
            inputType="radio"
            onChange={certUsageInputOnChange}
            isChecked={certUsageIsChecked}
          />
        )}
    </div>
  );
}
