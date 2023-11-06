import { Button, Card, Form, Input } from "antd";
import { useContext } from "react";
import { Link } from "../components/Link";
import { CertPolicyRefTable } from "./CertPolicyRefTable";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { AdminApi, NamespaceKind } from "../generated";
import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi";

type MemberOfGroupFormState = {
  groupId: string;
};

function MemberOfGroupForm() {
  const [form] = Form.useForm<MemberOfGroupFormState>();
  const api = useAuthedClient(AdminApi);
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
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

export default function NamespacePage() {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  //  const adminApi = useAuthedClient(AdminApi);

  return (
    <>
      <h1>{namespaceIdentifier}</h1>
      <div>{namespaceKind}</div>
      <Card
        title="Certificate Policies"
        extra={
          <Link to="./cert-policy/_create">Create certificate policy</Link>
        }
      >
        <CertPolicyRefTable routePrefix="./cert-policy/" />
      </Card>
      {namespaceKind === NamespaceKind.NamespaceKindUser && (
        <Card title="Sync group membership">
          <MemberOfGroupForm />
        </Card>
      )}
    </>
  );
}
