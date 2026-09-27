"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { ChevronDown, HelpCircle, ArrowRight } from "lucide-react";
import { cn } from "@/lib/utils";

interface FAQItem {
  q: string;
  a: string;
  category: "General" | "Tokens & Booking" | "Grading & Exams";
}

export default function FAQPage() {
  const [openIndex, setOpenIndex] = useState<number | null>(0);
  const [activeCategory, setActiveCategory] = useState<string>("All");

  const faqs: FAQItem[] = [
    {
      category: "General",
      q: "Do I need prior martial arts experience to start?",
      a: "No! All of our beginner classes (Traditional Kung Fu Level 1, Karate Foundations, Chen Style Tai Chi, and Ba Duan Jin Qigong) are designed specifically for new practitioners with zero prior experience. Our instructors break down stances and footwork step-by-step.",
    },
    {
      category: "Tokens & Booking",
      q: "How does the token booking system work?",
      a: "Each token is equivalent to one class session (1 to 1.5 hours). You purchase token packages (e.g. 7, 14, 28 tokens) in the student portal, then reserve your specific in-person or live-stream class occurrences. 1 token is deducted per booking.",
    },
    {
      category: "Tokens & Booking",
      q: "Can I cancel a booked class and get my token refunded?",
      a: "Yes! If you cancel your class booking prior to the occurrence start time, your token is automatically credited back to your account ledger immediately.",
    },
    {
      category: "Tokens & Booking",
      q: "Can I use my tokens for both Kung Fu and Tai Chi?",
      a: "Yes! Your token balance is flexible and valid across all regular adult programs: Kung Fu, Karate, Kobudo, Tai Chi, and Qigong.",
    },
    {
      category: "Grading & Exams",
      q: "How often are belt grading exams conducted?",
      a: "Formal belt grading exams occur at the conclusion of each seasonal term (Spring, Summer, Fall, Winter). Instructors evaluate stance depth, form sequence, coordination, and martial spirit. Passing scores grant official Academy certificates and colored belts.",
    },
    {
      category: "General",
      q: "How do Zoom live stream classes work?",
      a: "When you book a class with live-stream mode, the Zoom link and passcode appear in your student portal dashboard under 'Upcoming Classes'. Instructors monitor the live floor camera and video feed to provide real-time verbal corrections.",
    },
  ];

  const categories = ["All", "General", "Tokens & Booking", "Grading & Exams"];

  const filteredFaqs = faqs.filter(
    (f) => activeCategory === "All" || f.category === activeCategory
  );

  return (
    <div className="py-12 sm:py-16 space-y-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="gold">Common Questions</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Frequently Asked Questions
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Find quick answers regarding token registrations, school uniforms, online live streams, and belt advancement.
        </p>

        {/* Category Filters */}
        <div className="flex flex-wrap justify-center items-center gap-2 pt-4">
          {categories.map((cat) => (
            <button
              key={cat}
              onClick={() => setActiveCategory(cat)}
              className={cn(
                "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
                activeCategory === cat
                  ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                  : "bg-slate-900 text-slate-400 hover:text-slate-200 border border-slate-800"
              )}
            >
              {cat}
            </button>
          ))}
        </div>
      </div>

      {/* Accordion List */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 space-y-4">
        {filteredFaqs.map((faq, idx) => {
          const isOpen = openIndex === idx;
          return (
            <Card
              key={idx}
              className={cn(
                "transition-all duration-200 overflow-hidden",
                isOpen ? "border-red-900/60 bg-slate-900" : "hover:border-slate-700 bg-slate-950/80"
              )}
            >
              <button
                type="button"
                onClick={() => setOpenIndex(isOpen ? null : idx)}
                className="w-full p-5 text-left flex items-center justify-between gap-4 cursor-pointer"
              >
                <div className="flex items-center gap-3">
                  <HelpCircle className={cn("h-5 w-5 shrink-0", isOpen ? "text-red-500" : "text-slate-500")} />
                  <span className="font-bold text-slate-100 text-sm sm:text-base">{faq.q}</span>
                </div>
                <ChevronDown
                  className={cn(
                    "h-4 w-4 text-slate-400 shrink-0 transition-transform duration-200",
                    isOpen && "rotate-180 text-red-400"
                  )}
                />
              </button>

              {isOpen && (
                <div className="px-5 pb-5 pt-0 text-xs sm:text-sm text-slate-300 leading-relaxed border-t border-slate-800/40 mt-1">
                  <p className="pt-3">{faq.a}</p>
                </div>
              )}
            </Card>
          );
        })}

        <div className="p-8 rounded-2xl bg-gradient-to-r from-red-950/60 via-slate-900 to-amber-950/60 border border-red-900/40 text-center space-y-4 mt-12">
          <h3 className="text-xl font-bold text-slate-100">Still have questions?</h3>
          <p className="text-xs text-slate-400 max-w-md mx-auto">
            Our school coordinators and instructors are always here to help you get started with the right program.
          </p>
          <Link href="/register">
            <Button variant="primary" size="md">
              Create Student Account
              <ArrowRight className="h-4 w-4" />
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
