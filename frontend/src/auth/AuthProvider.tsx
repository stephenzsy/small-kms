import {
  AccountInfo,
  AuthenticationResult,
  InteractionStatus,
  InteractionType,
  PublicClientApplication,
} from "@azure/msal-browser";
import {
  MsalProvider,
  useAccount,
  useMsal,
  useMsalAuthentication,
} from "@azure/msal-react";
import { useLatest, useMemoizedFn, useRequest } from "ahooks";
import {
  createContext,
  useContext,
  type PropsWithChildren,
  useEffect,
  useState,
} from "react";

const pca = new PublicClientApplication({
  auth: {
    clientId: import.meta.env.VITE_AZURE_CLIENT_ID,
    authority: `https://login.microsoftonline.com/${
      import.meta.env.VITE_AZURE_TENANT_ID
    }`,
  },
  cache: {
    cacheLocation: "sessionStorage",
  },
});

export interface IAppAuthContext {
  account?: AccountInfo;
  logout: () => void;
  acquireToken: () => Promise<AuthenticationResult | void>;
}

export const AppAuthContext = createContext<IAppAuthContext>({
  logout: () => {},
  acquireToken: () => Promise.resolve(undefined),
});

const redirectUri = window.location.protocol + "//" + window.location.host;

function AuthContextProvider({ children }: PropsWithChildren<{}>) {
  const { instance, inProgress, accounts } = useMsal();
  const account = useAccount(accounts[0] || {}) ?? undefined;

  const logout = useMemoizedFn(() => {
    instance.logoutRedirect();
  });
  const accountRef = useLatest(account);
  const acquireToken = useMemoizedFn(
    (): Promise<AuthenticationResult | void> => {
      return accountRef.current
        ? instance.acquireTokenSilent({
            scopes: [import.meta.env.VITE_API_SCOPE],
            account: accountRef.current,
          })
        : instance.loginRedirect({
            scopes: [import.meta.env.VITE_API_SCOPE],
            redirectUri,
          });
    }
  );
  useEffect(() => {
    if (inProgress !== InteractionStatus.None) {
      return;
    }
    if (!accountRef.current) {
      const t = setTimeout(() => {
        if (!accountRef.current) {
          acquireToken();
        }
      }, 5000);
      return () => clearTimeout(t);
    }
  }, [account, inProgress]);
  return (
    inProgress !== InteractionStatus.Startup && (
      <AppAuthContext.Provider
        value={{
          account,
          logout,
          acquireToken,
        }}
      >
        {children}
      </AppAuthContext.Provider>
    )
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
