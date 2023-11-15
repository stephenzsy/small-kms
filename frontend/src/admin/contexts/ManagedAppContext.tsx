import { useMemoizedFn, useRequest } from "ahooks";
import { PropsWithChildren, createContext } from "react";
import { useParams } from "react-router-dom";
import { AdminApi, type ManagedAppRef } from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";

export type ManagedAppContextValue = {
  managedApp?: ManagedAppRef;
  syncApp: () => void;
};
export const ManagedAppContext = createContext<ManagedAppContextValue>({
  syncApp: () => {},
});

export function ManagedAppContextProvider({ children }: PropsWithChildren<{}>) {
  const { appId } = useParams() as { appId: string };
  const adminApi = useAuthedClient(AdminApi);
  const { data: managedApp, run } = useRequest(
    (sync?: boolean) => {
      if (sync) {
        return adminApi.syncManagedApp({ managedAppId: appId });
      }
      return adminApi.getManagedApp({ managedAppId: appId });
    },
    {
      refreshDeps: [appId],
    }
  );

  const syncApp = useMemoizedFn(() => {
    run(true);
  });

  return (
    <ManagedAppContext.Provider
      value={{
        managedApp,
        syncApp,
      }}
    >
      {children}
    </ManagedAppContext.Provider>
  );
}
