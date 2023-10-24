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
  defaultLocator: string | undefined;
  value?: CertPolicyRef;
  availableNamespaceProfiles: ProfileRef[] | undefined | null;
  onChange?: (value: CertPolicyRef | undefined) => void;
};

type NamespaceProfileOptionType = DefaultOptionType & {
  profile: ProfileRef;
};

export function CertificatePolicySelect({
  defaultLocator,
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
        value: v5(
          `https://example.com/v1/r/${profile.resourceKind}/${profile.resourceIdentifier}`,
          v5.URL
        ),
        profile,
      })
    );
  }, [availableNamespaceProfiles]);
  const [selectedProfileStorageNID, setSelectedProfileStorageNID] =
    useState<string>();

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    async () => {
      const selectedProfile = namespaceOptions.find(
        (option) => option.value === selectedProfileStorageNID
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
      refreshDeps: [selectedProfileStorageNID],
    }
  );

  const certPolicyOptions = useMemo(() => {
    return certPolicies?.map((policy) => ({
      label: (
        <span>
          {policy.displayName} ({policy.resourceIdentifier})
        </span>
      ),
      value: policy.id,
    }));
  }, [certPolicies]);

  const onChangeDerived = useMemoizedFn((value: string) => {
    const selected = certPolicies?.find((policy) => policy.id === value);
    onChange?.(selected);
  });

  const [defaultNID, defaultRID] = useMemo(() => {
    if (defaultLocator) {
      const [nid, rid] = defaultLocator.split(":");
      return [nid, rid];
    }
    return [];
  }, [defaultLocator]);

  return (
    <>
      <Form.Item label="Select certificate policy namespace">
        <Select
          options={namespaceOptions}
          value={selectedProfileStorageNID ?? defaultNID}
          onChange={setSelectedProfileStorageNID}
        />
      </Form.Item>
      <Form.Item label="Select certificate policy">
        <Select
          options={certPolicyOptions}
          value={value?.id ?? defaultRID}
          onChange={onChangeDerived}
        />
      </Form.Item>
    </>
  );
}
