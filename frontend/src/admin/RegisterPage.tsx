import { XCircleIcon } from "@heroicons/react/24/outline";
import { useRequest } from "ahooks";
import classNames from "classnames";
import { useState } from "react";
import { useParams } from "react-router-dom";
import {
  AdminApi,
  CreateManagedApplicationProfileRequest,
  CreateManagedApplicationProfileRequestFromJSON,
  CreateProfileRequestType,
  NamespaceKind,
} from "../generated3";
import { useAuthedClient } from "../utils/useCertsApi3";
import { InputField } from "./InputField";

interface RegiserDirectoryObjectInput {
  objectId: string;
}

export default function RegisterPage() {
  const { profileType } = useParams() as {
    profileType: NamespaceKind;
  };
  const [objectId, setObjectId] = useState<string>("");
  const [managedApplication, setManagedApplicationName] = useState<string>("");
  const client = useAuthedClient(AdminApi);

  const { run: registerNs, loading: registerNsLoading } = useRequest(
    async (oid: string) => {
      await client.syncProfile({
        namespaceId: oid,
        namespaceKind: profileType,
      });
      return oid;
    },
    { manual: true }
  );

  const {
    run: createManagedApplication,
    loading: createManagedApplicationLoading,
  } = useRequest(
    async (managedAppliationName: string) => {
      await client.createProfile({
        createProfileRequest: {
          type: CreateProfileRequestType.ProfileTypeManagedApplication,
          name: managedAppliationName,
        } as CreateManagedApplicationProfileRequest,
        namespaceKind: profileType,
      });
    },
    { manual: true }
  );

  return (
    <>
      <h1 className="font-semibold text-4xl">Register</h1>
      <form
        className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
        onSubmit={(e) => {
          e.preventDefault();
          registerNs(objectId);
        }}
      >
        <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
          <InputField
            labelContent="Azure AD Object ID"
            value={objectId}
            onChange={setObjectId}
            required
          />
        </div>

        <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center justify-end gap-x-6">
          <button
            type="submit"
            disabled={registerNsLoading}
            className={classNames(
              "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
              "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
            )}
          >
            {registerNsLoading ? "Registering..." : "Register"}
          </button>
        </div>
      </form>
      {profileType === NamespaceKind.NamespaceKindApplication && (
        <form
          className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
          onSubmit={(e) => {
            e.preventDefault();
            createManagedApplication(managedApplication);
          }}
        >
          <h2>Create managed application</h2>
          <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
            <InputField
              labelContent="Application display name"
              value={managedApplication}
              onChange={setManagedApplicationName}
              required
            />
          </div>

          <div className="bg-neutral-50 px-4 py-4 sm:px-6 flex align-center justify-end gap-x-6">
            <button
              type="submit"
              disabled={createManagedApplicationLoading}
              className={classNames(
                "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
                "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
              )}
            >
              {registerNsLoading ? "Creating..." : "Creating"}
            </button>
          </div>
        </form>
      )}
    </>
  );
}
