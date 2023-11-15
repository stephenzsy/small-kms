import classNames from "classnames";
import { PropsWithChildren } from "react";

export function PageHeader({ title }: { title: React.ReactNode }) {
  return (
    <header>
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold leading-tight tracking-tight text-gray-900">
          {title}
        </h1>
      </div>
    </header>
  );
}

export function PageContent({
  children,
  className,
}: PropsWithChildren<{
  className?: string;
}>) {
  return (
    <main
      className={classNames("mx-auto max-w-7xl sm:px-6 lg:px-8", className)}
    >
      {children}
    </main>
  );
}
