import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import {
  Flame,
  Calendar,
  Award,
  Coins,
  ArrowRight,
  ShieldCheck,
  Zap,
  CheckCircle2,
  Sword,
} from "lucide-react";

export default function HomePage() {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden pt-12 pb-24 lg:pt-20 lg:pb-32 border-b border-slate-800/80 bg-gradient-to-b from-[#090d16] via-slate-950 to-[#020617]">
        {/* Background glow ambient effects */}
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-red-600/10 rounded-full blur-3xl pointer-events-none" />
        <div className="absolute top-1/3 right-10 w-[400px] h-[400px] bg-amber-600/10 rounded-full blur-3xl pointer-events-none" />

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="max-w-3xl space-y-6">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-red-950/80 border border-red-800/60 text-red-300 text-xs font-semibold uppercase tracking-wider">
              <Flame className="h-3.5 w-3.5 text-red-500" />
              <span>Traditional Martial Arts & Fitness Academy</span>
            </div>

            <h1 className="text-4xl sm:text-6xl lg:text-7xl font-black tracking-tight text-slate-100 uppercase leading-[1.05]">
              Master Traditional{" "}
              <span className="bg-gradient-to-r from-red-500 via-amber-400 to-red-600 bg-clip-text text-transparent">
                Martial Arts
              </span>{" "}
              & Internal Flow
            </h1>

            <p className="text-base sm:text-lg text-slate-300 leading-relaxed font-light">
              Experience the disciplines, athletic conditioning, and mental focus of authentic martial arts. Comprehensive training in Kung Fu, Karate, Kobudo, Tai Chi, and Qigong with in-person dojo classes and live interactive streams.
            </p>

            <div className="flex flex-wrap gap-4 pt-4">
              <Link href="/schedule">
                <Button variant="primary" size="lg" className="shadow-red-900/50 shadow-xl">
                  <Calendar className="h-5 w-5" />
                  View Class Schedule
                </Button>
              </Link>

              <Link href="/register">
                <Button variant="outline" size="lg" className="border-slate-700 hover:border-red-600/60">
                  Join as New Student
                  <ArrowRight className="h-5 w-5" />
                </Button>
              </Link>
            </div>

            {/* Live Badges */}
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 pt-8 border-t border-slate-800/60 text-xs text-slate-400">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4 text-red-500 shrink-0" />
                <span>5 Core Disciplines</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4 text-amber-500 shrink-0" />
                <span>In-Person & Virtual Hybrid</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4 text-emerald-500 shrink-0" />
                <span>Flexible Token Ledgers</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Program Tracks Showcase */}
      <section className="py-20 bg-slate-950/60 border-b border-slate-800/80">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-2xl mx-auto mb-16 space-y-3">
            <Badge variant="primary">Core Disciplines</Badge>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-100 uppercase tracking-wide">
              Martial Arts Training Programs
            </h2>
            <p className="text-sm text-slate-400">
              Comprehensive curriculums designed for beginners, intermediate practitioners, and advanced martial artists.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-6">
            {/* Adult Kung Fu */}
            <Card glow="red" className="flex flex-col justify-between hover:-translate-y-1 transition-transform">
              <CardHeader>
                <div className="flex justify-between items-start mb-2">
                  <Badge variant="primary">Kung Fu 功夫</Badge>
                  <span className="text-xs text-red-400 font-mono">+15 Pts</span>
                </div>
                <CardTitle className="text-base">Traditional Kung Fu</CardTitle>
                <CardDescription className="text-xs">
                  Fist forms, low stances, kicking agility, conditioning, and martial applications.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2 text-[11px] text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Zap className="h-3.5 w-3.5 text-red-500" />
                  <span>Stance & Agility Drills</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5 text-amber-500" />
                  <span>Belt Rank Curriculum</span>
                </div>
              </CardContent>
            </Card>

            {/* Karate */}
            <Card glow="gold" className="flex flex-col justify-between hover:-translate-y-1 transition-transform">
              <CardHeader>
                <div className="flex justify-between items-start mb-2">
                  <Badge variant="gold">Karate 空手道</Badge>
                  <span className="text-xs text-amber-400 font-mono">+15 Pts</span>
                </div>
                <CardTitle className="text-base">Traditional Karate</CardTitle>
                <CardDescription className="text-xs">
                  Classical kata, kihon fundamentals, linear striking power, and controlled kumite sparring.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2 text-[11px] text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Zap className="h-3.5 w-3.5 text-amber-500" />
                  <span>Kihon Strikes & Blocks</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5 text-amber-500" />
                  <span>Kata Sequence Mastery</span>
                </div>
              </CardContent>
            </Card>

            {/* Kobudo */}
            <Card className="flex flex-col justify-between hover:-translate-y-1 transition-transform border-slate-800">
              <CardHeader>
                <div className="flex justify-between items-start mb-2">
                  <Badge variant="default">Kobudo 古武道</Badge>
                  <span className="text-xs text-slate-300 font-mono">+15 Pts</span>
                </div>
                <CardTitle className="text-base">Kobudo Weapons</CardTitle>
                <CardDescription className="text-xs">
                  Traditional weaponry: Bo Staff, Sai, Tonfa, Nunchaku, and Kama forms.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2 text-[11px] text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Sword className="h-3.5 w-3.5 text-slate-300" />
                  <span>Bo Staff & Sai Drills</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5 text-amber-500" />
                  <span>Weapon Form Katas</span>
                </div>
              </CardContent>
            </Card>

            {/* Tai Chi */}
            <Card glow="gold" className="flex flex-col justify-between hover:-translate-y-1 transition-transform">
              <CardHeader>
                <div className="flex justify-between items-start mb-2">
                  <Badge variant="gold">Tai Chi 太极</Badge>
                  <span className="text-xs text-amber-400 font-mono">+10 Pts</span>
                </div>
                <CardTitle className="text-base">Tai Chi Flow</CardTitle>
                <CardDescription className="text-xs">
                  Chen and Yang forms, silk reeling exercises, structural balance, and joint health.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2 text-[11px] text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Zap className="h-3.5 w-3.5 text-amber-500" />
                  <span>Silk Reeling & Balance</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5 text-amber-500" />
                  <span>Form Sets 18 & 56</span>
                </div>
              </CardContent>
            </Card>

            {/* Qi Gong */}
            <Card className="flex flex-col justify-between hover:-translate-y-1 transition-transform border-slate-800">
              <CardHeader>
                <div className="flex justify-between items-start mb-2">
                  <Badge variant="success">Qigong 气功</Badge>
                  <span className="text-xs text-emerald-400 font-mono">+8 Pts</span>
                </div>
                <CardTitle className="text-base">Qigong & Breath</CardTitle>
                <CardDescription className="text-xs">
                  Eight Brocades (Ba Duan Jin), standing posture meditation, and breath harmony.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2 text-[11px] text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Zap className="h-3.5 w-3.5 text-emerald-500" />
                  <span>Ba Duan Jin Regimen</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5 text-emerald-500" />
                  <span>Stress Reduction & Focus</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Modern Dojo Platform Features */}
      <section className="py-20 bg-[#090d16]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div className="space-y-6">
              <Badge variant="gold">Student Platform</Badge>
              <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-100 uppercase tracking-tight">
                Modern Technology Meets Martial Heritage
              </h2>
              <p className="text-sm text-slate-300 leading-relaxed">
                Our digital student portal gives martial arts practitioners full transparency into their training journey:
              </p>

              <div className="space-y-4 text-sm">
                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-red-950/80 text-red-400 border border-red-800/40">
                    <Coins className="h-5 w-5" />
                  </div>
                  <div>
                    <h4 className="font-semibold text-slate-200">Flexible Token Booking</h4>
                    <p className="text-xs text-slate-400">
                      Never lose money on missed classes. Use tokens across all Kung Fu, Karate, Kobudo, Tai Chi, and Qigong occurrences.
                    </p>
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-amber-950/80 text-amber-400 border border-amber-800/40">
                    <Award className="h-5 w-5" />
                  </div>
                  <div>
                    <h4 className="font-semibold text-slate-200">Gamified Attendance & Tier Ranks</h4>
                    <p className="text-xs text-slate-400">
                      Earn weekly attendance points, unlock consistency bonuses, and climb from Novice Disciple to Senior Martial Artist.
                    </p>
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-slate-800 text-slate-300 border border-slate-700">
                    <ShieldCheck className="h-5 w-5" />
                  </div>
                  <div>
                    <h4 className="font-semibold text-slate-200">Grading Exam Syllabus</h4>
                    <p className="text-xs text-slate-400">
                      Track stance mastery, review instructor evaluation feedback, and advance through official rank requirements.
                    </p>
                  </div>
                </div>
              </div>

              <div className="pt-2">
                <Link href="/portal">
                  <Button variant="accent">
                    Enter Student Command Portal
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>

            {/* Interactive Preview Card */}
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-red-600/20 to-amber-600/20 rounded-3xl blur-2xl -z-10" />
              <Card glow="gold" className="p-6 space-y-6">
                <div className="flex items-center justify-between pb-4 border-b border-slate-800">
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-xl bg-red-600 flex items-center justify-center font-bold text-white shadow-md">
                      武道
                    </div>
                    <div>
                      <h4 className="text-sm font-bold text-slate-100">Live Dojo System Status</h4>
                      <p className="text-xs text-emerald-400 font-mono flex items-center gap-1">
                        <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
                        Backend Connected (8081)
                      </p>
                    </div>
                  </div>
                  <Badge variant="gold">Summer 2026 Term</Badge>
                </div>

                <div className="grid grid-cols-2 gap-3 text-xs">
                  <div className="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                    <span className="text-slate-400 block mb-1">Active Disciplines</span>
                    <span className="text-slate-100 font-bold text-xs">Kung Fu, Karate, Kobudo, Tai Chi, Qigong</span>
                  </div>
                  <div className="p-3 rounded-xl bg-slate-950/80 border border-slate-800">
                    <span className="text-slate-400 block mb-1">Training Modes</span>
                    <span className="text-amber-400 font-bold text-sm">In-Person & Virtual</span>
                  </div>
                </div>

                <div className="p-4 rounded-xl bg-red-950/30 border border-red-800/40 text-xs text-slate-300">
                  <p className="font-semibold text-red-300 mb-1">Ready for Class Registration</p>
                  <p className="text-slate-400 leading-relaxed">
                    Create your student profile today to reserve your spots in Toronto training dojos or join live virtual streams.
                  </p>
                </div>
              </Card>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
