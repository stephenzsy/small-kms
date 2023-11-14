"use client";
import { PropsWithChildren, createContext, useContext, useMemo } from "react";
import { AdminApi, Configuration } from "./generated";
import { useMsal } from "@azure/msal-react";

const AdminApiContext = createContext<AdminApi | undefined>(undefined);

export function AdminApiProvider(props: PropsWithChildren<{}>) {
  const { instance } = useMsal();
  const client = useMemo(() => {
    return new AdminApi(
      new Configuration({
        basePath: process.env.NEXT_PUBLIC_API_BASE_PATH,
        accessToken: async () => {
          const result = await instance.acquireTokenSilent({
            scopes: [process.env.NEXT_PUBLIC_API_SCOPE!],
          });
          return result?.accessToken || "";
        },
      })
    );
  }, [instance]);
  console.log(client);
  return (
    <AdminApiContext.Provider value={client}>
      {props.children}
    </AdminApiContext.Provider>
  );
}

export function useAdminApi(): AdminApi {
  return useContext(AdminApiContext) as AdminApi;
}
