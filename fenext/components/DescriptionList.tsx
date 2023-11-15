import { PropsWithChildren } from "react";

export function DescriptionList({ children }: PropsWithChildren<{}>) {
  return <dl className="divide-y divide-gray-100">{children}</dl>;
}

export function DescriptionListItem({
  term,
  children,
}: PropsWithChildren<{ term: React.ReactNode }>) {
  return (
    <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
      <dt className="text-sm font-medium text-gray-900">{term}</dt>
      <dd className="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
        {children}
      </dd>
    </div>
  );
}
