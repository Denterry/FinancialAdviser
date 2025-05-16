// File: src/api/chat/route.ts
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";
const TOKEN_COOKIE_NAME = "token";

export async function POST(req: NextRequest) {
  try {
    const cookieStore = cookies();
    const token = (await cookieStore).get(TOKEN_COOKIE_NAME)?.value;

    if (!token) {
      return new Response("Unauthorized", { status: 401 });
    }

    const body = await req.text();

    const backendRes = await fetch(`${BACKEND_API_URL}/api/chat`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body,
    });

    if (!backendRes.ok || !backendRes.body) {
      return new Response("Failed to stream chat completion", {
        status: backendRes.status,
      });
    }

    return new Response(backendRes.body, {
      status: 200,
      headers: {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
      },
    });
  } catch (err) {
    console.error("[chat] POST error:", err);
    return new Response("Internal Server Error", { status: 500 });
  }
}
