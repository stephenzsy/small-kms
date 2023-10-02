import React from "react";
import {
  NamespaceProfile,
  Policy,
  PolicyParameters
} from "../generated";

export interface IPolicyContext {
  policyId: string;
  namespaceId: string;
  policy: Policy | null | undefined;
  namespaceProfile: NamespaceProfile | null | undefined;
  deletePolicy: (purge?: boolean) => void;
  applyPolicy: () => void;
  putPolicy: (p: PolicyParameters) => void;
  putPolicyError?: any;
}

export const PolicyContext = React.createContext<IPolicyContext>({
  namespaceId: "invalid",
  namespaceProfile: null,
  policyId: "invalid",
  policy: null,
  deletePolicy: () => {},
  applyPolicy: () => {},
  putPolicy: () => {},
});

/*
export function PolicyContextProvider(
  props: React.PropsWithChildren<{ policyId: string; namespaceId: string }>
) {
  const { policyId, namespaceId } = props;
  const policyApi = useAuthedClient(PolicyApi);
  const nsApi = useAuthedClient(DirectoryApi);
  // retrieve existing policy
  const { data: policy, mutate: mutatePolicy } = useRequest(
    async () => {
      try {
        return await policyApi.getPolicyV1({
          namespaceId: namespaceId,
          policyIdentifier: policyId!,
        });
      } catch (e) {
        if (e instanceof ResponseError && e.response.status === 404) {
          return null;
        }
        throw e;
      }
    },
    { refreshDeps: [policyId, namespaceId] }
  );
  const { data: ns } = useRequest(
    async () => {
      try {
        return await nsApi.getNamespaceProfileV1({
          namespaceId: namespaceId!,
        });
      } catch (e) {
        if (e instanceof ResponseError && e.response.status === 404) {
          return null;
        }
        throw e;
      }
    },
    { refreshDeps: [namespaceId] }
  );
  const { run: deletePolicy } = useRequest(
    async (purge?: boolean) => {
      try {
        const updated = await policyApi.deletePolicyV1({
          namespaceId,
          policyIdentifier: policyId,
          purge,
        });
        mutatePolicy(updated);
      } catch {
        if (purge) {
          mutatePolicy(null);
        }
      }
    },
    { manual: true }
  );

  const { run: applyPolicy } = useRequest(
    async () => {
      await policyApi.applyPolicyV1({
        namespaceId: namespaceId,
        policyId: policyId,
        applyPolicyRequest: {},
      });
    },
    { manual: true }
  );
  const { run: putPolicy, error: putPolicyError } = useRequest(
    async (p: PolicyParameters) => {
      const updated = await policyApi.putPolicyV1({
        namespaceId: namespaceId,
        policyIdentifier: policyId,
        policyParameters: p,
      });
      mutatePolicy(updated);
    },
    { manual: true }
  );

  return (
    <PolicyContext.Provider
      value={{
        namespaceId,
        policyId,
        namespaceProfile: ns,
        policy,
        deletePolicy,
        applyPolicy,
        putPolicy,
        putPolicyError,
      }}
    >
      {props.children}
    </PolicyContext.Provider>
  );
}
*/
