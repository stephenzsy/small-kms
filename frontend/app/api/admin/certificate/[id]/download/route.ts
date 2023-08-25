import { getMsAuth } from "@/utils/aadAuthUtils";
import { NextRequest, NextResponse } from "next/server";

export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const auth = getMsAuth();
  const resp = await fetch(
    `${process.env.BACKEND_URL_BASE}/admin/certificate/${
      params.id
    }/download?format=${request.nextUrl.searchParams.get("format") || "pem"}`,
    {
      method: "GET",
      headers: {
        "X-Ms-Client-Principal-Name": auth.principalName!,
        "X-Ms-Client-Principal-Id": auth.principalId!,
        "X-Ms-Client-Roles": auth.isAdmin ? "App.Admin" : "",
      },
    }
  );
  return new NextResponse(resp.body, {
    status: resp.status,
    headers: resp.headers,
  });
}