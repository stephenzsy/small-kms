import {
  AccountInfo,
  AuthenticationResult,
  InteractionStatus,
  PublicClientApplication,
} from "@azure/msal-browser";
import { MsalProvider, useAccount, useMsal } from "@azure/msal-react";
import { useLatest, useMemoizedFn } from "ahooks";
import {
  createContext,
  useContext,
  useEffect,
  type PropsWithChildren,
  useMemo,
} from "react";

const pca = new PublicClientApplication({
  auth: {
    clientId: import.meta.env.VITE_AZURE_CLIENT_ID,
    authority: `https://login.microsoftonline.com/${
      import.meta.env.VITE_AZURE_TENANT_ID
    }`,
    redirectUri: import.meta.env.VITE_MSAL_REDIRECT_URI,
  },
  cache: {
    cacheLocation: "sessionStorage",
  },
});

export interface IAppAuthContext {
  account?: AccountInfo;
  logout: () => void;
  acquireToken: () => Promise<AuthenticationResult | void>;
  readonly isAdmin: boolean;
}

export const AppAuthContext = createContext<IAppAuthContext>({
  logout: () => {},
  acquireToken: () => Promise.resolve(undefined),
  isAdmin: false,
});

function AuthContextProvider({ children }: PropsWithChildren<{}>) {
  const { instance, inProgress, accounts } = useMsal();
  const account = useAccount(accounts?.[0] ?? undefined);
  const logout = useMemoizedFn(() => {
    instance.logoutRedirect();
  });

  const accountRef = useLatest(account);
  const acquireToken = useMemoizedFn(
    async (): Promise<AuthenticationResult | void> => {
      try {
        if (accountRef.current) {
          return await instance.acquireTokenSilent({
            scopes: [import.meta.env.VITE_API_SCOPE],
            account: accountRef.current,
          });
        }
        instance.setActiveAccount(accountRef.current);
      } catch {}
      return instance.loginRedirect({
        scopes: [import.meta.env.VITE_API_SCOPE],
        redirectUri: import.meta.env.VITE_MSAL_REDIRECT_URI,
      });
    }
  );
  useEffect(() => {
    if (inProgress === InteractionStatus.None && accounts.length === 0) {
      acquireToken();
    }
  }, [inProgress, account]);
  const isAdmin = useMemo(
    () => !!account?.idTokenClaims?.roles?.includes("App.Admin"),
    [account]
  );

  return (
    <AppAuthContext.Provider
      value={{
        account: account ?? undefined,
        logout,
        acquireToken,
        isAdmin,
      }}
    >
      {children}
    </AppAuthContext.Provider>
  );
}

export function AuthProvider({ children }: PropsWithChildren<{}>) {
  return (
    <MsalProvider instance={pca}>
      <AuthContextProvider>{children}</AuthContextProvider>
    </MsalProvider>
  );
}

export function useAppAuthContext() {
  return useContext(AppAuthContext);
}
