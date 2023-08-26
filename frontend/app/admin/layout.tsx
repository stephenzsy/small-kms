import { getMsAuth } from "@/utils/aadAuthUtils";
import { NextResponse } from "next/server";
import type { PropsWithChildren } from "react";

export default function AdminLayout(props: PropsWithChildren<{}>) {
  const auth = getMsAuth();
  if (!auth.isAdmin) {
    return (
      <div className="max-w-7xl p-6 mx-auto">
        You must have administrator access
      </div>
    );
  }
  return <div className="max-w-7xl p-6 mx-auto">{props.children}</div>;
}
