"use client";

import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { Flame, Calendar, BookOpen, ShoppingBag, User, ShieldAlert } from "lucide-react";

export function Navbar() {
  const { user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-40 w-full border-b border-slate-800/80 bg-[#090d16]/90 backdrop-blur-md">
      <div className="max-w-7xl mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Brand Logo */}
        <Link href="/" className="flex items-center gap-3 group">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-red-600 to-amber-600 text-white shadow-lg shadow-red-950/60 group-hover:scale-105 transition-transform">
            <Flame className="h-6 w-6" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-black tracking-wider text-slate-100 text-lg uppercase">
                Martial Arts Academy
              </span>
              <span className="text-xs font-semibold px-1.5 py-0.5 rounded bg-red-950/80 text-red-400 border border-red-800/40">
                武道
              </span>
            </div>
            <p className="text-[10px] text-slate-400 font-medium tracking-wide">
              Kung Fu • Karate • Kobudo • Tai Chi • Qigong
            </p>
          </div>
        </Link>

        {/* Desktop Navigation Links */}
        <nav className="hidden md:flex items-center gap-6 text-sm font-medium text-slate-300">
          <Link
            href="/classes"
            className="hover:text-red-400 transition-colors flex items-center gap-1.5"
          >
            <Flame className="h-4 w-4 text-red-500" />
            Programs
          </Link>
          <Link
            href="/schedule"
            className="hover:text-amber-400 transition-colors flex items-center gap-1.5"
          >
            <Calendar className="h-4 w-4 text-amber-500" />
            Schedule
          </Link>
          <Link
            href="/tuition"
            className="hover:text-slate-100 transition-colors"
          >
            Tuition & Tokens
          </Link>
          <Link
            href="/store"
            className="hover:text-slate-100 transition-colors flex items-center gap-1.5"
          >
            <ShoppingBag className="h-4 w-4 text-slate-400" />
            Store
          </Link>
          <Link
            href="/guides"
            className="hover:text-slate-100 transition-colors flex items-center gap-1.5"
          >
            <BookOpen className="h-4 w-4 text-slate-400" />
            Guides
          </Link>
        </nav>

        {/* Auth CTA Actions */}
        <div className="flex items-center gap-3">
          {user ? (
            <div className="flex items-center gap-3">
              {user.role === "admin" || user.role === "instructor" ? (
                <Link href="/admin">
                  <Button variant="accent" size="sm" className="hidden sm:inline-flex">
                    <ShieldAlert className="h-4 w-4" />
                    Admin Hub
                  </Button>
                </Link>
              ) : null}
              <Link href="/portal">
                <Button variant="primary" size="sm">
                  <User className="h-4 w-4" />
                  Student Portal
                </Button>
              </Link>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => logout()}
                className="text-xs text-slate-400 hover:text-red-400"
              >
                Sign Out
              </Button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link href="/login">
                <Button variant="ghost" size="sm">
                  Sign In
                </Button>
              </Link>
              <Link href="/register">
                <Button variant="primary" size="sm">
                  Join School
                </Button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
