import Title from "antd/es/typography/Title";
import {
  AdminApi,
  CreateManagedAppRequest,
  ManagedAppRef,
  ProfileRef,
  ResourceKind,
  SystemAppName,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, TableColumnType } from "antd";
import { useForm } from "antd/es/form/Form";
import { useMemo } from "react";
import { Link } from "../components/Link";
import { ImportProfileForm } from "./ImportProfileForm";
import { SyncManagedApplicationForm } from "./forms/SyncManagedApplicationForm";
import { useSystemAppRequest } from "./forms/useSystemAppRequest";

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

function useProfileRefColumns(profileKind: ResourceKind) {
  return useMemo<TableColumnType<ProfileRef>[]>(
    () => [
      {
        title: "ID",
        render: (r: ProfileRef) => <span className="font-mono">{r.id}</span>,
      },
      {
        title: "Display name",
        render: (r: ProfileRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r: ProfileRef) => (
          <Link to={`/app/${profileKind}/${r.id}`}>View</Link>
        ),
      },
    ],
    [profileKind]
  );
}

function SystemAppsEntry({ systemAppName }: { systemAppName: SystemAppName }) {
  const adminApi = useAuthedClient(AdminApi);
  const { data: systemApp, run: syncSystemApp } =
    useSystemAppRequest(systemAppName);
  return (
    <div className="ring-1 ring-neutral-400 rounded-md p-4 space-y-2">
      <div>
        Type: <span className="font-mono">{systemAppName}</span>
      </div>
      <div>Display name: {systemApp?.displayName}</div>
      <div>
        Application ID (Client ID):{" "}
        <span className="font-mono">{systemApp?.id}</span>
      </div>
      <div className="flex items-center gap-4">
        {systemApp && <Link to={`/app/system/${systemAppName}`}>View</Link>}
        <Button
          type="link"
          onClick={() => {
            syncSystemApp(true);
          }}
        >
          Sync
        </Button>
      </div>
    </div>
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

  const { data: servicePrincipals, run: refreshServicePrincipals } = useRequest(
    () => {
      return adminApi.listProfiles({
        profileResourceKind: ResourceKind.ProfileResourceKindServicePrincipal,
      });
    },
    {
      refreshDeps: [],
    }
  );

  const managedAppColumns = useManagedAppsColumns();
  const spColumns = useProfileRefColumns(
    ResourceKind.ProfileResourceKindServicePrincipal
  );

  return (
    <>
      <Title>Applications</Title>
      <Card title="System applications">
        <div className="space-y-4">
          <SystemAppsEntry systemAppName={SystemAppName.SystemAppNameAPI} />
          <SystemAppsEntry systemAppName={SystemAppName.SystemAppNameBackend} />
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
      <Card title="Service principals">
        <Table<ProfileRef>
          columns={spColumns}
          dataSource={servicePrincipals}
          rowKey={(r) => r.id}
        />
      </Card>
      <Card title="Import service principal">
        <ImportProfileForm
          onCreated={refreshServicePrincipals}
          profileKind={ResourceKind.ProfileResourceKindServicePrincipal}
        />
      </Card>
    </>
  );
}
