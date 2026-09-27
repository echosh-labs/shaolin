"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import {
  Flame,
  Zap,
  Award,
  Calendar,
  ArrowRight,
  Shield,
  Heart,
  Sparkles,
  CheckCircle,
  Sword,
  Target,
} from "lucide-react";
import { cn } from "@/lib/utils";

export default function ProgramsPage() {
  const [activeTab, setActiveTab] = useState<"kungfu" | "karate" | "kobudo" | "taichi" | "qigong">("kungfu");

  const programs = {
    kungfu: {
      title: "Traditional Kung Fu (中国功夫)",
      subtitle: "Dynamic Power, Flexibility, Stance Agility & Empty Hand Forms",
      pointsBadge: "+15 Weekly Pts",
      description:
        "Traditional Kung Fu develops explosive martial conditioning, full-body coordination, and disciplined focus. Our structured curriculum teaches foundational five-step stances, classical fist forms, and weapon applications.",
      beltTracks: [
        "White Belt (Fundamentals & 5-Step Stances 五步拳)",
        "Yellow Belt (Lian Huan Quan & Agility Sets)",
        "Orange Belt (Tong Bi Quan & Staff Fundamentals)",
        "Green Belt (Xiao Hong Quan - Core Classical Form)",
        "Blue Belt (Staff Mastery & Dynamic Acrobatics)",
        "Red Belt (Advanced Classical Fist & Broadsword)",
        "Black Sash (Comprehensive Empty Hand & Weapon Mastery)",
      ],
      features: [
        "Five Fundamental Stances: Horse (Ma Bu), Bow (Gong Bu), Drop (Pu Bu), Empty (Xu Bu), Rest (Xie Bu)",
        "Classical Fist Forms & Practical Self-Defense Applications",
        "Traditional Weaponry: Waxwood Staff, Single Broadsword, Spear",
        "High-Intensity Functional Conditioning & Full Body Flexibility",
      ],
      schedulePreview: "Classes: Monday, Wednesday, Friday, Saturday (Evening & Weekend Morning slots)",
    },
    karate: {
      title: "Traditional Karate (空手道)",
      subtitle: "Linear Power, Crisp Kihon, Classical Kata & Controlled Kumite",
      pointsBadge: "+15 Weekly Pts",
      description:
        "Our Traditional Karate program focuses on maximum kinetic striking power, absolute body alignment, and rapid blocking sequences. Students progress through standardized belt katas (Heian/Pinan series), disciplined kihon drills, and respectful partner sparring.",
      beltTracks: [
        "White / Yellow Belt (Kihon Basics: Stances, Punches, Low & High Blocks)",
        "Orange Belt (Heian Shodan / Nidan Kata Sequences)",
        "Green Belt (Heian Sandan / Yondan & One-Step Kumite)",
        "Blue / Purple Belt (Heian Godan & Tekki Shodan)",
        "Brown Belt (Bassai Dai, Jion, Advanced Partner Kumite)",
        "Black Belt (Dan Rank Master Katas: Kanku Dai, Empi, Free Sparring)",
      ],
      features: [
        "Kihon Fundamentals: Gyaku-Zuki, Oi-Zuki, Age-Uke, Gedan-Barai, Mae-Geri",
        "Classical Kata Performance: Rhythm, Speed, Kime (Focus Point), and Bunkai",
        "Yakusoku & Jiyu Kumite: Safe, distance-controlled sparring drills",
        "Target Pad Striking: Developing instantaneous snap and hip torque",
      ],
      schedulePreview: "Classes: Tuesday, Thursday evenings & Saturday afternoons",
    },
    kobudo: {
      title: "Kobudo Traditional Weapons (古武道)",
      subtitle: "Okinawan Weaponry: Bo Staff, Sai, Tonfa, Nunchaku & Kama",
      pointsBadge: "+15 Weekly Pts",
      description:
        "Kobudo is the classical art of martial weaponry. Students learn to extend their biomechanical movement through traditional implements, cultivating wrist strength, grip endurance, spatial awareness, and precision weapon katas.",
      beltTracks: [
        "Novice Weaponry: Rokushaku Bo (6ft Staff Fundamentals & Striking Angles)",
        "Intermediate Level: Sai (Twin Iron Tridents: Trapping & Thrusting)",
        "Intermediate Level: Tonfa (Baton Spinning, Arm Shielding & Snapping)",
        "Advanced Level: Nunchaku (Speed Flow, Diagonal Switches, Wrist Rolls)",
        "Senior Weaponry: Kama (Sickle Forms) & Multi-Weapon Bunkai Combatives",
      ],
      features: [
        "Rokushaku Bo: Fundamental 8-angle strikes, thrusts, and overhead blocks",
        "Sai & Tonfa: Dual-weapon coordination, joint reinforcement, and parries",
        "Controlled Speed Flow: Smooth transitions, finger rolls, and stance shifts",
        "Weapon Safety & Maintenance: Preserving hardwood and steel equipment",
      ],
      schedulePreview: "Classes: Wednesday evenings & Sunday mornings",
    },
    taichi: {
      title: "Tai Chi & Internal Arts (太极拳)",
      subtitle: "Internal Energy, Silk Reeling, Flowing Grace & Meditative Strength",
      pointsBadge: "+10 Weekly Pts",
      description:
        "Combining the internal martial arts of Chen and Yang styles, our Tai Chi curriculum cultivates rooted balance, joint decompression, neuromuscular coordination, and internal vitality.",
      beltTracks: [
        "Level 1: Foundation Stances, Stepping & Silk Reeling (Chan Si Gong)",
        "Level 2: Chen Style 18 Essential Form Set",
        "Level 3: Yang Style 24 Form",
        "Level 4: Chen Style 56 Competition Routine",
        "Level 5: Traditional Tai Chi Straight Sword (Taiji Jian)",
      ],
      features: [
        "Spiral Silk Reeling (Chan Si Gong) for fascia and joint mobility",
        "Rooted postural alignment, grounding, and breath synchronization",
        "Low-impact resistance training that restores knees, hips, and spine",
        "Internal power generation (Fa Jin) and pushing hands (Tui Shou)",
      ],
      schedulePreview: "Classes: Tuesday, Thursday evenings & Saturday mornings",
    },
    qigong: {
      title: "Qigong & Breath Harmony (气功)",
      subtitle: "Ancient Health Cultivation, Eight Brocades & Mindful Vitality",
      pointsBadge: "+8 Weekly Pts",
      description:
        "Qigong focuses on diaphragmatic breathing, posture alignment, and internal energy cultivation. Proven over centuries to reduce stress, boost immune resilience, and cultivate mental stillness.",
      beltTracks: [
        "Module 1: Eight Pieces of Brocade (Ba Duan Jin 八段锦)",
        "Module 2: Muscle-Tendon Change Classic (Yi Jin Jing 易筋经)",
        "Module 3: Six Healing Sounds (Liu Zi Jue 六字诀)",
        "Module 4: Standing Posture Meditation (Zhan Zhuang 站桩)",
      ],
      features: [
        "Ba Duan Jin (8 Brocades): Restores internal organ vitality",
        "Yi Jin Jing (Muscle-Tendon Classic): Strengthens tendons and posture",
        "Diaphragmatic breathing and emotional de-stressing",
        "Gentle, accessible practice suitable for all ages and physical conditions",
      ],
      schedulePreview: "Classes: Tuesday, Thursday evenings & Saturday mid-day",
    },
  };

  const active = programs[activeTab];

  return (
    <div className="py-12 sm:py-16 space-y-16">
      {/* Header Banner */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="primary">5 Core Disciplines</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Martial Arts Training Programs
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto leading-relaxed">
          From dynamic physical combat arts to meditative internal practices. Select a discipline below to explore syllabus tracks and requirements.
        </p>

        {/* Tab Switcher */}
        <div className="flex flex-wrap justify-center items-center gap-2 pt-6">
          <button
            onClick={() => setActiveTab("kungfu")}
            className={cn(
              "px-4 py-2 rounded-xl font-bold text-xs uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
              activeTab === "kungfu"
                ? "bg-red-600 text-white shadow-lg shadow-red-950/60 scale-105 border border-red-500"
                : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
            )}
          >
            <Flame className="h-4 w-4" />
            Kung Fu
          </button>
          <button
            onClick={() => setActiveTab("karate")}
            className={cn(
              "px-4 py-2 rounded-xl font-bold text-xs uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
              activeTab === "karate"
                ? "bg-amber-600 text-white shadow-lg shadow-amber-950/60 scale-105 border border-amber-500"
                : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
            )}
          >
            <Target className="h-4 w-4" />
            Karate
          </button>
          <button
            onClick={() => setActiveTab("kobudo")}
            className={cn(
              "px-4 py-2 rounded-xl font-bold text-xs uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
              activeTab === "kobudo"
                ? "bg-slate-700 text-white shadow-lg scale-105 border border-slate-600"
                : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
            )}
          >
            <Sword className="h-4 w-4" />
            Kobudo
          </button>
          <button
            onClick={() => setActiveTab("taichi")}
            className={cn(
              "px-4 py-2 rounded-xl font-bold text-xs uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
              activeTab === "taichi"
                ? "bg-amber-600 text-white shadow-lg shadow-amber-950/60 scale-105 border border-amber-500"
                : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
            )}
          >
            <Sparkles className="h-4 w-4" />
            Tai Chi
          </button>
          <button
            onClick={() => setActiveTab("qigong")}
            className={cn(
              "px-4 py-2 rounded-xl font-bold text-xs uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
              activeTab === "qigong"
                ? "bg-emerald-600 text-white shadow-lg shadow-emerald-950/60 scale-105 border border-emerald-500"
                : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
            )}
          >
            <Heart className="h-4 w-4" />
            Qigong
          </button>
        </div>
      </div>

      {/* Program Detail Overview */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Info Col */}
          <div className="lg:col-span-2 space-y-6">
            <Card glow={activeTab === "kungfu" ? "red" : activeTab === "karate" || activeTab === "taichi" ? "gold" : "none"}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant={activeTab === "kungfu" ? "primary" : activeTab === "qigong" ? "success" : "gold"}>
                    {active.pointsBadge}
                  </Badge>
                  <span className="text-xs text-slate-400 font-mono">In-Person & Virtual Available</span>
                </div>
                <CardTitle className="text-2xl pt-2">{active.title}</CardTitle>
                <CardDescription className="text-base text-slate-300">
                  {active.subtitle}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6 text-sm text-slate-300 leading-relaxed">
                <p>{active.description}</p>

                <div className="space-y-3 pt-2">
                  <h4 className="font-bold text-slate-100 text-sm uppercase tracking-wider flex items-center gap-2">
                    <Zap className="h-4 w-4 text-red-500" />
                    Key Training Components:
                  </h4>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs text-slate-300">
                    {active.features.map((feat, idx) => (
                      <div
                        key={idx}
                        className="flex items-start gap-2.5 p-3 rounded-xl bg-slate-950/80 border border-slate-800/80"
                      >
                        <CheckCircle className="h-4 w-4 text-red-400 shrink-0 mt-0.5" />
                        <span>{feat}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between text-xs">
                  <div className="flex items-center gap-2 text-slate-300">
                    <Calendar className="h-4 w-4 text-amber-500" />
                    <span>{active.schedulePreview}</span>
                  </div>
                  <Link href="/schedule">
                    <Button variant="ghost" size="sm" className="text-red-400">
                      Full Schedule &rarr;
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Syllabus Belt Track Card */}
          <div className="space-y-6">
            <Card glow="gold">
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Award className="h-5 w-5 text-amber-500" />
                  Progression Syllabus
                </CardTitle>
                <CardDescription>
                  Structured rank requirements from novice fundamentals to advanced mastery.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ol className="relative border-l border-slate-800 ml-3 space-y-4 text-xs">
                  {active.beltTracks.map((belt, idx) => (
                    <li key={idx} className="mb-2 ml-4">
                      <div className="absolute w-2.5 h-2.5 bg-amber-500 rounded-full -left-1.5 border border-slate-900 mt-1.5" />
                      <p className="font-semibold text-slate-200">{belt}</p>
                    </li>
                  ))}
                </ol>

                <div className="mt-6 pt-4 border-t border-slate-800">
                  <Link href="/register">
                    <Button variant="primary" size="md" className="w-full">
                      Start Your Training
                      <ArrowRight className="h-4 w-4" />
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
