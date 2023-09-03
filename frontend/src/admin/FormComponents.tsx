import { useId } from "react";
import { UseFormRegister } from "react-hook-form";

export interface InputFieldProps<T extends {}> {
  labelContent: React.ReactNode;
  defaultValue?: string;
  required?: boolean;
  register: UseFormRegister<T>;
  inputKey: keyof T;
  placeholder?: string;
  type?: "text" | "number";
}

export function InputField<T extends {}>({
  labelContent,
  defaultValue = "",
  required = false,
  register,
  inputKey,
  placeholder,
  type = "text",
}: InputFieldProps<T>) {
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
            type={type}
            defaultValue={defaultValue}
            id={id}
            {...register(inputKey as any, { required })}
            className="block flex-1 border-0 bg-transparent py-1.5 px-em text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
            placeholder={placeholder}
          />
        </div>
      </div>
    </div>
  );
}
