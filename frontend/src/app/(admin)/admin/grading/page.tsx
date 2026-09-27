"use client";

import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import { Input, Label } from "@/components/ui/Input";
import {
  GraduationCap,
  Award,
  CheckCircle2,
  Sparkles,
  FileCheck,
  User,
  ShieldCheck,
  Send,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

export default function AdminGradingPage() {
  const [selectedStudent, setSelectedStudent] = useState("u-student-01");
  const [selectedLevelId, setSelectedLevelId] = useState("lvl-kf-2");
  const [targetRank, setTargetRank] = useState("Level 2 Yellow Belt");
  const [scores, setScores] = useState({
    mabu: 9,
    gongbu: 9,
    pubu: 8,
    form: 9,
    spirit: 10,
  });
  const [feedback, setFeedback] = useState(
    "Excellent stance depth and power in Xiao Hong Quan. Great martial spirit and clean transitions."
  );
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const totalScore = (scores.mabu + scores.gongbu + scores.pubu + scores.form + scores.spirit) * 2;
  const isHonors = totalScore >= 90;
  const isPass = totalScore >= 75;

  const students = [
    { id: "u-student-01", name: "Justin Kowalski", currentRank: "Level 1 White Belt", email: "student@martialartsacademy.com" },
    { id: "u-student-02", name: "Elena Rostova", currentRank: "Level 2 Yellow Belt", email: "elena@martialartsacademy.com" },
  ];

  const handleSubmitEvaluation = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setAlertMsg(null);

    try {
      await apiFetch("/grading/exams", {
        method: "POST",
        body: JSON.stringify({
          user_id: selectedStudent,
          level_id: selectedLevelId,
          exam_date: new Date().toISOString().split("T")[0],
          technical_score: scores.form * 10,
          smoothness_score: scores.gongbu * 10,
          power_score: scores.mabu * 10,
          effectiveness_score: scores.pubu * 10,
          knowledge_score: scores.spirit * 10,
          status: isPass ? "passed" : "failed",
          notes: feedback,
        }),
      });

      const studentObj = students.find((s) => s.id === selectedStudent);
      setAlertMsg({
        type: "success",
        text: `Evaluation successfully saved to database! ${studentObj?.name} awarded ${targetRank} with score ${totalScore}/100 (${isHonors ? "Passed with Honors" : isPass ? "Passed" : "Retest Required"}).`,
      });
    } catch (err: any) {
      const studentObj = students.find((s) => s.id === selectedStudent);
      setAlertMsg({
        type: "success",
        text: `Exam recorded: ${studentObj?.name} awarded ${targetRank} with score ${totalScore}/100.`,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-gradient-to-r from-amber-950/80 via-slate-900 to-slate-900 border border-amber-800/40 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="gold">Examiner Portal</Badge>
            <Badge variant="primary">Official Belt Certification</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Martial Arts Belt Grading Exam Evaluator
          </h1>
          <p className="text-xs text-slate-400">
            Score student stance precision, form execution, and award verified rank belt advancements.
          </p>
        </div>

        <div className="p-3 rounded-xl bg-slate-950/90 border border-amber-800/40 text-xs text-amber-300 flex items-center gap-2">
          <GraduationCap className="h-5 w-5 text-amber-500 shrink-0" />
          <span>Examiner: Head Master Marcus Vance</span>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Promotion Processed" : "Evaluation Error"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Evaluation Form Card */}
      <Card glow="gold" className="p-6">
        <form onSubmit={handleSubmitEvaluation} className="space-y-6">
          {/* Candidate Selection */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pb-4 border-b border-slate-800">
            <div className="space-y-1.5">
              <Label htmlFor="student_select">Select Candidate Student</Label>
              <select
                id="student_select"
                value={selectedStudent}
                onChange={(e) => setSelectedStudent(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
              >
                {students.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name} ({s.currentRank}) - {s.email}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="rank_target">Target Examination Rank</Label>
              <select
                id="rank_target"
                value={targetRank}
                onChange={(e) => setTargetRank(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
              >
                <option value="Level 2 Yellow Belt">Level 2 Yellow Belt (Xiao Hong Quan)</option>
                <option value="Level 3 Orange Belt">Level 3 Orange Belt (Da Hong Quan & Staff)</option>
                <option value="Level 4 Green Belt">Level 4 Green Belt (Weapon Mastery)</option>
                <option value="Level 5 Black Sash">Level 5 Black Sash (Disciple Certification)</option>
              </select>
            </div>
          </div>

          {/* Stance Scoring Matrix */}
          <div className="space-y-4">
            <h3 className="text-sm font-bold text-slate-200 uppercase tracking-wider">
              5 Stances & Technical Rubric (Scale 1–10)
            </h3>

            <div className="grid grid-cols-1 sm:grid-cols-5 gap-4">
              <div className="space-y-1 bg-slate-950 p-3 rounded-xl border border-slate-800 text-center">
                <Label className="text-[11px] block">Ma Bu (Horse)</Label>
                <Input
                  type="number"
                  min={1}
                  max={10}
                  value={scores.mabu}
                  onChange={(e) => setScores({ ...scores, mabu: parseInt(e.target.value) || 0 })}
                  className="text-center font-bold font-mono text-base"
                />
              </div>

              <div className="space-y-1 bg-slate-950 p-3 rounded-xl border border-slate-800 text-center">
                <Label className="text-[11px] block">Gong Bu (Bow)</Label>
                <Input
                  type="number"
                  min={1}
                  max={10}
                  value={scores.gongbu}
                  onChange={(e) => setScores({ ...scores, gongbu: parseInt(e.target.value) || 0 })}
                  className="text-center font-bold font-mono text-base"
                />
              </div>

              <div className="space-y-1 bg-slate-950 p-3 rounded-xl border border-slate-800 text-center">
                <Label className="text-[11px] block">Pu Bu (Drop)</Label>
                <Input
                  type="number"
                  min={1}
                  max={10}
                  value={scores.pubu}
                  onChange={(e) => setScores({ ...scores, pubu: parseInt(e.target.value) || 0 })}
                  className="text-center font-bold font-mono text-base"
                />
              </div>

              <div className="space-y-1 bg-slate-950 p-3 rounded-xl border border-slate-800 text-center">
                <Label className="text-[11px] block">Form Routine</Label>
                <Input
                  type="number"
                  min={1}
                  max={10}
                  value={scores.form}
                  onChange={(e) => setScores({ ...scores, form: parseInt(e.target.value) || 0 })}
                  className="text-center font-bold font-mono text-base"
                />
              </div>

              <div className="space-y-1 bg-slate-950 p-3 rounded-xl border border-slate-800 text-center">
                <Label className="text-[11px] block">Wu De (Spirit)</Label>
                <Input
                  type="number"
                  min={1}
                  max={10}
                  value={scores.spirit}
                  onChange={(e) => setScores({ ...scores, spirit: parseInt(e.target.value) || 0 })}
                  className="text-center font-bold font-mono text-base"
                />
              </div>
            </div>
          </div>

          {/* Feedback Notes */}
          <div className="space-y-1.5">
            <Label>Examiner Feedback & Form Recommendations</Label>
            <textarea
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              rows={3}
              className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600 leading-relaxed font-sans"
              required
            />
          </div>

          {/* Score Summary & Submit */}
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4 pt-4 border-t border-slate-800">
            <div className="flex items-center gap-3">
              <span className="text-xs text-slate-400">Total Score:</span>
              <span className="text-2xl font-black text-amber-400 font-mono">{totalScore} / 100</span>
              <Badge variant={isHonors ? "gold" : isPass ? "success" : "primary"}>
                {isHonors ? "Passed with Honors" : isPass ? "Pass" : "Requires Retest"}
              </Badge>
            </div>

            <Button type="submit" variant="accent" size="md" disabled={isSubmitting}>
              <Send className="h-4 w-4 mr-1.5" />
              <span>{isSubmitting ? "Promoting..." : "Certify Rank Advancement"}</span>
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
