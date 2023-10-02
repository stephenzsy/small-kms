import { useRequest } from "ahooks";
import { useId } from "react";
import { useParams } from "react-router-dom";
import { Button } from "../components/Button";
import {
  AdminApi,
  CertificateUsage,
  DirectoryApi,
  NamespaceTypeShortName,
  PolicyApi,
  ResponseError
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

interface IssuerNamespaceSelectorProps {
  requesterNamespace: string;
  client: DirectoryApi;
  selectedIssuerId: string;
  onChange: (issuerId: string) => void;
}

interface IssuerPolicySelectorProps {
  issuerNamespaceId: string;
  client: PolicyApi;
  selectedPolicyId: string;
  onChange: (policyId: string) => void;
}

function IssuerPolicySelector({
  issuerNamespaceId,
  client,
  selectedPolicyId,
  onChange,
}: IssuerPolicySelectorProps) {
  const { data: policies } = useRequest(
    async () => {
      return client.listPoliciesV1({ namespaceId: issuerNamespaceId });
    },
    {
      refreshDeps: [issuerNamespaceId],
    }
  );
  const selectId = useId();
  return (
    <div>
      <label
        htmlFor={selectId}
        className="block text-sm font-medium leading-6 text-gray-900"
      >
        Issuer policy
      </label>
      <select
        id={selectId}
        name="location"
        className="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6"
        value={selectedPolicyId}
        onChange={(e) => onChange(e.target.value)}
      >
        <option disabled value="">
          Select issuer policy
        </option>
        {policies?.map((p) => (
          <option key={p.id} value={p.id}>
            {p.id} ({p.policyType})
          </option>
        ))}
      </select>
    </div>
  );
}

interface CertCreatePolicyFormInput {
  subjectCN: string;
  subjectOU?: string;
  subjectO?: string;
  subjectC?: string;
  validityInMonths?: number;
  keyStorePath: string;
  usage: CertificateUsage;
}

export default function CertificatePage() {
  const { namespaceId, nsType, templateId, certId } = useParams() as {
    namespaceId: string;
    nsType: NamespaceTypeShortName;
    templateId: string;
    certId: string;
  };

  const adminApi = useAuthedClient(AdminApi);
  const {
    data: cert,
    loading,
    error: certError,
    run,
  } = useRequest(
    async (create?: boolean) => {
      try {
        if (create) {
          return await adminApi.issueCertificateByTemplateV2({
            namespaceId,
            namespaceType: nsType,
            templateId,
          });
        }
        return await adminApi.getCertificateTemplateV2({
          namespaceId,
          namespaceType: nsType,
          templateId,
        });
      } catch (e) {
        if (e instanceof ResponseError) {
          if (e.response.status === 404) {
            return null;
          }
        }
        throw e;
      }
    },
    {
      refreshDeps: [namespaceId, nsType, templateId, certId],
    }
  );
  return (
    <>
      <h1>
        {nsType}/{namespaceId}/certificate-templates/{templateId}/certificates/
        {certId}
      </h1>
      <div className="p-6 bg-white rounded-lg overflow-hidden shadow-sm">
        {certError && <div>Fetch cert has error</div>}
        {loading ? (
          <div>Loading...</div>
        ) : cert ? (
          <pre>{JSON.stringify(cert, undefined, 2)}</pre>
        ) : (
          "No cert"
        )}
      </div>
      <div className="p-6 bg-white rounded-lg overflow-hidden shadow-sm">
        <Button
          variant="primary"
          onClick={() => {
            run(true);
          }}
        >
          Apply Template
        </Button>
      </div>
    </>
  );
}
