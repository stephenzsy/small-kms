import { getMsAuth } from "@/utils/aadAuthUtils";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const auth = getMsAuth();
  const forward = await request.blob();
  const resp = await fetch(
    `${process.env.BACKEND_URL_BASE}/admin/certificate`,
    {
      method: "POST",
      body: forward,
      headers: {
        "X-Ms-Client-Principal-Name": auth.principalName!,
        "X-Ms-Client-Principal-Id": auth.principalId!,
        "X-Ms-Client-Roles": auth.isAdmin ? "App.Admin" : "",
      },
    }
  );
  return NextResponse.json(await resp.json(), {
    status: resp.status,
  });
}
