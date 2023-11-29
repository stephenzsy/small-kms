import { useRequest } from "ahooks";
import { useAdminApi } from "../../utils/useCertsApi";

export function useSystemAppRequest(systemAppName: string) {
  const adminApi = useAdminApi();
  return useRequest(
    async (isSync?: boolean) => {
      if (isSync) {
        return adminApi?.syncSystemApp({
          id: systemAppName,
        });
      }
      return adminApi?.getSystemApp({
        id: systemAppName,
      });
    },
    {
      refreshDeps: [systemAppName],
    }
  );
}
