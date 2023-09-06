import { useRequest } from "ahooks";
import { useAppAuthContext } from "./auth/AuthProvider";
import { DirectoryApi } from "./generated";
import { useAuthedClient } from "./utils/useCertsApi";

export function MainPage() {
  const client = useAuthedClient(DirectoryApi);
  const { account } = useAppAuthContext();
  useRequest(
    async () => {
      if (account?.tenantId && account.localAccountId) {
        client.registerNamespaceV1({ namespaceId: account.localAccountId });
      }
    },
    { manual: true }
  );

  return (
    <>
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold tracking-tight text-gray-900">
            Coming Soon
          </h1>
        </div>
      </header>
      <main>
        <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
          {/* Your content */}
        </div>
      </main>
    </>
  );
}
