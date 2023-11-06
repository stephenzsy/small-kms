import React from "react";
import {
  CertificateRuleIssuer,
  CertificateRuleMsEntraClientCredential,
  NamespaceKind,
} from "../../generated";

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
  entraClientCred: CertificateRuleMsEntraClientCredential | undefined;
  setEntraClientCred: (rule: CertificateRuleMsEntraClientCredential) => void;
};

export const NamespaceConfigContext =
  React.createContext<NamespaceConfigContextValue>({
    issuer: undefined,
    setIssuer: () => {},
    entraClientCred: undefined,
    setEntraClientCred: () => {},
  });
