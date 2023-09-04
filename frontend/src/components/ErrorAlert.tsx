import { XCircleIcon } from "@heroicons/react/24/solid";

export function ErrorAlert(props: React.PropsWithChildren<{ error: any }>) {
  return (
    <div className="rounded-md bg-red-50 p-4">
      <div className="flex">
        <div className="flex-shrink-0">
          <XCircleIcon className="h-5 w-5 text-red-400" aria-hidden="true" />
        </div>
        <div className="ml-3">
          <h3 className="text-sm font-medium text-yellow-800">Error</h3>
          <div className="mt-2 text-sm text-yellow-700">
            <p>
              {props.error.message}
              {props.children}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
