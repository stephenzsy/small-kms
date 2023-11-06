import { AdminApi, ImportProfileRequest, ResourceKind } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import { Button, Form, Input } from "antd";
import { useForm } from "antd/es/form/Form";

type ImportProfileFormState = {
  objectId?: string;
};
export function ImportProfileForm({
  onCreated, profileKind,
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
