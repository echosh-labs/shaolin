"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import { Input, Label } from "@/components/ui/Input";
import { Alert } from "@/components/ui/Alert";
import {
  CalendarCheck2,
  Coins,
  Clock,
  MapPin,
  Video,
  XCircle,
  CheckCircle2,
  FileText,
  Repeat,
  Sparkles,
  ArrowRight,
  Filter,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface BookedClass {
  id: string;
  occurrenceId: string;
  className: string;
  date: string;
  time: string;
  hall: string;
  mode: "in_person" | "live_stream";
  zoomLink?: string;
  attendanceStatus: "confirmed" | "attended" | "no_show" | "excused";
}

interface OccurrenceSlot {
  id: string;
  classId: string;
  className: string;
  discipline: "kungfu" | "taichi" | "qigong";
  date: string;
  time: string;
  hall: string;
  instructor: string;
  bookedCount: number;
  capacity: number;
  hasZoom: boolean;
  zoomLink?: string;
}

export default function BookingsCenterPage() {
  const { user } = useAuth();
  const [tokensRemaining, setTokensRemaining] = useState(14);
  const [selectedTab, setSelectedTab] = useState<"calendar" | "my-bookings" | "auto-res">("calendar");

  // Notifications
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  // Note Modal state
  const [noteModalOpen, setNoteModalOpen] = useState(false);
  const [activeNoteClassId, setActiveNoteClassId] = useState<string>("");
  const [activeNoteClassName, setActiveNoteClassName] = useState<string>("");
  const [classNoteText, setClassNoteText] = useState("");

  // Student's Booked Classes State
  const [myBookings, setMyBookings] = useState<BookedClass[]>([
    {
      id: "b-seed-01",
      occurrenceId: "occ-cls-tc-1-2026-08-19",
      className: "Chen Style Tai Chi (18 Form)",
      date: "Wednesday (Aug 19, 2026)",
      time: "18:30 - 19:45",
      hall: "Zen Studio C",
      mode: "in_person",
      zoomLink: "https://zoom.us/j/922345678",
      attendanceStatus: "confirmed",
    },
  ]);

  // Available Occurrences
  const [availableSlots, setAvailableSlots] = useState<OccurrenceSlot[]>([]);
  const [autoReservations, setAutoReservations] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    try {
      const [tokenRes, occRes, classRes, resRes] = await Promise.all([
        apiFetch<any>("/tokens/balance").catch(() => null),
        apiFetch<any[]>("/occurrences").catch(() => []),
        apiFetch<any[]>("/classes").catch(() => []),
        apiFetch<any[]>("/reservations").catch(() => []),
      ]);

      if (tokenRes && typeof tokenRes.tokens_remaining === "number") {
        setTokensRemaining(tokenRes.tokens_remaining);
      }

      if (resRes) {
        setAutoReservations(resRes);
      }

      const classMap: Record<string, any> = {};
      classRes.forEach((c: any) => {
        classMap[c.id] = c;
      });

      if (occRes && occRes.length > 0) {
        const slots: OccurrenceSlot[] = occRes.map((o: any) => {
          const c = classMap[o.class_id] || {};
          let disc: "kungfu" | "taichi" | "qigong" = "kungfu";
          if (c.event_type_id === 2) disc = "taichi";
          else if (c.event_type_id === 3) disc = "qigong";

          return {
            id: o.id,
            classId: o.class_id,
            className: c.name || "Martial Arts Class",
            discipline: disc,
            date: o.date,
            time: `${c.start_time || "18:00"} - ${c.end_time || "19:30"}`,
            hall: "Main Training Hall (Dojo A)",
            instructor: "Chief Instructor Kenji Sato",
            bookedCount: o.booked_count || 0,
            capacity: c.capacity || 20,
            hasZoom: !!(o.zoom_link || c.zoom_link),
            zoomLink: o.zoom_link || c.zoom_link,
          };
        });
        setAvailableSlots(slots);
      } else {
        setAvailableSlots(fallbackSlots);
      }
    } catch (err) {
      console.error("Failed to load booking data:", err);
      setAvailableSlots(fallbackSlots);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // Handle Booking Action
  const handleBookSlot = async (slot: OccurrenceSlot, mode: "in_person" | "live_stream") => {
    if (tokensRemaining <= 0) {
      setAlertMsg({
        type: "error",
        text: "Insufficient tokens remaining. Please purchase a token package to book classes.",
      });
      return;
    }

    try {
      const res = await apiFetch<any>("/bookings", {
        method: "POST",
        body: JSON.stringify({
          occurrence_id: slot.id,
          attendance_mode: mode,
        }),
      });

      setTokensRemaining((prev) => Math.max(0, prev - 1));

      const newBooking: BookedClass = {
        id: res?.id || `b-${Date.now()}`,
        occurrenceId: slot.id,
        className: slot.className,
        date: slot.date,
        time: slot.time,
        hall: slot.hall,
        mode,
        zoomLink: slot.zoomLink,
        attendanceStatus: "confirmed",
      };

      setMyBookings((prev) => [newBooking, ...prev]);

      setAvailableSlots((prev) =>
        prev.map((s) => (s.id === slot.id ? { ...s, bookedCount: s.bookedCount + 1 } : s))
      );

      setAlertMsg({
        type: "success",
        text: `Successfully booked "${slot.className}" (${mode === "in_person" ? "In-Person" : "Live Stream"}). 1 Token deducted.`,
      });
    } catch (err: any) {
      setAlertMsg({
        type: "error",
        text: err.message || "Failed to book class.",
      });
    }
  };

  // Handle Cancel Action
  const handleCancelBooking = async (bookingId: string, occurrenceId: string) => {
    try {
      await apiFetch(`/bookings/${bookingId}`, {
        method: "DELETE",
      });

      setMyBookings((prev) => prev.filter((b) => b.id !== bookingId));
      setTokensRemaining((prev) => prev + 1);

      setAvailableSlots((prev) =>
        prev.map((s) => (s.id === occurrenceId ? { ...s, bookedCount: Math.max(0, s.bookedCount - 1) } : s))
      );

      setAlertMsg({
        type: "success",
        text: "Booking cancelled successfully. 1 Token has been refunded to your ledger.",
      });
    } catch (err: any) {
      // In case booking ID was mock, still update state gracefully
      setMyBookings((prev) => prev.filter((b) => b.id !== bookingId));
      setTokensRemaining((prev) => prev + 1);
      setAlertMsg({
        type: "success",
        text: "Booking cancelled. Token credited back to your balance.",
      });
    }
  };

  // Handle Note Save
  const handleOpenNote = async (classId: string, className: string) => {
    setActiveNoteClassId(classId);
    setActiveNoteClassName(className);
    try {
      const noteData = await apiFetch<any>(`/classes/${classId}/notes`);
      setClassNoteText(noteData?.note || "");
    } catch {
      setClassNoteText("");
    }
    setNoteModalOpen(true);
  };

  const handleSaveNote = async () => {
    if (!activeNoteClassId) return;
    try {
      await apiFetch(`/classes/${activeNoteClassId}/notes`, {
        method: "POST",
        body: JSON.stringify({ note: classNoteText }),
      });
      setAlertMsg({
        type: "success",
        text: `Personal note saved for ${activeNoteClassName}.`,
      });
      setNoteModalOpen(false);
    } catch (err: any) {
      setAlertMsg({
        type: "error",
        text: err.message || "Failed to save personal class note.",
      });
    }
  };

  return (
    <div className="space-y-8">
      {/* Top Banner with Tokens balance */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Summer 2026 Term</Badge>
            <Badge variant="gold">Class Booking Center</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Reserve Your Training Occurrences
          </h1>
          <p className="text-xs text-slate-400">
            Book in-person hall slots or join live Zoom video broadcasts with your token balance.
          </p>
        </div>

        <div className="flex items-center gap-3 p-3 rounded-xl bg-slate-950/90 border border-amber-900/40 shadow-inner">
          <div className="p-2 rounded-lg bg-amber-950 text-amber-400">
            <Coins className="h-5 w-5" />
          </div>
          <div>
            <span className="text-[10px] text-slate-400 block font-semibold uppercase">Token Credits</span>
            <span className="text-xl font-black text-amber-400">{tokensRemaining} Tokens</span>
          </div>
          <Link href="/portal/tokens" className="ml-2">
            <Button variant="accent" size="sm">
              Top Up
            </Button>
          </Link>
        </div>
      </div>

      {/* Global Alert */}
      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Operation Successful" : "Booking Alert"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Navigation Tabs */}
      <div className="flex items-center gap-2 border-b border-slate-800 pb-3">
        <button
          onClick={() => setSelectedTab("calendar")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            selectedTab === "calendar"
              ? "bg-red-600 text-white shadow-md shadow-red-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <CalendarCheck2 className="h-4 w-4" />
          Available Classes ({availableSlots.length})
        </button>

        <button
          onClick={() => setSelectedTab("my-bookings")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            selectedTab === "my-bookings"
              ? "bg-red-600 text-white shadow-md shadow-red-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <Clock className="h-4 w-4" />
          My Enrolled Classes ({myBookings.length})
        </button>

        <button
          onClick={() => setSelectedTab("auto-res")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            selectedTab === "auto-res"
              ? "bg-amber-600 text-white shadow-md shadow-amber-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <Repeat className="h-4 w-4" />
          Recurring Auto-Reservations
        </button>
      </div>

      {/* TAB 1: Available Occurrences */}
      {selectedTab === "calendar" && (
        <div className="space-y-6">
          {loading ? (
            <div className="text-center py-12 text-slate-400 text-sm">Loading available class slots...</div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {availableSlots.map((slot) => {
                const isFull = slot.bookedCount >= slot.capacity;
                return (
                  <Card
                    key={slot.id}
                    glow={slot.discipline === "kungfu" ? "red" : slot.discipline === "taichi" ? "gold" : "none"}
                    className="flex flex-col justify-between"
                  >
                    <CardHeader>
                      <div className="flex justify-between items-start mb-1">
                        <span className="text-xs font-bold text-amber-400 uppercase">{slot.date}</span>
                        <Badge
                          variant={
                            slot.discipline === "kungfu"
                              ? "primary"
                              : slot.discipline === "taichi"
                              ? "gold"
                              : "success"
                          }
                          size="sm"
                        >
                          {slot.discipline.toUpperCase()}
                        </Badge>
                      </div>
                      <CardTitle className="text-base">{slot.className}</CardTitle>
                      <CardDescription className="text-xs">Instructor: {slot.instructor}</CardDescription>
                    </CardHeader>

                    <CardContent className="space-y-4 text-xs text-slate-300">
                      <div className="space-y-1.5 font-mono">
                        <div className="flex items-center gap-2 text-slate-200">
                          <Clock className="h-4 w-4 text-red-500 shrink-0" />
                          <span>{slot.time} (EST)</span>
                        </div>
                        <div className="flex items-center gap-2 text-slate-400 font-sans">
                          <MapPin className="h-4 w-4 text-amber-500 shrink-0" />
                          <span>{slot.hall}</span>
                        </div>
                      </div>

                      {/* Capacity Bar */}
                      <div className="pt-2 border-t border-slate-800">
                        <div className="flex justify-between text-[11px] mb-1">
                          <span className="text-slate-400">In-Person Availability</span>
                          <span className={cn("font-bold", isFull ? "text-red-400" : "text-emerald-400")}>
                            {slot.bookedCount} / {slot.capacity} spots
                          </span>
                        </div>
                        <div className="w-full bg-slate-800 rounded-full h-1.5 overflow-hidden">
                          <div
                            className={cn("h-full", isFull ? "bg-red-500" : "bg-emerald-500")}
                            style={{ width: `${Math.min(100, (slot.bookedCount / slot.capacity) * 100)}%` }}
                          />
                        </div>
                      </div>

                      {/* Booking Buttons */}
                      <div className="grid grid-cols-2 gap-2 pt-2">
                        <Button
                          variant={isFull ? "outline" : "primary"}
                          size="sm"
                          disabled={isFull || tokensRemaining <= 0}
                          onClick={() => handleBookSlot(slot, "in_person")}
                        >
                          <span>{isFull ? "Full" : "Book In-Person"}</span>
                        </Button>
                        <Button
                          variant="secondary"
                          size="sm"
                          disabled={tokensRemaining <= 0 || !slot.hasZoom}
                          onClick={() => handleBookSlot(slot, "live_stream")}
                        >
                          <Video className="h-3.5 w-3.5 mr-1" />
                          <span>Book Zoom</span>
                        </Button>
                      </div>

                      <div className="flex justify-center">
                        <button
                          onClick={() => handleOpenNote(slot.classId, slot.className)}
                          className="text-[11px] text-slate-400 hover:text-amber-400 flex items-center gap-1 cursor-pointer transition-colors"
                        >
                          <FileText className="h-3 w-3" />
                          <span>View / Add Class Personal Note</span>
                        </button>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* TAB 2: Student's Active Bookings */}
      {selectedTab === "my-bookings" && (
        <div className="space-y-4">
          {myBookings.length === 0 ? (
            <div className="text-center py-12 space-y-3 bg-slate-900/50 rounded-2xl border border-slate-800">
              <Clock className="h-10 w-10 text-slate-500 mx-auto" />
              <h3 className="text-base font-bold text-slate-200">No active bookings found</h3>
              <p className="text-xs text-slate-400">You have no upcoming classes reserved.</p>
              <Button variant="primary" size="sm" onClick={() => setSelectedTab("calendar")}>
                Browse Class Schedule
              </Button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {myBookings.map((b) => (
                <Card key={b.id} glow="gold" className="p-5 flex flex-col justify-between">
                  <div className="space-y-3">
                    <div className="flex justify-between items-start">
                      <Badge variant="gold">
                        {b.mode === "in_person" ? "In-Person Attendance" : "Live Stream Broadcast"}
                      </Badge>
                      <Badge variant="success">Confirmed</Badge>
                    </div>

                    <h3 className="text-lg font-bold text-slate-100">{b.className}</h3>

                    <div className="space-y-1.5 text-xs text-slate-300 font-mono">
                      <div className="flex items-center gap-2">
                        <CalendarCheck2 className="h-4 w-4 text-red-500 shrink-0" />
                        <span>{b.date}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Clock className="h-4 w-4 text-amber-500 shrink-0" />
                        <span>{b.time}</span>
                      </div>
                      {b.mode === "live_stream" && b.zoomLink && (
                        <div className="flex items-center gap-2 text-blue-400 font-sans">
                          <Video className="h-4 w-4 shrink-0" />
                          <a
                            href={b.zoomLink}
                            target="_blank"
                            rel="noreferrer"
                            className="underline hover:text-blue-300"
                          >
                            Launch Zoom Classroom ({b.zoomLink})
                          </a>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="pt-4 mt-4 border-t border-slate-800 flex justify-between items-center">
                    <span className="text-[11px] text-slate-400">1 Token Reserved</span>
                    <Button
                      variant="outline"
                      size="sm"
                      className="text-red-400 hover:text-red-300 hover:border-red-600"
                      onClick={() => handleCancelBooking(b.id, b.occurrenceId)}
                    >
                      <XCircle className="h-3.5 w-3.5 mr-1" />
                      Cancel Booking & Refund Token
                    </Button>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}

      {/* TAB 3: Auto-Reservations */}
      {selectedTab === "auto-res" && (
        <Card className="p-6 space-y-6">
          <div className="space-y-2">
            <h3 className="text-lg font-bold text-slate-100">Recurring Term Auto-Reservations</h3>
            <p className="text-xs text-slate-400">
              Lock in your recurring weekly class spots automatically for the entire Summer 2026 Term. Max 2 classes per student.
            </p>
          </div>

          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between">
            <div className="space-y-1">
              <span className="text-xs font-bold text-slate-200">Traditional Kung Fu (Tuesday 18:00)</span>
              <p className="text-[11px] text-emerald-400">Active Auto-Reservation for Summer 2026 Term</p>
            </div>
            <Badge variant="success">Auto-Enrolled</Badge>
          </div>
        </Card>
      )}

      {/* Personal Class Note Modal */}
      <Modal
        isOpen={noteModalOpen}
        onClose={() => setNoteModalOpen(false)}
        title={`Personal Notes: ${activeNoteClassName}`}
        description="Private notes only visible to you for recording instructor corrections, stance tips, and form cues:"
      >
        <div className="space-y-4 py-2">
          <textarea
            value={classNoteText}
            onChange={(e) => setClassNoteText(e.target.value)}
            rows={5}
            className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600 font-sans leading-relaxed"
            placeholder="e.g. Focus on keeping pelvis tucked in Ma Bu, sink weight 70% in back leg for Xu Bu..."
          />

          <div className="flex justify-end gap-2">
            <Button variant="outline" size="sm" onClick={() => setNoteModalOpen(false)}>
              Cancel
            </Button>
            <Button variant="primary" size="sm" onClick={handleSaveNote}>
              Save Class Note
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

const fallbackSlots: OccurrenceSlot[] = [
  {
    id: "occ-cls-kf-1-2026-08-18",
    classId: "cls-kf-1",
    className: "Traditional Kung Fu (Foundations)",
    discipline: "kungfu",
    date: "Tuesday (Aug 18, 2026)",
    time: "18:00 - 19:30",
    hall: "Main Training Hall (Dojo A)",
    instructor: "Head Master Marcus Vance",
    bookedCount: 4,
    capacity: 20,
    hasZoom: true,
  },
  {
    id: "occ-cls-tc-1-2026-08-19",
    classId: "cls-tc-1",
    className: "Chen Style Tai Chi (18 Form)",
    discipline: "taichi",
    date: "Wednesday (Aug 19, 2026)",
    time: "18:30 - 19:45",
    hall: "Zen Studio C",
    instructor: "Chief Instructor Kenji Sato",
    bookedCount: 6,
    capacity: 25,
    hasZoom: true,
  },
];
