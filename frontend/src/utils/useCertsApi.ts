import { useCreation } from "ahooks";
import { useAppAuthContext } from "../auth/AuthProvider";
import { BaseAPI, Configuration } from "../generated";

export function useAuthedClient<T extends BaseAPI>(ClientType: {
  new (configuration: Configuration): T;
}): T {
  const { account, acquireToken } = useAppAuthContext();

  return useCreation(() => {
    return new ClientType(
      new Configuration({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await acquireToken();
          return result?.accessToken || "";
        },
      })
    );
  }, [account, acquireToken]);
}
