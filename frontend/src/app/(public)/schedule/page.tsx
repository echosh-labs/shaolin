"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import {
  Calendar as CalendarIcon,
  Clock,
  MapPin,
  User,
  Users,
  Video,
  CheckCircle2,
  AlertCircle,
  Sparkles,
  ArrowRight,
  Filter,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface ScheduleItem {
  id: string;
  name: string;
  discipline: "kungfu" | "karate" | "kobudo" | "taichi" | "qigong";
  day: number; // 0 = Sun, 1 = Mon, ..., 6 = Sat
  dayName: string;
  startTime: string;
  endTime: string;
  hall: string;
  instructor: string;
  capacity: number;
  bookedCount: number;
  hasZoom: boolean;
}

export default function SchedulePage() {
  const [selectedDay, setSelectedDay] = useState<number | "all">("all");
  const [selectedDiscipline, setSelectedDiscipline] = useState<string>("all");
  const [selectedMode, setSelectedMode] = useState<"all" | "in_person" | "zoom">("all");

  const [scheduleData, setScheduleData] = useState<ScheduleItem[]>(defaultScheduleData);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadSchedule() {
      try {
        const [classesRes, occsRes] = await Promise.all([
          apiFetch<any[]>("/classes").catch(() => []),
          apiFetch<any[]>("/occurrences").catch(() => []),
        ]);

        if (classesRes && classesRes.length > 0) {
          const dayNames = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
          
          const items: ScheduleItem[] = classesRes.map((c: any) => {
            const classOccs = occsRes?.filter((o: any) => o.class_id === c.id) || [];
            const bookedCount = classOccs.length > 0 
              ? Math.max(...classOccs.map((o: any) => o.booked_count || 0), 0)
              : 0;

            const disc = c.name.toLowerCase().includes("karate") 
              ? "karate" 
              : c.name.toLowerCase().includes("kobudo") || c.name.toLowerCase().includes("staff") || c.name.toLowerCase().includes("weapon")
              ? "kobudo"
              : c.name.toLowerCase().includes("tai chi") 
              ? "taichi" 
              : c.name.toLowerCase().includes("qigong") || c.name.toLowerCase().includes("baduanjin") 
              ? "qigong" 
              : "kungfu";

            return {
              id: c.id,
              name: c.name,
              discipline: disc,
              day: c.day_of_week,
              dayName: dayNames[c.day_of_week] || "Monday",
              startTime: c.start_time,
              endTime: c.end_time,
              hall: "Main Training Hall (Dojo A)",
              instructor: "Chief Instructor Kenji Sato",
              capacity: c.capacity || 20,
              bookedCount,
              hasZoom: !!c.zoom_link,
            };
          });
          setScheduleData(items);
        } else {
          setScheduleData(defaultScheduleData);
        }
      } catch (err) {
        console.error("Failed to load live schedule:", err);
      } finally {
        setLoading(false);
      }
    }
    loadSchedule();
  }, []);

  const filteredSchedule = scheduleData.filter((item) => {
    if (selectedDay !== "all" && item.day !== selectedDay) return false;
    if (selectedDiscipline !== "all" && item.discipline !== selectedDiscipline) return false;
    if (selectedMode === "zoom" && !item.hasZoom) return false;
    return true;
  });

  const daysFilter = [
    { label: "All Days", value: "all" },
    { label: "Mon", value: 1 },
    { label: "Tue", value: 2 },
    { label: "Wed", value: 3 },
    { label: "Thu", value: 4 },
    { label: "Fri", value: 5 },
    { label: "Sat", value: 6 },
    { label: "Sun", value: 0 },
  ];

  const disciplineFilter = [
    { label: "All Disciplines", value: "all" },
    { label: "Kung Fu", value: "kungfu" },
    { label: "Karate", value: "karate" },
    { label: "Kobudo", value: "kobudo" },
    { label: "Tai Chi", value: "taichi" },
    { label: "Qigong", value: "qigong" },
  ];

  return (
    <div className="py-12 sm:py-16 space-y-12">
      {/* Header */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="primary">Summer 2026 Term Active</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Weekly Class Schedule
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Explore our weekly in-person and live-stream schedule across all 5 disciplines. Reserve spots dynamically using your student tokens.
        </p>

        {/* Live Status Alert */}
        <div className="max-w-xl mx-auto p-3 rounded-xl bg-amber-950/40 border border-amber-800/40 text-xs text-amber-300 flex items-center justify-center gap-2">
          <Sparkles className="h-4 w-4 text-amber-400 shrink-0" />
          <span>Real-time in-person capacity counters updated directly from database occurrences.</span>
        </div>
      </div>

      {/* Filters Bar */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="p-4 rounded-2xl bg-slate-900/90 border border-slate-800 shadow-xl space-y-4">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            {/* Days Tabs */}
            <div className="flex items-center gap-1.5 overflow-x-auto w-full md:w-auto pb-2 md:pb-0">
              {daysFilter.map((d) => (
                <button
                  key={d.label}
                  onClick={() => setSelectedDay(d.value as any)}
                  className={cn(
                    "px-3 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider transition-all cursor-pointer whitespace-nowrap",
                    selectedDay === d.value
                      ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                      : "bg-slate-800/60 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                  )}
                >
                  {d.label}
                </button>
              ))}
            </div>

            {/* Discipline Dropdown */}
            <div className="flex items-center gap-2 w-full md:w-auto overflow-x-auto">
              <div className="flex items-center gap-1 bg-slate-950 p-1 rounded-xl border border-slate-800">
                {disciplineFilter.map((disc) => (
                  <button
                    key={disc.value}
                    onClick={() => setSelectedDiscipline(disc.value)}
                    className={cn(
                      "px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap",
                      selectedDiscipline === disc.value
                        ? "bg-slate-800 text-slate-100"
                        : "text-slate-400 hover:text-slate-300"
                    )}
                  >
                    {disc.label}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Mode Switcher */}
          <div className="flex items-center justify-between border-t border-slate-800/80 pt-3 text-xs text-slate-400">
            <div className="flex items-center gap-2">
              <Filter className="h-3.5 w-3.5" />
              <span>Training Mode:</span>
              <div className="flex items-center gap-1">
                {(["all", "in_person", "zoom"] as const).map((mode) => (
                  <button
                    key={mode}
                    onClick={() => setSelectedMode(mode)}
                    className={cn(
                      "px-2.5 py-0.5 rounded-md uppercase font-medium transition-colors cursor-pointer",
                      selectedMode === mode
                        ? "bg-amber-950 text-amber-400 border border-amber-800/60"
                        : "hover:text-slate-200"
                    )}
                  >
                    {mode === "all" ? "All Modes" : mode === "in_person" ? "Dojo Only" : "Zoom Only"}
                  </button>
                ))}
              </div>
            </div>

            <span className="font-mono text-slate-500">
              Showing {filteredSchedule.length} classes
            </span>
          </div>
        </div>
      </div>

      {/* Schedule Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {filteredSchedule.length === 0 ? (
          <div className="text-center py-16 bg-slate-900/40 rounded-3xl border border-slate-800 space-y-3">
            <CalendarIcon className="h-10 w-10 text-slate-600 mx-auto" />
            <h3 className="text-lg font-bold text-slate-300">No classes match selected filters</h3>
            <p className="text-xs text-slate-500">Try selecting another day or discipline filter.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredSchedule.map((item) => {
              const spotsLeft = item.capacity - item.bookedCount;
              const isFull = spotsLeft <= 0;
              const isAlmostFull = spotsLeft > 0 && spotsLeft <= 3;

              return (
                <Card
                  key={item.id}
                  glow={item.discipline === "kungfu" ? "red" : item.discipline === "karate" || item.discipline === "taichi" ? "gold" : "none"}
                  className="flex flex-col justify-between hover:border-slate-700 transition-all group"
                >
                  <div className="p-6 space-y-4">
                    {/* Header Top Tag */}
                    <div className="flex items-center justify-between">
                      <Badge
                        variant={
                          item.discipline === "kungfu"
                            ? "primary"
                            : item.discipline === "karate" || item.discipline === "taichi"
                            ? "gold"
                            : item.discipline === "qigong"
                            ? "success"
                            : "default"
                        }
                      >
                        {item.discipline.toUpperCase()}
                      </Badge>

                      <div className="flex items-center gap-1.5 text-xs text-slate-400 font-mono">
                        <Clock className="h-3.5 w-3.5 text-slate-500" />
                        <span>
                          {item.startTime} - {item.endTime}
                        </span>
                      </div>
                    </div>

                    {/* Title */}
                    <div>
                      <h3 className="text-lg font-bold text-slate-100 group-hover:text-amber-400 transition-colors">
                        {item.name}
                      </h3>
                      <p className="text-xs text-amber-500 font-semibold mt-0.5">{item.dayName}</p>
                    </div>

                    {/* Hall & Instructor */}
                    <div className="space-y-2 text-xs text-slate-400 pt-2 border-t border-slate-800/60">
                      <div className="flex items-center gap-2">
                        <MapPin className="h-4 w-4 text-red-500 shrink-0" />
                        <span className="truncate">{item.hall}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-amber-500 shrink-0" />
                        <span>{item.instructor}</span>
                      </div>
                    </div>

                    {/* Capacity & Zoom Flags */}
                    <div className="flex items-center justify-between pt-2 text-xs">
                      <div className="flex items-center gap-1.5">
                        <Users className="h-3.5 w-3.5 text-slate-500" />
                        <span
                          className={cn(
                            "font-semibold",
                            isFull
                              ? "text-red-400"
                              : isAlmostFull
                              ? "text-amber-400"
                              : "text-emerald-400"
                          )}
                        >
                          {isFull ? "Full (Waitlist)" : `${spotsLeft} spots available`}
                        </span>
                      </div>

                      {item.hasZoom && (
                        <div className="flex items-center gap-1 text-sky-400 font-medium">
                          <Video className="h-3.5 w-3.5" />
                          <span>Zoom Live</span>
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Footer Action */}
                  <div className="p-4 bg-slate-950/60 border-t border-slate-800/80 flex items-center justify-between">
                    <span className="text-[11px] text-slate-400 font-mono">1 Token Deduction</span>
                    <Link href="/portal/bookings">
                      <Button variant="primary" size="sm" className="cursor-pointer">
                        Book Class
                        <ArrowRight className="h-3.5 w-3.5" />
                      </Button>
                    </Link>
                  </div>
                </Card>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}

const defaultScheduleData: ScheduleItem[] = [
  {
    id: "cls-kf-1",
    name: "Traditional Kung Fu (Foundations)",
    discipline: "kungfu",
    day: 2,
    dayName: "Tuesday",
    startTime: "18:00",
    endTime: "19:30",
    hall: "Main Training Hall (Dojo A)",
    instructor: "Head Master Marcus Vance",
    capacity: 20,
    bookedCount: 4,
    hasZoom: true,
  },
  {
    id: "cls-karate-1",
    name: "Traditional Karate (Kata & Kihon)",
    discipline: "karate",
    day: 2,
    dayName: "Tuesday",
    startTime: "19:30",
    endTime: "21:00",
    hall: "Martial Arts Studio B",
    instructor: "Chief Instructor Kenji Sato",
    capacity: 20,
    bookedCount: 5,
    hasZoom: true,
  },
  {
    id: "cls-kobudo-1",
    name: "Kobudo Weapons (Bo Staff & Sai)",
    discipline: "kobudo",
    day: 3,
    dayName: "Wednesday",
    startTime: "18:00",
    endTime: "19:30",
    hall: "Martial Arts Studio B",
    instructor: "Chief Instructor Kenji Sato",
    capacity: 16,
    bookedCount: 6,
    hasZoom: true,
  },
  {
    id: "cls-tc-1",
    name: "Chen Style Tai Chi (18 Form)",
    discipline: "taichi",
    day: 3,
    dayName: "Wednesday",
    startTime: "19:30",
    endTime: "20:45",
    hall: "Zen Studio C",
    instructor: "Master Elena Rostova",
    capacity: 25,
    bookedCount: 8,
    hasZoom: true,
  },
  {
    id: "cls-qg-1",
    name: "Baduanjin & Yi Jin Jing Qigong Flow",
    discipline: "qigong",
    day: 6,
    dayName: "Saturday",
    startTime: "09:30",
    endTime: "10:45",
    hall: "Zen Studio C",
    instructor: "Master Elena Rostova",
    capacity: 30,
    bookedCount: 12,
    hasZoom: true,
  },
  {
    id: "cls-kf-youth",
    name: "Youth Martial Arts (Dragons)",
    discipline: "kungfu",
    day: 6,
    dayName: "Saturday",
    startTime: "11:00",
    endTime: "12:15",
    hall: "Main Training Hall (Dojo A)",
    instructor: "Chief Instructor Kenji Sato",
    capacity: 18,
    bookedCount: 6,
    hasZoom: true,
  },
];
