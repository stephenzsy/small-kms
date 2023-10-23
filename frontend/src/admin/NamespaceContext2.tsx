import React from "react";
import { useParams } from "react-router-dom";
import { NamespaceKind } from "../generated";

export const NamespaceContext = React.createContext<{
  namespaceKind: NamespaceKind;
  namespaceId: string;
}>({ namespaceKind: "" as never, namespaceId: "" as never });

export function NamespaceContextProvider(props: React.PropsWithChildren<{}>) {
  const { nsKind, nsId } = useParams() as {
    nsKind: NamespaceKind;
    nsId: string;
  };

  return (
    <NamespaceContext.Provider
      value={{ namespaceId: nsId, namespaceKind: nsKind }}
    >
      {props.children}
    </NamespaceContext.Provider>
  );
}
