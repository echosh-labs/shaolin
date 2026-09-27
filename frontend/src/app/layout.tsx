import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/context/AuthContext";
import { CartProvider } from "@/context/CartContext";

export const metadata: Metadata = {
  title: "Martial Arts Academy | Kung Fu, Karate, Kobudo, Tai Chi & Qigong",
  description:
    "Authentic traditional martial arts training featuring Kung Fu, Karate, Kobudo weapons, Tai Chi, and Qigong. Structured curriculum, rank progression, and modern student portal.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="h-full bg-[#090d16] text-slate-100 antialiased dark">
      <body className="min-h-full flex flex-col bg-[#090d16] text-slate-100 selection:bg-red-600 selection:text-white">
        <AuthProvider>
          <CartProvider>{children}</CartProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
