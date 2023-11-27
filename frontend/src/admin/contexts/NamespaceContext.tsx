import React from "react";
import { CertificateRuleIssuer, NamespaceKind } from "../../generated";

export type NamespaceContextValue = {
  namespaceKind: NamespaceKind;
  namespaceId: string;
};

export const NamespaceContext = React.createContext<NamespaceContextValue>({
  namespaceKind: "" as never,
  namespaceId: "" as never,
});

type NamespaceConfigContextValue = {
  issuer: CertificateRuleIssuer | undefined;
  setIssuer: (rule: CertificateRuleIssuer) => void;
};

export const NamespaceConfigContext =
  React.createContext<NamespaceConfigContextValue>({
    issuer: undefined,
    setIssuer: () => {},
  });
