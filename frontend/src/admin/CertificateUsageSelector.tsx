import { CertificateUsage } from "../generated";
import { certUsageNames } from "./displayConstants";
import { useId } from "react";

export const configurableCertificateUsages: CertificateUsage[] = [
  CertificateUsage.Usage_ServerAndClient,
  CertificateUsage.Usage_ServerOnly,
  CertificateUsage.Usage_ClientOnly,
];

export function CertificateUsageSelector({
  inputType,
  availableUsages = configurableCertificateUsages,
  onChange,
  isChecked,
  label = "Certificate Usage",
}: {
  inputType: "radio" | "checkbox";
  availableUsages?: CertificateUsage[];
  onChange: React.ChangeEventHandler<HTMLInputElement>;
  isChecked(usage: CertificateUsage): boolean;
  label?: React.ReactNode;
}) {
  const idBase = useId();
  return (
    <fieldset className="space-y-4">
      <legend className="text-base font-semibold text-gray-900">{label}</legend>
      <div className="space-y-4">
        {availableUsages.map((usage) => (
          <div className="flex items-center" key={usage}>
            <input
              id={`${idBase}:${usage}`}
              type={inputType}
              onChange={onChange}
              value={usage}
              checked={isChecked(usage)}
              className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
            />
            <label
              htmlFor={`${idBase}:${usage}`}
              className="ml-3 block text-sm font-medium leading-6 text-gray-900"
            >
              {certUsageNames[usage]}
            </label>
          </div>
        ))}
      </div>
    </fieldset>
  );
}
