import { EyeIcon } from "@heroicons/react/24/outline";
import { useBoolean, useRequest } from "ahooks";
import { Button, Card, Input, Typography } from "antd";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";

export default function SecretPage() {
  const { namespaceKind, namespaceId: namespaceIdentifier } =
    useContext(NamespaceContext);

  const { id } = useParams() as { id: string };

  const adminApi = useAuthedClient(AdminApi);
  const {
    data: secret,
    run: getWithValue,
    loading,
  } = useRequest((withValue?: boolean) => {
    return adminApi.getSecret({
      resourceId: id,
      namespaceId: namespaceIdentifier,
      namespaceKind: namespaceKind,
      withValue: withValue,
    });
  }, {});

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
      <Typography.Title>Certificate</Typography.Title>
      <Card title="Certificate">
        <dl>
          <div>
            <dt className="font-medium">ID</dt>
            <dd className="font-mono">{secret?.id}</dd>
          </div>
          <div>
            <dt className="font-medium">Secret Key Vault ID</dt>
            <dd className="font-mono">{secret?.sid}</dd>
          </div>
        </dl>
      </Card>
      {
        <Card title="Actions">
          <div className="flex flex-row gap-4">
            <Button
              onClick={() => {
                getWithValue(true);
              }}
              type="primary"
              loading={loading}
            >
              Get secret value
            </Button>
            <Input
              readOnly
              value={secret?.value}
              type={reviewSecret ? "text" : "password"}
            />
            <Button
              type="text"
              onClick={toggle}
              icon={<EyeIcon className="h-em w-em" />}
            >
              Review secret
            </Button>
          </div>
        </Card>
      }
    </>
  );
}
