import { useRequest } from "ahooks";
import {
  DirectoryApi,
  NamespaceType,
  PolicyApi,
  PolicyType,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import React, { useContext, useMemo } from "react";
import { policyTypeNames } from "./displayConstants";
import { PolicyContext } from "./PolicyContext";
import { WellknownId } from "../constants";

export interface ISelectorItem {
  value: string;
  title: React.ReactNode;
}

export function BaseSelector<T extends ISelectorItem = ISelectorItem>({
  value,
  onChange,
  label,
  placeholder,
  items,
}: {
  value: string;
  onChange: (v: string) => void;
  label: React.ReactNode;
  placeholder: React.ReactNode;
  items: readonly T[] | undefined;
}) {
  const selectId = React.useId();
  return (
    <div>
      <label
        htmlFor={selectId}
        className="block text-sm font-medium leading-6 text-gray-900"
      >
        {label}
      </label>
      <select
        id={selectId}
        name="location"
        className="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      >
        <option disabled value="">
          {placeholder}
        </option>
        {items?.map((item) => (
          <option key={item.value} value={item.value}>
            {item.title}
          </option>
        ))}
      </select>
    </div>
  );
}

export function PolicySelector({
  namespaceId,
  policyType,
  selectedPolicyId,
  onChange,
  label,
}: {
  namespaceId: string;
  policyType: PolicyType;
  selectedPolicyId: string;
  onChange: (policyId: string) => void;
  label?: React.ReactNode;
}) {
  const policyApi = useAuthedClient(PolicyApi);
  const { data: policies } = useRequest(
    async () => {
      const policies = await policyApi.listPoliciesV1({
        namespaceId: namespaceId,
      });
      return policies.filter((p) => p.policyType === policyType);
    },
    {
      refreshDeps: [namespaceId, policyType],
    }
  );
  const items = useMemo<ISelectorItem[] | undefined>(
    () =>
      policies?.map((p) => ({
        value: p.id,
        title: (
          <>
            {p.id} ({policyTypeNames[p.policyType]})
          </>
        ),
      })),
    [policies]
  );
  return (
    <BaseSelector
      items={items}
      label={label ?? "Select Policy"}
      placeholder="Select Policy"
      value={selectedPolicyId}
      onChange={onChange}
    />
  );
}

interface IssuerNamespaceSelectorProps {
  selectedIssuerId: string;
  onChange: (issuerId: string) => void;
}

export function IssuerSelector({
  selectedIssuerId,
  onChange,
}: IssuerNamespaceSelectorProps) {
  const { namespaceId, namespaceProfile } = useContext(PolicyContext);
  const client = useAuthedClient(DirectoryApi);
  const { data: issuers } = useRequest(
    async () => {
      switch (namespaceProfile?.objectType) {
        case NamespaceType.NamespaceType_BuiltInCaRoot:
          // force itself
          return [namespaceProfile];
        case NamespaceType.NamespaceType_BuiltInCaInt:
          const l = await client.listNamespacesV1({
            namespaceType: NamespaceType.NamespaceType_BuiltInCaRoot,
          });
          if (namespaceId === WellknownId.nsTestIntCa) {
            return l.filter((x) => x.id === WellknownId.nsTestRootCa);
          } else {
            return l.filter((x) => x.id === WellknownId.nsRootCa);
          }
      }
      return await client.listNamespacesV1({
        namespaceType: NamespaceType.NamespaceType_BuiltInCaInt,
      });
    },
    {
      ready: !!namespaceProfile,
      refreshDeps: [namespaceProfile?.id],
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
      value={selectedIssuerId}
      onChange={onChange}
    />
  );
}
