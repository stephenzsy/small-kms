import { useMemoizedFn, useRequest, useSet } from "ahooks";
import React, { useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import { Button } from "../components/Button";
import { WellknownId, uuidNil } from "../constants";
import {
  AdminApi as AdminApiOld,
  NamespaceTypeShortName,
  ResponseError,
} from "../generated";
import {
  AdminApi,
  CertificateTemplate,
  CertificateTemplateParameters,
  CertificateUsage,
  NamespaceKind,
} from "../generated3";
import {
  ValueState,
  ValueStateMayBeFixed,
  useFixedValueState,
  useValueState,
} from "../utils/formStateUtils";
import { useAuthedClient as useAuthedClientOld } from "../utils/useCertsApi";
import { useAuthedClient } from "../utils/useCertsApi3";
import { CertificateUsageSelector } from "./CertificateUsageSelector";
import { InputField } from "./InputField";
import { RefsTable3 } from "./RefsTable";
import { BaseSelector } from "./Selectors";

export interface CertificateTemplateFormState {
  issuerNamespaceId: ValueStateMayBeFixed<string>;
  issuerTemplateId: ValueState<string>;
  subjectCN: ValueState<string>;
  validityInMonths: number;
  setValidityInMonths: (value: number) => void;
  keyStorePath: ValueStateMayBeFixed<string>;
  certUsages: ValueStateMayBeFixed<ReadonlySet<CertificateUsage>>;
}

export function useCertificateTemplateFormState(
  certTemplate: CertificateTemplate | undefined,
  nsType: NamespaceKind | undefined,
  nsId: string,
  templateId: string
): CertificateTemplateFormState {
  const [validityInMonths, setValidityInMonths] = React.useState<number>(0);

  const randKeyStoreSuffix = useMemo(() => {
    return uuidv4().substring(0, 8);
  }, [nsType, nsId, templateId]);

  const fixedIssuerNamespaceId = useMemo(() => {
    switch (nsType) {
      case NamespaceKind.NamespaceKindCaRoot:
        return nsId;
      case NamespaceKind.NamespaceKindCaInt:
        return nsId === WellknownId.nsTestIntCa
          ? WellknownId.nsTestRootCa
          : WellknownId.nsRootCa;
    }
    return undefined;
  }, [nsType, nsId]);

  const state: CertificateTemplateFormState = {
    issuerNamespaceId: useFixedValueState(
      useValueState(""),
      fixedIssuerNamespaceId
    ),
    issuerTemplateId: useValueState("default"),
    subjectCN: useValueState(""),
    validityInMonths,
    setValidityInMonths,
    keyStorePath: useFixedValueState(
      useValueState(`${nsType}-${randKeyStoreSuffix}`),
      nsType === NamespaceTypeShortName.NSType_Group ? "" : undefined
    ),
    certUsages: useFixedValueState(
      useValueState(
        (): ReadonlySet<CertificateUsage> =>
          new Set([
            CertificateUsage.CertUsageServerAuth,
            CertificateUsage.CertUsageClientAuth,
          ])
      ),
      nsType == NamespaceKind.NamespaceKindCaRoot
        ? new Set([
            CertificateUsage.CertUsageCA,
            CertificateUsage.CertUsageCARoot,
          ])
        : nsType == NamespaceKind.NamespaceKindCaInt
        ? new Set([CertificateUsage.CertUsageCA])
        : templateId == "default-ms-entra-client-creds"
        ? new Set([
            CertificateUsage.CertUsageServerAuth,
            CertificateUsage.CertUsageClientAuth,
          ])
        : undefined
    ),
  };

  React.useEffect(() => {
    if (certTemplate) {
      // state.issuerNamespaceId.onChange?.(certTemplate.issuer.namespaceId);
      // state.issuerTemplateId.onChange(
      //   certTemplate.issuer.templateId ?? uuidNil
      // );
      state.subjectCN.onChange(certTemplate.subjectCommonName);
      setValidityInMonths(certTemplate.validityMonths ?? 0);
      state.keyStorePath.onChange?.(certTemplate.keyStorePath ?? "");
      state.certUsages.onChange?.(new Set(certTemplate.usages));
    }
  }, [certTemplate]);

  return state;
}

type CertificateTemplateFormProps = CertificateTemplateFormState & {
  nsType: NamespaceTypeShortName;
  templateId: string;
  adminApi: AdminApiOld;
};

export function CertificateIssuerSelector({
  value,
  onChange,
  adminApi,
}: {
  adminApi: AdminApiOld;
  value: string;
  onChange: (value: string) => void;
}) {
  const { data: issuers } = useRequest(
    () => {
      return adminApi.listNamespacesByTypeV2({
        namespaceType: NamespaceTypeShortName.NSType_IntCA,
      });
    },
    {
      refreshDeps: [],
    }
  );
  const items = useMemo(
    () =>
      issuers?.map((issuer) => ({
        value: issuer.id,
        title: issuer.displayName || issuer.id,
      })),
    [issuers]
  );
  return (
    <BaseSelector
      items={items}
      label={"Issuer namespace"}
      placeholder="Select issuer namespace"
      value={value}
      onChange={onChange}
    />
  );
}

export function CertificateIssuerTemplateSelector({
  value,
  onChange,
  adminApi,
  issuerNsId,
}: {
  adminApi: AdminApiOld;
  value: string;
  onChange: (value: string) => void;
  issuerNsId: string;
}) {
  const { data: issuers } = useRequest(
    () => {
      return adminApi.listCertificateTemplatesV2({
        namespaceId: issuerNsId,
      });
    },
    {
      refreshDeps: [issuerNsId],
    }
  );
  const items = useMemo(
    () =>
      issuers
        ?.filter((issuer) => issuer.id !== uuidNil)
        .map((issuer) => ({
          value: issuer.id,
          title: issuer.displayName || issuer.id,
        })),
    [issuers]
  );
  return (
    <BaseSelector
      items={items}
      label={"Issuer template"}
      placeholder="Select issuer template"
      value={value}
      onChange={onChange}
      defaultItem={{ value: uuidNil, title: "default" }}
    />
  );
}

function SANList({
  sansState: [itemsSet, { add: addSan, remove: removeSan }],
  title,
}: {
  sansState: ReturnType<typeof useSet<string>>;
  title: React.ReactNode;
}) {
  const [tobeAdded, setToBeAdded] = useState("");
  const items = useMemo(() => [...itemsSet], [itemsSet]);
  return (
    <div className="p-6 space-y-4">
      <h4 className="text-lg font-medium">{title}</h4>
      {items.length > 0 ? (
        <ul
          role="list"
          className="divide-y mt-4 divide-neutral-200 ring-1 ring-neutral-200 "
        >
          {items.map((item) => (
            <li key={item} className="flex justify-between gap-x-6 p-4">
              <span>{item}</span>
              <Button
                onClick={() => {
                  removeSan(item);
                }}
              >
                Remove
              </Button>
            </li>
          ))}
        </ul>
      ) : (
        <div> No items </div>
      )}
      <div className="flex flex-row gap-8">
        <div className="flex-1 flex rounded-md shadow-sm ring-1 ring-inset ring-neutral-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
          <input
            className="block flex-1 border-0 bg-transparent py-1.5 px-em text-neutral-900 placeholder:text-neutral-400 focus:ring-0 sm:text-sm sm:leading-6"
            value={tobeAdded}
            onChange={(e) => {
              setToBeAdded(e.target.value);
            }}
          />
        </div>
        <Button
          variant="primary"
          className="flex-0"
          onClick={() => {
            if (tobeAdded) {
              addSan(tobeAdded);
              setToBeAdded("");
            }
          }}
        >
          Add
        </Button>
      </div>
    </div>
  );
}

export function CertificateTemplatesForm(props: CertificateTemplateFormProps) {
  const {
    validityInMonths,
    setValidityInMonths,
    certUsages,
    nsType,
    adminApi,
  } = props;
  const certUsageInputOnChange = useMemoizedFn<
    (e: React.ChangeEvent<HTMLInputElement>) => void
  >((e) => {
    certUsages.onChange?.((prev) => {
      const s = new Set([...prev]);
      if (e.target.checked) {
        s.add(e.target.value as CertificateUsage);
      } else {
        s.delete(e.target.value as CertificateUsage);
      }
      return s;
    });
  });

  const certUsageIsChecked = useMemoizedFn((usage: CertificateUsage) => {
    return certUsages.value.has(usage);
  });
  return (
    <div className="space-y-6 ">
      <h2 className="text-2xl font-semibold">Certificate template</h2>
      {nsType !== NamespaceTypeShortName.NSType_RootCA && (
        <div className="pt-6 space-y-4">
          {props.issuerNamespaceId.onChange && (
            <CertificateIssuerSelector
              adminApi={adminApi}
              value={props.issuerNamespaceId.value}
              onChange={props.issuerNamespaceId.onChange}
            />
          )}
          {props.issuerNamespaceId.value && (
            <CertificateIssuerTemplateSelector
              issuerNsId={props.issuerNamespaceId.value}
              adminApi={adminApi}
              value={props.issuerTemplateId.value}
              onChange={props.issuerTemplateId.onChange}
            />
          )}
        </div>
      )}
      <InputField
        labelContent="Subject common name (CN)"
        placeholder="Sample Common Name"
        required
        value={props.subjectCN.value}
        onChange={props.subjectCN.onChange}
      />
      <div className="space-y-2">
        <InputField
          labelContent="Validity in months"
          type="number"
          inputMode="numeric"
          placeholder="12"
          value={validityInMonths}
          onChange={setValidityInMonths as any}
        />
        <div className="text-neutral-500">0 as default</div>
      </div>
      {props.keyStorePath.onChange && (
        <InputField
          labelContent="Key Store Path"
          required={nsType !== NamespaceTypeShortName.NSType_Group}
          value={props.keyStorePath.value}
          onChange={props.keyStorePath.onChange}
        />
      )}
      {nsType !== NamespaceTypeShortName.NSType_RootCA &&
        nsType !== NamespaceTypeShortName.NSType_IntCA && (
          <CertificateUsageSelector
            inputType="checkbox"
            onChange={certUsageInputOnChange}
            isChecked={certUsageIsChecked}
          />
        )}
    </div>
  );
}

export default function CertificateTemplatePage() {
  const { namespaceId, templateId, profileType } = useParams() as {
    namespaceId: string;
    templateId: string;
    profileType: NamespaceKind;
  };

  const adminApiOld = useAuthedClientOld(AdminApiOld);
  const adminApi = useAuthedClient(AdminApi);

  const { data, loading, run } = useRequest(
    async (p?: CertificateTemplateParameters) => {
      if (!p) {
        try {
          return await adminApi.getCertificateTemplate({
            profileType,
            profileId: namespaceId,
            templateId,
          });
        } catch (e) {
          if (e instanceof ResponseError && e.response.status === 404) {
            return undefined;
          }
          throw e;
        }
      } else {
        adminApi;
        return await adminApi.putCertificateTemplate({
          profileId: namespaceId,
          templateId,
          profileType,
          certificateTemplateParameters: p,
        });
      }
    },
    { refreshDeps: [namespaceId, templateId] }
  );

  const { run: issueCert } = useRequest(
    async () => {
      await adminApi.issueCertificateFromTemplate({
        profileId: namespaceId,
        templateId,
        profileType,
      });
    },
    { manual: true }
  );

  const state = useCertificateTemplateFormState(
    data,
    profileType,
    namespaceId,
    templateId
  );

  const onSubmit = useMemoizedFn<React.FormEventHandler<HTMLFormElement>>(
    (e) => {
      e.preventDefault();

      run({
        subjectCommonName: state.subjectCN.value,
        usages: [...state.certUsages.value],
        keyStorePath: state.keyStorePath.value,
        validityMonths: state.validityInMonths,
      });
    }
  );

  const { data: issuedCertificates } = useRequest(
    () => {
      return adminApi.listCertificatesByTemplate({
        profileId: namespaceId,
        profileType,
        templateId: templateId,
      });
    },
    { refreshDeps: [namespaceId, templateId] }
  );

  return (
    <>
      <h1>
        {namespaceId}/certificate-templates/{templateId}
      </h1>
      <RefsTable3
        items={issuedCertificates}
        title="Issued certificates"
        tableActions={
          <div>
            <Link
              to={`/admin/${namespaceId}/certificate-templates/${templateId}/certificates/${uuidNil}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              View latest certificate
            </Link>
          </div>
        }
        columns={[
          {
            header: "Thumbprint",
            render: (r) => r.thumbprint ?? "",
            columnKey: "thumbprint",
          },
        ]}
        refActions={(ref) => (
          <Link
            to={`/admin/${namespaceId}/certificate/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
      />
      <div className="rounded-lg bg-white shadow p-6 space-y-6">
        <Button variant="primary" onClick={issueCert}>
          Request certificate
        </Button>
      </div>
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
          nsType={"root-ca"}
          adminApi={adminApiOld}
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
