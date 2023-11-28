import React from "react";
import { Profile } from "../../generated/apiv2";

type AgentContextValue = {
  agent: Profile | undefined;
};

export const AgentContext = React.createContext<AgentContextValue>({
  agent: undefined,
});
