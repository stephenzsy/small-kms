import { useContext, useMemo, type PropsWithChildren } from "react";
import { AppAuthContext } from "../auth/AuthProvider";

export function AdminLayout(props: PropsWithChildren) {
  const { account } = useContext(AppAuthContext);

  const isAdmin = useMemo(
    () => !!account?.idTokenClaims?.roles?.includes("App.Admin"),
    [account]
  );

  if (!isAdmin) {
    return (
      <main className="grid min-h-full place-items-center px-6 py-24 sm:py-32 lg:px-8">
        <div className="text-center">
          <h1 className="mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-5xl">
            Admin permission required
          </h1>
        </div>
      </main>
    );
  }
  return props.children;
}
