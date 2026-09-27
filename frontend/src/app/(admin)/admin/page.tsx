"use client";

import React from "react";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import {
  ShieldAlert,
  ClipboardCheck,
  CalendarDays,
  GraduationCap,
  PackageCheck,
  UserCheck,
  ArrowRight,
  Sparkles,
} from "lucide-react";

export default function AdminOperationsHub() {
  const { user } = useAuth();

  return (
    <div className="space-y-8">
      {/* Staff Operations Banner */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-amber-950/80 via-slate-900 to-slate-900 border border-amber-800/40 p-6 sm:p-8 shadow-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 relative z-10">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <Badge variant="gold">Staff & Instructors Only</Badge>
              <Badge variant="primary">Admin Role: {user?.role || "Staff"}</Badge>
            </div>
            <h1 className="text-2xl sm:text-3xl font-black text-slate-100 tracking-tight">
              School Operations Center
            </h1>
            <p className="text-xs sm:text-sm text-slate-300">
              Manage class attendance, curriculum stances, grading evaluations, and student ledgers.
            </p>
          </div>

          <Link href="/admin/attendance">
            <Button variant="accent" size="md">
              <ClipboardCheck className="h-4 w-4" />
              Launch Attendance Scanner
            </Button>
          </Link>
        </div>
      </div>

      {/* Admin Modules Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Attendance Sheet */}
        <Card glow="red" className="flex flex-col justify-between">
          <CardHeader>
            <div className="flex justify-between items-start mb-2">
              <div className="p-2 rounded-lg bg-red-950/80 text-red-400 border border-red-800/40">
                <ClipboardCheck className="h-5 w-5" />
              </div>
              <Badge variant="primary">High Priority</Badge>
            </div>
            <CardTitle className="text-lg">Class Attendance Check-In</CardTitle>
            <CardDescription>
              Select today&apos;s occurrence, toggle student attendance (in-person or livestream), award gamified points, and trigger rank promotions.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/admin/attendance">
              <Button variant="primary" size="sm" className="w-full">
                Open Attendance Sheet
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Schedule & Breaks */}
        <Card glow="gold" className="flex flex-col justify-between">
          <CardHeader>
            <div className="flex justify-between items-start mb-2">
              <div className="p-2 rounded-lg bg-amber-950/80 text-amber-400 border border-amber-800/40">
                <CalendarDays className="h-5 w-5" />
              </div>
              <Badge variant="gold">Calendar</Badge>
            </div>
            <CardTitle className="text-lg">Schedules & Holiday Breaks</CardTitle>
            <CardDescription>
              Add term break blackout periods (e.g. Summer Mid-Term Break), adjust class occurrence Zoom links, and update capacity limits.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/admin/schedules">
              <Button variant="accent" size="sm" className="w-full">
                Manage Class Schedules
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Grading Exam Evaluator */}
        <Card className="flex flex-col justify-between border-slate-800">
          <CardHeader>
            <div className="flex justify-between items-start mb-2">
              <div className="p-2 rounded-lg bg-slate-800 text-slate-300">
                <GraduationCap className="h-5 w-5" />
              </div>
              <Badge variant="default">Syllabus</Badge>
            </div>
            <CardTitle className="text-lg">Grading Exam Evaluator</CardTitle>
            <CardDescription>
              Record student belt exam evaluations, score stance precision, and award official Academy certification belts.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/admin/grading">
              <Button variant="secondary" size="sm" className="w-full">
                Evaluate Student Exams
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Store Order Fulfillment */}
        <Card className="flex flex-col justify-between border-slate-800">
          <CardHeader>
            <div className="flex justify-between items-start mb-2">
              <div className="p-2 rounded-lg bg-slate-800 text-slate-300">
                <PackageCheck className="h-5 w-5" />
              </div>
              <Badge variant="default">E-Commerce</Badge>
            </div>
            <CardTitle className="text-lg">Store Orders & Inventory</CardTitle>
            <CardDescription>
              Process store checkouts, verify e-transfer payments, and mark uniform silks and weapons ready for dojo pickup.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/admin/store">
              <Button variant="secondary" size="sm" className="w-full">
                Process Orders
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Student Records */}
        <Card className="flex flex-col justify-between border-slate-800">
          <CardHeader>
            <div className="flex justify-between items-start mb-2">
              <div className="p-2 rounded-lg bg-slate-800 text-slate-300">
                <UserCheck className="h-5 w-5" />
              </div>
              <Badge variant="default">Members</Badge>
            </div>
            <CardTitle className="text-lg">Student Directory</CardTitle>
            <CardDescription>
              Look up student account records, assign annual membership exemptions, and review token balance audit ledgers.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/admin/students">
              <Button variant="secondary" size="sm" className="w-full">
                View Student Roster
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
