import { ArrowPathIcon } from "@heroicons/react/24/outline";
import type { Group } from "@microsoft/microsoft-graph-types";
import { useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Form, Input, Select, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { DefaultOptionType } from "antd/es/select";
import classNames from "classnames";
import { useMemo } from "react";
import { Link } from "../components/Link";
import {
  AdminApi,
  NamespaceProvider,
  Ref,
  SyncMemberOfRequest,
} from "../generated/apiv2";
import { useAuthedClientV2, useGraphClient } from "../utils/useCertsApi";
import { useNamespace } from "./contexts/useNamespace";
import { NamespacePoliciesTableCard } from "./tables/NamespacePoliciesTableCard";
import { ResourceRefsTable } from "./tables/ResourceRefsTable";

function GroupMembershipSyncForm({ userId }: { userId: string }) {
  const [form] = useForm<{ groupId: string }>();
  const gclient = useGraphClient();
  const {
    data: dirObjects,
    run: getDirectoryObjects,
    loading: dirObjLoading,
  } = useRequest(
    async () => {
      return (
        await gclient
          .api(`/users/${userId}/memberOf`)
          .select(["id", "displayName"])
          .get()
      ).value as Group[];
    },
    { manual: true }
  );

  const api = useAuthedClientV2(AdminApi);
  const { run: syncMemberOf } = useRequest(
    async (req: SyncMemberOfRequest) => {
      return await api.syncMemberOf(req);
    },
    { manual: true }
  );

  const dirOpjOptions = useMemo<DefaultOptionType[] | undefined>(() => {
    return dirObjects?.map((obj: Group) => ({
      label: obj.displayName,
      value: obj.id,
    }));
  }, [dirObjects]);
  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(values) => {
        if (values.groupId) {
          return syncMemberOf({
            namespaceId: userId,
            id: values.groupId,
            namespaceProvider: NamespaceProvider.NamespaceProviderUser,
          });
        }
      }}
    >
      <Form.Item
        name="groupId"
        label={
          <div className="inline-flex items-center gap-4">
            <span>Select Graph Object</span>
            <Button
              type="link"
              size="small"
              onClick={getDirectoryObjects}
              className="inline-flex items-center gap-2"
            >
              <ArrowPathIcon
                className={classNames(
                  "h-em w-em",
                  dirObjLoading && "animate-spin"
                )}
              />
              <span>Get List from Microsoft Graph API</span>
            </Button>
          </div>
        }
      >
        <Select options={dirOpjOptions} loading={dirObjLoading} />
      </Form.Item>
      <Form.Item label="Enter Group ID" name="groupId">
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Sync
        </Button>
      </Form.Item>
    </Form>
  );
}

export default function NamespacePage() {
  const { namespaceId, namespaceProvider } = useNamespace();
  const showProfile = [
    NamespaceProvider.NamespaceProviderServicePrincipal,
    NamespaceProvider.NamespaceProviderGroup,
    NamespaceProvider.NamespaceProviderUser,
  ].some((np) => np === namespaceProvider);
  const showCertPolicies = [
    NamespaceProvider.NamespaceProviderRootCA,
    NamespaceProvider.NamespaceProviderIntermediateCA,
    NamespaceProvider.NamespaceProviderServicePrincipal,
    NamespaceProvider.NamespaceProviderGroup,
  ].some((np) => np === namespaceProvider);

  const adminApi = useAuthedClientV2(AdminApi);

  const { data: profile } = useRequest(
    () => {
      return adminApi.getProfile({
        namespaceId: namespaceId,
        namespaceProvider: namespaceProvider,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceProvider],
      ready: showProfile,
    }
  );

  const { data: userCerts, loading: userCertsLoading } = useRequest(
    () => {
      return adminApi.listCertificates({
        namespaceId: namespaceId,
        namespaceProvider: namespaceProvider,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceProvider],
      ready: namespaceProvider === NamespaceProvider.NamespaceProviderUser,
    }
  );

  const renderUserCertActions = useMemoizedFn((ref: Ref) => {
    return (
      <div className="flex flex-row gap-2">
        <Link to={`./certificates/${ref.id}`}>View</Link>
      </div>
    );
  });
  return (
    <>
      <Typography.Title>
        Namespace: {profile?.displayName ?? namespaceId}
      </Typography.Title>
      {showProfile && (
        <Card title="Profile">
          <dl className="dl">
            <div>
              <dt>ID</dt>
              <dd className="font-mono">{profile?.id}</dd>
            </div>
            <div>
              <dt>Display Name</dt>
              <dd className="font-mono">{profile?.displayName}</dd>
            </div>
            {namespaceProvider === NamespaceProvider.NamespaceProviderUser && (
              <div>
                <dt>User Principal Name</dt>
                <dd className="font-mono">{profile?.userPrincipalName}</dd>
              </div>
            )}
          </dl>
        </Card>
      )}
      {showCertPolicies && <NamespacePoliciesTableCard type="cert" />}
      {namespaceProvider ===
        NamespaceProvider.NamespaceProviderServicePrincipal && (
        <NamespacePoliciesTableCard type="key" />
      )}
      {namespaceProvider === NamespaceProvider.NamespaceProviderUser && (
        <Card title="Certificates">
          <ResourceRefsTable
            renderActions={renderUserCertActions}
            loading={userCertsLoading}
            dataSource={userCerts}
          />
        </Card>
      )}
      {namespaceProvider === NamespaceProvider.NamespaceProviderUser && (
        <Card title="Sync group membership">
          <GroupMembershipSyncForm userId={namespaceId} />
        </Card>
      )}
      {/* {namespaceKind === NamespaceKind.NamespaceKindUser && (
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
      )} */}
    </>
  );
}
