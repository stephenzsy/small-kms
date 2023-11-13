import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Table } from "antd";
import { ColumnsType } from "antd/es/table";
import { useContext, useMemo } from "react";
import { Link } from "../components/Link";
import { AdminApi, NamespaceKind, ResourceReference } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { CertPolicyRefTable } from "./CertPolicyRefTable";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { CertificatesTableCard } from "./tables/CertificatesTableCard";

type MemberOfGroupFormState = {
  groupId: string;
};

function MemberOfGroupForm() {
  const [form] = Form.useForm<MemberOfGroupFormState>();
  const api = useAuthedClient(AdminApi);
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const { run: addMember } = useRequest(
    (groupId: string) => {
      return api.syncGroupMemberOf({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
        groupId: groupId,
      });
    },
    {
      manual: true,
      onSuccess: () => {
        form.resetFields();
      },
    }
  );
  return (
    <Form<MemberOfGroupFormState>
      form={form}
      onFinish={(values) => {
        if (values.groupId) {
          addMember(values.groupId);
        }
      }}
    >
      <Form.Item<MemberOfGroupFormState>
        label="Group Microsoft Entra Object ID"
        name="groupId"
        required
      >
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Add
        </Button>
      </Form.Item>
    </Form>
  );
}

function useGroupMemberOfColumns() {
  return useMemo((): ColumnsType<ResourceReference> => {
    return [
      {
        title: "ID",
        render: (r: ResourceReference) => (
          <span className="font-mono">{r.id}</span>
        ),
      },
      {
        title: "Actions",
        render: (r: ResourceReference) => (
          <Link to={`/entra/group/${r.id}`}>View</Link>
        ),
      },
    ];
  }, []);
}

export default function NamespacePage() {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: groupMemberOf } = useRequest(
    async () => {
      return await adminApi.listGroupMemberOf({
        namespaceId,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceKind],
      ready: namespaceKind === NamespaceKind.NamespaceKindUser,
    }
  );

  const groupMemberOfColumns = useGroupMemberOfColumns();

  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{namespaceKind}</div>
      {namespaceKind !== NamespaceKind.NamespaceKindUser && (
        <Card
          title="Certificate Policies"
          extra={
            <Link to="./cert-policy/_create">Create certificate policy</Link>
          }
        >
          <CertPolicyRefTable routePrefix="./cert-policy/" />
        </Card>
      )}
      {namespaceKind === NamespaceKind.NamespaceKindUser && (
        <>
          <CertificatesTableCard />
          <Card title="Listed group memberships">
            <Table<ResourceReference>
              dataSource={groupMemberOf}
              columns={groupMemberOfColumns}
              rowKey="id"
            />
          </Card>
          <Card title="Sync group membership">
            <MemberOfGroupForm />
          </Card>
        </>
      )}
    </>
  );
}
