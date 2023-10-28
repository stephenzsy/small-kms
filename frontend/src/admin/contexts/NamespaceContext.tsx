import React from "react";
import {
  CertificateRuleIssuer,
  CertificateRuleMsEntraClientCredential,
  NamespaceKind,
} from "../../generated";

export const NamespaceContext = React.createContext<{
  namespaceKind: NamespaceKind;
  namespaceIdentifier: string;
}>({ namespaceKind: "" as never, namespaceIdentifier: "" as never });

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
