import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { DeviceGroupInstall } from "./DeviceGroupInstall";

export function DeviceServicePrincipalLink(props: { namespaceId: string }) {
  const adminApi = useAuthedClient(AdminApi);

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
