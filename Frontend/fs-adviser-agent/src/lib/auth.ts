import { cookies } from "next/headers"
import jwt, { JwtPayload } from "jsonwebtoken"

const JWT_SECRET = process.env.JWT_SECRET!
const TOKEN_COOKIE_NAME = "token"
const API_BASE_URL = "/api/auth"

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

class AuthError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'AuthError';
  }
}

// Read and verify the JWT from cookies; returns user or null
export async function getCurrentUser(): Promise<User | null> {
  const cookieStore = await cookies()
  const token = cookieStore.get(TOKEN_COOKIE_NAME)?.value
  if (!token) {
    return null
  }

  try {
    const payload = jwt.verify(token, JWT_SECRET) as JwtPayload
    return {
      id: payload.sub as string,
      email: payload.email as string,
      name: payload.name as string,
      role: payload.role as string,
    }
  } catch {
    return null
  }
}

// Call the login endpoint; on success this sets the HTTP-only cookie
export async function signIn(email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/signin`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const { error } = await response.json()
    throw new AuthError(error || "Failed to sign in")
  }
}

// Call the signup endpoint; on success sets cookie too
export async function signUp(email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const { error } = await response.json()
    throw new AuthError(error || "Failed to sign up")
  }
}

// Call the logout endpoint; on success clears the HTTP-only cookie
export async function signOut(): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/signout`, {
    method: "POST",
  })

  if (!response.ok) {
    const { error } = await response.json()
    throw new AuthError(error || "Failed to sign out")
  }
}
