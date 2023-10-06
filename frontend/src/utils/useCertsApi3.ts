import { useCreation } from "ahooks";
import { useAppAuthContext } from "../auth/AuthProvider";
import {
  BaseAPI as BaseAPI3,
  Configuration as Configuration3,
} from "../generated3";

export function useAuthedClient<T extends BaseAPI3>(ClientType: {
  new (configuration: Configuration3): T;
}): T {
  const { account, acquireToken } = useAppAuthContext();

  return useCreation(() => {
    return new ClientType(
      new Configuration3({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await acquireToken();
          return result?.accessToken || "";
        },
      })
    );
  }, [account, acquireToken]);
}
