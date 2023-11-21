import { Button, Card, Typography } from "antd";
import { useSystemAppRequest } from "../admin/forms/useSystemAppRequest";
import { Link } from "../components/Link";

function SystemAppsEntry({ systemAppName }: { systemAppName: string }) {
  const { data: systemApp, run: syncSystemApp } =
    useSystemAppRequest(systemAppName);
  return (
    <div className="space-y-4">
      <dl className="dl">
        <div className="dl-item">
          <dt>Application ID (Client ID)</dt>
          <dd className="font-mono">{systemApp?.id}</dd>
        </div>
        <div className="dl-item">
          <dt>Display Name</dt>
          <dd>{systemApp?.displayName}</dd>
        </div>
        <div className="dl-item">
          <dt>Service Principal ID</dt>
          <dd className="font-mono">{systemApp?.servicePrincipalId}</dd>
        </div>
      </dl>
      <div className="flex items-center gap-4">
        <Link to={`/service-principal/${systemApp?.servicePrincipalId}`}>View</Link>
        <Button
          type="link"
          onClick={() => {
            syncSystemApp(true);
          }}
        >
          Sync
        </Button>
      </div>
    </div>
  );
}

export default function SystemAppsPage() {
  return (
    <>
      <Typography.Title>System applications</Typography.Title>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Card title="API">
          <SystemAppsEntry systemAppName="api" />
        </Card>
        <Card title="Backend">
          <SystemAppsEntry systemAppName="backend" />
        </Card>
      </div>
    </>
  );
}
