"use client";

import {
  BrowserCacheLocation,
  Configuration,
  InteractionType,
  PublicClientApplication,
} from "@azure/msal-browser";
import { MsalAuthenticationTemplate, MsalProvider } from "@azure/msal-react";
import { useRequest } from "ahooks";
import { PropsWithChildren } from "react";

const msalConfig: Configuration = {
  auth: {
    clientId: process.env.NEXT_PUBLIC_MSAL_CLIENT_ID!,
    redirectUri: "http://localhost:3000/admin",
    authority: `https://login.microsoftonline.com/${process.env.NEXT_PUBLIC_MSAL_TENANT_ID}/`,
  },
  cache: {
    cacheLocation: BrowserCacheLocation.SessionStorage,
    storeAuthStateInCookie: false,
  },
};
const msalInstance = new PublicClientApplication(msalConfig);

export function LoginProvider(props: PropsWithChildren<{}>) {
  useRequest(async () => {
    await msalInstance.initialize();
    msalInstance.setActiveAccount(msalInstance.getAllAccounts()[0]);
  });

  return (
    <MsalProvider instance={msalInstance}>
      <MsalAuthenticationTemplate interactionType={InteractionType.Redirect}>
        {props.children}
      </MsalAuthenticationTemplate>
    </MsalProvider>
  );
}
