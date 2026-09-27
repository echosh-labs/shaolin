import { NextRequest, NextResponse } from "next/server";
import { backendFetch } from "@/lib/api";
import { cookies } from "next/headers";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const { email, password, first_name, last_name, date_of_birth } = body;

    if (!email || !password || !first_name || !last_name || !date_of_birth) {
      return NextResponse.json(
        { error: "All fields are required" },
        { status: 400 }
      );
    }

    // 1. Create account
    await backendFetch("/auth/register", {
      method: "POST",
      body: JSON.stringify({
        email,
        password,
        first_name,
        last_name,
        date_of_birth,
      }),
    });

    // 2. Auto login
    const authData = await backendFetch<{ token: string; user: any }>(
      "/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }
    );

    // 3. Set cookie
    const cookieStore = await cookies();
    cookieStore.set("auth_token", authData.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24 * 7,
    });

    return NextResponse.json({
      user: authData.user,
      message: "Registration and login successful",
    });
  } catch (err: any) {
    return NextResponse.json(
      { error: err.message || "Failed to register" },
      { status: err.status || 500 }
    );
  }
}
