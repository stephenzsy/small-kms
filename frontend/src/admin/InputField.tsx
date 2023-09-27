import classNames from "classnames";
import { useId } from "react";

export interface InputFieldProps {
  labelContent: React.ReactNode;
  defaultValue?: string;
  required?: boolean;
  placeholder?: string;
  type?: "text" | "number";
  value: string | number | undefined;
  onChange: (value: string) => void;
  inputMode?: "text" | "numeric";
  className?: string;
}

export function InputField({
  labelContent,
  required = false,
  placeholder,
  type = "text",
  value,
  onChange,
  inputMode,
  className,
}: InputFieldProps) {
  const id = useId();
  return (
    <div className={classNames("sm:col-span-4", className)}>
      <label
        htmlFor={id}
        className="block text-sm font-medium leading-6 text-neutral-900"
      >
        {labelContent}
        {required && <span className="ml-1 text-red-500">*</span>}
      </label>
      <div className="mt-2">
        <div className="flex rounded-md shadow-sm ring-1 ring-inset ring-neutral-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
          <input
            type={type}
            id={id}
            value={value}
            className="block flex-1 border-0 bg-transparent py-1.5 px-em text-neutral-900 placeholder:text-neutral-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder={placeholder}
            onChange={(e) => onChange(e.target.value)}
            inputMode={inputMode}
          />
        </div>
      </div>
    </div>
  );
}
