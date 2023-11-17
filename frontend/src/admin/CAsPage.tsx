import { useMemoizedFn, useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  Select,
  Table,
  TableColumnType,
} from "antd";
import { DefaultOptionType } from "antd/es/cascader";
import { useForm } from "antd/es/form/Form";
import Title from "antd/es/typography/Title";
import { useMemo } from "react";
import { Link } from "../components/Link";
import {
  AdminApi,
  NamespaceKind,
  ProfileParameters,
  ProfileRef,
  PutProfileRequest,
  ResourceKind,
} from "../generated";
import {
  AdminApi as AdminApiV2,
  NamespaceProvider,
  Ref,
} from "../generated/apiv2";
import { useAuthedClient, useAuthedClientV2 } from "../utils/useCertsApi";
import { ResourceRefsTable } from "./tables/ResourceRefsTable";

type ProfileTypeMapRecord<T> = Record<
  | typeof ResourceKind.ProfileResourceKindRootCA
  | typeof ResourceKind.ProfileResourceKindIntermediateCA,
  T
>;

type CreateProfileFormState = {
  name?: string;
  profileType?: keyof ProfileTypeMapRecord<any>;
  displayName?: string;
};

function useProfileTypeSelectOptions(): Array<DefaultOptionType> {
  return useMemo(
    () => [
      {
        label: "Root CA",
        value: ResourceKind.ProfileResourceKindRootCA,
      },
      {
        label: "Intermediate CA",
        value: ResourceKind.ProfileResourceKindIntermediateCA,
      },
    ],
    []
  );
}

function CreateProfileForm({
  onCreated,
}: {
  onCreated: ProfileTypeMapRecord<() => void>;
}) {
  const [form] = useForm<CreateProfileFormState>();

  const adminApi = useAuthedClient(AdminApi);

  const { run } = useRequest(
    async (
      profileType: keyof ProfileTypeMapRecord<any>,
      name: string,
      params: ProfileParameters
    ) => {
      const req: PutProfileRequest = {
        namespaceId: name,
        profileParameters: params,
        profileResourceKind: profileType,
      };
      switch (profileType) {
        case ResourceKind.ProfileResourceKindRootCA:
        case ResourceKind.ProfileResourceKindIntermediateCA:
          await adminApi.putProfile(req);
          break;
      }
      onCreated[profileType]();
      form.resetFields();
    },
    { manual: true }
  );

  const typeOptions = useProfileTypeSelectOptions();

  return (
    <Form
      form={form}
      onFinish={(values) => {
        if (values.profileType && values.name?.trim()) {
          return run(values.profileType, values.name.trim(), {
            displayName: values.displayName?.trim(),
          });
        }
      }}
    >
      <Form.Item<CreateProfileFormState>
        label="Select profile type"
        name="profileType"
        required
      >
        <Select options={typeOptions} />
      </Form.Item>
      <Form.Item<CreateProfileFormState> name="name" label="Name" required>
        <Input />
      </Form.Item>
      <Form.Item<CreateProfileFormState>
        name="displayName"
        label="Display name"
      >
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Create
        </Button>
      </Form.Item>
    </Form>
  );
}

export default function CAsPage() {
  const adminApi = useAuthedClientV2(AdminApiV2);

  const { data: rootCAs, run: listRootCAs } = useRequest(
    () => {
      return adminApi.listProfiles({
        namespaceProvider: NamespaceProvider.NamespaceProviderRootCA,
      });
    },
    {
      refreshDeps: [],
    }
  );
  const { data: intermediateCAs, run: listIntermediateCAs } = useRequest(
    () => {
      return adminApi.listProfiles({
        namespaceProvider: NamespaceProvider.NamespaceProviderIntermediateCA,
      });
    },
    {
      refreshDeps: [],
    }
  );

  const onProfileUpsert: ProfileTypeMapRecord<() => void> = useMemo(() => {
    return {
      [ResourceKind.ProfileResourceKindRootCA]: listRootCAs,
      [ResourceKind.ProfileResourceKindIntermediateCA]: listIntermediateCAs,
    };
  }, [listRootCAs]);

  const renderRootCaActions = useMemoizedFn((item: Ref) => {
    return <Link to={`/${NamespaceProvider.NamespaceProviderRootCA}/${item.id}`}>View</Link>;
  });

  const renderIntCaactions = useMemoizedFn((item: Ref) => {
    return <Link to={`/${NamespaceProvider.NamespaceProviderIntermediateCA}/${item.id}`}>View</Link>;
  });

  return (
    <>
      <Title>Certificate Authorities</Title>
      <Card title="Root certificate authorities">
        <ResourceRefsTable<Ref>
          renderActions={renderRootCaActions}
          dataSource={rootCAs}
        />
      </Card>
      <Card title="Intermediate certificate authorities">
        <ResourceRefsTable<Ref>
          renderActions={renderIntCaactions}
          dataSource={intermediateCAs}
        />
      </Card>
      <Card title="Create certificate authority profile">
        <CreateProfileForm onCreated={onProfileUpsert} />
      </Card>
    </>
  );
}
