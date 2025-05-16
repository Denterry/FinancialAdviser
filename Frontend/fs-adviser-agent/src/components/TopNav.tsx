"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { name: "Chat", path: "/chat" },
  { name: "Agent", path: "/agent" },
  { name: "Articles", path: "/articles" },
];

export default function TopNav() {
  const pathname = usePathname();
  return (
    <header className="p-4 bg-gray-800 text-white flex justify-between items-center">
      <h1 className="text-xl font-semibold">FS ChatBot</h1>
      <nav className="space-x-4">
        {navItems.map((item) => (
          <Link
            key={item.path}
            href={item.path}
            className={`px-3 py-2 rounded ${
              pathname === item.path ? "bg-blue-700" : "hover:bg-gray-700"
            }`}
          >
            {item.name}
          </Link>
        ))}
      </nav>
    </header>
  );
} 