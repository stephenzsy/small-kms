import { useMemoizedFn, useRequest, useSet } from "ahooks";
import React, { useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { v4 as uuidv4, version as uuidversion } from "uuid";
import { Button } from "../components/Button";
import { WellknownId, uuidNil } from "../constants";
import {
  AdminApi,
  CertificateTemplate,
  CertificateTemplateParameters,
  CertificateUsage,
  NamespaceTypeShortName,
  ResponseError,
} from "../generated";
import {
  ValueState,
  ValueStateMayBeFixed,
  useFixedValueState,
  useValueState,
} from "../utils/formStateUtils";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertificateUsageSelector } from "./CertificateUsageSelector";
import { InputField } from "./InputField";
import { RefsTable } from "./RefsTable";
import { BaseSelector } from "./Selectors";

export interface CertificateTemplateFormState {
  displayName: ValueStateMayBeFixed<string>;
  issuerNamespaceId: ValueStateMayBeFixed<string>;
  issuerTemplateId: ValueState<string>;
  subjectCN: ValueState<string>;
  subjectOU: string;
  setSubjectOU: (value: string) => void;
  subjectO: string;
  setSubjectO: (value: string) => void;
  subjectC: string;
  setSubjectC: (value: string) => void;
  validityInMonths: number;
  setValidityInMonths: (value: number) => void;
  keyStorePath: ValueStateMayBeFixed<string>;
  certUsage: CertificateUsage;
  setCertUsage: (value: CertificateUsage) => void;
  sanURIs: ReturnType<typeof useSet<string>>;
}

export function useCertificateTemplateFormState(
  certTemplate: CertificateTemplate | undefined,
  nsType: NamespaceTypeShortName,
  nsId: string,
  templateId: string
): CertificateTemplateFormState {
  const [subjectOU, setSubjectOU] = React.useState<string>("");
  const [subjectO, setSubjectO] = React.useState<string>("");
  const [subjectC, setSubjectC] = React.useState<string>("");
  const [validityInMonths, setValidityInMonths] = React.useState<number>(0);
  const [certUsage, setCertUsage] = React.useState<CertificateUsage>(
    CertificateUsage.Usage_ServerAndClient
  );

  const randKeyStoreSuffix = useMemo(() => {
    return uuidv4().substring(0, 8);
  }, [nsType, nsId, templateId]);

  const fixedIssuerNamespaceId = useMemo(() => {
    switch (nsType) {
      case NamespaceTypeShortName.NSType_RootCA:
        return nsId;
      case NamespaceTypeShortName.NSType_IntCA:
        return nsId === WellknownId.nsTestIntCa
          ? WellknownId.nsTestRootCa
          : WellknownId.nsRootCa;
    }
    return undefined;
  }, [nsType, nsId]);

  const state = {
    displayName: useFixedValueState(
      useValueState(""),
      uuidversion(templateId) === 4 ? undefined : "default"
    ),
    issuerNamespaceId: useFixedValueState(
      useValueState(""),
      fixedIssuerNamespaceId
    ),
    issuerTemplateId: useValueState(uuidNil),
    subjectCN: useValueState(""),
    subjectOU,
    setSubjectOU,
    subjectO,
    setSubjectO,
    subjectC,
    setSubjectC,
    validityInMonths,
    setValidityInMonths,
    keyStorePath: useFixedValueState(
      useValueState(`${nsType}-${randKeyStoreSuffix}`),
      nsType === NamespaceTypeShortName.NSType_Group ? "" : undefined
    ),
    certUsage,
    setCertUsage,
    sanURIs: useSet<string>(),
  };

  React.useEffect(() => {
    if (certTemplate) {
      state.displayName.onChange?.(
        certTemplate.ref.metadata?.["displayName"] ?? ""
      );
      state.issuerNamespaceId.onChange?.(certTemplate.issuer.namespaceId);
      state.issuerTemplateId.onChange(
        certTemplate.issuer.templateId ?? uuidNil
      );
      state.subjectCN.onChange(certTemplate.subject.cn);
      setSubjectOU(certTemplate.subject.ou ?? "");
      setSubjectO(certTemplate.subject.o ?? "");
      setSubjectC(certTemplate.subject.c ?? "");
      setValidityInMonths(certTemplate.validityMonths ?? 0);
      state.keyStorePath.onChange?.(certTemplate.keyStorePath ?? "");
      setCertUsage(certTemplate.usage);
      const [sanUriSet, { add: addSanUri, remove: removeSanUri }] =
        state.sanURIs;
      for (const san of sanUriSet) {
        if (!certTemplate.subjectAlternativeNames?.uris?.includes(san)) {
          removeSanUri(san);
        }
      }
      for (const san of certTemplate.subjectAlternativeNames?.uris ?? []) {
        addSanUri(san);
      }
    }
  }, [certTemplate]);

  return state;
}

type CertificateTemplateFormProps = CertificateTemplateFormState & {
  nsType: NamespaceTypeShortName;
  templateId: string;
  adminApi: AdminApi;
};

export function CertificateIssuerSelector({
  value,
  onChange,
  adminApi,
}: {
  adminApi: AdminApi;
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
        title: issuer.metadata?.["displayName"] || issuer.id,
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
  issuerNsType,
  issuerNsId,
}: {
  adminApi: AdminApi;
  value: string;
  onChange: (value: string) => void;
  issuerNsType: NamespaceTypeShortName;
  issuerNsId: string;
}) {
  const { data: issuers } = useRequest(
    () => {
      return adminApi.listCertificateTemplatesV2({
        namespaceType: issuerNsType,
        namespaceId: issuerNsId,
      });
    },
    {
      refreshDeps: [issuerNsType, issuerNsId],
    }
  );
  const items = useMemo(
    () =>
      issuers
        ?.filter((issuer) => issuer.id !== uuidNil)
        .map((issuer) => ({
          value: issuer.id,
          title: issuer.metadata?.["displayName"] || issuer.id,
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
    subjectOU,
    setSubjectOU,
    subjectO,
    setSubjectO,
    subjectC,
    setSubjectC,
    validityInMonths,
    setValidityInMonths,
    certUsage,
    setCertUsage,
    nsType,
    adminApi,
  } = props;
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
      {props.displayName.onChange && (
        <InputField
          className="pt-6"
          labelContent="Display name"
          placeholder=""
          required
          value={props.displayName.value}
          onChange={props.displayName.onChange}
        />
      )}
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
              issuerNsType={
                nsType === "intermediate-ca" ? "root-ca" : "intermediate-ca"
              }
              issuerNsId={props.issuerNamespaceId.value}
              adminApi={adminApi}
              value={props.issuerTemplateId.value}
              onChange={props.issuerTemplateId.onChange}
            />
          )}
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
          value={props.subjectCN.value}
          onChange={props.subjectCN.onChange}
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
      <div>
        <h3 className="text-base font-semibold leading-7 text-gray-900">
          Subject Alternative Names (SANs)
        </h3>
        <SANList sansState={props.sanURIs} title="URIs" />
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
        {props.keyStorePath.onChange && (
          <InputField
            labelContent="Key Store Path"
            required={nsType !== NamespaceTypeShortName.NSType_Group}
            value={props.keyStorePath.value}
            onChange={props.keyStorePath.onChange}
          />
        )}
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
    async (p?: CertificateTemplateParameters) => {
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
        });
      }
    },
    { refreshDeps: [nsType, namespaceId, templateId] }
  );
  const state = useCertificateTemplateFormState(
    data,
    nsType,
    namespaceId,
    templateId
  );

  const onSubmit = useMemoizedFn<React.FormEventHandler<HTMLFormElement>>(
    (e) => {
      e.preventDefault();
      const [sanUrisSet] = state.sanURIs;
      const uris = sanUrisSet.size > 0 ? [...sanUrisSet] : undefined;
      run({
        displayName: state.displayName.value,
        issuer: {
          namespaceId: state.issuerNamespaceId.value,
          namespaceType:
            nsType === "root-ca" || nsType === "intermediate-ca"
              ? NamespaceTypeShortName.NSType_RootCA
              : NamespaceTypeShortName.NSType_IntCA,
          templateId: state.issuerTemplateId.value,
        },
        subject: {
          cn: state.subjectCN.value,
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
        keyStorePath: state.keyStorePath.value,
        subjectAlternativeNames: uris && {
          uris,
        },
        validityMonths: state.validityInMonths,
      });
    }
  );

  const { data: issuedCertificates } = useRequest(
    () => {
      return adminApi.listCertificatesV2({
        namespaceId,
        namespaceType: nsType,
        templateId: templateId,
      });
    },
    { refreshDeps: [nsType, namespaceId, templateId] }
  );

  return (
    <>
      <h1>
        {nsType}/{namespaceId}/certificate-templates/{templateId}
      </h1>
      <RefsTable
        items={issuedCertificates}
        title="Issued certificates"
        tableActions={
          <div>
            <Link
              to={`/admin/${nsType}/${namespaceId}/certificate-templates/${templateId}/certificates/${uuidNil}`}
              className="text-indigo-600 hover:text-indigo-900"
            >
              View latest certificate
            </Link>
          </div>
        }
        itemTitleMetadataKey="thumbprint"
        refActions={(ref) => (
          <Link
            to={`/admin/${ref.namespaceType}/${ref.namespaceId}/certificate-templates/${ref.id}`}
            className="text-indigo-600 hover:text-indigo-900"
          >
            View
          </Link>
        )}
      />
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
          adminApi={adminApi}
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
