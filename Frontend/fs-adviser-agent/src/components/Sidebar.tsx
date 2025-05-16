import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { name: "Chat", path: "/chat" },
  { name: "Agent", path: "/agent" },
  { name: "Articles", path: "/articles" },
];

export default function Sidebar() {
  const pathname = usePathname();
  return (
    <aside className="w-60 bg-gray-900 text-white flex flex-col h-full p-4">
      <div className="mb-8 text-2xl font-bold">FS ChatBot</div>
      <nav className="flex-1">
        {navItems.map((item) => (
          <Link
            key={item.path}
            href={item.path}
            className={`block px-4 py-2 rounded mb-2 ${
              pathname === item.path ? "bg-blue-700" : "hover:bg-gray-700"
            }`}
          >
            {item.name}
          </Link>
        ))}
      </nav>
    </aside>
  );
} 