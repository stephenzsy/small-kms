import type { Result } from "ahooks/lib/useRequest/src/types";
import { createContext } from "react";

export type UseRequestContextValue<
  TData,
  TParams extends unknown[] = []
> = Result<TData, TParams>;

class NotContextErrorResult<TData, TParams extends unknown[] | []>
  implements Result<TData, TParams[]>
{
  constructor() {}
  public readonly loading = false;
  public get data(): TData | undefined {
    return undefined;
  }
  public readonly error = new Error("No service in context provider");
  public readonly params = [] as [];
  public cancel() {}
  public refresh() {
    throw this.error;
  }
  public refreshAsync() {
    return Promise.reject(this.error);
  }
  public run() {
    throw this.error;
  }
  public runAsync() {
    return Promise.reject(this.error);
  }
  public mutate() {
    throw this.error;
  }
}

export function createUseRequestContext<
  TData,
  TParams extends unknown[] | [] = []
>() {
  return createContext<Result<TData, TParams>>(
    new NotContextErrorResult<TData, TParams>()
  );
}
