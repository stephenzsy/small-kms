import { getMsAuth } from "@/utils/aadAuthUtils";
import { DefaultAzureCredential } from "@azure/identity";
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
      headers: await auth.getAuthHeaders(),
    }
  );
  return new NextResponse(resp.body, {
    status: resp.status,
    headers: resp.headers,
  });
}
