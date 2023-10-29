import { PropsWithChildren, createContext } from "react";
import { useParams } from "react-router-dom";
import { useAuthedClient } from "../../utils/useCertsApi";
import { AdminApi, type ManagedAppRef, SystemAppName } from "../../generated";
import { useMemoizedFn, useRequest } from "ahooks";

export type ManagedAppContextValue = {
  managedApp?: ManagedAppRef;
  syncApp: () => void;
};
export const ManagedAppContext = createContext<ManagedAppContextValue>({
  syncApp: () => {},
});

export function ManagedAppContextProvider({
  children,
  isSystemApp = false,
}: PropsWithChildren<{
  isSystemApp?: boolean;
}>) {
  const { appId } = useParams() as { appId: string };
  const adminApi = useAuthedClient(AdminApi);
  const { data: managedApp, run } = useRequest(
    (sync?: boolean) => {
      if (isSystemApp) {
        if (sync) {
          return adminApi.syncSystemApp({
            systemAppName: appId as SystemAppName,
          });
        }
        return adminApi.getSystemApp({ systemAppName: appId as SystemAppName });
      }
      if (sync) {
        return adminApi.syncManagedApp({ managedAppId: appId });
      }
      return adminApi.getManagedApp({ managedAppId: appId });
    },
    {
      refreshDeps: [appId, isSystemApp],
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
