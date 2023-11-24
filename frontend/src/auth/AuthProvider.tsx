import {
  AccountInfo,
  AuthenticationResult,
  PublicClientApplication,
} from "@azure/msal-browser";
import {
  MsalProvider,
  useAccount,
  useIsAuthenticated,
  useMsal,
} from "@azure/msal-react";
import { useCreation, useMemoizedFn } from "ahooks";
import { createContext, useMemo, type PropsWithChildren } from "react";
import { AdminApi, Configuration } from "../generated/apiv2";

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
  readonly account: AccountInfo | null;
  readonly isAuthenticated: boolean;
  login: () => void;
  logout: () => void;
  acquireToken: (scopes?: string[]) => Promise<AuthenticationResult | void>;
  readonly isAdmin: boolean;
  readonly api: AdminApi | undefined;
}

export const AppAuthContext = createContext<IAppAuthContext>({
  account: null,
  isAuthenticated: false,
  login: () => {},
  logout: () => {},
  acquireToken: () => Promise.resolve(undefined),
  isAdmin: false,
  api: undefined,
});

function useActiveAccount() {
  const { accounts } = useMsal();
  let account: AccountInfo | undefined;
  if (accounts && accounts.length > 0) {
    account = accounts[0];
  }
  return useAccount(account);
}

function AuthContextProvider({ children }: PropsWithChildren) {
  const { instance } = useMsal();
  const account = useActiveAccount();
  const logout = useMemoizedFn(() => {
    instance.logoutRedirect();
  });

  const acquireToken = useMemoizedFn(
    async (
      scopes: string[] = [import.meta.env.VITE_API_SCOPE]
    ): Promise<AuthenticationResult | void> => {
      return await instance.acquireTokenSilent({
        scopes,
        account: account ?? undefined,
      });
    }
  );

  const login = useMemoizedFn(() => {
    instance.loginRedirect({
      scopes: [import.meta.env.VITE_API_SCOPE],
      extraScopesToConsent: ["https://graph.microsoft.com/Directory.Read.All"],
      redirectUri: import.meta.env.VITE_MSAL_REDIRECT_URI,
    });
  });

  const isAuthenticated = useIsAuthenticated(account ?? undefined);

  const isAdmin = useMemo(
    () => !!account?.idTokenClaims?.roles?.includes("App.Admin"),
    [account]
  );

  const api = useCreation(() => {
    return new AdminApi(
      new Configuration({
        basePath: import.meta.env.VITE_API_BASE_PATH,
        accessToken: async () => {
          const result = await acquireToken();
          return result?.accessToken || "";
        },
      })
    );
  }, [account, acquireToken]);

  return (
    <AppAuthContext.Provider
      value={{
        account,
        isAuthenticated,
        login,
        logout,
        acquireToken,
        isAdmin,
        api,
      }}
    >
      {children}
    </AppAuthContext.Provider>
  );
}

export function AuthProvider({ children }: PropsWithChildren) {
  return (
    <MsalProvider instance={pca}>
      <AuthContextProvider>{children}</AuthContextProvider>
    </MsalProvider>
  );
}
