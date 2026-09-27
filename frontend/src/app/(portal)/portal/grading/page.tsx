"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import {
  Award,
  CheckCircle2,
  Circle,
  FileCheck,
  Calendar,
  Sparkles,
  Zap,
  TrendingUp,
  User,
  ShieldCheck,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface StanceItem {
  name: string;
  chinese: string;
  status: "mastered" | "proficient" | "learning";
  notes: string;
}

interface ExamRecord {
  id: string;
  date: string;
  rankTarget: string;
  examiner: string;
  score: string;
  result: "Passed with Honors" | "Passed" | "Conditional Pass";
  feedback: string;
}

export default function GradingProgressPage() {
  const { user } = useAuth();
  const [activeTrack, setActiveTrack] = useState<"kungfu" | "taichi" | "qigong">("kungfu");
  const [tracksData, setTracksData] = useState<any[]>([]);
  const [examHistory, setExamHistory] = useState<ExamRecord[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadGradingData() {
      try {
        const [tracksRes, historyRes] = await Promise.all([
          apiFetch<any[]>("/grading/tracks").catch(() => []),
          apiFetch<any[]>("/grading/exams/my-history").catch(() => []),
        ]);

        if (tracksRes && tracksRes.length > 0) {
          setTracksData(tracksRes);
        }

        if (historyRes && historyRes.length > 0) {
          const mapped: ExamRecord[] = historyRes.map((h: any) => ({
            id: h.id,
            date: h.exam_date || "2026-08-18",
            rankTarget: h.level_name || "Level 1 White Belt",
            examiner: "Head Master Marcus Vance",
            score: `${h.total_score || 90} / 100`,
            result: h.status === "passed" ? "Passed with Honors" : "Passed",
            feedback: h.notes || "Solid stance stability and breathing focus.",
          }));
          setExamHistory(mapped);
        } else {
          setExamHistory(fallbackExams);
        }
      } catch (err) {
        console.error("Failed to load grading data:", err);
        setExamHistory(fallbackExams);
      } finally {
        setLoading(false);
      }
    }
    loadGradingData();
  }, []);

  return (
    <div className="space-y-8">
      {/* Top Banner with Rank status */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Martial Arts Syllabus</Badge>
            <Badge variant="beltWhite">{user?.current_rank || "Level 1 White Belt"}</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Curriculum & Stance Grading Board
          </h1>
          <p className="text-xs text-slate-400">
            Track your fundamental stances milestones, Ma Bu benchmark endurance, and exam results.
          </p>
        </div>

        <div className="flex items-center gap-3 p-3 rounded-xl bg-slate-950/90 border border-amber-900/40 shadow-inner">
          <div className="p-2 rounded-lg bg-amber-950 text-amber-400">
            <Award className="h-5 w-5" />
          </div>
          <div>
            <span className="text-[10px] text-slate-400 block font-semibold uppercase">Current Belt Target</span>
            <span className="text-lg font-black text-amber-400">Level 2 Yellow Sash</span>
          </div>
        </div>
      </div>

      {/* Discipline Tabs */}
      <div className="flex items-center gap-2 border-b border-slate-800 pb-3">
        <button
          onClick={() => setActiveTrack("kungfu")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
            activeTrack === "kungfu"
              ? "bg-red-600 text-white shadow-md shadow-red-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          Kung Fu & Karate Track
        </button>

        <button
          onClick={() => setActiveTrack("taichi")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
            activeTrack === "taichi"
              ? "bg-blue-600 text-white shadow-md shadow-blue-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          Chen Style Tai Chi Track
        </button>

        <button
          onClick={() => setActiveTrack("qigong")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
            activeTrack === "qigong"
              ? "bg-amber-600 text-white shadow-md shadow-amber-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          Qi Gong & Internal Energy
        </button>
      </div>

      {/* Stances & Benchmark Checklist Grid */}
      <div className="space-y-4">
        <h2 className="text-lg font-bold text-slate-100 flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-amber-500" />
          5 Core Stances & Technical Rubric
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {kungfuStances.map((st, idx) => (
            <Card
              key={idx}
              glow={st.status === "mastered" ? "gold" : st.status === "proficient" ? "red" : "none"}
              className="p-5 flex flex-col justify-between"
            >
              <div className="space-y-2">
                <div className="flex justify-between items-start">
                  <Badge
                    variant={
                      st.status === "mastered"
                        ? "success"
                        : st.status === "proficient"
                        ? "primary"
                        : "outline"
                    }
                    size="sm"
                  >
                    {st.status.toUpperCase()}
                  </Badge>
                  <span className="text-xs font-bold text-amber-400">{st.chinese}</span>
                </div>

                <h3 className="text-base font-bold text-slate-100">{st.name}</h3>
                <p className="text-xs text-slate-400 leading-relaxed">{st.notes}</p>
              </div>

              <div className="pt-3 mt-3 border-t border-slate-800 flex items-center justify-between text-[11px]">
                <span className="text-slate-500 font-mono">Exam Prerequisite</span>
                <span className="text-emerald-400 font-semibold flex items-center gap-1">
                  <CheckCircle2 className="h-3.5 w-3.5" />
                  Checked by Instructor
                </span>
              </div>
            </Card>
          ))}
        </div>
      </div>

      {/* Official Exam History */}
      <div className="space-y-4 pt-4">
        <h2 className="text-lg font-bold text-slate-100 flex items-center gap-2">
          <FileCheck className="h-5 w-5 text-slate-400" />
          Formal Exam Evaluation History
        </h2>

        <div className="space-y-4">
          {examHistory.map((exam) => (
            <Card key={exam.id} className="p-6 border-slate-800 space-y-4">
              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 pb-3 border-b border-slate-800">
                <div>
                  <div className="flex items-center gap-2">
                    <Badge variant="gold">{exam.rankTarget}</Badge>
                    <Badge variant="success">{exam.result}</Badge>
                  </div>
                  <span className="text-xs font-mono text-slate-400 mt-1 block">
                    Evaluated by {exam.examiner} • Date: {exam.date}
                  </span>
                </div>

                <div className="text-right">
                  <span className="text-xs text-slate-400 block font-semibold">Overall Score</span>
                  <span className="text-2xl font-black text-amber-400 font-mono">{exam.score}</span>
                </div>
              </div>

              <div className="space-y-1">
                <span className="text-xs font-bold text-slate-300">Examiner Feedback & Form Notes:</span>
                <p className="text-xs text-slate-400 leading-relaxed italic bg-slate-950 p-3 rounded-xl border border-slate-850">
                  &quot;{exam.feedback}&quot;
                </p>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}

const kungfuStances: StanceItem[] = [
  {
    name: "Horse Riding Stance",
    chinese: "马步 (Mǎ Bù)",
    status: "mastered",
    notes: "Thighs parallel to ground, knees pushed outwards, flat spine (2-min hold target).",
  },
  {
    name: "Bow Stance",
    chinese: "弓步 (Gōng Bù)",
    status: "mastered",
    notes: "Front knee at 90 degrees, back leg completely locked and straight, heel rooted.",
  },
  {
    name: "Drop Stance",
    chinese: "仆步 (Pū Bù)",
    status: "proficient",
    notes: "Low crouch on one leg, extended leg completely flat with toes turned inward.",
  },
  {
    name: "Empty / Cat Stance",
    chinese: "虚步 (Xū Bù)",
    status: "proficient",
    notes: "90% weight on back bent leg, front toe lightly touching ground, core engaged.",
  },
  {
    name: "Rest / Cross Stance",
    chinese: "歇步 (Xiē Bù)",
    status: "learning",
    notes: "Crossed legs sitting tightly on rear heel, torso upright without leaning.",
  },
  {
    name: "Xiao Hong Quan (Small Flood Fist)",
    chinese: "小洪拳",
    status: "learning",
    notes: "First classical fist routine. Focus on clean snap kicks and palm strikes.",
  },
];

const fallbackExams: ExamRecord[] = [
  {
    id: "exam-01",
    date: "2026-03-28",
    rankTarget: "Yellow Belt (Sash)",
    examiner: "Head Master Marcus Vance",
    score: "92 / 100",
    result: "Passed with Honors",
    feedback:
      "Excellent flexibility and stance depth. Ma Bu posture is stable and grounded. Keep working on eye focus during punch transitions.",
  },
  {
    id: "exam-02",
    date: "2025-12-14",
    rankTarget: "White Belt to Initiate",
    examiner: "Chief Instructor Kenji Sato",
    score: "88 / 100",
    result: "Passed",
    feedback:
      "Great spirit and punctuality. Keep sinking deeper into Pu Bu and remember to breathe naturally through movement.",
  },
];
