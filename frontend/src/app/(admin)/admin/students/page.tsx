"use client";

import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import { Input, Label } from "@/components/ui/Input";
import { Alert } from "@/components/ui/Alert";
import {
  UserCheck,
  Coins,
  ShieldCheck,
  Search,
  Award,
  Edit2,
  Sparkles,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface StudentRecord {
  id: string;
  name: string;
  email: string;
  rank: string;
  tier: string;
  tokens: number;
  overallPoints: number;
  isExempt: boolean;
}

export default function AdminStudentsPage() {
  const [search, setSearch] = useState("");
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const [students, setStudents] = useState<StudentRecord[]>([
    {
      id: "u-student-01",
      name: "Justin Kowalski",
      email: "student@martialartsacademy.com",
      rank: "Level 1 White Belt",
      tier: "Iron Disciple (铁弟子)",
      tokens: 14,
      overallPoints: 380,
      isExempt: false,
    },
    {
      id: "u-student-02",
      name: "Elena Rostova",
      email: "elena@martialartsacademy.com",
      rank: "Level 2 Yellow Belt",
      tier: "Bronze Monk (铜和尚)",
      tokens: 28,
      overallPoints: 720,
      isExempt: true,
    },
  ]);

  // Token adjustment modal state
  const [tokenModalOpen, setTokenModalOpen] = useState(false);
  const [selectedStudent, setSelectedStudent] = useState<StudentRecord | null>(null);
  const [adjustAmount, setAdjustAmount] = useState<number>(5);
  const [adjustReason, setAdjustReason] = useState("Belt grading exam prerequisite bonus");

  useEffect(() => {
    async function loadStudents() {
      try {
        const boardData = await apiFetch<any[]>("/users/leaderboard/overall");
        if (boardData && boardData.length > 0) {
          const mapped: StudentRecord[] = boardData.map((b: any, idx: number) => ({
            id: b.user_id,
            name: b.first_name ? `${b.first_name} ${b.last_name}` : `Disciple #${idx + 1}`,
            email: b.email || `disciple${idx + 1}@martialartsacademy.com`,
            rank: b.current_rank || "Level 1 White Belt",
            tier: b.level_tier || "Novice Disciple",
            tokens: 14,
            overallPoints: b.overall_points || 0,
            isExempt: idx === 1,
          }));
          setStudents(mapped);
        }
      } catch (err) {
        console.error("Failed to load students leaderboard:", err);
      }
    }
    loadStudents();
  }, []);

  const handleToggleExemption = (studentId: string) => {
    setStudents(
      students.map((st) => {
        if (st.id === studentId) {
          const next = !st.isExempt;
          setAlertMsg({
            type: "success",
            text: `${st.name} annual membership fee exemption ${next ? "GRANTED" : "REVOKED"}.`,
          });
          return { ...st, isExempt: next };
        }
        return st;
      })
    );
  };

  const handleConfirmTokenAdjustment = () => {
    if (!selectedStudent) return;
    setStudents(
      students.map((st) => {
        if (st.id === selectedStudent.id) {
          return { ...st, tokens: st.tokens + adjustAmount };
        }
        return st;
      })
    );
    setAlertMsg({
      type: "success",
      text: `Adjusted tokens for ${selectedStudent.name}: +${adjustAmount} Tokens (${adjustReason}).`,
    });
    setTokenModalOpen(false);
  };

  const filteredStudents = students.filter(
    (s) =>
      s.name.toLowerCase().includes(search.toLowerCase()) ||
      s.email.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-8">
      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Operations Console</Badge>
            <Badge variant="gold">Student Roster & Ledgers</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Disciple Directory & Membership Management
          </h1>
          <p className="text-xs text-slate-400">
            Inspect student registration status, grant charitable fee exemptions, and adjust token ledger balances.
          </p>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Operation Saved" : "Notice"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Filter and Search Bar */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="h-4 w-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by student name or email..."
            className="pl-9"
          />
        </div>
      </div>

      {/* Students Table */}
      <Card className="overflow-hidden border-slate-800">
        <div className="divide-y divide-slate-800">
          {filteredStudents.map((st) => (
            <div
              key={st.id}
              className="p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 hover:bg-slate-900/40 transition-colors"
            >
              <div className="space-y-1.5">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-slate-100">{st.name}</span>
                  <Badge variant="beltWhite" size="sm">{st.rank}</Badge>
                  <Badge variant="gold" size="sm">{st.tier}</Badge>
                  {st.isExempt && (
                    <Badge variant="primary" size="sm" className="flex items-center gap-1">
                      <ShieldCheck className="h-3 w-3" />
                      Fee Exempt
                    </Badge>
                  )}
                </div>

                <div className="flex items-center gap-4 text-xs text-slate-400 font-mono">
                  <span>{st.email}</span>
                  <span className="text-amber-400 font-bold">{st.tokens} Tokens Available</span>
                  <span className="text-slate-300">{st.overallPoints} Overall Points</span>
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setSelectedStudent(st);
                    setTokenModalOpen(true);
                  }}
                >
                  <Coins className="h-3.5 w-3.5 mr-1 text-amber-400" />
                  Adjust Tokens
                </Button>

                <Button
                  variant={st.isExempt ? "primary" : "secondary"}
                  size="sm"
                  onClick={() => handleToggleExemption(st.id)}
                >
                  {st.isExempt ? "Revoke Exemption" : "Grant Exemption"}
                </Button>
              </div>
            </div>
          ))}
        </div>
      </Card>

      {/* Token Adjustment Modal */}
      <Modal
        isOpen={tokenModalOpen}
        onClose={() => setTokenModalOpen(false)}
        title={`Adjust Tokens: ${selectedStudent?.name}`}
        description="Manually credit or deduct tokens from this student's account ledger:"
      >
        <div className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label>Tokens Delta (+ / -)</Label>
            <Input
              type="number"
              value={adjustAmount}
              onChange={(e) => setAdjustAmount(parseInt(e.target.value) || 0)}
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label>Audit Reason & Reference</Label>
            <Input
              value={adjustReason}
              onChange={(e) => setAdjustReason(e.target.value)}
              placeholder="e.g. Belt grading exam prerequisite or administrative correction"
              required
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" size="sm" onClick={() => setTokenModalOpen(false)}>
              Cancel
            </Button>
            <Button variant="accent" size="sm" onClick={handleConfirmTokenAdjustment}>
              Save Ledger Adjustment
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
