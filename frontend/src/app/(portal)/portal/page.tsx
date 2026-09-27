"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import {
  Flame,
  Coins,
  Award,
  CalendarCheck2,
  TrendingUp,
  ShoppingBag,
  Clock,
  ArrowRight,
  Sparkles,
  ShieldCheck,
} from "lucide-react";
import { StudentOverallStats, StudentWeeklyStats } from "@/types";
import { apiFetch } from "@/lib/api";

export default function StudentPortalDashboard() {
  const { user } = useAuth();
  const [tokensRemaining, setTokensRemaining] = useState<number>(14);
  const [weeklyPoints, setWeeklyPoints] = useState<number>(45);
  const [tierTitle, setTierTitle] = useState<string>("Iron Disciple (铁弟子)");
  const [cumulativePoints, setCumulativePoints] = useState<number>(380);
  const [nextClassInfo, setNextClassInfo] = useState<{ title: string; time: string }>({
    title: "Traditional Kung Fu (Foundations)",
    time: "18:00 - 19:30",
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadDashboardData() {
      try {
        const [tokenData, weeklyData, occurrences] = await Promise.all([
          apiFetch<any>("/tokens/balance").catch(() => null),
          apiFetch<any[]>("/users/leaderboard/weekly").catch(() => []),
          apiFetch<any[]>("/occurrences").catch(() => []),
        ]);

        if (tokenData && typeof tokenData.tokens_remaining === "number") {
          setTokensRemaining(tokenData.tokens_remaining);
        }

        if (weeklyData && weeklyData.length > 0 && user) {
          const userStat = weeklyData.find((w: any) => w.user_id === user.id);
          if (userStat) {
            setWeeklyPoints(userStat.weekly_points || 45);
          }
        }

        if (occurrences && occurrences.length > 0) {
          const firstOcc = occurrences[0];
          setNextClassInfo({
            title: firstOcc.notes || "Martial Arts Training",
            time: `${firstOcc.date} (Scheduled)`,
          });
        }
      } catch (err) {
        console.error("Dashboard data load error:", err);
      } finally {
        setLoading(false);
      }
    }
    loadDashboardData();
  }, [user]);

  return (
    <div className="space-y-8">
      {/* Student Welcome Header */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-red-950/80 via-slate-900 to-slate-900 border border-red-900/40 p-6 sm:p-8 shadow-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 relative z-10">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <Badge variant="primary">Disciple Portal</Badge>
              <Badge variant="beltWhite">{user?.current_rank || "Level 1 White Belt"}</Badge>
            </div>
            <h1 className="text-2xl sm:text-3xl font-black text-slate-100 tracking-tight">
              Welcome back, {user?.first_name || "Disciple"}!
            </h1>
            <p className="text-xs sm:text-sm text-slate-300">
              Toronto Downtown Dojo • Summer 2026 Term Active
            </p>
          </div>

          <div className="flex items-center gap-3">
            <Link href="/portal/bookings">
              <Button variant="primary" size="md">
                <CalendarCheck2 className="h-4 w-4" />
                Book Next Class
              </Button>
            </Link>
            <Link href="/portal/tokens">
              <Button variant="accent" size="md">
                <Coins className="h-4 w-4" />
                Buy Tokens
              </Button>
            </Link>
          </div>
        </div>
      </div>

      {/* Gamification & Token Metrics Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Tokens Card */}
        <Card glow="gold" className="p-5">
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs text-slate-400 font-semibold uppercase tracking-wider">
              Token Balance
            </span>
            <div className="p-2 rounded-lg bg-amber-950/80 text-amber-400 border border-amber-800/40">
              <Coins className="h-4 w-4" />
            </div>
          </div>
          <div className="text-3xl font-black text-slate-100 mb-1">
            {tokensRemaining} <span className="text-xs text-slate-400 font-normal">Tokens</span>
          </div>
          <p className="text-[11px] text-amber-400 flex items-center gap-1">
            <Sparkles className="h-3 w-3" />
            Valid across all Kung Fu & Tai Chi
          </p>
        </Card>

        {/* Weekly Points & Consistency */}
        <Card glow="red" className="p-5">
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs text-slate-400 font-semibold uppercase tracking-wider">
              Weekly Points
            </span>
            <div className="p-2 rounded-lg bg-red-950/80 text-red-400 border border-red-800/40">
              <TrendingUp className="h-4 w-4" />
            </div>
          </div>
          <div className="text-3xl font-black text-slate-100 mb-1">
            {weeklyPoints} <span className="text-xs text-slate-400 font-normal">Pts</span>
          </div>
          <div className="w-full bg-slate-800 rounded-full h-1.5 mt-2 overflow-hidden">
            <div
              className="bg-red-500 h-full rounded-full transition-all duration-500"
              style={{ width: `${Math.min(100, (weeklyPoints / 60) * 100)}%` }}
            />
          </div>
          <p className="text-[11px] text-slate-400 mt-1">Consistency tier active</p>
        </Card>

        {/* Overall Rank Tier */}
        <Card className="p-5 border-slate-800">
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs text-slate-400 font-semibold uppercase tracking-wider">
              Discipleship Tier
            </span>
            <div className="p-2 rounded-lg bg-slate-800 text-slate-300">
              <Award className="h-4 w-4" />
            </div>
          </div>
          <div className="text-base font-bold text-amber-400 mb-1">{tierTitle}</div>
          <p className="text-[11px] text-slate-400 font-mono">{cumulativePoints} Overall Cumulative Pts</p>
        </Card>

        {/* Next Scheduled Class */}
        <Card className="p-5 border-slate-800">
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs text-slate-400 font-semibold uppercase tracking-wider">
              Next Enrolled Slot
            </span>
            <div className="p-2 rounded-lg bg-slate-800 text-slate-300">
              <Clock className="h-4 w-4" />
            </div>
          </div>
          <div className="text-sm font-bold text-slate-100 mb-0.5 truncate">{nextClassInfo.title}</div>
          <p className="text-[11px] text-emerald-400 font-medium">{nextClassInfo.time}</p>
        </Card>
      </div>

      {/* Quick Navigation Action Hub */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="hover:border-slate-700 transition-colors">
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <CalendarCheck2 className="h-5 w-5 text-red-500" />
              Class Booking Center
            </CardTitle>
            <CardDescription>
              Reserve your spot in daily Kung Fu, Tai Chi, and Qigong occurrences or manage recurring auto-reservations.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/portal/bookings">
              <Button variant="secondary" size="sm" className="w-full">
                View Class Calendar
                <ArrowRight className="h-4 w-4 ml-1" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="hover:border-slate-700 transition-colors">
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Award className="h-5 w-5 text-amber-500" />
              Grading Syllabus & Stances
            </CardTitle>
            <CardDescription>
              Track your fundamental stances (Ma Bu, Gong Bu, Pu Bu, Zenkutsu-dachi) benchmark timings, exam rubrics, and feedback.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/portal/grading">
              <Button variant="secondary" size="sm" className="w-full">
                View Grading Syllabus
                <ArrowRight className="h-4 w-4 ml-1" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="hover:border-slate-700 transition-colors">
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <ShoppingBag className="h-5 w-5 text-emerald-500" />
              Online Store & Tokens
            </CardTitle>
            <CardDescription>
              Top up class token bundles, order martial arts Feiyue footwear, uniforms, and weapon gear.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="/portal/tokens">
              <Button variant="secondary" size="sm" className="w-full">
                Manage Tokens & Store
                <ArrowRight className="h-4 w-4 ml-1" />
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
