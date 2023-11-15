import classNames from "classnames";
import { PropsWithChildren } from "react";

export function Card({
  title,
  className,
  children,
}: PropsWithChildren<{ title?: React.ReactNode; className?: string }>) {
  return (
    <div className={"overflow-hidden rounded-lg bg-white shadow"}>
      {title && (
        <div className="border-b border-gray-200 bg-white px-4 py-5 sm:px-6">
          <h3 className="text-base font-semibold leading-6 text-gray-900">
            {title}
          </h3>
        </div>
      )}
      <div className={classNames("px-4 py-5 sm:p-6", className)}>
        {children}
      </div>
    </div>
  );
}
