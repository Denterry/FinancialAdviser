import { NextRequest, NextResponse } from "next/server"
import jwt from "jsonwebtoken"

const JWT_SECRET = process.env.JWT_SECRET!
const BACKEND_URL = process.env.BACKEND_URL!
const TOKEN_COOKIE_NAME = "token"
const COOKIE_MAX_AGE = 7 * 24 * 60 * 60 // 7 days in seconds

interface BackendAuthResponse {
    id: string;
    name: string;
    email: string;
    role: string;
}

export async function POST(request: NextRequest) {
    try {
        const { email, password } = await request.json()

        // Forward the request to the backend
        const response = await fetch(`${BACKEND_URL}/v1/auth/signin`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password }),
        })

        if (!response.ok) {
            const { error } = await response.json()
            return NextResponse.json(
                { error: error || "Invalid credentials" }, 
                { status: 401 }
            )
        }

        const userData: BackendAuthResponse = await response.json()

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
        console.error("Sign-in error:", error)
        return NextResponse.json(
            { error: "Internal server error" }, 
            { status: 500 }
        )
    }
}