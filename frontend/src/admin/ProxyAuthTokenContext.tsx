import { useMemoizedFn, useSessionStorageState } from "ahooks";
import {
  PropsWithChildren,
  createContext,
  useEffect,
  useMemo,
  useState,
} from "react";

type ProxyAuthTokenContextValue = {
  getAccessToken(): string;
  hasToken: boolean;
  instanceId: string;
  setAccessToken: (token: string) => void;
};
export const ProxyAuthTokenContext = createContext<ProxyAuthTokenContextValue>({
  getAccessToken: () => "",
  hasToken: false,
  instanceId: "",
  setAccessToken: () => {},
});
export function ProxyAuthTokenContextProvider({
  instanceId,
  children,
}: PropsWithChildren<{ instanceId: string }>) {
  const [storedToken, setStoredToken] = useSessionStorageState<string>(
    `proxy-token-${instanceId}`
  );

  const [token, exp] = useMemo((): [string, number | undefined] => {
    if (!storedToken) {
      return ["", undefined];
    }
    const [, payload] = storedToken.split(".");
    const decoded = JSON.parse(atob(payload));
    return [storedToken, decoded.exp];
  }, [storedToken]);

  const [evictToken, setEvictToken] = useState<string>();

  useEffect(() => {
    setStoredToken((t) => {
      if (t === evictToken) {
        return "";
      }
      return t ?? "";
    });
  }, [evictToken]);

  const getAccessToken = useMemoizedFn((): string => {
    const now = Date.now() / 1000;
    if (!token) {
      throw new Error("Not authorized");
    }
    if (now > exp!) {
      setEvictToken(token);
      throw new Error("Not authorized");
    }
    return token;
  });

  return (
    <ProxyAuthTokenContext.Provider
      value={{
        getAccessToken,
        hasToken: !!token,
        instanceId,
        setAccessToken: setStoredToken,
      }}
    >
      {children}
    </ProxyAuthTokenContext.Provider>
  );
}
