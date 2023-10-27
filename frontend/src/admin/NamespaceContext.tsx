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
      {nsKind === NamespaceKind.NamespaceKindRootCA ||
      nsKind === NamespaceKind.NamespaceKindIntermediateCA ? (
        <CertificateIssuerContextProvider
          namespaceKind={nsKind}
          namespaceIdentifier={nsId}
        >
          {props.children}
        </CertificateIssuerContextProvider>
      ) : (
        props.children
      )}
    </NamespaceContext.Provider>
  );
}
