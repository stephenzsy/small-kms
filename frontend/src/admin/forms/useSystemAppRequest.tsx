import { useRequest } from "ahooks";
import { AdminApi } from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";

export function useSystemAppRequest(systemAppName: string) {
  const adminApi = useAuthedClientV2(AdminApi);
  return useRequest(
    (isSync?: boolean) => {
      if (isSync) {
        return adminApi.syncSystemApp({
          id: systemAppName,
        });
      }
      return adminApi.getSystemApp({
        id: systemAppName,
      });
    },
    {
      refreshDeps: [systemAppName],
    }
  );
}
