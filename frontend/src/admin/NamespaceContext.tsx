import React from "react";
import { useParams } from "react-router-dom";
import { NamespaceKind } from "../generated";
import { CertificateIssuerContextProvider } from "./CertIssuerContext";

export const NamespaceContext = React.createContext<{
  namespaceKind: NamespaceKind;
  namespaceIdentifier: string;
}>({ namespaceKind: "" as never, namespaceIdentifier: "" as never });

export function NamespaceContextProvider(props: React.PropsWithChildren<{}>) {
  const { nsKind, nsId } = useParams() as {
    nsKind: NamespaceKind;
    nsId: string;
  };

  return (
    <NamespaceContext.Provider
      value={{ namespaceIdentifier: nsId, namespaceKind: nsKind }}
    >
      <CertificateIssuerContextProvider
        namespaceKind={nsKind}
        namespaceIdentifier={nsId}
        ruleIssuer={
          nsKind === NamespaceKind.NamespaceKindRootCA ||
          nsKind === NamespaceKind.NamespaceKindIntermediateCA
        }
        ruleEntraClientCred={
          nsKind === NamespaceKind.NamespaceKindServicePrincipal
        }
      >
        {props.children}
      </CertificateIssuerContextProvider>
    </NamespaceContext.Provider>
  );
}
