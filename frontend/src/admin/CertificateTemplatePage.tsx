import React from "react";
import {
  AdminApi,
  CertificateTemplate,
  CertificateTemplateParameters,
  CertificateUsage,
  NamespaceType,
  NamespaceTypeShortName,
  PolicyType,
  ResponseError,
} from "../generated";
import { InputField } from "./InputField";
import { PolicyContext } from "./PolicyContext";
import { IssuerSelector, PolicySelector } from "./Selectors";
import { useMemoizedFn, useRequest } from "ahooks";
import { CertificateUsageSelector } from "./CertificateUsageSelector";
import { useParams } from "react-router-dom";
import { useAuthedClient } from "../utils/useCertsApi";
import { uuidNil } from "../constants";
import { Button } from "../components/Button";

export interface CertificateTemplateFormState {
  displayName: string;
  setDisplayName: (value: string) => void;
  issuerNamespaceId: string;
  setIssuerNamespaceId: (value: string) => void;
  issuerCertificateTemplateId: string | undefined;
  setIssuerCertificateTemplateId: (value: string | undefined) => void;
  subjectCN: string;
  setSubjectCN: (value: string) => void;
  subjectOU: string;
  setSubjectOU: (value: string) => void;
  subjectO: string;
  setSubjectO: (value: string) => void;
  subjectC: string;
  setSubjectC: (value: string) => void;
  validityInMonths: number;
  setValidityInMonths: (value: number) => void;
  keyStorePath: string;
  setKeyStorePath: (value: string) => void;
  certUsage: CertificateUsage;
  setCertUsage: (value: CertificateUsage) => void;
}

export function useCertificateTemplateFormState(
  certTemplate: CertificateTemplate | undefined
): CertificateTemplateFormState {
  const [displayName, setDisplayName] = React.useState<string>("");
  const [issuerNamespaceId, setIssuerNamespaceId] = React.useState<string>("");
  const [issuerCertificateTemplateId, setIssuerCertificateTemplateId] =
    React.useState<string | undefined>();
  const [subjectCN, setSubjectCN] = React.useState<string>("");
  const [subjectOU, setSubjectOU] = React.useState<string>("");
  const [subjectO, setSubjectO] = React.useState<string>("");
  const [subjectC, setSubjectC] = React.useState<string>("");
  const [validityInMonths, setValidityInMonths] = React.useState<number>(0);
  const [keyStorePath, setKeyStorePath] = React.useState<string>("");
  const [certUsage, setCertUsage] = React.useState<CertificateUsage>(
    CertificateUsage.Usage_ServerAndClient
  );

  React.useEffect(() => {
    if (certTemplate) {
      setDisplayName(certTemplate.ref.displayName);
      setIssuerNamespaceId(certTemplate.issuer.namespaceId);
      if (certTemplate.issuer.templateId) {
        setIssuerCertificateTemplateId(certTemplate.issuer.templateId);
      } else {
        setIssuerCertificateTemplateId(undefined);
      }
      setSubjectCN(certTemplate.subject.cn);
      setSubjectOU(certTemplate.subject.ou ?? "");
      setSubjectO(certTemplate.subject.o ?? "");
      setSubjectC(certTemplate.subject.c ?? "");
      setValidityInMonths(certTemplate.validityMonths ?? 0);
      setKeyStorePath(certTemplate.keyStorePath ?? "");
      setCertUsage(certTemplate.usage);
    }
  }, [certTemplate]);

  return {
    displayName,
    setDisplayName,
    issuerNamespaceId,
    setIssuerNamespaceId,
    issuerCertificateTemplateId,
    setIssuerCertificateTemplateId,
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

type CertificateTemplateFormProps = CertificateTemplateFormState & {
  nsType: NamespaceTypeShortName;
  templateId: string;
};

export function CertificateTemplatesForm({
  displayName,
  setDisplayName,
  issuerNamespaceId,
  setIssuerNamespaceId,
  issuerCertificateTemplateId,
  setIssuerCertificateTemplateId,
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
  nsType,
  templateId,
}: CertificateTemplateFormProps) {
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
      <h2 className="text-2xl font-semibold">Certificate template</h2>
      {templateId !== uuidNil && (
        <InputField
          className="pt-6"
          labelContent="Display name"
          placeholder=""
          required
          value={displayName}
          onChange={setDisplayName}
        />
      )}
      {nsType !== NamespaceTypeShortName.NSType_RootCA && (
        <div className="pt-6 space-y-4">
          <IssuerSelector
            selectedIssuerId={issuerNamespaceId ?? ""}
            onChange={setIssuerNamespaceId}
          />
          {issuerNamespaceId /*
            <PolicySelector
              namespaceId={issuerNamespaceId}
              policyType={PolicyType.PolicyType_CertRequest}
              selectedPolicyId={issuerPolicyId ?? ""}
              onChange={setIssuerPolicyId}
              label="Select certificate policy"
      /> */ && null}
        </div>
      )}
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
          required={nsType !== NamespaceTypeShortName.NSType_Group}
          value={keyStorePath}
          onChange={setKeyStorePath}
        />
      </div>
      {nsType !== NamespaceTypeShortName.NSType_RootCA &&
        nsType !== NamespaceTypeShortName.NSType_IntCA && (
          <CertificateUsageSelector
            inputType="radio"
            onChange={certUsageInputOnChange}
            isChecked={certUsageIsChecked}
          />
        )}
    </div>
  );
}

export default function CertificateTemplatePage() {
  const { nsType, namespaceId, templateId } = useParams() as {
    nsType: NamespaceTypeShortName;
    namespaceId: string;
    templateId: string;
  };

  const adminApi = useAuthedClient(AdminApi);
  const { data, loading, run } = useRequest(
    async (p?: CertificateTemplateParameters, updateDisplayName?: string) => {
      if (!p) {
        try {
          return await adminApi.getCertificateTemplateV2({
            namespaceType: nsType,
            namespaceId,
            templateId,
          });
        } catch (e) {
          if (e instanceof ResponseError && e.response.status === 404) {
            return undefined;
          }
          throw e;
        }
      } else {
        return await adminApi.putCertificateTemplateV2({
          namespaceType: nsType,
          namespaceId,
          templateId,
          certificateTemplateParameters: p,
          displayName:
            templateId === uuidNil ? "default" : updateDisplayName ?? "",
        });
      }
    },
    { refreshDeps: [nsType, namespaceId, templateId] }
  );
  const state = useCertificateTemplateFormState(data);

  const onSubmit = useMemoizedFn<React.FormEventHandler<HTMLFormElement>>(
    (e) => {
      e.preventDefault();
      run(
        {
          issuer: {
            namespaceId:
              nsType === "root-ca" ? namespaceId : state.issuerNamespaceId,
            namespaceType:
              nsType === "root-ca" || nsType === "intermediate-ca"
                ? NamespaceTypeShortName.NSType_RootCA
                : NamespaceTypeShortName.NSType_IntCA,
          },
          subject: {
            cn: state.subjectCN,
            ou: state.subjectOU || undefined,
            o: state.subjectO || undefined,
            c: state.subjectC || undefined,
          },
          usage:
            nsType === "root-ca"
              ? CertificateUsage.Usage_RootCA
              : nsType === "intermediate-ca"
              ? CertificateUsage.Usage_IntCA
              : state.certUsage,
          keyStorePath: state.keyStorePath,
          validityMonths: state.validityInMonths,
        },
        state.displayName
      );
      /*

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
        validityMonths: certReqState.validityInMonths
          ? parseInt(certReqState.validityInMonths.toString())
          : 0,
      }; */
    }
  );
  return (
    <>
      <h1>
        {nsType}/{namespaceId}/certificate-templates/{templateId}
      </h1>
      <div className="rounded-lg bg-white shadow p-6 space-y-6">
        <h2>Current policy</h2>
        {loading ? (
          <div>Loading...</div>
        ) : data ? (
          <pre>{JSON.stringify(data, undefined, 2)}</pre>
        ) : (
          <div>Not found</div>
        )}
      </div>
      <form
        className="divide-y-2 divide-neutral-500 overflow-hidden rounded-lg bg-white shadow p-6 space-y-6"
        onSubmit={onSubmit}
      >
        <CertificateTemplatesForm
          templateId={templateId}
          nsType={nsType}
          {...state}
        />
        <div className="pt-6 flex flex-row items-center gap-x-6 justify-end">
          <Button variant="primary" type="submit">
            Create or Update
          </Button>
        </div>
      </form>
    </>
  );
}
