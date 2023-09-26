import { XCircleIcon } from "@heroicons/react/24/outline";
import { useBoolean, useRequest } from "ahooks";
import classNames from "classnames";
import { useMemo, useState, useId } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router-dom";
import { WellknownId } from "../constants";
import {
  CertificateUsage,
  DirectoryApi,
  NamespaceType,
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
  policyTypeNames,
  isRootCaNamespace,
  nsDisplayNames,
} from "./displayConstants";
import { Button } from "../components/Button";
import AlertDialog from "../components/AlertDialog";

interface IssuerNamespaceSelectorProps {
  requesterNamespace: string;
  client: DirectoryApi;
  selectedIssuerId: string;
  onChange: (issuerId: string) => void;
}

function IssuerSelector({
  requesterNamespace,
  client,
  selectedIssuerId,
  onChange,
}: IssuerNamespaceSelectorProps) {
  const { data: issuers } = useRequest(
    async () => {
      // if is intCA or root ca, query root ca namespaces
      if (
        isRootCaNamespace(requesterNamespace) ||
        IsIntCaNamespace(requesterNamespace)
      ) {
        const l = await client.listNamespacesV1({
          namespaceType: NamespaceType.NamespaceType_BuiltInCaRoot,
        });
        if (isRootCaNamespace(requesterNamespace)) {
          return [l.find((x) => x.id === requesterNamespace)];
        } else {
          if (requesterNamespace === WellknownId.nsTestIntCa) {
            return [l.find((x) => x.id === WellknownId.nsTestRootCa)];
          } else {
            return [l.find((x) => x.id === WellknownId.nsRootCa)];
          }
        }
      } else {
        return await client.listNamespacesV1({
          namespaceType: NamespaceType.NamespaceType_BuiltInCaInt,
        });
      }
    },
    {
      refreshDeps: [requesterNamespace],
    }
  );
  const selectId = useId();
  return (
    <div>
      <label
        htmlFor={selectId}
        className="block text-sm font-medium leading-6 text-gray-900"
      >
        Issuer namespace
      </label>
      <select
        id={selectId}
        name="location"
        className="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6"
        value={selectedIssuerId}
        onChange={(e) => onChange(e.target.value)}
      >
        <option disabled value="">
          Select issuer namespace
        </option>
        {issuers?.map((issuer) => (
          <option key={issuer?.id} value={issuer?.id}>
            {issuer?.displayName}
          </option>
        ))}
      </select>
    </div>
  );
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
  const directoryApi = useAuthedClient(DirectoryApi);

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
      usage: CertificateUsage.Usage_ClientOnly,
    },
  });

  const [formInvalid, { setTrue: setFormInvalid, setFalse: clearFormInvalid }] =
    useBoolean(false);

  const { run: updatePolicy, loading: updatePolicyLoading } = useRequest(
    async (policyParameters: PolicyParameters) => {
      const policy = await client.putPolicyV1({
        namespaceId,
        policyIdentifier: policyId,
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

  const [issuerNamespaceIdOverride, issuerPolicyIdOverride] = useMemo(() => {
    switch (namespaceId) {
      case WellknownId.nsRootCa:
        return [WellknownId.nsRootCa, policyId];
      case WellknownId.nsTestRootCa:
        return [WellknownId.nsTestRootCa, policyId];
    }
    return [];
  }, [namespaceId, policyId]);

  const [selectedIssuerNamespaceId, setSelectedIssuerNamespaceId] =
    useState("");
  const [selectedIssuerPolicyId, setSelectedIssuerPolicyId] = useState("");

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
        issuerNamespaceId:
          issuerNamespaceIdOverride ?? selectedIssuerNamespaceId,
        issuerPolicyIdentifier:
          issuerPolicyIdOverride ?? selectedIssuerPolicyId,
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
          : input.usage,
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
        {!isRootCaNamespace(namespaceId) && (
          <div className="border-neutral-900/10 space-y-6">
            <h2 className="text-base font-semibold leading-7 text-gray-900">
              Issuer
            </h2>
            <IssuerSelector
              requesterNamespace={namespaceId}
              client={directoryApi}
              selectedIssuerId={selectedIssuerNamespaceId}
              onChange={(value) => {
                console.log(value);
                setSelectedIssuerNamespaceId(value);
              }}
            />
            {selectedIssuerNamespaceId && (
              <IssuerPolicySelector
                issuerNamespaceId={selectedIssuerNamespaceId}
                client={client}
                selectedPolicyId={selectedIssuerPolicyId}
                onChange={setSelectedIssuerPolicyId}
              />
            )}
          </div>
        )}
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
        {!isRootCaNamespace(namespaceId) && !IsIntCaNamespace(namespaceId) && (
          <fieldset>
            <legend className="text-base font-semibold text-gray-900">
              Certificate Usage
            </legend>
            <div className="space-y-4">
              <div className="flex items-center">
                <input
                  id="usage-client-and-server"
                  type="radio"
                  {...register("usage", {})}
                  value={CertificateUsage.Usage_ServerAndClient}
                  className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
                />
                <label
                  htmlFor="usage-client-and-server"
                  className="ml-3 block text-sm font-medium leading-6 text-gray-900"
                >
                  Server and client
                </label>
              </div>{" "}
              <div className="flex items-center">
                <input
                  id="usage-server-only"
                  type="radio"
                  {...register("usage", {})}
                  value={CertificateUsage.Usage_ServerOnly}
                  className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
                />
                <label
                  htmlFor="usage-server-only"
                  className="ml-3 block text-sm font-medium leading-6 text-gray-900"
                >
                  Server only
                </label>
              </div>
              <div className="flex items-center">
                <input
                  id="usage-client-only"
                  type="radio"
                  {...register("usage", {})}
                  value={CertificateUsage.Usage_ClientOnly}
                  className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
                />
                <label
                  htmlFor="usage-client-only"
                  className="ml-3 block text-sm font-medium leading-6 text-gray-900"
                >
                  Client only
                </label>
              </div>
            </div>
          </fieldset>
        )}
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
          policyIdentifier: policyId!,
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

  const { run: deletePolicy } = useRequest(
    async (purge?: boolean) => {
      const updated = await client.deletePolicyV1({
        namespaceId: namespaceId!,
        policyIdentifier: policyId!,
        purge: purge ?? false,
      });
      mutate(updated);
    },
    { manual: true }
  );

  const [
    disableDialogOpen,
    { setFalse: closeDisableDialog, setTrue: openDisableDialog },
  ] = useBoolean(false);
  const [
    deleteDialogOpen,
    { setFalse: closeDeleteDialog, setTrue: openDeleteDialog },
  ] = useBoolean(false);
  return (
    <>
      <h1 className="text-4xl font-semibold">Policy</h1>
      <div>{nsDisplayNames[namespaceId!] ?? namespaceId}</div>
      <div>{policyTypeNames[certCategory]}</div>
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
        {fetchedPolicy?.deleted ? (
          <>
            <Button color="danger" variant="soft" onClick={openDeleteDialog}>
              Delete
            </Button>
            <AlertDialog
              title="Delete policy"
              description="Are you sure you want to delete this policy? This action cannot be undone."
              open={deleteDialogOpen}
              onClose={closeDeleteDialog}
              okText="Delete policy"
              onOk={() => {
                deletePolicy(true);
                closeDeleteDialog();
              }}
            />
          </>
        ) : (
          <>
            <Button color="danger" variant="soft" onClick={openDisableDialog}>
              Disable
            </Button>
            <AlertDialog
              title="Disable policy"
              description="Are you sure you want to disable this policy?"
              open={disableDialogOpen}
              onClose={closeDisableDialog}
              okText="Disable policy"
              onOk={() => {
                deletePolicy();
                closeDisableDialog();
              }}
            />
          </>
        )}
        <button
          type="button"
          className={classNames(
            "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
            "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
          )}
          onClick={toggleForm}
        >
          {formOpen ? "Cancel" : "Update"}
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
