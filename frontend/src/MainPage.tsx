import { useRequest } from "ahooks";
import classNames from "classnames";
import { useAppAuthContext } from "./auth/AuthProvider";
import { DirectoryApi } from "./generated";
import { useAuthedClient } from "./utils/useCertsApi";
import { Link } from "react-router-dom";
/*
async function genKeypair() {
  const subtle = new SubtleCrypto();
  const kp = await subtle.generateKey(
    {
      name: "RSASSA-PKCS1-v1_5",
      modulusLength: 2048,
      publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
      hash: "SHA-384",
    },
    true,
    ["encrypt", "decrypt", "sign", "verify"]
  );
  subtle.exportKey("jwk", kp.publicKey);
}
*/
export default function MainPage() {
  const client = useAuthedClient(DirectoryApi);

  const { account } = useAppAuthContext();

  const { data: profiles } = useRequest(() => client.getMyProfilesV1(), {
    ready: !!account,
  });

  const { run: syncProfiles, loading: syncProfilesLoading } = useRequest(
    () => client.syncMyProfilesV1(),
    {
      manual: true,
    }
  );

  const { data: managedDevices, loading: managedDevicesLoading } = useRequest(
    () =>
      client.myHasPermissionV1({
        permissionKey: "allowEnrollDeviceCertificate",
      }),
    {}
  );
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
        <div className="gap-x-4 p-6 bg-white rounded-md shadow-sm">
          <h2 className="text-2xl font-semibold">My managed devices</h2>
          <table className="mt-6 min-w-full divide-y divide-gray-300">
            <thead>
              <tr className="text-left">
                <th>ID</th>
                <th>Name</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {managedDevices?.map((d) => (
                <tr key={d.id}>
                  <td>{d.id}</td>
                  <td>{d.displayName}</td>
                  <td>
                    <Link to={"/my/enroll"}>Enroll certificate</Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </>
  );
}
