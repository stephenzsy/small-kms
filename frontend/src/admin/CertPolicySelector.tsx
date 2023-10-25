import { Form } from "antd";
import Select, { DefaultOptionType } from "antd/es/select";
import { useMemo, useState, useEffect } from "react";
import { v5 } from "uuid";
import {
  AdminApi,
  CertPolicyRef,
  NamespaceKind,
  ProfileRef,
  ResourceLocator1,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useMemoizedFn, useRequest } from "ahooks";

export type CertificatePolicySelectorProps = {
  value?: ResourceLocator1;
  availableNamespaceProfiles: ProfileRef[] | undefined | null;
  onChange?: (value: ResourceLocator1 | undefined) => void;
};

type NamespaceProfileOptionType = DefaultOptionType & {
  profile: ProfileRef;
};

export function CertificatePolicySelect({
  value,
  availableNamespaceProfiles,
  onChange,
}: CertificatePolicySelectorProps) {
  const namespaceOptions = useMemo(() => {
    if (!availableNamespaceProfiles) {
      return [];
    }
    return availableNamespaceProfiles.map(
      (profile): NamespaceProfileOptionType => ({
        label: (
          <span>
            {profile.displayName} ({profile.resourceKind}:{" "}
            {profile.resourceIdentifier})
          </span>
        ),
        value: `${profile.resourceKind}/${profile.resourceIdentifier}`,
        profile,
      })
    );
  }, [availableNamespaceProfiles]);
  const [_selectedProfileStorageNID, setSelectedProfileStorageNID] =
    useState<string>();

  const selectedProfileValue =
    _selectedProfileStorageNID ??
    (value ? `${value.namespaceKind}/${value.namespaceIdentifier}` : undefined);

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    async () => {
      const selectedProfile = namespaceOptions.find(
        (option) => option.value === selectedProfileValue
      )?.profile;
      if (!selectedProfile) {
        return;
      }
      return await adminApi.listCertPolicies({
        namespaceKind: selectedProfile.resourceKind as unknown as NamespaceKind,
        namespaceIdentifier: selectedProfile.resourceIdentifier,
      });
    },
    {
      refreshDeps: [selectedProfileValue],
    }
  );

  const certPolicyOptions = useMemo(() => {
    return certPolicies?.map((policy) => ({
      label: (
        <span>
          {policy.displayName} ({policy.resourceIdentifier})
        </span>
      ),
      value: policy.resourceIdentifier,
    }));
  }, [certPolicies]);

  const onChangeDerived = useMemoizedFn((v: string) => {
    const selected = certPolicies?.find(
      (policy) => policy.resourceIdentifier === v
    );
    onChange?.(selected);
  });

  return (
    <>
      <Form.Item label="Select certificate policy namespace">
        <Select
          options={namespaceOptions}
          value={selectedProfileValue}
          onChange={setSelectedProfileStorageNID}
        />
      </Form.Item>
      <Form.Item label="Select certificate policy">
        <Select
          options={certPolicyOptions}
          value={value?.resourceIdentifier}
          onChange={onChangeDerived}
        />
      </Form.Item>
    </>
  );
}
