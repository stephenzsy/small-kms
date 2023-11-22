import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Divider, Form, Select, Typography } from "antd";
import { useEffect, useMemo } from "react";
import { useParams } from "react-router-dom";
import { ResourceRefsTable } from "../../admin/tables/ResourceRefsTable";
import { Link } from "../../components/Link";
import {
  AdminApi,
  AgentConfigIdentity,
  AgentConfigIdentityFields,
  AgentConfigName,
  CreateAgentConfigRequest,
  NamespaceProvider,
  Ref,
  Ref as ResourceReference,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

function useAgent(agentId: string | undefined) {
  const api = useAuthedClientV2(AdminApi);
  return useRequest(
    async () => {
      if (agentId) {
        return api.getAgent({
          id: agentId,
        });
      }
    },
    {
      refreshDeps: [agentId],
    }
  );
}

export default function AgentPage() {
  const { id } = useParams<{ id: string }>();
  const { data: agent } = useAgent(id);
  const spId = agent?.servicePrincipalId;
  const api = useAuthedClientV2(AdminApi);
  const { data: certPolicies, loading: certPoliciesLoading } = useRequest(
    async () => {
      if (spId) {
        return api.listCertificatePolicies({
          namespaceProvider:
            NamespaceProvider.NamespaceProviderServicePrincipal,
          namespaceId: spId,
        });
      }
    },
    {
      refreshDeps: [spId],
    }
  );
  const renderActions = useMemoizedFn((item: ResourceReference) => {
    return (
      <div className="flex flex-row gap-2">
        <Link to={`/service-principal/${spId}/cert-policies/${item.id}`}>
          View
        </Link>
      </div>
    );
  });

  return (
    <>
      <Typography.Title>Agent</Typography.Title>
      <Card title="Agent info">
        <dl className="dl">
          <div className="dl-item">
            <dt>ID (Client ID)</dt>
            <dd className="font-mono">{agent?.id}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Display name</dt>
            <dd className="dd">{agent?.displayName}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Application ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.applicationId}</dd>
          </div>
          <div className="dl-item">
            <dt className="dt">Service Principal ID (Object ID)</dt>
            <dd className="dd font-mono">{agent?.servicePrincipalId}</dd>
          </div>
        </dl>
      </Card>
      <Divider />
      <Typography.Title level={2}>Service Principal</Typography.Title>
      <Card title="Certificate Policies">
        <ResourceRefsTable
          dataSource={certPolicies}
          loading={certPoliciesLoading}
          renderActions={renderActions}
        />
      </Card>
      {agent?.servicePrincipalId && (
        <AgentConfigurationsCard
          namespaceId={agent.servicePrincipalId}
          certificatePolicies={certPolicies}
        />
      )}
    </>
  );
}

function AgentIdentityForm({
  certificatePolicies,
  value,
  onPutConfig,
}: {
  value: AgentConfigIdentity | undefined;
  certificatePolicies: Ref[] | undefined;
  onPutConfig: (config: AgentConfigIdentityFields) => void;
}) {
  const [formInstance] = Form.useForm<AgentConfigIdentityFields>();

  const selectOpts = useMemo(() => {
    return certificatePolicies?.map((cp) => {
      return {
        label: (
          <span>
            {cp.displayName} ({cp.id})
          </span>
        ),
        value: cp.id,
      };
    });
  }, [certificatePolicies]);

  useEffect(() => {
    if (value) {
      formInstance.setFieldsValue(value);
    }
  }, [value]);

  return (
    <Form form={formInstance} layout="vertical" onFinish={onPutConfig}>
      <Form.Item<AgentConfigIdentityFields>
        name="keyCredentialCertificatePolicyId"
        label="Select Key Credential Certificate Policy"
      >
        <Select options={selectOpts} />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

function AgentConfigurationsCard({
  namespaceId,
  certificatePolicies,
}: {
  namespaceId: string;
  certificatePolicies: Ref[] | undefined;
}) {
  const api = useAuthedClientV2(AdminApi);
  const { run: updateIdentityConfig, data: identityConfig } = useRequest(
    async (req?: CreateAgentConfigRequest) => {
      if (req) {
        return api.putAgentConfig({
          configName: AgentConfigName.AgentConfigNameIdentity,
          namespaceId: namespaceId,
          createAgentConfigRequest: req,
        });
      }
      return await api.getAgentConfig({
        configName: AgentConfigName.AgentConfigNameIdentity,
        namespaceId: namespaceId,
      });
    },
    {
      refreshDeps: [namespaceId],
    }
  );
  return (
    <Card title="Agent configurations">
      <AgentIdentityForm
        certificatePolicies={certificatePolicies}
        onPutConfig={updateIdentityConfig}
        value={identityConfig}
      />
    </Card>
  );
}
