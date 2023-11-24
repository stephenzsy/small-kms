import { Client } from "@microsoft/microsoft-graph-client";
import { useCreation } from "ahooks";
import { useContext } from "react";
import { AppAuthContext } from "../auth/AuthProvider";
import { BaseAPI, Configuration } from "../generated";
import {
  BaseAPI as BaseAPI2,
  Configuration as Configuration2,
} from "../generated/apiv2";

export function useAuthedClient<T extends BaseAPI>(ClientType: {
  new (configuration: Configuration): T;
}): T {
  const { account, acquireToken } = useContext(AppAuthContext);

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

export function useAuthedClientV2<T extends BaseAPI2>(ClientType: {
  new (configuration: Configuration2): T;
}): T {
  const { account, acquireToken } = useContext(AppAuthContext);

  return useCreation(() => {
    return new ClientType(
      new Configuration2({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await acquireToken();
          return result?.accessToken || "";
        },
      })
    );
  }, [account, acquireToken]);
}

export function useAdminApi() {
  const { api } = useContext(AppAuthContext);
  return api;
}

export function useGraphClient(): Client {
  const { account, acquireToken } = useContext(AppAuthContext);

  return useCreation(() => {
    return Client.initWithMiddleware({
      authProvider: {
        getAccessToken: async () => {
          const result = await acquireToken([
            "https://graph.microsoft.com/Directory.Read.All",
          ]);
          return result?.accessToken || "";
        },
      },
    });
  }, [account, acquireToken]);
}
