import { useCreation } from "ahooks";
import { BaseAPI, CertsApi, Configuration } from "../generated";
import { useMsal } from "@azure/msal-react";
import { IPublicClientApplication } from "@azure/msal-browser";

function getDevAuthHeaders(instance: IPublicClientApplication) {
  const activeAccount = instance.getActiveAccount();
  if (!activeAccount) {
    return undefined;
  }

  return {
    "x-ms-client-principal-id": activeAccount.localAccountId,
    "x-ms-client-principal-name": activeAccount.username,
    "x-ms-client-principal": btoa(
      JSON.stringify({
        claims:
          activeAccount.idTokenClaims?.roles?.map((v) => ({
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
  const { instance } = useMsal();

  return useCreation(() => {
    const headers =
      import.meta.env.VITE_USE_DEV_AUTH_HEADERS === "true"
        ? getDevAuthHeaders(instance)
        : undefined;
    return new ClientType(
      new Configuration({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await instance.acquireTokenSilent({
            scopes: [import.meta.env.VITE_API_SCOPE],
          });
          return result.accessToken;
        },
        headers,
      })
    );
  }, [instance]);
}

export function useCertsApi() {
  return useAuthedClient(CertsApi);
}
