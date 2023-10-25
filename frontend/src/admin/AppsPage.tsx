import Title from "antd/es/typography/Title";
import {
  AdminApi,
  CreateManagedAppRequest,
  ImportProfileRequest,
  ManagedAppRef,
  ProfileRef,
  ResourceKind,
} from "../generated";
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

function useManagedAppsColumns() {
  return useMemo<TableColumnType<ProfileRef>[]>(
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
        render: (r) => <Link to={`/app/managed/${r.appId}`}>View</Link>,
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
        render: (r: ProfileRef) => (
          <Link to={`/app/${profileKind}/${r.resourceIdentifier}`}>View</Link>
        ),
      },
    ],
    [profileKind]
  );
}

type ImportProfileFormState = {
  objectId?: string;
};

function ImportProfileForm({
  onCreated,
  profileKind,
}: {
  onCreated: () => void;
  profileKind: ResourceKind;
}) {
  const [form] = useForm<ImportProfileFormState>();

  const adminApi = useAuthedClient(AdminApi);

  const { run } = useRequest(
    async (req: ImportProfileRequest) => {
      await adminApi.importProfile(req);
      onCreated();
      form.resetFields();
    },
    { manual: true }
  );

  return (
    <Form
      form={form}
      onFinish={(values) => {
        const objectId = values.objectId?.trim();
        if (objectId) {
          return run({
            profileResourceKind: profileKind,
            namespaceIdentifier: objectId,
          });
        }
      }}
    >
      <Form.Item<ImportProfileFormState>
        name="objectId"
        label="Microsoft Entra object ID"
        required
      >
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Import
        </Button>
      </Form.Item>
    </Form>
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
      <Card title="Managed applications">
        <Table<ProfileRef>
          columns={managedAppColumns}
          dataSource={managedApps}
          rowKey={(r) => r.resourceIdentifier}
        />
      </Card>
      <Card title="Create managed application">
        <CreateManagedAppForm onCreated={listManagedApps} />
      </Card>
      <Card title="Service principals">
        <Table<ProfileRef>
          columns={spColumns}
          dataSource={servicePrincipals}
          rowKey={(r) => r.resourceIdentifier}
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