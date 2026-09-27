import Link from "next/link";
import { Flame, MapPin, Phone, Mail, Clock } from "lucide-react";

export function Footer() {
  return (
    <footer className="w-full border-t border-slate-800/80 bg-slate-950 text-slate-400 text-sm mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          {/* Brand Col */}
          <div className="space-y-4 md:col-span-1">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-red-600 text-white shadow-md shadow-red-950/60">
                <Flame className="h-5 w-5" />
              </div>
              <span className="font-bold tracking-wider text-slate-100 text-base uppercase">
                Martial Arts Academy
              </span>
            </div>
            <p className="text-xs text-slate-400 leading-relaxed">
              Dedicated to authentic traditional Kung Fu, Karate, Kobudo, Tai Chi, and Qigong education, physical vitality, and martial discipline.
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="font-semibold text-slate-200 mb-3 text-xs uppercase tracking-wider">
              Training Programs
            </h4>
            <ul className="space-y-2 text-xs">
              <li>
                <Link href="/classes" className="hover:text-red-400 transition-colors">
                  Traditional Kung Fu
                </Link>
              </li>
              <li>
                <Link href="/classes" className="hover:text-amber-400 transition-colors">
                  Karate (Kata & Kihon)
                </Link>
              </li>
              <li>
                <Link href="/classes" className="hover:text-red-400 transition-colors">
                  Kobudo (Traditional Weapons)
                </Link>
              </li>
              <li>
                <Link href="/classes" className="hover:text-amber-400 transition-colors">
                  Tai Chi & Internal Arts
                </Link>
              </li>
              <li>
                <Link href="/classes" className="hover:text-red-400 transition-colors">
                  Qigong (Ba Duan Jin)
                </Link>
              </li>
              <li>
                <Link href="/tuition" className="hover:text-slate-200 transition-colors">
                  Token Packages & Rates
                </Link>
              </li>
            </ul>
          </div>

          {/* Community & Student */}
          <div>
            <h4 className="font-semibold text-slate-200 mb-3 text-xs uppercase tracking-wider">
              Student Community
            </h4>
            <ul className="space-y-2 text-xs">
              <li>
                <Link href="/portal" className="hover:text-slate-200 transition-colors">
                  Student Dashboard
                </Link>
              </li>
              <li>
                <Link href="/guides" className="hover:text-slate-200 transition-colors">
                  Uniform & Belt Guide
                </Link>
              </li>
              <li>
                <Link href="/store" className="hover:text-slate-200 transition-colors">
                  Dojo Store & Equipment
                </Link>
              </li>
              <li>
                <Link href="/faq" className="hover:text-slate-200 transition-colors">
                  Frequently Asked Questions
                </Link>
              </li>
            </ul>
          </div>

          {/* Location & Hall Info */}
          <div className="space-y-3 text-xs">
            <h4 className="font-semibold text-slate-200 mb-3 uppercase tracking-wider">
              Training Centre
            </h4>
            <div className="flex items-start gap-2 text-slate-400">
              <MapPin className="h-4 w-4 text-red-500 shrink-0 mt-0.5" />
              <span>Downtown Toronto Training Centre, ON</span>
            </div>
            <div className="flex items-center gap-2 text-slate-400">
              <Clock className="h-4 w-4 text-amber-500 shrink-0" />
              <span>Training: Mon–Sat (In-Person & Virtual)</span>
            </div>
            <div className="flex items-center gap-2 text-slate-400">
              <Mail className="h-4 w-4 text-slate-500 shrink-0" />
              <span>contact@martialartsacademy.com</span>
            </div>
          </div>
        </div>

        <div className="border-t border-slate-900 pt-6 flex flex-col sm:flex-row items-center justify-between text-xs text-slate-500 gap-4">
          <p>© {new Date().getFullYear()} Martial Arts Academy. All rights reserved.</p>
          <div className="flex items-center gap-4">
            <Link href="/waiver" className="hover:text-slate-400 transition-colors">
              Liability Waiver
            </Link>
            <Link href="/terms" className="hover:text-slate-400 transition-colors">
              Terms & Refund Policy
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
