import { useParams } from "react-router-dom";
import { InputField } from "./InputField";
import React from "react";
import {
  DirectoryApi,
  NamespacePermissions,
  PutPermissionsV1Request,
} from "../generated";
import { useRequest, useSetState } from "ahooks";
import { Button } from "../components/Button";
import { useAuthedClient } from "../utils/useCertsApi";

const keys: [keyof NamespacePermissions, string][] = [
  ["allowEnrollDeviceCertificate", "Allow enroll device certificate"],
];

function PutPermissionsForm(props: { namespaceId: string }) {
  const [objectId, setObjectId] = React.useState("");
  const [permissions, setPermissions] = useSetState<NamespacePermissions>({});

  const client = useAuthedClient(DirectoryApi);
  const { run: putPermissions } = useRequest(
    (r: PutPermissionsV1Request) => {
      return client.putPermissionsV1(r);
    },
    { manual: true }
  );
  return (
    <form
      className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow p-6 space-y-4"
      onSubmit={(e) => {
        e.preventDefault();
        if (!objectId) {
          return;
        }
        putPermissions({
          namespaceId: props.namespaceId,
          objectId,
          namespacePermissions: permissions,
        });
      }}
    >
      <h2 className="font-semibold text-2xl">Add permission</h2>
      <InputField
        className="pt-4"
        labelContent="AAD Object ID"
        required
        value={objectId}
        onChange={setObjectId}
      />
      <div className="pt-4">
        <h3 className="font-medium text-large">Permissions</h3>
        <ul className="pt-4">
          {keys.map(([key, label]) => (
            <li key={key}>
              <label>
                <input
                  type="checkbox"
                  checked={permissions[key]}
                  onChange={(e) => {
                    setPermissions({ [key]: e.target.checked });
                  }}
                />{" "}
                {label}
              </label>
            </li>
          ))}
        </ul>
      </div>
      <div className="pt-4">
        <Button variant="primary" type="submit">
          Set permission
        </Button>
      </div>
    </form>
  );
}

export default function PermissionsPage() {
  const { namespaceId } = useParams();
  const client = useAuthedClient(DirectoryApi);

  const { data: permissions } = useRequest(
    () => {
      return client.hasPermissionV1({
        namespaceId: namespaceId!,
        permissionKey: "allowEnrollDeviceCertificate",
      });
    },
    { refreshDeps: [namespaceId] }
  );
  return (
    <>
      <h1 className="font-semibold text-4xl">Permissions</h1>
      <div className="bg-white">
        Current permissions of allowEnrollDeviceCertificate:
        <pre>{JSON.stringify(permissions, undefined, 2)}</pre>
      </div>
      <PutPermissionsForm namespaceId={namespaceId!} />
    </>
  );
}
