"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import {
  ShieldAlert,
  ClipboardCheck,
  CalendarDays,
  GraduationCap,
  PackageCheck,
  UserCheck,
  LogOut,
  ArrowLeft,
} from "lucide-react";
import { cn } from "@/lib/utils";

export function AdminNavbar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();

  const adminNav = [
    { href: "/admin", label: "Operations", icon: ShieldAlert },
    { href: "/admin/attendance", label: "Attendance Sheet", icon: ClipboardCheck },
    { href: "/admin/schedules", label: "Schedules & Breaks", icon: CalendarDays },
    { href: "/admin/grading", label: "Exam Evaluator", icon: GraduationCap },
    { href: "/admin/store", label: "Store Orders", icon: PackageCheck },
    { href: "/admin/students", label: "Students", icon: UserCheck },
  ];

  return (
    <header className="sticky top-0 z-40 w-full border-b border-amber-900/40 bg-slate-950/95 backdrop-blur-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-3">
            <Link href="/portal" className="text-slate-400 hover:text-slate-200 flex items-center gap-1.5 text-xs mr-2">
              <ArrowLeft className="h-3.5 w-3.5" />
              <span className="hidden sm:inline">Student View</span>
            </Link>

            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-600 text-white font-bold text-xs shadow-md shadow-amber-950/50">
                ADM
              </div>
              <Badge variant="gold" size="sm">
                Staff Operations
              </Badge>
            </div>
          </div>

          <nav className="hidden md:flex items-center gap-1">
            {adminNav.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors",
                    isActive
                      ? "bg-amber-950/70 text-amber-300 border border-amber-800/50"
                      : "text-slate-400 hover:text-slate-100 hover:bg-slate-900"
                  )}
                >
                  <Icon className={cn("h-3.5 w-3.5", isActive ? "text-amber-400" : "text-slate-500")} />
                  {item.label}
                </Link>
              );
            })}
          </nav>

          <div className="flex items-center gap-3">
            {user && (
              <span className="text-xs font-medium text-amber-400 hidden sm:inline">
                {user.first_name} ({user.role})
              </span>
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
