import React from "react";
import { useParams } from "react-router-dom";
import { NamespaceKind } from "../generated";

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
      {props.children}
    </NamespaceContext.Provider>
  );
}
