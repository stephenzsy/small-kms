import { useRequest } from "ahooks";
import React, { useContext } from "react";
import {
  AdminApi,
  CertificateRuleIssuer,
  CertificateRuleMsEntraClientCredential,
} from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";
import { NamespaceConfigContext, NamespaceContext } from "./NamespaceContext";

export function NamespaceConfigContextProvider(
  props: React.PropsWithChildren<{
    ruleIssuer?: boolean;
    ruleEntraClientCred?: boolean;
  }>
) {
  const { namespaceKind, namespaceId: namespaceIdentifier } = useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: issuer, run: setIssuer } = useRequest(
    (params?: CertificateRuleIssuer) => {
      if (params) {
        return adminApi.putCertificateRuleIssuer({
          namespaceId: namespaceIdentifier,
          namespaceKind,
          certificateRuleIssuer: params,
        });
      }
      return adminApi.getCertificateRuleIssuer({
        namespaceId: namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!props.ruleIssuer && !!namespaceIdentifier,
    }
  );

  const { data: msEntraClientCred, run: setMsEntraClientCred } = useRequest(
    (params?: CertificateRuleMsEntraClientCredential) => {
      if (params) {
        return adminApi.putCertificateRuleMsEntraClientCredential({
          namespaceId: namespaceIdentifier,
          namespaceKind,
          certificateRuleMsEntraClientCredential: params,
        });
      }
      return adminApi.getCertificateRuleMsEntraClientCredential({
        namespaceId: namespaceIdentifier,
        namespaceKind,
      });
    },
    {
      refreshDeps: [namespaceIdentifier, namespaceKind],
      ready: !!props.ruleEntraClientCred && !!namespaceIdentifier,
    }
  );

  return (
    <NamespaceConfigContext.Provider
      value={{
        issuer: issuer,
        setIssuer: setIssuer,
        entraClientCred: msEntraClientCred,
        setEntraClientCred: setMsEntraClientCred,
      }}
    >
      {props.children}
    </NamespaceConfigContext.Provider>
  );
}
