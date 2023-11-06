import { Button, Table, Tag, type TableColumnType } from "antd";
import { AdminApi, CertPolicyRef, NamespaceKind } from "../generated";
import { useMemo, useContext } from "react";
import { Link } from "../components/Link";
import {
  NamespaceContext,
  NamespaceConfigContext,
} from "./contexts/NamespaceContext";
import { useAuthedClient } from "../utils/useCertsApi";
import { useMemoizedFn, useRequest } from "ahooks";

export function usePolicyRefTableColumns(
  routePrefix: string,
  onRenderTags?: (r: CertPolicyRef) => React.ReactNode
) {
  return useMemo<TableColumnType<CertPolicyRef>[]>(
    () => [
      {
        title: "App ID",
        render: (r: CertPolicyRef) => (
          <>
            <span className="font-mono">{r.id}</span>
            {onRenderTags?.(r)}
          </>
        ),
      },
      {
        title: "Display name",
        render: (r: CertPolicyRef) => r.displayName,
      },

      {
        title: "Actions",
        render: (r: CertPolicyRef) => (
          <>
            <Link to={`${routePrefix}${r.id}`}>View</Link>
          </>
        ),
      },
    ],
    [routePrefix, onRenderTags]
  );
}

function useColumns(
  routePrefix: string,
  activeIssuerPolicyId: string | undefined
) {
  return usePolicyRefTableColumns(
    routePrefix,
    useMemoizedFn((r: CertPolicyRef) => {
      return (
        r.id === activeIssuerPolicyId && (
          <Tag className="ml-2" color="blue">
            Current issuer
          </Tag>
        )
      );
    })
  );
}

export function useCertPolicies() {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    () => {
      return adminApi.listCertPolicies({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier,
    }
  );

  return certPolicies;
}

export function useSecretPolicies() {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data } = useRequest(
    () => {
      return adminApi.listSecretPolicies({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier,
    }
  );

  return data;
}

export function CertPolicyRefTable({ routePrefix }: { routePrefix: string }) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: certPolicies } = useRequest(
    () => {
      return adminApi.listCertPolicies({
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!namespaceIdentifier,
    }
  );

  const { issuer: issuerRule } = useContext(NamespaceConfigContext);

  const columns = useColumns(routePrefix, issuerRule?.policyId);
  return (
    <Table<CertPolicyRef>
      columns={columns}
      dataSource={certPolicies}
      rowKey={(r) => r.id}
    />
  );
}
