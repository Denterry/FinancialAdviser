import { NextResponse } from "next/server"

const TOKEN_COOKIE_NAME = "token"

export async function POST() {
    try {
        const result = NextResponse.json({ success: true })
        result.headers.set(
            "Set-Cookie",
            `${TOKEN_COOKIE_NAME}=; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=0;`
        )
        return result
    } catch (error) {
        console.error("Sign-out error:", error)
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        )
    }
}