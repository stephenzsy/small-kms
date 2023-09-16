import { useRequest } from "ahooks";
import { useAppAuthContext } from "./auth/AuthProvider";
import { DirectoryApi, NamespaceType } from "./generated";
import { useAuthedClient } from "./utils/useCertsApi";
import { useMemo } from "react";
import classNames from "classnames";

export default function MainPage() {
  const client = useAuthedClient(DirectoryApi);

  const { account } = useAppAuthContext();

  const { data: profiles } = useRequest(() => client.getMyProfilesV1(), {
    ready: !!account,
  });

  const { run: syncProfiles, loading: syncProfilesLoading } = useRequest(() => client.syncMyProfilesV1(), {
    manual: true,
  });

  return (
    <>
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold tracking-tight text-gray-900">
            My profiles
          </h1>
        </div>
      </header>
      <main className="space-y-6 p-6">
        <div className="flex flex-row items-center gap-x-4 p-6 bg-white rounded-md shadow-sm">
          <button
            type="button"
            disabled={syncProfilesLoading}
            onClick={syncProfiles}
            className={classNames(
              "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
              "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
            )}
          >
            {syncProfilesLoading ? "Syncing..." : "Sync all profiles"}
          </button>
        </div>
        {profiles === undefined ? (
          <div>Loading...</div>
        ) : profiles.length ? (
          profiles.map((profile) => (
            <div
              className="gap-x-4 p-6 bg-white rounded-md shadow-sm"
              key={profile.id}
            >
              <dl className="space-y-2">
                <div>
                  <dt className="font-semibold">Display name</dt>
                  <dd>{profile.displayName}</dd>
                </div>
                <div>
                  <dt className="font-semibold">Type</dt>
                  <dd>{profile.objectType}</dd>
                </div>
              </dl>
            </div>
          ))
        ) : (
          <div>No profiles</div>
        )}
      </main>
    </>
  );
}
