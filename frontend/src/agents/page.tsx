import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { AdminApi, CreateAgentRequest, Ref } from "../generated/apiv2";

import { useAuthedClientV2 } from "../utils/useCertsApi";
import { ResourceRefsTable } from "../admin/tables/ResourceRefsTable";
import { useMemo } from "react";
import { Link } from "../components/Link";

function useListAgents() {
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    () => {
      return api.listAgents();
    },
    {
      refreshDeps: [],
    }
  );
}

function useCreateAgents(onSuccess?: () => void) {
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    (req: CreateAgentRequest) => {
      return api.createAgent({
        createAgentRequest: req,
      });
    },
    {
      manual: true,
      onSuccess: onSuccess,
    }
  );
}

function CreateAgentForm({
  isCreate,
  onSuccess,
}: {
  isCreate: boolean;
  onSuccess?: () => void;
}) {
  const [form] = useForm<CreateAgentRequest>();
  const { run, loading } = useCreateAgents(onSuccess);
  return (
    <Form
      form={form}
      onFinish={(values) => {
        if (values.appId || values.displayName) {
          run(values);
        }
      }}
      layout="vertical"
    >
      {!isCreate && (
        <Form.Item<CreateAgentRequest>
          label="Application ID (Client ID)"
          name="appId"
        >
          <Input />
        </Form.Item>
      )}
      {isCreate && (
        <Form.Item<CreateAgentRequest> label="Display Name" name="displayName">
          <Input />
        </Form.Item>
      )}
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading}>
          {isCreate ? "Create" : "Import"}
        </Button>
      </Form.Item>
    </Form>
  );
}


export default function AgentsPage() {
  const { data: agents, run: refreshAgents } = useListAgents();
  const renderActions = useMemoizedFn((item: Ref) => {
    return (
      <div className="flex flex-row gap-2">
        <Link to={`/agents/${item.id}`}>View</Link>
      </div>
    );
  });
  return (
    <>
      <Typography.Title>Agents</Typography.Title>
      <Card title="List of agents">
        <ResourceRefsTable resourceRefs={agents} renderActions={renderActions} />
      </Card>
      <Card title="Import agent application">
        <CreateAgentForm onSuccess={refreshAgents} isCreate={false} />
      </Card>
      <Card title="Create new agent application">
        <CreateAgentForm onSuccess={refreshAgents} isCreate={true} />
      </Card>
    </>
  );
}
