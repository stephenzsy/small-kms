"use client";

import { useAdminApi } from "@/client/useClient";
import { Card } from "@/components/Card";
import {
  DescriptionList,
  DescriptionListItem,
} from "@/components/DescriptionList";
import { useRequest } from "ahooks";

export default function SystemAppInfoCard({ id }: { id: string }) {
  const api = useAdminApi();

  const { data, loading, run } = useRequest(
    (sync?: boolean) => {
      if (sync) {
        return api.syncSystemApp({ id: id });
      }
      return api.getSystemApp({
        id: id,
      });
    },
    {
      refreshDeps: [id],
    }
  );

  console.log(data);
  return (
    <Card title="Application Information" className="space-y-4">
      <DescriptionList>
        <DescriptionListItem term="ID">{data?.id}</DescriptionListItem>
        <DescriptionListItem term="Display Name">
          {data?.displayName}
        </DescriptionListItem>
      </DescriptionList>
      <div>
        <button
          className="btn btn-primary btn-lg"
          disabled={loading}
          onClick={() => {
            run(true);
          }}
        >
          Sync
        </button>
      </div>
    </Card>
  );
}
