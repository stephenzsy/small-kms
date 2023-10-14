import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { AdminApi, NamespaceKind } from "../generated";
import {
  useAuthedClient,
  useAuthedClient as useAuthedClientOld,
} from "../utils/useCertsApi";
import { DeviceGroupInstall } from "./DeviceGroupInstall";

export function DeviceServicePrincipalLink(props: { namespaceId: string }) {
  const { data, run } = useRequest((apply?: boolean) => {
    return Promise.resolve({}); // TODO: fix this
  }, {});
  return (
    <>
      <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
        <Button
          variant="primary"
          onClick={() => {
            run(true);
          }}
        >
          Link service principal
        </Button>
      </section>
      <DeviceGroupInstall namespaceId={props.namespaceId} />
    </>
  );
}

export function ApplicationServicePrincipalLink(props: {
  namespaceId: string;
}) {
  const adminApi = useAuthedClient(AdminApi);

  const { data, run } = useRequest(
    async (apply?: boolean) => {
      if (apply) {
        return await adminApi.createManagedNamespace({
          namespaceKind: NamespaceKind.NamespaceKindApplication,
          namespaceId: props.namespaceId,
          targetNamespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
        });
      }
    },
    {
      manual: true,
    }
  );
  return (
    <>
      <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
        <Button
          variant="primary"
          onClick={() => {
            run(true);
          }}
        >
          Create or link service principal
        </Button>
      </section>
    </>
  );
}
