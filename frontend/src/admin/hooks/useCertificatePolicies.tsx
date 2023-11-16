import { useRequest } from "ahooks";
import { AdminApi, NamespaceProvider } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

export function useCertificatePolicies(
  namespaceProvider: NamespaceProvider,
  namespaceId?: string
) {
  const api = useAuthedClientV2(AdminApi);

  return useRequest(
    async () => {
      if (namespaceId) {
        return api.listCertificatePolicies({ namespaceProvider, namespaceId });
      }
    },
    {
      refreshDeps: [api, namespaceProvider, namespaceId],
    }
  );
}
