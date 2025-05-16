import { NextRequest, NextResponse } from "next/server"
import jwt from "jsonwebtoken"

const JWT_SECRET = process.env.JWT_SECRET!
const BACKEND_URL = process.env.BACKEND_URL!
const TOKEN_COOKIE_NAME = "token"
const COOKIE_MAX_AGE = 7 * 24 * 60 * 60 // 7 days in seconds

interface BackendSignUpResponse {
    id: string;
    email: string;
    name: string;
    role: string;
}

export async function POST(request: NextRequest) {
    try {
        const { email, password, name } = await request.json()

        // Base validation
        if (!email || !password || !name) {
            return NextResponse.json(
                { error: "Email, password, and name are required" },
                { status: 400 }
            )
        }

        const response = await fetch(`${BACKEND_URL}/v1/auth/signup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password, name }),
        })

        if (!response.ok) {
            const { error } = await response.json()
            return NextResponse.json(
                { error: error || "Registration failed" },
                { status: 400 }
            )
        }

        const userData: BackendSignUpResponse = await response.json()

        // Issue own JWT with complete user data
        const token = jwt.sign(
            {
                sub: userData.id,
                email: userData.email,
                name: userData.name,
                role: userData.role
            },
            JWT_SECRET,
            { expiresIn: "1d" }
        )

        // Set the HTTP-only cookie
        const result = NextResponse.json({ success: true })
        result.headers.set(
            "Set-Cookie",
            `${TOKEN_COOKIE_NAME}=${token}; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=${COOKIE_MAX_AGE};`
        )
        return result
    } catch (error) {
        console.error("Sign-up error:", error)
        return NextResponse.json(
            { error: "Internal server error" }, 
            { status: 500 }
        )
    }
}