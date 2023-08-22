import type { IPublicClientApplication } from "@azure/msal-browser";
import { useMsal } from "@azure/msal-react";
import { useCreation, useRequest } from "ahooks";
import { AdminApi, Configuration } from "../generated";

function useAdminClient(msalInstance: IPublicClientApplication) {
  return useCreation(() => {
    return new AdminApi(
      new Configuration({
        basePath: "http://localhost:8080",
        accessToken: () =>
          msalInstance
            .acquireTokenSilent({
              scopes: ["api://253cbd3f-312e-4657-af79-f954bb6877e8/api.admin"],
            })
            .then((result) => result.accessToken),
      })
    );
  }, [msalInstance]);
}

export default function AdminPage() {
  const { instance } = useMsal();
  const client = useAdminClient(instance);
  const { data } = useRequest(
    async () => {
      return await client.adminGetCAMetadata({ id: "root" });
    },
    {
      refreshDeps: [client],
    }
  );

  return <div>{JSON.stringify(data, undefined, 2)}</div>;
}
