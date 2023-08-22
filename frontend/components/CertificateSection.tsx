"use client";

import { AdminApi, Configuration } from "@/generated";
import { IPublicClientApplication } from "@azure/msal-browser";
import { useMsal } from "@azure/msal-react";
import { useCreation, useRequest } from "ahooks";

function useAdminClient(msalInstance: IPublicClientApplication) {
  return useCreation(() => {
    return new AdminApi(
      new Configuration({
        basePath: process.env.NEXT_PUBLIC_BACKEND_URL_BASE!,
        accessToken: () =>
          msalInstance
            .acquireTokenSilent({ scopes: [] })
            .then((result) => result.idToken),
      })
    );
  }, [msalInstance]);
}

export function CertificateSection() {
  const { instance, accounts } = useMsal();
  const idToken = accounts[0]?.idToken;
  const client = useAdminClient(instance);

  const { data } = useRequest(
    async () => {
      const result = await client?.adminGetCAMetadata({ id: "root" });
      return result;
    },
    {
      ready: !!client,
      refreshDeps: [client],
    }
  );

  return <div>{data && JSON.stringify(data, undefined, 2)}</div>;
}
