import { Card } from "@/components/Card";
import { PageContent, PageHeader } from "@/components/PageLayout";
import Link from "next/link";

export default function SystemAppsPage() {
  return (
    <>
      <PageHeader title="System Applications" />
      <PageContent className="space-y-4 mt-6">
        <Card>
          <ul>
            <li>
              <Link
                className="text-sm leading-6 link font-semibold"
                href="/system/api"
              >
                API
              </Link>
            </li>
            <li>
              <Link
                className="text-sm leading-6 link font-semibold"
                href="/system/backend"
              >
                Backend
              </Link>
            </li>
          </ul>
        </Card>
      </PageContent>
    </>
  );
}
