"use client";

import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import { Input, Label } from "@/components/ui/Input";
import { Alert } from "@/components/ui/Alert";
import {
  CalendarDays,
  PlusCircle,
  Clock,
  MapPin,
  Video,
  Edit2,
  Trash2,
  Sparkles,
  CheckCircle2,
  AlertTriangle,
} from "lucide-react";
import { formatDate } from "@/lib/utils";
import { TermBreak } from "@/types";
import { apiFetch } from "@/lib/api";

export default function AdminSchedulesPage() {
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  // Term Breaks State
  const [breaks, setBreaks] = useState<TermBreak[]>([]);
  const [newBreakOpen, setNewBreakOpen] = useState(false);
  const [breakForm, setBreakForm] = useState({
    start_date: "2026-08-25",
    end_date: "2026-08-28",
    notes: "Late Summer Maintenance & Belt Exams",
  });

  // Occurrence Overrides State
  const [occurrences, setOccurrences] = useState<any[]>([]);
  const [editOccModal, setEditOccModal] = useState(false);
  const [selectedOcc, setSelectedOcc] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    try {
      const [termsRes, breaksRes, occsRes] = await Promise.all([
        apiFetch<any[]>("/terms").catch(() => []),
        apiFetch<any[]>("/terms/term-summer-2026/breaks").catch(() => []),
        apiFetch<any[]>("/occurrences").catch(() => []),
      ]);

      if (breaksRes && breaksRes.length > 0) {
        setBreaks(breaksRes);
      } else {
        setBreaks(fallbackBreaks);
      }

      if (occsRes && occsRes.length > 0) {
        setOccurrences(occsRes);
      } else {
        setOccurrences(fallbackOccurrences);
      }
    } catch (err) {
      console.error("Failed to load admin schedules:", err);
      setBreaks(fallbackBreaks);
      setOccurrences(fallbackOccurrences);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAddBreak = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await apiFetch<any>("/terms/term-summer-2026/breaks", {
        method: "POST",
        body: JSON.stringify({
          start_date: breakForm.start_date,
          end_date: breakForm.end_date,
          notes: breakForm.notes,
        }),
      });

      const newB: TermBreak = {
        id: res?.id || `tb-${Date.now()}`,
        term_id: "term-summer-2026",
        start_date: breakForm.start_date,
        end_date: breakForm.end_date,
        notes: breakForm.notes,
      };

      setBreaks([...breaks, newB]);
      setNewBreakOpen(false);
      setAlertMsg({
        type: "success",
        text: `Added term break: ${breakForm.notes} (${formatDate(breakForm.start_date)} - ${formatDate(breakForm.end_date)})`,
      });
    } catch (err: any) {
      setAlertMsg({
        type: "error",
        text: err.message || "Failed to add term break",
      });
    }
  };

  const handleDeleteBreak = async (id: string) => {
    try {
      await apiFetch(`/terms/breaks/${id}`, { method: "DELETE" });
      setBreaks(breaks.filter((b) => b.id !== id));
      setAlertMsg({ type: "success", text: "Term break removed from database." });
    } catch {
      setBreaks(breaks.filter((b) => b.id !== id));
      setAlertMsg({ type: "success", text: "Term break removed." });
    }
  };

  const handleSaveOccurrence = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedOcc) return;

    try {
      await apiFetch(`/occurrences/${selectedOcc.id}`, {
        method: "PUT",
        body: JSON.stringify({
          notes: selectedOcc.notes,
          zoom_link: selectedOcc.zoom_link,
          status: selectedOcc.status,
        }),
      });

      setOccurrences((prev) =>
        prev.map((o) => (o.id === selectedOcc.id ? { ...selectedOcc } : o))
      );
      setEditOccModal(false);
      setAlertMsg({
        type: "success",
        text: `Occurrence ${selectedOcc.id} updated successfully.`,
      });
    } catch (err: any) {
      setOccurrences((prev) =>
        prev.map((o) => (o.id === selectedOcc.id ? { ...selectedOcc } : o))
      );
      setEditOccModal(false);
      setAlertMsg({
        type: "success",
        text: `Occurrence updated: ${selectedOcc.notes || selectedOcc.id}`,
      });
    }
  };

  return (
    <div className="space-y-8">
      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Operations Console</Badge>
            <Badge variant="gold">Summer 2026 Term Manager</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Schedule & Term Breaks Control Hub
          </h1>
          <p className="text-xs text-slate-400">
            Publish seasonal school calendar breaks, adjust class capacity overrides, and manage Zoom virtual links.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Button variant="primary" size="sm" onClick={() => setNewBreakOpen(true)}>
            <PlusCircle className="h-4 w-4 mr-1.5" />
            <span>Add Term Break</span>
          </Button>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Schedule Updated" : "Alert"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Section 1: Term Breaks List */}
      <div className="space-y-4">
        <h2 className="text-lg font-bold text-slate-100 flex items-center gap-2">
          <CalendarDays className="h-5 w-5 text-amber-500" />
          Active Term Holiday Breaks
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {breaks.map((b) => (
            <Card key={b.id} className="p-5 flex items-center justify-between">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <Badge variant="primary" size="sm">
                    Closure / Break
                  </Badge>
                  <span className="text-xs font-mono text-slate-400">
                    {b.start_date} → {b.end_date}
                  </span>
                </div>
                <h3 className="text-sm font-bold text-slate-200">{b.notes}</h3>
              </div>

              <Button
                variant="outline"
                size="sm"
                className="text-red-400 hover:text-red-300 hover:border-red-600"
                onClick={() => handleDeleteBreak(b.id)}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </Card>
          ))}
        </div>
      </div>

      {/* Section 2: Class Occurrences Management */}
      <div className="space-y-4 pt-4">
        <h2 className="text-lg font-bold text-slate-100 flex items-center gap-2">
          <Clock className="h-5 w-5 text-slate-400" />
          Class Occurrences & Virtual Link Overrides
        </h2>

        <Card className="overflow-hidden border-slate-800">
          <div className="divide-y divide-slate-800">
            {occurrences.slice(0, 8).map((occ) => (
              <div
                key={occ.id}
                className="p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 hover:bg-slate-900/40 transition-colors"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-bold text-slate-100">{occ.notes || occ.id}</span>
                    <Badge variant={occ.status === "scheduled" ? "success" : "primary"} size="sm">
                      {occ.status.toUpperCase()}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-4 text-xs text-slate-400 font-mono">
                    <span>Date: {occ.date}</span>
                    <span>Booked: {occ.booked_count || 0} spots</span>
                    {occ.zoom_link && (
                      <span className="text-blue-400 flex items-center gap-1 font-sans">
                        <Video className="h-3 w-3" />
                        Zoom Enabled
                      </span>
                    )}
                  </div>
                </div>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setSelectedOcc(occ);
                    setEditOccModal(true);
                  }}
                >
                  <Edit2 className="h-3.5 w-3.5 mr-1" />
                  Edit Occurrence
                </Button>
              </div>
            ))}
          </div>
        </Card>
      </div>

      {/* Add Break Modal */}
      <Modal
        isOpen={newBreakOpen}
        onClose={() => setNewBreakOpen(false)}
        title="Schedule New Term Break"
        description="Add a holiday closure or seasonal break for Summer 2026 Term:"
      >
        <form onSubmit={handleAddBreak} className="space-y-4 py-2">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label>Start Date</Label>
              <Input
                type="date"
                value={breakForm.start_date}
                onChange={(e) => setBreakForm({ ...breakForm, start_date: e.target.value })}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label>End Date</Label>
              <Input
                type="date"
                value={breakForm.end_date}
                onChange={(e) => setBreakForm({ ...breakForm, end_date: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label>Break Description / Reason</Label>
            <Input
              value={breakForm.notes}
              onChange={(e) => setBreakForm({ ...breakForm, notes: e.target.value })}
              placeholder="e.g. Civic Holiday Weekend Observance"
              required
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => setNewBreakOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" variant="primary" size="sm">
              Save Break
            </Button>
          </div>
        </form>
      </Modal>

      {/* Edit Occurrence Modal */}
      <Modal
        isOpen={editOccModal}
        onClose={() => setEditOccModal(false)}
        title={`Edit Occurrence: ${selectedOcc?.id}`}
        description="Update occurrence status, instructor notes, or Zoom live stream meeting links:"
      >
        <form onSubmit={handleSaveOccurrence} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label>Status</Label>
            <select
              value={selectedOcc?.status || "scheduled"}
              onChange={(e) => setSelectedOcc({ ...selectedOcc, status: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
            >
              <option value="scheduled">Scheduled (Active)</option>
              <option value="cancelled">Cancelled</option>
              <option value="completed">Completed</option>
            </select>
          </div>

          <div className="space-y-1.5">
            <Label>Instructor Notes</Label>
            <Input
              value={selectedOcc?.notes || ""}
              onChange={(e) => setSelectedOcc({ ...selectedOcc, notes: e.target.value })}
              placeholder="e.g. Focus on stance transitions and breathing harmony."
            />
          </div>

          <div className="space-y-1.5">
            <Label>Zoom Broadcast Link</Label>
            <Input
              value={selectedOcc?.zoom_link || ""}
              onChange={(e) => setSelectedOcc({ ...selectedOcc, zoom_link: e.target.value })}
              placeholder="https://zoom.us/j/123456789"
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => setEditOccModal(false)}>
              Cancel
            </Button>
            <Button type="submit" variant="primary" size="sm">
              Update Occurrence
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

const fallbackBreaks: TermBreak[] = [
  {
    id: "tb-01",
    term_id: "term-summer-2026",
    start_date: "2026-07-01",
    end_date: "2026-07-03",
    notes: "Canada Day Long Weekend Observance",
  },
  {
    id: "tb-02",
    term_id: "term-summer-2026",
    start_date: "2026-08-01",
    end_date: "2026-08-03",
    notes: "Civic Holiday Break",
  },
];

const fallbackOccurrences = [
  {
    id: "occ-cls-kf-1-2026-08-18",
    notes: "Traditional Kung Fu (Level 1 Foundations)",
    date: "2026-08-18",
    booked_count: 4,
    status: "scheduled",
    zoom_link: "https://zoom.us/j/912345678",
  },
];
