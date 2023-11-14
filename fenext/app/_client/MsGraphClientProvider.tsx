"use client";

import {
  type PropsWithChildren,
  createContext,
  useMemo,
  useContext,
} from "react";
import { Client } from "@microsoft/microsoft-graph-client";
import { useMsal } from "@azure/msal-react";

const GraphClientContext = createContext<Client | undefined>(undefined);

export function GraphClientProvider(props: PropsWithChildren<{}>) {
  const { instance } = useMsal();
  const client = useMemo((): Client => {
    return Client.initWithMiddleware({
      authProvider: {
        getAccessToken: async (opts) => {
          const authResult = await instance.acquireTokenSilent({
            scopes: opts?.scopes || ["User.Read"],
          });
          return authResult.accessToken;
        },
      },
    });
  }, []);
  return (
    <GraphClientContext.Provider value={client}>
      {props.children}
    </GraphClientContext.Provider>
  );
}

export function useGraphClient(): Client {
  return useContext(GraphClientContext) as Client;
}
