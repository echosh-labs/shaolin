import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";

const BACKEND_API_BASE =
  process.env.BACKEND_API_URL || "http://localhost:8081/api";

async function proxyRequest(
  req: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  try {
    const { path } = await params;
    const targetPath = path.join("/");
    const searchParams = req.nextUrl.search;
    const targetUrl = `${BACKEND_API_BASE}/${targetPath}${searchParams}`;

    // Extract auth token from httpOnly cookie
    const cookieStore = await cookies();
    const token = cookieStore.get("auth_token")?.value;

    const forwardHeaders: Record<string, string> = {
      "Content-Type": req.headers.get("Content-Type") || "application/json",
    };

    if (token) {
      forwardHeaders["Authorization"] = `Bearer ${token}`;
    }

    const fetchOptions: RequestInit = {
      method: req.method,
      headers: forwardHeaders,
      cache: "no-store",
    };

    if (req.method !== "GET" && req.method !== "HEAD") {
      const bodyText = await req.text();
      if (bodyText) {
        fetchOptions.body = bodyText;
      }
    }

    const response = await fetch(targetUrl, fetchOptions);

    if (response.status === 204) {
      return new NextResponse(null, { status: 204 });
    }

    const contentType = response.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      const data = await response.json();
      return NextResponse.json(data, { status: response.status });
    } else {
      const textData = await response.text();
      return new NextResponse(textData, {
        status: response.status,
        headers: { "Content-Type": contentType || "text/plain" },
      });
    }
  } catch (err: any) {
    return NextResponse.json(
      { error: err.message || "Failed to reach Go backend API" },
      { status: 502 }
    );
  }
}

export async function GET(
  req: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(req, context);
}

export async function POST(
  req: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(req, context);
}

export async function PUT(
  req: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(req, context);
}

export async function DELETE(
  req: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(req, context);
}
