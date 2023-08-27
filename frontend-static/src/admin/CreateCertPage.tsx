import {
  PhotoIcon,
  UserCircleIcon,
  XCircleIcon,
} from "@heroicons/react/24/solid";
import { useBoolean, useRequest } from "ahooks";
import { useEffect, useId, useMemo } from "react";
import { useForm, type UseFormRegister } from "react-hook-form";
import { redirect, useNavigate, useParams } from "react-router-dom";
import {
  CertificateUsage,
  CreateCertificateV1Request,
  WellKnownNamespaceId,
} from "../generated";
import { useCertsApi } from "../utils/useCertsApi";
import {
  AdminBreadcrumb,
  type BreadcrumbPageMetadata,
} from "./AdminBreadcrumb";
import { caBreadcrumPages } from "./CaPage";
import { NIL as uuidNil } from "uuid";
import classNames from "classnames";

interface CreateReactFormInput {
  subjectCN: string;
  subjectOU?: string;
  subjectO?: string;
  subjectC?: string;
}

interface TextInputFieldProps {
  labelContent: React.ReactNode;
  defaultValue?: string;
  required?: boolean;
  register: UseFormRegister<CreateReactFormInput>;
  inputKey: keyof CreateReactFormInput;
  placeholder?: string;
}

function TextInputField({
  labelContent,
  defaultValue = "",
  required = false,
  register,
  inputKey,
  placeholder,
}: TextInputFieldProps) {
  const id = useId();
  return (
    <div className="sm:col-span-4">
      <label
        htmlFor={id}
        className="block text-sm font-medium leading-6 text-gray-900"
      >
        {labelContent}
        {required && <span className="ml-1 text-red-500">*</span>}
      </label>
      <div className="mt-2">
        <div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
          <input
            type="text"
            defaultValue={defaultValue}
            id={id}
            {...register(inputKey, { required })}
            className="block flex-1 border-0 bg-transparent py-1.5 px-em text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder={placeholder}
          />
        </div>
      </div>
    </div>
  );
}

export default function CreateCertPage() {
  const titleId = useId();
  const client = useCertsApi();
  const navigate = useNavigate();
  const { loading: createCertInProgress, runAsync: createCertAsync } =
    useRequest(
      async (request: CreateCertificateV1Request) => {
        return client.createCertificateV1(request);
      },
      { manual: true }
    );

  const { register, handleSubmit } = useForm<CreateReactFormInput>();

  const { namespaceId } = useParams();

  const breadcrumPages: BreadcrumbPageMetadata[] = useMemo(() => {
    switch (namespaceId) {
      case WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA:
        return [...caBreadcrumPages, { name: "Create root CA", to: "#" }];
    }
    return [{ name: "Create certificate", to: "#" }];
  }, [namespaceId]);

  const [formInvalid, { setTrue: setFormInvalid, setFalse: clearFormInvalid }] =
    useBoolean(false);

  const onSubmit = async (data: CreateReactFormInput) => {
    clearFormInvalid();
    try {
      const result = await createCertAsync({
        namespaceId: namespaceId!,
        createCertificateParameters: {
          usage: CertificateUsage.Usage_RootCA,
          issuerNamespace: namespaceId!,
          issuer: uuidNil,
          subject: {
            cn: data.subjectCN,
            ou: data.subjectOU,
            o: data.subjectO,
            c: data.subjectC,
          },
        },
      });
      console.log(result);
      navigate("/admin/ca", { replace: true });
    } catch {
      // do nothing in onSubmit handler
    }
  };
  return (
    <>
      <AdminBreadcrumb pages={breadcrumPages} />
      <form
        className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow"
        onSubmit={handleSubmit(onSubmit, setFormInvalid)}
      >
        <h1 id={titleId} className="px-4 py-5 sm:px-6 text-lg font-semibold">
          {namespaceId === WellKnownNamespaceId.WellKnownNamespaceIDStr_RootCA
            ? "Create root CA"
            : "Create certificate"}
        </h1>

        <div className="px-4 py-5 sm:p-6 space-y-12 [&>*+*]:border-t [&>*+*]:pt-6">
          <div className="border-neutral-900/10 space-y-6">
            <h2 className="text-base font-semibold leading-7 text-gray-900">
              Subject
            </h2>

            <TextInputField
              inputKey="subjectCN"
              labelContent="Common Name (CN)"
              register={register}
              placeholder="Sample Internal Root CA"
              required
            />
            <TextInputField
              inputKey="subjectOU"
              labelContent="Organizational Unit (OU)"
              register={register}
              placeholder="Sample Organizational Unit"
            />

            <TextInputField
              inputKey="subjectO"
              labelContent="Organization (O)"
              register={register}
              placeholder="Sample Organization"
            />

            <TextInputField
              inputKey="subjectC"
              labelContent="Country or Region (C)"
              register={register}
              placeholder="US"
            />
          </div>
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
            disabled={createCertInProgress}
            className={classNames(
              "rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
              "disabled:bg-neutral-100 disabled:text-neutral-400 disabled:shadow-none"
            )}
          >
            {createCertInProgress ? "Creating..." : "Create"}
          </button>
        </div>
      </form>
    </>
  );
}
