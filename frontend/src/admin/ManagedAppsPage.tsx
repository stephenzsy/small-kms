import Title from "antd/es/typography/Title";
import { AdminApi, CreateManagedAppRequest, ManagedAppRef } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, TableColumnType } from "antd";
import { useForm } from "antd/es/form/Form";
import { useMemo } from "react";
import { Link } from "../components/Link";

type CreateManagedAppFormState = {
  displayName?: string;
};

function CreateManagedAppForm({ onCreated }: { onCreated: () => void }) {
  const [form] = useForm<CreateManagedAppFormState>();

  const adminApi = useAuthedClient(AdminApi);

  const { run } = useRequest(
    async (req: CreateManagedAppRequest) => {
      await adminApi.createManagedApp(req);
      onCreated();
      form.resetFields();
    },
    { manual: true }
  );

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
      <Form.Item name="displayName" label="Display name" required>
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
  return useMemo<TableColumnType<ManagedAppRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: ManagedAppRef) => (
          <span className="font-mono">{r.appId}</span>
        ),
      },
      {
        title: "Display name",
        render: (r: ManagedAppRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r) => <Link to={`/apps/${r.appId}`}>View</Link>,
      },
    ],
    []
  );
}

export default function ManagedAppsPage() {
  const adminApi = useAuthedClient(AdminApi);

  const { data, run: listApps } = useRequest(
    () => {
      return adminApi.listManagedApps();
    },
    {
      refreshDeps: [],
    }
  );

  const columns = useColumns();

  return (
    <>
      <Title>Managed Applications</Title>
      <Card title="Managed Applications">
        <Table<ManagedAppRef>
          columns={columns}
          dataSource={data}
          rowKey={(r) => r.resourceIdentifier}
        />
      </Card>
      <Card title="Create managed application">
        <CreateManagedAppForm onCreated={listApps} />
      </Card>
    </>
  );
}
