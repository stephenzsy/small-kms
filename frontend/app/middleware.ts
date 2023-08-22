import { NextRequest, NextResponse } from "next/server";

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  console.log(request.url);
  const modifiedPath = request.nextUrl.pathname.substring(4);
  const forwardUrl = new URL(
    `${process.env.BACKEND_URL_BASE}/${modifiedPath}}`
  );
  forwardUrl.search = request.nextUrl.search;
  console.log(forwardUrl);
  return NextResponse.rewrite(forwardUrl);
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ["/api/.*"],
};
