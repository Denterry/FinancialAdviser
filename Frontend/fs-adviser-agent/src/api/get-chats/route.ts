import { cookies } from "next/headers";
import { NextResponse } from "next/server";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";
const TOKEN_COOKIE_NAME = "token";

export async function GET() {
  try {
    const cookieStore = cookies();
    const token = (await cookieStore).get(TOKEN_COOKIE_NAME)?.value;

    if (!token) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
    }

    const backendRes = await fetch(`${BACKEND_API_URL}/api/get-chats`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!backendRes.ok) {
      const error = await backendRes.text();
      return NextResponse.json({ error }, { status: backendRes.status });
    }

    const chats = await backendRes.json();
    return NextResponse.json(chats, { status: 200 });
  } catch (err: Error | unknown) {
    console.error("[get-chats] GET error:", err);
    return NextResponse.json({ error: "Internal Server Error" }, { status: 500 });
  }
}
