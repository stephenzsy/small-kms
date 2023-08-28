import { AccountInfo } from "@azure/msal-browser";
import { useCreation } from "ahooks";
import { useAppAuthContext } from "../auth/AuthProvider";
import { BaseAPI, CertsApi, Configuration } from "../generated";

function getDevAuthHeaders(account: AccountInfo | undefined) {
  if (!account) return undefined;
  return {
    "x-ms-client-principal-id": account.localAccountId,
    "x-ms-client-principal-name": account.username,
    "x-ms-client-principal": btoa(
      JSON.stringify({
        claims:
          account.idTokenClaims?.roles?.map((v) => ({
            typ: "roles",
            val: v,
          })) ?? [],
      })
    ),
  };
}

export function useAuthedClient<T extends BaseAPI>(ClientType: {
  new (configuration: Configuration): T;
}): T {
  const { account, acquireToken } = useAppAuthContext();

  return useCreation(() => {
    const headers =
      import.meta.env.VITE_USE_DEV_AUTH_HEADERS === "true"
        ? getDevAuthHeaders(account)
        : undefined;
    return new ClientType(
      new Configuration({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await acquireToken();
          return result?.accessToken || "";
        },
        headers,
      })
    );
  }, [account, acquireToken]);
}

export function useCertsApi() {
  return useAuthedClient(CertsApi);
}
