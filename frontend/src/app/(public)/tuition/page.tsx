"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import {
  Coins,
  Check,
  Zap,
  ShieldCheck,
  Users,
  Sparkles,
  ArrowRight,
  Info,
  Calendar,
} from "lucide-react";
import { cn, formatCurrency } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

export default function TuitionPage() {
  const [ageGroup, setAgeGroup] = useState<"adult" | "youth">("adult");
  const [packages, setPackages] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadTuitionPackages() {
      try {
        const data = await apiFetch<any[]>("/tokens/packages");
        if (data && data.length > 0) {
          const mapped = data.map((p: any) => ({
            id: p.id,
            name: p.name,
            tokens: p.tokens_count,
            adultCents: p.price_cents,
            youthCents: Math.round(p.price_cents * 0.8),
            description: p.description || `${p.tokens_count} class tokens.`,
            popular: p.tokens_count === 14,
            isUnlimited: p.tokens_count > 100,
          }));
          setPackages(mapped);
        } else {
          setPackages(fallbackPackages);
        }
      } catch (err) {
        console.error("Failed to load tuition packages:", err);
        setPackages(fallbackPackages);
      } finally {
        setLoading(false);
      }
    }
    loadTuitionPackages();
  }, []);

  return (
    <div className="py-12 sm:py-16 space-y-16">
      {/* Page Title */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="gold">Tuition & Token Packages</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Flexible Martial Arts Training
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Our token-based system gives you complete freedom. Tokens never expire during the term and can be used interchangeably across Kung Fu, Tai Chi, and Qi Gong.
        </p>

        {/* Age Toggle */}
        <div className="flex justify-center items-center gap-2 pt-4">
          <div className="p-1 rounded-xl bg-slate-900 border border-slate-800 flex items-center gap-1">
            <button
              onClick={() => setAgeGroup("adult")}
              className={cn(
                "px-4 py-2 rounded-lg text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
                ageGroup === "adult"
                  ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                  : "text-slate-400 hover:text-slate-200"
              )}
            >
              Adult Programs (14+)
            </button>
            <button
              onClick={() => setAgeGroup("youth")}
              className={cn(
                "px-4 py-2 rounded-lg text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
                ageGroup === "youth"
                  ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                  : "text-slate-400 hover:text-slate-200"
              )}
            >
              Youth / Kids (5–13)
            </button>
          </div>
        </div>
      </div>

      {/* Package Pricing Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {packages.slice(0, 6).map((pkg) => {
            const price = ageGroup === "adult" ? pkg.adultCents : pkg.youthCents;
            const pricePerClass = pkg.isUnlimited
              ? null
              : Math.round(price / pkg.tokens);

            return (
              <Card
                key={pkg.id}
                glow={pkg.popular ? "gold" : pkg.isUnlimited ? "red" : "none"}
                className={cn(
                  "flex flex-col justify-between relative",
                  pkg.popular && "border-amber-500/60 bg-gradient-to-b from-slate-900 via-slate-900 to-slate-950"
                )}
              >
                {pkg.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                    <Badge variant="gold" className="shadow-lg shadow-amber-950/60">
                      Most Popular
                    </Badge>
                  </div>
                )}

                <CardHeader className="space-y-3">
                  <div className="flex justify-between items-start">
                    <span className="text-xs text-amber-400 font-bold uppercase tracking-wider">
                      {pkg.name}
                    </span>
                    <Badge variant="primary" size="sm">
                      {pkg.isUnlimited ? "All Year" : `${pkg.tokens} Tokens`}
                    </Badge>
                  </div>

                  <div>
                    <div className="text-4xl font-black text-slate-100 font-mono">
                      {formatCurrency(price)}
                    </div>
                    <span className="text-xs text-slate-400 block mt-1">
                      {pkg.isUnlimited
                        ? "Unlimited classes for 3 terms (1 Year)"
                        : `~${formatCurrency(pricePerClass || 0)} per class session (+ HST)`}
                    </span>
                  </div>

                  <CardDescription className="text-xs text-slate-300 leading-relaxed">
                    {pkg.description}
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-6">
                  <div className="space-y-2 text-xs text-slate-300">
                    <div className="flex items-center gap-2">
                      <Check className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>Valid across Kung Fu, Tai Chi & Qi Gong</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Check className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>Includes in-person & Zoom livestream</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Check className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>Instant token credit upon booking cancellation</span>
                    </div>
                  </div>

                  <Link href="/portal/tokens" className="block w-full">
                    <Button
                      variant={pkg.popular ? "accent" : "primary"}
                      size="md"
                      className="w-full"
                    >
                      <span>Select Package</span>
                      <ArrowRight className="h-4 w-4 ml-1" />
                    </Button>
                  </Link>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>
    </div>
  );
}

const fallbackPackages = [
  {
    id: "pkg-1",
    name: "Single Class Drop-in",
    tokens: 1,
    adultCents: 2500,
    youthCents: 2000,
    description: "Great for visitors, trial sessions, or occasional schedule top-ups.",
    popular: false,
  },
  {
    id: "pkg-7",
    name: "7 Tokens Flex Pack",
    tokens: 7,
    adultCents: 15400,
    youthCents: 12600,
    description: "Flexible training pack (~$22/class). Valid across all disciplines.",
    popular: false,
  },
  {
    id: "pkg-14",
    name: "14 Tokens Term Pack",
    tokens: 14,
    adultCents: 29400,
    youthCents: 23800,
    description: "Recommended for 1 class per week over a 14-week term (~$21/class).",
    popular: true,
  },
];
