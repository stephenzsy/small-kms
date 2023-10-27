import { useRequest } from "ahooks";
import React from "react";
import {
  AdminApi,
  CertificateRuleIssuer,
  CertificateRuleMsEntraClientCredential,
  NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

type CertificateRulesContextValue = {
  issuer: CertificateRuleIssuer | undefined;
  setIssuer: (rule: CertificateRuleIssuer) => void;
  entraClientCred: CertificateRuleMsEntraClientCredential | undefined;
  setEntraClientCred: (rule: CertificateRuleMsEntraClientCredential) => void;
};

export const CertificateIssuerContext =
  React.createContext<CertificateRulesContextValue>({
    issuer: undefined,
    setIssuer: () => {},
    entraClientCred: undefined,
    setEntraClientCred: () => {},
  });

export function CertificateIssuerContextProvider(
  props: React.PropsWithChildren<{
    namespaceKind: NamespaceKind;
    namespaceIdentifier: string;
    ruleIssuer?: boolean;
    ruleEntraClientCred?: boolean;
  }>
) {
  const { namespaceKind, namespaceIdentifier } = props;

  const adminApi = useAuthedClient(AdminApi);
  const { data: issuer, run: setIssuer } = useRequest(
    (params?: CertificateRuleIssuer) => {
      if (params) {
        return adminApi.putCertificateRuleIssuer({
          namespaceIdentifier,
          namespaceKind,
          certificateRuleIssuer: params,
        });
      }
      return adminApi.getCertificateRuleIssuer({
        namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!props.ruleIssuer,
    }
  );

  const { data: msEntraClientCred, run: setMsEntraClientCred } = useRequest(
    (params?: CertificateRuleMsEntraClientCredential) => {
      if (params) {
        return adminApi.putCertificateRuleMsEntraClientCredential({
          namespaceIdentifier,
          namespaceKind,
          certificateRuleMsEntraClientCredential: params,
        });
      }
      return adminApi.getCertificateRuleMsEntraClientCredential({
        namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!props.ruleEntraClientCred,
    }
  );

  return (
    <CertificateIssuerContext.Provider
      value={{
        issuer: issuer,
        setIssuer: setIssuer,
        entraClientCred: msEntraClientCred,
        setEntraClientCred: setMsEntraClientCred,
      }}
    >
      {props.children}
    </CertificateIssuerContext.Provider>
  );
}
