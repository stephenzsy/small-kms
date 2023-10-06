import { useRequest } from "ahooks";
import React from "react";
import { useParams } from "react-router-dom";
import { AdminApi, Profile, ProfileType } from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";

export const NamespaceContext = React.createContext<{
  nsInfo: Profile | undefined;
}>({ nsInfo: undefined });

export function NamespaceContextProvider(props: React.PropsWithChildren<{}>) {
  const { namespaceId, profileType } = useParams() as {
    namespaceId: string;
    profileType: ProfileType;
  };
  const adminApi = useAuthedClient(AdminApi);
  const { data: namespaceInfo } = useRequest(
    async () => {
      return await adminApi.getProfile({
        profileType: profileType,
        profileId: namespaceId,
      });
    },
    {
      refreshDeps: [namespaceId, profileType],
    }
  );

  return (
    <NamespaceContext.Provider value={{ nsInfo: namespaceInfo }}>
      {props.children}
    </NamespaceContext.Provider>
  );
}
