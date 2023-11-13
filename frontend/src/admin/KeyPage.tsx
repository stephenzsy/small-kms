import { useBoolean, useRequest } from "ahooks";
import { Button, Card, Input, Typography } from "antd";
import { useContext, useState } from "react";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import { EyeIcon } from "@heroicons/react/24/outline";

export default function KeyPage() {
  const { namespaceKind, namespaceId: namespaceIdentifier } =
    useContext(NamespaceContext);

  const { id } = useParams() as { id: string };

  const adminApi = useAuthedClient(AdminApi);
  const { data: key, loading } = useRequest(
    () => {
      return adminApi.getKey({
        resourceId: id,
        namespaceId: namespaceIdentifier,
        namespaceKind: namespaceKind,
      });
    },
    {
      refreshDeps: [id, namespaceIdentifier, namespaceKind],
      ready: !!id && !!namespaceIdentifier && !!namespaceKind,
    }
  );

  // const {
  //   data: deleted,
  //   loading: deleteLoading,
  //   run: deleteCert,
  // } = useRequest(
  //   async () => {
  //     await adminApi.deleteCertificate({
  //       resourceId: certId,
  //       namespaceId: namespaceIdentifier,
  //       namespaceKind,
  //     });
  //     return true;
  //   },
  //   { manual: true }
  // );
  const [reviewSecret, { toggle }] = useBoolean();
  return (
    <>
      <Typography.Title>Key</Typography.Title>
      <Card title="Key" loading={loading}>
        <dl>
          <div>
            <dt className="font-medium">ID</dt>
            <dd className="font-mono">{key?.id}</dd>
          </div>
          <div>
            <dt className="font-medium">Key Vault ID</dt>
            <dd className="font-mono">{key?.kid}</dd>
          </div>
        </dl>
      </Card>
      {
        <Card title="Actions">
          <div className="flex flex-row gap-4"></div>
        </Card>
      }
    </>
  );
}
