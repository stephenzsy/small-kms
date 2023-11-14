import { useRequest } from "ahooks";
import { AdminApi, SystemAppName } from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";

export function useSystemAppRequest(systemAppName: SystemAppName) {
  const adminApi = useAuthedClient(AdminApi);
  return useRequest(
    (isSync?: boolean) => {
      if (isSync) {
        return adminApi.syncSystemApp({
          systemAppName,
        });
      }
      return adminApi.getSystemApp({
        systemAppName,
      });
    },
    {
      refreshDeps: [systemAppName],
    }
  );
}
