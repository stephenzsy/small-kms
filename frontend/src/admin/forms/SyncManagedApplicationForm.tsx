import {
  AdminApi,
  ImportProfileRequest,
  ResourceKind,
  SyncManagedAppRequest,
} from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";
import { useRequest } from "ahooks";
import { Button, Form, Input } from "antd";
import { useForm } from "antd/es/form/Form";

type SyncAppFormState = {
  clientId?: string;
};
export function SyncManagedApplicationForm({
  onSynced,
}: {
  onSynced: () => void;
}) {
  const [form] = useForm<SyncAppFormState>();

  const adminApi = useAuthedClient(AdminApi);

  const { run } = useRequest(
    async (req: SyncManagedAppRequest) => {
      await adminApi.syncManagedApp(req);
      onSynced();
      form.resetFields();
    },
    { manual: true }
  );

  return (
    <Form
      form={form}
      onFinish={(values) => {
        const objectId = values.clientId?.trim();
        if (objectId) {
          return run({
            managedAppId: objectId,
          });
        }
      }}
    >
      <Form.Item<SyncAppFormState>
        name="clientId"
        label="Microsoft Entra application (client) ID"
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
