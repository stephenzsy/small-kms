import { useContext } from "react";
import { NamespaceContext } from "./NamespaceContext";
import { NamespaceProvider } from "../../generated/apiv2";


export function useNamespace(): {
  namespaceId: string;
  namespaceProvider: NamespaceProvider;
} {
  const ctx = useContext(NamespaceContext);
  return {
    namespaceId: ctx.namespaceId,
    namespaceProvider: ctx.namespaceKind as NamespaceProvider,
  };
}
