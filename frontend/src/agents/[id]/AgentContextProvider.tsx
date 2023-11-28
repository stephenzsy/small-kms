import { useRequest } from "ahooks";
import { PropsWithChildren } from "react";
import { useParams } from "react-router-dom";
import { useAdminApi } from "../../utils/useCertsApi";
import { AgentContext } from "../../admin/contexts/AgentContext";

function useAgent(agentId: string | undefined) {
  const api = useAdminApi();
  return useRequest(
    async () => {
      if (agentId) {
        return api?.getAgent({
          id: agentId,
        });
      }
    },
    {
      refreshDeps: [agentId],
    }
  );
}

export function AgentContextProvider({ children }: PropsWithChildren) {
  const { id } = useParams<{ id: string }>();
  const { data: agent } = useAgent(id);

  return (
    <AgentContext.Provider
      value={{
        agent,
      }}
    >
      {children}
    </AgentContext.Provider>
  );
}
