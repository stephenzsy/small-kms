import { useRequest } from "ahooks";
import { CertsApi } from "../generated";
/*
export function useCaList(certsApi: CertsApi, namespaceId: string | undefined) {
  return useRequest(
    async () => {
      if (!namespaceId) {
        return undefined;
      }
      return certsApi.listCertificatesV1({
        namespaceId,
      });
    },
    { refreshDeps: [namespaceId] }
  );
}
*/