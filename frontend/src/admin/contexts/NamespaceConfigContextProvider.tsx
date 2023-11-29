import { useRequest } from "ahooks";
import React, { useContext } from "react";
import { AdminApi } from "../../generated";
import { useAuthedClient } from "../../utils/useCertsApi";
import { NamespaceConfigContext, NamespaceContext } from "./NamespaceContext";

export function NamespaceConfigContextProvider(
  props: React.PropsWithChildren<{
    ruleIssuer?: boolean;
    ruleEntraClientCred?: boolean;
  }>
) {
  const { namespaceKind, namespaceId: namespaceIdentifier } =
    useContext(NamespaceContext);

  const adminApi = useAuthedClient(AdminApi);
  const { data: issuer, run: setIssuer } = useRequest(
    () => {
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

  return (
    <NamespaceConfigContext.Provider
      value={{
        issuer: issuer,
        setIssuer: setIssuer,
      }}
    >
      {props.children}
    </NamespaceConfigContext.Provider>
  );
}
