import { NextRequest } from "next/server";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();

    const backendResponse = await fetch(`${BACKEND_API_URL}/api/get-messages`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    if (!backendResponse.ok) {
      const error = await backendResponse.text();
      return new Response(JSON.stringify({ error }), { status: backendResponse.status });
    }

    const messages = await backendResponse.json();
    return new Response(JSON.stringify(messages), {
      headers: { "Content-Type": "application/json" },
      status: 200,
    });
  } catch (err: Error | unknown) {
    console.error("[get-messages] POST error:", err);
    return new Response(JSON.stringify({ error: "Internal Server Error" }), { status: 500 });
  }
}
