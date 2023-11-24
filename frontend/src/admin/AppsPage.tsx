import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, TableColumnType } from "antd";
import { useForm } from "antd/es/form/Form";
import Title from "antd/es/typography/Title";
import { useMemo } from "react";
import { Link } from "../components/Link";
import {
  AdminApi,
  CreateManagedAppRequest,
  ManagedAppRef,
  ProfileRef
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { SyncManagedApplicationForm } from "./forms/SyncManagedApplicationForm";

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

function useManagedAppsColumns() {
  return useMemo<TableColumnType<ProfileRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: ManagedAppRef) => <span className="font-mono">{r.id}</span>,
      },
      {
        title: "Display name",
        render: (r: ManagedAppRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r) => <Link to={`/app/managed/${r.id}`}>View</Link>,
      },
    ],
    []
  );
}

export default function ManagedAppsPage() {
  const adminApi = useAuthedClient(AdminApi);

  const { data: managedApps, run: listManagedApps } = useRequest(
    () => {
      return adminApi.listManagedApps();
    },
    {
      refreshDeps: [],
    }
  );

  const managedAppColumns = useManagedAppsColumns();

  return (
    <>
      <Title>Applications</Title>
      <Card title="System applications">
        <div className="space-y-4">
          <Link className="block" to={`/app/system/default/provision-agent`}>
            Agent global configurations
          </Link>
          <Link className="block" to={`/app/system/default/radius-config`}>
            Radius global configurations
          </Link>
        </div>
      </Card>
      <Card title="Managed applications">
        <Table<ProfileRef>
          columns={managedAppColumns}
          dataSource={managedApps}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Create new managed application">
        <CreateManagedAppForm onCreated={listManagedApps} />
      </Card>
      <Card title="Sync existing managed application">
        <SyncManagedApplicationForm onSynced={listManagedApps} />
      </Card>
    </>
  );
}
