import React, { useContext } from "react";
import { useParams } from "react-router-dom";
import { NamespaceKind } from "../../generated";
import { NamespaceConfigContextProvider } from "./NamespaceConfigContextProvider";
import { NamespaceContext } from "./NamespaceContext";
import { NamespaceProvider } from "../../generated/apiv2";

export function NamespaceContextRouteProvider(
  props: React.PropsWithChildren<{}>
) {
  const { nsKind, nsId } = useParams() as {
    nsKind: NamespaceKind;
    nsId: string;
  };

  return (
    <NamespaceContext.Provider
      value={{ namespaceId: nsId, namespaceKind: nsKind }}
    >
      <NamespaceConfigContextProvider
        ruleIssuer={
          nsKind === NamespaceKind.NamespaceKindRootCA ||
          nsKind === NamespaceKind.NamespaceKindIntermediateCA
        }
        ruleEntraClientCred={
          nsKind === NamespaceKind.NamespaceKindServicePrincipal
        }
      >
        {props.children}
      </NamespaceConfigContextProvider>
    </NamespaceContext.Provider>
  );
}

export function NamespaceContextRouteProvider2(
  props: React.PropsWithChildren<{}>
) {
  const { nsKind, nsId } = useParams() as {
    nsKind: NamespaceProvider;
    nsId: string;
  };

  return (
    <NamespaceContext.Provider
      value={{ namespaceId: nsId, namespaceKind: nsKind as any }}
    >
      {props.children}
    </NamespaceContext.Provider>
  );
}
export function useNamespace(): {
  namespaceId: string;
  namespaceProvider: NamespaceProvider;
} {
  const ctx = useContext(NamespaceContext);
  return {
    namespaceId: ctx.namespaceId,
    namespaceProvider: ctx.namespaceKind as any,
  };
}
