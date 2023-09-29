import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

export function DeviceServicePrincipalLink(props: { namespaceId: string }) {
  const adminApi = useAuthedClient(AdminApi);

  const { data, run } = useRequest(
    (apply?: boolean) => {
      return adminApi.getDeviceServicePrincipalLinkV2({
        namespaceId: props.namespaceId,
        apply,
      });
    },
    {
      manual: true,
    }
  );
  return (
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
  );
}
