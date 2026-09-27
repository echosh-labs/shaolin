import { NextRequest, NextResponse } from "next/server";
import { backendFetch } from "@/lib/api";
import { cookies } from "next/headers";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const { email, password } = body;

    if (!email || !password) {
      return NextResponse.json(
        { error: "Email and password are required" },
        { status: 400 }
      );
    }

    // Call Go API
    const authData = await backendFetch<{ token: string; user: any }>(
      "/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }
    );

    const cookieStore = await cookies();
    cookieStore.set("auth_token", authData.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24 * 7, // 7 days
    });

    return NextResponse.json({
      user: authData.user,
      message: "Login successful",
    });
  } catch (err: any) {
    return NextResponse.json(
      { error: err.message || "Failed to authenticate" },
      { status: err.status || 500 }
    );
  }
}
