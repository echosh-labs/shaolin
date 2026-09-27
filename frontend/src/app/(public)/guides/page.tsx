"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import {
  BookOpen,
  Sparkles,
  Shirt,
  Flame,
  Award,
  CheckCircle2,
  HeartHandshake,
  ArrowRight,
} from "lucide-react";
import { apiFetch } from "@/lib/api";

interface GuideItem {
  id: string;
  title: string;
  summary: string;
  content_markdown: string;
  target_audience: string;
  category: string;
  display_order: number;
}

export default function GuidesPage() {
  const [guides, setGuides] = useState<GuideItem[]>([]);
  const [selectedGuide, setSelectedGuide] = useState<GuideItem | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadGuides() {
      try {
        const data = await apiFetch<GuideItem[]>("/guides");
        if (data && data.length > 0) {
          setGuides(data);
        } else {
          setGuides(fallbackGuides);
        }
      } catch (err) {
        console.error("Failed to load guides:", err);
        setGuides(fallbackGuides);
      } finally {
        setLoading(false);
      }
    }
    loadGuides();
  }, []);

  return (
    <div className="py-12 sm:py-16 space-y-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="primary">Student & Parent Resource Center</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Martial Arts Training Guides
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Essential information on dojo etiquette, uniform preparation, belt promotions, and training philosophies.
        </p>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {loading ? (
          <div className="text-center py-12 text-slate-400 text-sm">Loading guides...</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {guides.map((g) => (
              <Card key={g.id} glow="red" className="flex flex-col justify-between hover:border-slate-700 transition-all">
                <CardHeader>
                  <div className="flex items-center justify-between mb-2">
                    <Badge variant="primary">{g.category.replace("_", " ").toUpperCase()}</Badge>
                    <Badge variant="outline">{g.target_audience.toUpperCase()}</Badge>
                  </div>
                  <CardTitle className="text-base font-bold">{g.title}</CardTitle>
                  <CardDescription className="text-xs text-slate-400 leading-relaxed">
                    {g.summary}
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-4">
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full"
                    onClick={() => setSelectedGuide(g)}
                  >
                    <BookOpen className="h-4 w-4 mr-1.5" />
                    Read Complete Guide
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* Guide Reader Modal */}
      <Modal
        isOpen={!!selectedGuide}
        onClose={() => setSelectedGuide(null)}
        title={selectedGuide?.title || "Guide Details"}
        description={selectedGuide?.summary}
      >
        <div className="space-y-4 py-2 max-h-[65vh] overflow-y-auto pr-2">
          <div className="prose prose-invert prose-sm max-w-none text-slate-300 whitespace-pre-wrap font-sans leading-relaxed">
            {selectedGuide?.content_markdown}
          </div>

          <div className="pt-4 border-t border-slate-800 flex justify-end">
            <Button variant="primary" size="sm" onClick={() => setSelectedGuide(null)}>
              Close Guide
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

const fallbackGuides: GuideItem[] = [
  {
    id: "guide-01",
    title: "Martial Arts Belt & Rank Progression Guide",
    summary: "Complete handbook detailing curriculum milestones from Level 1 Foundations to Advanced Mastery.",
    content_markdown: "Traditional martial arts training at the Academy follows comprehensive rank testing and standardized forms.",
    target_audience: "all",
    category: "grading_exams",
    display_order: 1,
  },
  {
    id: "guide-02",
    title: "Stance Conditioning & Endurance Handbook",
    summary: "Daily training regimen to build unbreakable leg power, core stability, and mental fortitude.",
    content_markdown: "To build mastery that lasts ten thousand days, one must first lay foundations three yards deep.",
    target_audience: "student",
    category: "training_philosophy",
    display_order: 2,
  },
];
