import { getMsAuth } from "@/utils/aadAuthUtils";
import { NextRequest, NextResponse } from "next/server";

// create certificate
export async function POST(request: NextRequest) {
  const auth = getMsAuth();
  const forward = await request.blob();
  const resp = await fetch(
    `${process.env.BACKEND_URL_BASE}/admin/certificate`,
    {
      method: "POST",
      body: forward,
      headers: await auth.getAuthHeaders(),
    }
  );
  return NextResponse.json(await resp.json(), {
    status: resp.status,
  });
}
