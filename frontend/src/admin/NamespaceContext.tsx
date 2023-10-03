import { useRequest } from "ahooks";
import React from "react";
import { useParams } from "react-router-dom";
import { AdminApi, NamespaceInfo } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

export const NamespaceContext = React.createContext<{
  nsInfo: NamespaceInfo | undefined;
}>({ nsInfo: undefined });

export function NamespaceContextProvider(props: React.PropsWithChildren<{}>) {
  const { namespaceId } = useParams() as {
    namespaceId: string;
  };
  const adminApi = useAuthedClient(AdminApi);
  const { data: namespaceInfo } = useRequest(
    async () => {
      return await adminApi.getNamespaceInfoV2({ namespaceId });
    },
    {
      refreshDeps: [namespaceId],
    }
  );

  return (
    <NamespaceContext.Provider value={{ nsInfo: namespaceInfo }}>
      {props.children}
    </NamespaceContext.Provider>
  );
}
