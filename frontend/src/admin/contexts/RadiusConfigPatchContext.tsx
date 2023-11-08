import { useRequest } from "ahooks";
import {
  AdminApi,
  AgentConfigRadius,
  AgentConfigRadiusFields,
  NamespaceKind,
} from "../../generated";
import { Result } from "ahooks/lib/useRequest/src/types";
import { PropsWithChildren, createContext, useContext } from "react";
import { useAuthedClient } from "../../utils/useCertsApi";

export type RadiusConfigPatchContextValue = Result<
  AgentConfigRadius,
  [AgentConfigRadiusFields?]
>;

export const RadiusConfigPatchContext =
  createContext<RadiusConfigPatchContextValue>(
    {} as RadiusConfigPatchContextValue
  );

export function RadiusConfigPatchProvider({
  children,
  namespaceId,
  namespaceKind,
}: PropsWithChildren<{ namespaceKind: NamespaceKind; namespaceId: string }>) {
  const api = useAuthedClient(AdminApi);
  const v = useRequest((patchOps?: AgentConfigRadiusFields) => {
    if (patchOps) {
      return api.patchAgentConfigRadius({
        agentConfigRadiusFields: patchOps,
        namespaceId,
        namespaceKind,
      });
    }
    return api.getAgentConfigRadius({
      namespaceId,
      namespaceKind,
    });
  }, {});
  return (
    <RadiusConfigPatchContext.Provider value={v}>
      {children}
    </RadiusConfigPatchContext.Provider>
  );
}

export function useRadiusConfigPatch() {
  return useContext(RadiusConfigPatchContext);
}
