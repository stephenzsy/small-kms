import { XCircleIcon } from "@heroicons/react/24/outline";
import { useBoolean, useRequest } from "ahooks";
import classNames from "classnames";
import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router-dom";
import { WellknownId } from "../constants";
import {
  CertificateUsage,
  Policy,
  PolicyApi,
  PolicyParameters,
  PolicyType,
  ResponseError,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { InputField } from "./FormComponents";
import {
  IsIntCaNamespace,
  certRequestPolicyNames,
  isRootCaNamespace,
  nsDisplayNames,
} from "./displayConstants";

interface CertCreatePolicyFormInput {
  subjectCN: string;
  subjectOU?: string;
  subjectO?: string;
  subjectC?: string;
  validityInMonths?: number;
  keyStorePath: string;
}

function CertCreatePolicyForm({
  client,
  namespaceId,
  policyId,
  onPolicyMutate,
}: {
  client: PolicyApi;
  namespaceId: string;
  policyId: string;
  onPolicyMutate?: (policy: Policy) => void;
}) {
  const defaultKeyStorePath = useMemo(() => {
    const id = crypto.randomUUID();
    let prefix = "cert-";
    switch (namespaceId) {
      case WellknownId.nsRootCa:
        prefix = "root-ca-";
        break;
      case WellknownId.nsTestRootCa:
        prefix = "test-root-ca-";
        break;
      case WellknownId.nsIntCaIntranet:
        prefix = "int-ca-intranet-";
        break;
      case WellknownId.nsTestIntCa:
        prefix = "test-int-ca-";
        break;
    }
    return prefix + id.substring(0, 6);
  }, [namespaceId]);

  const { register, handleSubmit } = useForm<CertCreatePolicyFormInput>({
    defaultValues: {
      keyStorePath: defaultKeyStorePath,
    },
  });

  const [formInvalid, { setTrue: setFormInvalid, setFalse: clearFormInvalid }] =
    useBoolean(false);

  const { run: updatePolicy, loading: updatePolicyLoading } = useRequest(
    async (policyParameters: PolicyParameters) => {
      const policy = await client.putPolicyV1({
        namespaceId,
        policyId,
        policyParameters,
      });
      onPolicyMutate?.(policy);
      return policy;
    },
    { manual: true }
  );

  const defaultValidityPlaceholder = useMemo(() => {
    if (isRootCaNamespace(namespaceId)) {
      return 120;
    }
    return 12;
  }, [namespaceId]);

  const onSubmit = (input: CertCreatePolicyFormInput) => {
    let validityMonths: number | undefined = undefined;
    try {
      validityMonths = parseInt(input.validityInMonths as any);
    } catch {}
    if ((validityMonths ?? 0) < 0 || (validityMonths ?? 0) > 120) {
      validityMonths = undefined;
    }
    updatePolicy({
      policyType: PolicyType.PolicyType_CertRequest,
      certRequest: {
        issuerNamespaceId: policyId,
        subject: {
          cn: input.subjectCN,
          ou: input.subjectOU?.trim() || undefined,
          o: input.subjectO?.trim() || undefined,
          c: input.subjectC?.trim() || undefined,
        },
        usage: isRootCaNamespace(namespaceId)
          ? CertificateUsage.Usage_RootCA
          : IsIntCaNamespace(namespaceId)
          ? CertificateUsage.Usage_IntCA
          : CertificateUsage.Usage_ClientOnly,
        validityMonths,
        keyStorePath: input.keyStorePath,
      },
    });
  };
  return (
    <form
      className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
      onSubmit={handleSubmit(onSubmit, setFormInvalid)}
    >
      <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
        <div className="border-neutral-900/10 space-y-6">
          <h2 className="text-base font-semibold leading-7 text-gray-900">
            Subject
          </h2>

          <InputField
            inputKey="subjectCN"
            labelContent="Common Name (CN)"
            register={register}
            placeholder="Sample Internal Root CA"
            required
          />
          <InputField
            inputKey="subjectOU"
            labelContent="Organizational Unit (OU)"
            register={register}
            placeholder="Sample Organizational Unit"
          />

          <InputField
            inputKey="subjectO"
            labelContent="Organization (O)"
            register={register}
            placeholder="Sample Organization"
          />

          <InputField
            inputKey="subjectC"
            labelContent="Country or Region (C)"
            register={register}
            placeholder="US"
          />
        </div>
        <InputField
          inputKey="validityInMonths"
          labelContent="Validity in months"
          register={register}
          type="number"
          placeholder={defaultValidityPlaceholder.toString()}
        />
        <InputField
          inputKey="keyStorePath"
          labelContent="Key Store Path"
          register={register}
          type="text"
          required
        />
      </div>
      {formInvalid && (
        <div className="bg-red-50 px-4 py-4 sm:px-6 ">
          <div className="flex items-center gap-x-2">
            <div className="flex-shrink-0">
              <XCircleIcon
                className="h-5 w-5 text-red-400"
                aria-hidden="true"
              />
            </div>
            <div>
              <h3 className="text-sm font-medium text-red-800">
                Invalid input, please correect before proceeding
              </h3>
            </div>
          </div>
        </div>
      )}

      <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center justify-end gap-x-6">
        <button
          type="submit"
          disabled={updatePolicyLoading}
          className={classNames(
            "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
            "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
          )}
        >
          {updatePolicyLoading ? "Creating..." : "Create or update"}
        </button>
      </div>
    </form>
  );
}

export default function PolicyPage() {
  const { namespaceId, policyId } = useParams();
  const certCategory = useMemo(() => {
    return PolicyType.PolicyType_CertRequest;
  }, []);
  const client = useAuthedClient(PolicyApi);
  const {
    data: fetchedPolicy,
    run: refresh,
    mutate,
  } = useRequest(
    async () => {
      try {
        return await client.getPolicyV1({
          namespaceId: namespaceId!,
          policyId: policyId!,
        });
      } catch (e) {
        if (e instanceof ResponseError && e.response.status === 404) {
          return null;
        }
        throw e;
      }
    },
    { refreshDeps: [] }
  );

  const { run: applyPolicy } = useRequest(
    async () => {
      await client.applyPolicyV1({
        namespaceId: namespaceId!,
        policyId: policyId!,
        applyPolicyRequest: {},
      });
    },
    { manual: true }
  );

  const [formOpen, { toggle: toggleForm, setFalse: closeForm }] =
    useBoolean(false);

  return (
    <>
      <h1 className="text-4xl font-semibold">Policy</h1>
      <div>{nsDisplayNames[namespaceId!] ?? namespaceId}</div>
      <div>{certRequestPolicyNames[certCategory]}</div>
      {fetchedPolicy !== undefined && fetchedPolicy ? (
        <div>
          <pre className="text-sm">
            {JSON.stringify(fetchedPolicy, undefined, 2)}
          </pre>
        </div>
      ) : (
        <div>No policy</div>
      )}
      <div className="flex flex-row items-center gap-x-6">
        <button
          type="button"
          className={classNames(
            "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
            "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
          )}
          onClick={toggleForm}
        >
          {formOpen ? "Candel" : "Update"}
        </button>
        {fetchedPolicy && (
          <button
            type="button"
            className={classNames(
              "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
              "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
            )}
            onClick={applyPolicy}
          >
            Apply policy
          </button>
        )}
      </div>
      {formOpen && certCategory === PolicyType.PolicyType_CertRequest && (
        <CertCreatePolicyForm
          onPolicyMutate={(p) => {
            mutate(p);
            closeForm();
          }}
          client={client}
          namespaceId={namespaceId!}
          policyId={policyId!}
        />
      )}
    </>
  );
}
