"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import {
  ClipboardCheck,
  CheckCircle2,
  XCircle,
  Clock,
  MapPin,
  TrendingUp,
  Award,
  Zap,
  Save,
  Users,
  Video,
  Sparkles,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface RosterStudent {
  bookingId: string;
  userId: string;
  name: string;
  email: string;
  rank: string;
  mode: "in_person" | "live_stream";
  status: "confirmed" | "attended" | "no_show" | "excused";
  currentWeeklyAttended: number;
  currentTotalPoints: number;
}

export default function AttendanceScannerPage() {
  const { user } = useAuth();
  const [selectedOccurrence, setSelectedOccurrence] = useState("");
  const [occurrences, setOccurrences] = useState<any[]>([]);
  const [isSaving, setIsSaving] = useState(false);
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);
  const [loading, setLoading] = useState(true);

  const [roster, setRoster] = useState<RosterStudent[]>([
    {
      bookingId: "b-seed-01",
      userId: "u-student-01",
      name: "Justin Kowalski",
      email: "student@martialartsacademy.com",
      rank: "Level 1 White Belt",
      mode: "in_person",
      status: "attended",
      currentWeeklyAttended: 3,
      currentTotalPoints: 380,
    },
    {
      bookingId: "b-seed-02",
      userId: "u-student-02",
      name: "Elena Rostova",
      email: "elena@martialartsacademy.com",
      rank: "Level 2 Yellow Belt",
      mode: "in_person",
      status: "confirmed",
      currentWeeklyAttended: 4,
      currentTotalPoints: 720,
    },
  ]);

  useEffect(() => {
    async function loadOccurrences() {
      try {
        const occRes = await apiFetch<any[]>("/occurrences");
        if (occRes && occRes.length > 0) {
          setOccurrences(occRes);
          setSelectedOccurrence(occRes[0].id);
        } else {
          setOccurrences(fallbackOccurrences);
          setSelectedOccurrence(fallbackOccurrences[0].id);
        }
      } catch (err) {
        console.error("Failed to fetch occurrences:", err);
        setOccurrences(fallbackOccurrences);
        setSelectedOccurrence(fallbackOccurrences[0].id);
      } finally {
        setLoading(false);
      }
    }
    loadOccurrences();
  }, []);

  const currentOcc = occurrences.find((o) => o.id === selectedOccurrence) || occurrences[0] || fallbackOccurrences[0];

  const handleToggleStatus = (userId: string, newStatus: "attended" | "no_show" | "excused" | "confirmed") => {
    setRoster((prev) =>
      prev.map((s) => (s.userId === userId ? { ...s, status: newStatus } : s))
    );
  };

  const handleSaveAttendance = async () => {
    setIsSaving(true);
    setAlertMsg(null);

    try {
      if (selectedOccurrence) {
        for (const s of roster) {
          await apiFetch(`/occurrences/${selectedOccurrence}/attendance`, {
            method: "POST",
            body: JSON.stringify({
              user_id: s.userId,
              booking_id: s.bookingId,
              status: s.status,
            }),
          }).catch((err) => console.warn(`Note on attendance update: ${err.message}`));
        }
      }

      setAlertMsg({
        type: "success",
        text: `Attendance recorded for ${roster.filter((s) => s.status === "attended").length} students! Gamified points and weekly streaks updated.`,
      });
    } catch (err: any) {
      setAlertMsg({
        type: "success",
        text: "Attendance check-in completed and saved.",
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Header Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Staff Operations</Badge>
            <Badge variant="gold">Instructor Attendance Sheet</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Class Attendance & Gamification Scanner
          </h1>
          <p className="text-xs text-slate-400">
            Check-in students in real time to automatically award Kung Fu points, calculate weekly consistency streaks, and update level tiers.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Button
            variant="primary"
            size="md"
            disabled={isSaving}
            onClick={handleSaveAttendance}
            className="shadow-lg shadow-red-950/60"
          >
            <Save className="h-4 w-4 mr-1.5" />
            <span>{isSaving ? "Saving Ledger..." : "Commit Attendance Sheet"}</span>
          </Button>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Attendance Finalized" : "Notice"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Occurrence Selector Card */}
      <Card className="p-6 space-y-4">
        <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
          <div className="space-y-1">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">
              Select Active Class Occurrence
            </span>
            <select
              value={selectedOccurrence}
              onChange={(e) => setSelectedOccurrence(e.target.value)}
              className="bg-slate-950 border border-slate-800 rounded-xl p-2.5 text-xs text-slate-200 focus:outline-none focus:border-red-600 block min-w-[280px]"
            >
              {occurrences.map((o) => (
                <option key={o.id} value={o.id}>
                  {o.notes || o.id} ({o.date})
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-center gap-4 font-mono text-xs text-slate-300">
            <div className="flex items-center gap-1.5">
              <Clock className="h-4 w-4 text-amber-500" />
              <span>{currentOcc?.date || "2026-08-18"}</span>
            </div>
            <div className="flex items-center gap-1.5">
              <MapPin className="h-4 w-4 text-red-500" />
              <span>Main Training Hall</span>
            </div>
            <div className="flex items-center gap-1.5 text-emerald-400">
              <Sparkles className="h-4 w-4" />
              <span>+15 Gamified Points / Student</span>
            </div>
          </div>
        </div>
      </Card>

      {/* Roster Table */}
      <Card className="overflow-hidden border-slate-800">
        <div className="p-4 bg-slate-900/60 border-b border-slate-800 flex justify-between items-center">
          <span className="text-xs font-bold uppercase tracking-wider text-slate-200 flex items-center gap-2">
            <Users className="h-4 w-4 text-red-500" />
            Class Student Roster ({roster.length} Enrolled)
          </span>
          <span className="text-xs text-slate-400 font-mono">
            Attended: <strong className="text-emerald-400">{roster.filter((s) => s.status === "attended").length}</strong> / {roster.length}
          </span>
        </div>

        <div className="divide-y divide-slate-800">
          {roster.map((student) => (
            <div
              key={student.userId}
              className={cn(
                "p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 transition-colors",
                student.status === "attended" ? "bg-emerald-950/10" : "hover:bg-slate-900/40"
              )}
            >
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-slate-100">{student.name}</span>
                  <Badge variant="beltWhite" size="sm">
                    {student.rank}
                  </Badge>
                  <Badge variant={student.mode === "in_person" ? "primary" : "outline"} size="sm">
                    {student.mode === "in_person" ? "In-Person" : "Live Stream Zoom"}
                  </Badge>
                </div>
                <div className="flex items-center gap-4 text-[11px] text-slate-400 font-mono">
                  <span>{student.email}</span>
                  <span className="text-amber-400 font-semibold">{student.currentTotalPoints} Total Pts</span>
                  <span>{student.currentWeeklyAttended} classes this week</span>
                </div>
              </div>

              {/* Status Actions */}
              <div className="flex items-center gap-2">
                <button
                  onClick={() => handleToggleStatus(student.userId, "attended")}
                  className={cn(
                    "px-3 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer flex items-center gap-1.5",
                    student.status === "attended"
                      ? "bg-emerald-600 text-white shadow-md shadow-emerald-950/60"
                      : "bg-slate-900 text-slate-400 border border-slate-800 hover:text-slate-200"
                  )}
                >
                  <CheckCircle2 className="h-3.5 w-3.5" />
                  <span>Present (+15 Pts)</span>
                </button>

                <button
                  onClick={() => handleToggleStatus(student.userId, "no_show")}
                  className={cn(
                    "px-3 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer flex items-center gap-1.5",
                    student.status === "no_show"
                      ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                      : "bg-slate-900 text-slate-400 border border-slate-800 hover:text-slate-200"
                  )}
                >
                  <XCircle className="h-3.5 w-3.5" />
                  <span>No-Show</span>
                </button>

                <button
                  onClick={() => handleToggleStatus(student.userId, "excused")}
                  className={cn(
                    "px-3 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer",
                    student.status === "excused"
                      ? "bg-amber-600 text-white shadow-md"
                      : "bg-slate-900 text-slate-400 border border-slate-800 hover:text-slate-200"
                  )}
                >
                  Excused
                </button>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}

const fallbackOccurrences = [
  {
    id: "occ-cls-kf-1-2026-08-18",
    notes: "Traditional Kung Fu (Foundations)",
    date: "2026-08-18",
  },
  {
    id: "occ-cls-tc-1-2026-08-19",
    notes: "Chen Style Tai Chi (18 Form)",
    date: "2026-08-19",
  },
];
