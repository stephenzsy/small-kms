import { useAppAuthContext } from "./auth/AuthProvider";
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
  const { account } = useAppAuthContext();

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
        <div className="flex flex-row items-center gap-x-4 p-6 bg-white rounded-md shadow-sm"></div>

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
              {/*managedDevices?.map((d) => (
                <tr key={d.id}>
                  <td>{d.id}</td>
                  <td>{d.displayName}</td>
                  <td>
                    <Link to={"/my/enroll"}>Enroll certificate</Link>
                  </td>
                </tr>
              ))*/}
            </tbody>
          </table>
        </div>
      </main>
    </>
  );
}
