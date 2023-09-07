import { useForm } from "react-hook-form";
import { InputField } from "./FormComponents";
import { useBoolean, useRequest } from "ahooks";
import { useAuthedClient } from "../utils/useCertsApi";
import { DirectoryApi } from "../generated";
import { XCircleIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";

interface RegiserDirectoryObjectInput {
  objectId: string;
}

export default function RegisterPage() {
  const { register, handleSubmit } = useForm<RegiserDirectoryObjectInput>({
    defaultValues: {
      objectId: "",
    },
  });

  const [formInvalid, { setTrue: setFormInvalid, setFalse: clearFormInvalid }] =
    useBoolean(false);

  const client = useAuthedClient(DirectoryApi);

  const { run: registerNs, loading: registerNsLoading } = useRequest(
    async (objectId: string) => {
      await client.registerNamespaceProfileV1({
        namespaceId: objectId,
      });
      return objectId;
    },
    { manual: true }
  );

  const onSubmit = (input: RegiserDirectoryObjectInput) => {
    registerNs(input.objectId);
  };
  return (
    <>
      <h1 className="font-semibold text-4xl">Register</h1>
      <form
        className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
        onSubmit={handleSubmit(onSubmit, setFormInvalid)}
      >
        <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
          <InputField
            inputKey="objectId"
            labelContent="Azure AD Object ID"
            register={register}
            required
          />
        </div>
        {formInvalid && (
          <div className="bg-red-50 px-4 py-4 sm:px-6 ">
            <div className="flex items-center gap-x-2">
              <div className="flex-shrink-0">
                <XCircleIcon
                  className="h-5 w-5 text-red-400"
                  aria-hidden="true"
                />
              </div>
              <div>
                <h3 className="text-sm font-medium text-red-800">
                  Invalid input, please correect before proceeding
                </h3>
              </div>
            </div>
          </div>
        )}

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
    </>
  );
}
