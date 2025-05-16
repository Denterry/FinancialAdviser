// File: src/app/page.tsx
import Link from "next/link";
import Image from "next/image";
import { Metadata } from "next";
import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth";

export const metadata: Metadata = {
  title: "FS Adviser Agent",
  description: "AI-driven financial advisory platform — start your journey in investing here.",
};

export default async function HomePage() {
  const user = await getCurrentUser();

  return (
    <div className="grid grid-rows-[1fr_auto] min-h-screen p-8 relative">
      {/* Top Right Buttons */}
      <div className="absolute top-8 right-8 space-x-4">
        {!user && (
          <>
            <Link href="/signin">
              <Button variant="outline" size="lg">
                Sign In
              </Button>
            </Link>
            <Link href="/signup">
              <Button variant="outline" size="lg">
                Sign Up
              </Button>
            </Link>
          </>
        )}
      </div>

      {/* Main Content */}
      <main className="flex flex-col items-center justify-center gap-8">
        <div className="text-center space-y-6 max-w-2xl">
          <h1 className="text-6xl font-extrabold tracking-tight bg-gradient-to-r from-primary to-primary-foreground bg-clip-text">
            FS Adviser Agent
          </h1>
          <p className="text-xl text-muted-foreground">
            Leverage AI-powered sentiment analysis and reinforcement learning
            to make smarter financial and real estate investments.
          </p>
          {user && (
            <Link href="/chat">
              <Button size="lg" className="mt-4">
                Go to Chat
              </Button>
            </Link>
          )}
        </div>
      </main>

      {/* Footer with Social Links and Brand Logo */}
      <footer className="flex justify-between items-center py-8">
        {/* Brand Logo */}
        <div className="flex items-center gap-2 cursor-pointer group">
          <div className="w-10 h-10 rounded-full bg-white flex items-center justify-center shadow-md overflow-hidden">
            <Image
              src="/fs-logo.svg"
              alt="FS Adviser Logo"
              width={32}
              height={32}
              className="object-contain"
            />
          </div>
          <span className="text-sm text-muted-foreground group-hover:text-foreground transition-colors">
            About Project
          </span>
        </div>

        {/* Social Links as Text */}
        <div className="flex gap-10 items-center text-lg font-semibold justify-center w-full absolute bottom-8 -ml-8">
          <a
            href="https://t.me/alreadygoat"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-blue-600 transition-colors"
          >
            Telegram
          </a>
          <a
            href="https://github.com/Denterry/FinancialAdviser"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-blue-600 transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://linkedin.com/in/denis-todorov-b5a308243/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-blue-600 transition-colors"
          >
            LinkedIn
          </a>
        </div>
      </footer>
    </div>
  );
}
