import { useRequest } from "ahooks";
import React from "react";
import { AdminApi, CertificateRuleIssuer, NamespaceKind } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

type CertificateIssuerContextValue = {
  rule: CertificateRuleIssuer | undefined;
  setRule: (rule: CertificateRuleIssuer) => void;
};

export const CertificateIssuerContext =
  React.createContext<CertificateIssuerContextValue>({
    rule: undefined,
    setRule: () => {},
  });

export function CertificateIssuerContextProvider(
  props: React.PropsWithChildren<{
    namespaceKind: NamespaceKind;
    namespaceIdentifier: string;
  }>
) {
  const { namespaceKind, namespaceIdentifier } = props;

  const adminApi = useAuthedClient(AdminApi);
  const { data, run } = useRequest(
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
    }
  );

  return (
    <CertificateIssuerContext.Provider value={{ rule: data, setRule: run }}>
      {props.children}
    </CertificateIssuerContext.Provider>
  );
}
