import SystemAppInfoCard from "@/app/system/[id]/SystemAppInfo";
import { PageContent, PageHeader } from "@/components/PageLayout";

export default function SystemAppPage({
  params: { id },
}: {
  params: { id: string };
}) {
  return (
    <>
      <PageHeader title="System Application" />
      <PageContent className="space-y-4 mt-6">
        <SystemAppInfoCard id={id} />
      </PageContent>
    </>
  );
}
