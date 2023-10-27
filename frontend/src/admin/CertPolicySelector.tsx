import { Form } from "antd";
import Select, { DefaultOptionType } from "antd/es/select";
import { useMemo, useState } from "react";
import { ProfileRef, ResourceKind } from "../generated";

export type CertificatePolicySelectorProps = {
  value?: string;
  availableNamespaceProfiles: ProfileRef[] | undefined | null;
  onChange?: (value: string | undefined) => void;
  profileKind: ResourceKind;
};

type NamespaceProfileOptionType = DefaultOptionType & {
  profile: ProfileRef;
};

export function CertificateIssuerNamespaceSelect({
  value,
  availableNamespaceProfiles,
  onChange,
  profileKind,
}: CertificatePolicySelectorProps) {
  const namespaceOptions = useMemo(() => {
    if (!availableNamespaceProfiles) {
      return [];
    }
    return availableNamespaceProfiles.map(
      (profile): NamespaceProfileOptionType => ({
        label: (
          <span>
            {profile.displayName} ({profileKind}:{profile.id})
          </span>
        ),
        value: profile.id,
        profile,
      })
    );
  }, [availableNamespaceProfiles]);
  const [_selectedProfileStorageNID, setSelectedProfileStorageNID] =
    useState<string>();

  return (
    <>
      <Form.Item label="Select issuer namespace">
        <Select options={namespaceOptions} value={value} onChange={onChange} />
      </Form.Item>
    </>
  );
}
