"use client";

import {
  AccountInfo,
  AuthenticationResult,
  EventType,
  PublicClientApplication,
} from "@azure/msal-browser";
import {
  AuthenticatedTemplate,
  MsalProvider,
  UnauthenticatedTemplate,
} from "@azure/msal-react";
import type { PropsWithChildren } from "react";

const msalInstance = new PublicClientApplication({
  auth: {
    clientId: process.env.NEXT_PUBLIC_AZURE_CLIENT_ID!,
    authority: `https://login.microsoftonline.com/${process.env.NEXT_PUBLIC_AZURE_TENANT_ID}`,
    redirectUri: process.env.NEXT_PUBLIC_MSAL_REDIRECT_URI,
  },
  cache: {
    cacheLocation: "sessionStorage",
  },
});

msalInstance.initialize().then(() => {
  // Account selection logic is app dependent. Adjust as needed for different use cases.
  const accounts = msalInstance.getAllAccounts();
  if (accounts.length > 0) {
    msalInstance.setActiveAccount(accounts[0]);
  }

  msalInstance.addEventCallback((event) => {
    if (
      event.eventType === EventType.LOGIN_SUCCESS &&
      (event.payload as AuthenticationResult).account
    ) {
      const account = (event.payload as AuthenticationResult).account;
      msalInstance.setActiveAccount(account);
    }
  });
});

export function AppMsalProvider(props: PropsWithChildren<{}>) {
  console.log(process.env.NEXT_PUBLIC_AZURE_TENANT_ID)
  return (
    <MsalProvider instance={msalInstance}>
      <AuthenticatedTemplate>{props.children}</AuthenticatedTemplate>
      <UnauthenticatedTemplate>
        You must authenticate{" "}
        <button
          onClick={() => {
            msalInstance.loginRedirect({
              scopes: [process.env.NEXT_PUBLIC_API_SCOPE!],
              redirectUri: process.env.NEXT_PUBLIC_MSAL_REDIRECT_URI,
            });
          }}
        >
          Login
        </button>
      </UnauthenticatedTemplate>
    </MsalProvider>
  );
}
