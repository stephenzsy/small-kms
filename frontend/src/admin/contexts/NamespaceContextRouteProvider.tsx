import React from "react";
import { useParams } from "react-router-dom";
import { NamespaceKind } from "../../generated";
import { NamespaceConfigContextProvider } from "./NamespaceConfigContextProvider";
import { NamespaceContext } from "./NamespaceContext";

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
