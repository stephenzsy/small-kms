import Title from "antd/es/typography/Title";
import {
  AdminApi,
  CreateManagedAppRequest,
  ManagedAppRef,
  ProfileRef,
  ResourceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import {
  Button,
  Card,
  Form,
  Input,
  Select,
  SelectProps,
  Table,
  TableColumnType,
} from "antd";
import { useForm } from "antd/es/form/Form";
import { useMemo } from "react";
import { Link } from "../components/Link";
import { DefaultOptionType } from "antd/es/cascader";

type ProfileTypeMapRecord<T> = Record<
  typeof ResourceKind.ProfileResourceKindRootCA,
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
    async (req: CreateManagedAppRequest) => {
      await adminApi.createManagedApp(req);
      //onCreated();
      form.resetFields();
    },
    { manual: true }
  );

  const typeOptions = useProfileTypeSelectOptions();

  return (
    <Form
      form={form}
      onFinish={(values) => {
        if (values.displayName) {
          return run({
            managedAppParameters: {
              displayName: values.displayName.trim(),
            },
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
        required
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

function useColumns() {
  return useMemo<TableColumnType<ProfileRef>[]>(
    () => [
      {
        title: "Name",
        render: (r: ProfileRef) => (
          <span className="font-mono">{r.resourceIdentifier}</span>
        ),
      },
      {
        title: "Display name",
        render: (r: ProfileRef) => r.displayName,
      },
      {
        title: "Actions",
        render: (r) => <Link to={`/apps/${r.appId}`}>View</Link>,
      },
    ],
    []
  );
}

export default function CAsPage() {
  const adminApi = useAuthedClient(AdminApi);

  const { data: rootCAs, run: listApps } = useRequest(
    () => {
      return adminApi.listRootCAs();
    },
    {
      refreshDeps: [],
    }
  );

  const columns = useColumns();
  const onProfileUpsert: ProfileTypeMapRecord<() => void> = useMemo(() => {
    return {
      [ResourceKind.ProfileResourceKindRootCA]: listApps,
    };
  }, [listApps]);

  return (
    <>
      <Title>Certificate Authorities</Title>
      <Card title="Certificate Authorities">
        <Table<ProfileRef>
          columns={columns}
          dataSource={rootCAs}
          rowKey={(r) => r.resourceIdentifier}
        />
      </Card>
      <Card title="Create certificate authority profile">
        <CreateProfileForm onCreated={onProfileUpsert} />
      </Card>
    </>
  );
}
