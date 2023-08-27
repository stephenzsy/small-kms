import { PublicClientApplication } from "@azure/msal-browser";
import { MsalProvider } from "@azure/msal-react";
import type { PropsWithChildren } from "react";

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

export function AuthProvider({ children }: PropsWithChildren<{}>) {
  return <MsalProvider instance={pca}>{children}</MsalProvider>;
}
