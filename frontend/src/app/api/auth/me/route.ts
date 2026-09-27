import { NextResponse } from "next/server";
import { backendFetch } from "@/lib/api";
import { cookies } from "next/headers";

export async function GET() {
  const cookieStore = await cookies();
  const token = cookieStore.get("auth_token")?.value;

  if (!token) {
    return NextResponse.json({ user: null }, { status: 200 });
  }

  try {
    const user = await backendFetch("/auth/me", {
      token,
    });
    return NextResponse.json({ user });
  } catch {
    // If backend reports token expired or invalid, clear cookie
    cookieStore.delete("auth_token");
    return NextResponse.json({ user: null }, { status: 200 });
  }
}
