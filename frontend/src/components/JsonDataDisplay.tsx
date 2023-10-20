import classNames from "classnames";
import { useMemo } from "react";

export function JsonDataDisplay<T>({
  loading,
  data,
  toJson,
  className,
}: {
  data: T | undefined;
  loading?: boolean;
  toJson?: (data?: T) => any;
  className?: string;
}) {
  const jsonText = useMemo(() => {
    if (!data) {
      return undefined;
    }
    if (toJson) {
      return JSON.stringify(toJson(data), undefined, 2);
    }
    return JSON.stringify(data, undefined, 2);
  }, [data, toJson]);
  return (
    <div
      className={classNames(
        "bg-neutral-100 ring-1 ring-neutral-500 px-4 overflow-auto",
        className
      )}
    >
      {loading ? (
        <div>Loading...</div>
      ) : data ? (
        <pre className="font-mono">{jsonText}</pre>
      ) : (
        <div>No data</div>
      )}
    </div>
  );
}
