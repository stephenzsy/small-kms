import { getMsAuth } from "@/utils/aadAuthUtils";
import { NextResponse } from "next/server";
import type { PropsWithChildren } from "react";

export default function AdminLayout(props: PropsWithChildren<{}>) {
  const auth = getMsAuth();
  if (!auth.isAdmin) {
    const e = new NextResponse(undefined, {
      status: 403,
    });
    return e;
  }
  return <div className="max-w-7xl p-6 mx-auto">{props.children}</div>;
}
