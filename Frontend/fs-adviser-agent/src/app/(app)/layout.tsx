import TopNav from "../../../components/ui/TopNav";
import "../globals.css";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-gray-950 min-h-screen flex flex-col">
        <TopNav />
        <main className="flex-1 p-8 overflow-auto">{children}</main>
      </body>
    </html>
  );
} 