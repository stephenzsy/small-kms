import { useRequest } from "ahooks";
import React, { useMemo } from "react";
import { PolicyApi, PolicyType } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { policyTypeNames } from "./displayConstants";

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
  defaultItem,
}: {
  value: string;
  onChange: (v: string) => void;
  label: React.ReactNode;
  placeholder: React.ReactNode;
  items: readonly T[] | undefined;
  defaultItem?: T;
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
        {defaultItem ? (
          <option value={defaultItem.value}>{defaultItem.title}</option>
        ) : (
          <option disabled value="">
            {placeholder}
          </option>
        )}
        {items?.map((item) => (
          <option key={item.value} value={item.value}>
            {item.title}
          </option>
        ))}
      </select>
    </div>
  );
}

export function GroupSelector({
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
