"use client";
import { useRequest } from "ahooks";

import {
  CreateCertificateParameters,
  CreateCertificateParametersToJSON,
  CreateCertificateRequest,
} from "@/generated";
import Link from "next/link";
import { useState } from "react";

function TextField({
  value,
  onChange,
  id,
  label,
}: {
  value: string | undefined;
  onChange: (nextValue: string) => void;
  id: string;
  label: React.ReactNode;
}) {
  return (
    <div className="sm:col-span-4">
      <label
        htmlFor={id}
        className="block text-sm font-medium leading-6 text-gray-900"
      >
        {label}
      </label>
      <div className="mt-2">
        <div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
          <input
            id={id}
            type="text"
            className="block flex-1 border-0 bg-transparent py-1.5 px-4 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder="Sample Internal Root CA"
            onChange={(e) => {
              onChange(e.target.value);
            }}
            value={value}
          />
        </div>
      </div>
    </div>
  );
}

export function CreateCeritificateForm() {
  const [subjectCN, setSubjectCN] = useState<string>();
  const [subjectOU, setSubjectOU] = useState<string>();
  const [subjectO, setSubjectO] = useState<string>();
  const [subjectC, setSubjectC] = useState<string>();

  const { run: sendRequest } = useRequest(
    async (params: CreateCertificateParameters) => {
      const resp = await fetch("/api/admin/certificate", {
        body: JSON.stringify(CreateCertificateParametersToJSON(params)),
        method: "POST",
      });
      return await resp.json();
    },
    { manual: true }
  );

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        if (!subjectCN) {
          return;
        }
        const params: CreateCertificateParameters = {
          category: "root-ca",
          name: "root",
          subject: {
            commonName: subjectCN,
            organizationUnit: subjectOU,
            organization: subjectO,
            country: subjectC,
          },
          options: {
            keepKeyVersion: true,
          },
        };
        sendRequest(params);
      }}
    >
      <div className="space-y-12">
        <div className="border-b border-gray-900/10 pb-12">
          <h2 className="text-base font-semibold leading-7 text-gray-900">
            Create CA Certificate
          </h2>

          <div className="mt-10">
            <h3 className="text-base font-medium leading-7 text-gray-900">
              Subject
            </h3>

            <div className="mt-6 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
              <TextField
                value={subjectCN}
                label="Common name"
                onChange={setSubjectCN}
                id="cn"
              />
              <TextField
                value={subjectOU}
                label="Organization unit"
                onChange={setSubjectOU}
                id="ou"
              />
              <TextField
                value={subjectO}
                label="Organization"
                onChange={setSubjectO}
                id="o"
              />
              <TextField
                value={subjectC}
                label="Country"
                onChange={setSubjectC}
                id="c"
              />
            </div>
          </div>
        </div>
      </div>
      <div className="mt-6 flex items-center justify-end gap-x-6">
        <Link
          type="button"
          className="text-sm font-semibold leading-6 text-gray-900"
          href="/admin/ca"
        >
          Cancel
        </Link>
        <button
          type="submit"
          className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        >
          Create
        </button>
      </div>
    </form>
  );
}
