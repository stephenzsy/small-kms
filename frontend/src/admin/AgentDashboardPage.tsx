import { useParams } from "react-router-dom";
import { Card, CardSection, CardTitle } from "../components/Card";
import { useRequest } from "ahooks";
import { Button } from "../components/Button";
import { useAuthedClient } from "../utils/useCertsApi";
import { AdminApi } from "../generated";

export default function AgentDashboardPage() {
  const { namespaceId } = useParams() as { namespaceId: string };

  const client = useAuthedClient(AdminApi);
  const { run } = useRequest(
    async () => {
      await client.getDockerInfo({
        namespaceId,
      });
    },
    { manual: true }
  );
  return (
    <>
      <h1>Agent Dashboard</h1>
      <Card>
        <CardTitle>Start</CardTitle>
        <CardSection>
          <Button onClick={run} variant="primary">
            Start
          </Button>
        </CardSection>
      </Card>
    </>
  );
}
