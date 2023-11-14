"use client";
import { useAdminApi } from "@/client/useClient";
import { useRequest } from "ahooks";

export default function AgentsTable() {
  const client = useAdminApi();
  const { data } = useRequest(() => {
    return client.listAgents();
  }, {});
  console.log("agents: ", data);
  return <div>{JSON.stringify(data)}</div>;
}
