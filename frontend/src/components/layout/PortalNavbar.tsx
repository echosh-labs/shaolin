"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import {
  Flame,
  LayoutDashboard,
  CalendarCheck2,
  Coins,
  Award,
  ShoppingBag,
  MessagesSquare,
  LogOut,
  ShieldCheck,
} from "lucide-react";
import { cn } from "@/lib/utils";

export function PortalNavbar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();

  const navItems = [
    { href: "/portal", label: "Dashboard", icon: LayoutDashboard },
    { href: "/portal/bookings", label: "Book Classes", icon: CalendarCheck2 },
    { href: "/portal/tokens", label: "Buy Tokens", icon: Coins },
    { href: "/portal/grading", label: "Grading & Progress", icon: Award },
    { href: "/store", label: "Dojo Store", icon: ShoppingBag },
    { href: "/portal/community", label: "Community", icon: MessagesSquare },
  ];

  return (
    <header className="sticky top-0 z-40 w-full border-b border-slate-800 bg-slate-950/95 backdrop-blur-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Brand Logo & Portal Tag */}
          <div className="flex items-center gap-4">
            <Link href="/" className="flex items-center gap-2.5 group">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-red-600 to-amber-600 text-white shadow-md shadow-red-950/50">
                <Flame className="h-5 w-5" />
              </div>
              <span className="font-extrabold tracking-wider text-slate-100 text-sm uppercase hidden sm:inline">
                Martial Arts Portal
              </span>
            </Link>

            {user?.role === "admin" || user?.role === "instructor" ? (
              <Link href="/admin">
                <Badge variant="gold" size="sm" className="hidden md:inline-flex cursor-pointer">
                  <ShieldCheck className="h-3 w-3" />
                  Staff View
                </Badge>
              </Link>
            ) : null}
          </div>

          {/* Desktop Nav Links */}
          <nav className="hidden lg:flex items-center gap-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors",
                    isActive
                      ? "bg-red-950/60 text-red-300 border border-red-800/40"
                      : "text-slate-400 hover:text-slate-100 hover:bg-slate-900"
                  )}
                >
                  <Icon className={cn("h-3.5 w-3.5", isActive ? "text-red-400" : "text-slate-500")} />
                  {item.label}
                </Link>
              );
            })}
          </nav>

          {/* User Profile Details & Sign Out */}
          <div className="flex items-center gap-3">
            {user && (
              <div className="flex items-center gap-2.5">
                <div className="text-right hidden sm:block">
                  <p className="text-xs font-semibold text-slate-200">
                    {user.first_name} {user.last_name}
                  </p>
                  <p className="text-[10px] text-slate-400 font-mono">{user.email}</p>
                </div>

                <Badge variant="beltWhite" size="sm">
                  {user.current_rank || "White Belt"}
                </Badge>
              </div>
            )}

            <Button
              variant="ghost"
              size="icon"
              onClick={() => logout()}
              title="Sign Out"
              className="text-slate-400 hover:text-red-400"
            >
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
}
