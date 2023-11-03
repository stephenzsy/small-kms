import classNames from "classnames";
import { useMemo } from "react";

export type JsonDataDisplayProps<T> = {
  data: T | undefined;
  loading?: boolean;
  toJson?: (data?: T) => any;
  className?: string;
};

export function JsonDataDisplay<T>({
  loading,
  data,
  toJson,
  className,
}: JsonDataDisplayProps<T>) {
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
        "bg-neutral-100 ring-1 ring-neutral-500 px-4 overflow-auto max-h-full",
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
