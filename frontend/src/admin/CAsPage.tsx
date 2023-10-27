import { useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  Select,
  Table,
  TableColumnType
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
  ResourceKind
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

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
      const req = {
        namespaceIdentifier: name,
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

function useColumns(nsKind: NamespaceKind) {
  return useMemo<TableColumnType<ProfileRef>[]>(
    () => [
      {
        title: "Name",
        render: (r: ProfileRef) => (
          <span className="font-mono">{r.id}</span>
        ),
      },
      {
        title: "Display name",
        render: (r: ProfileRef) => r.displayName,
      },
      {
        title: "Actions",
        render: (r: ProfileRef) => (
          <Link to={`/ca/${nsKind}/${r.id}`}>View</Link>
        ),
      },
    ],
    []
  );
}

export default function CAsPage() {
  const adminApi = useAuthedClient(AdminApi);

  const { data: rootCAs, run: listRootCAs } = useRequest(
    () => {
      return adminApi.listProfiles({
        profileResourceKind: ResourceKind.ProfileResourceKindRootCA,
      });
    },
    {
      refreshDeps: [],
    }
  );
  const { data: intermediateCAs, run: listIntermediateCAs } = useRequest(
    () => {
      return adminApi.listProfiles({
        profileResourceKind: ResourceKind.ProfileResourceKindIntermediateCA,
      });
    },
    {
      refreshDeps: [],
    }
  );
  const rootColumns = useColumns(NamespaceKind.NamespaceKindRootCA);
  const intColumns = useColumns(NamespaceKind.NamespaceKindIntermediateCA);
  const onProfileUpsert: ProfileTypeMapRecord<() => void> = useMemo(() => {
    return {
      [ResourceKind.ProfileResourceKindRootCA]: listRootCAs,
      [ResourceKind.ProfileResourceKindIntermediateCA]: listIntermediateCAs,
    };
  }, [listRootCAs]);

  return (
    <>
      <Title>Certificate Authorities</Title>
      <Card title="Root certificate authorities">
        <Table<ProfileRef>
          columns={rootColumns}
          dataSource={rootCAs}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Intermediate certificate authorities">
        <Table<ProfileRef>
          columns={intColumns}
          dataSource={intermediateCAs}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Create certificate authority profile">
        <CreateProfileForm onCreated={onProfileUpsert} />
      </Card>
    </>
  );
}
