import { useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi";
import { DiagnosticsApi } from "../generated";

export default function DiagnosticsPage() {
  const client = useAuthedClient(DiagnosticsApi);
  const { data: diagnosticsData } = useRequest(() => {
    return client.getDiagnosticsV1();
  }, {});
  return (
    <main className="min-h-full place-items-center p-6">
      <h1 className="mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-5xl">
        Diagnostics
      </h1>
      <div className="mt-10 rounded-md bg-white text-sm p-6 overflow-x-auto max-w-full">
        <pre>{JSON.stringify(diagnosticsData, undefined, 2)}</pre>
      </div>
    </main>
  );
}
