import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { AdminApi as AdminApiOld } from "../generated";
import { AdminApi, NamespaceKind } from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import { useAuthedClient as useAuthedClientOld } from "../utils/useCertsApi";
import { DeviceGroupInstall } from "./DeviceGroupInstall";

export function DeviceServicePrincipalLink(props: { namespaceId: string }) {
  const adminApi = useAuthedClientOld(AdminApiOld);

  const { data, run } = useRequest((apply?: boolean) => {
    if (apply) {
      return adminApi.createDeviceServicePrincipalLinkV2({
        namespaceId: props.namespaceId,
      });
    }
    return adminApi.getDeviceServicePrincipalLinkV2({
      namespaceId: props.namespaceId,
    });
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
      <DeviceGroupInstall namespaceId={props.namespaceId} linkInfo={data} />
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
