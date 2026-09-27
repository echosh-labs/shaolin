"use client";

import React, { Suspense, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Input, Label } from "@/components/ui/Input";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/Card";
import { Alert } from "@/components/ui/Alert";
import { Flame, Lock, Mail, ArrowRight, ShieldCheck } from "lucide-react";

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirectPath = searchParams.get("redirect") || "/portal";

  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login(email, password);
      router.push(redirectPath);
    } catch (err: any) {
      setError(err.message || "Failed to sign in. Please verify your credentials.");
    } finally {
      setLoading(false);
    }
  };

  const handleQuickLogin = (testEmail: string, testPass: string) => {
    setEmail(testEmail);
    setPassword(testPass);
  };

  return (
    <Card glow="red">
      <CardHeader>
        <CardTitle className="text-lg">Student & Member Sign In</CardTitle>
        <CardDescription>Enter your email and password to access your dojo account.</CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <Alert variant="error" title="Authentication Error">
              {error}
            </Alert>
          )}

          <div className="space-y-1.5">
            <Label htmlFor="email">Email Address</Label>
            <div className="relative">
              <Input
                id="email"
                type="email"
                placeholder="student@martialartsacademy.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="pl-9"
              />
              <Mail className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            </div>
          </div>

          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label htmlFor="password">Password</Label>
              <Link href="/forgot-password" className="text-xs text-amber-400 hover:underline">
                Forgot?
              </Link>
            </div>
            <div className="relative">
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="pl-9"
              />
              <Lock className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            </div>
          </div>

          <Button type="submit" variant="primary" className="w-full mt-2" isLoading={loading}>
            Sign In to Portal
            <ArrowRight className="h-4 w-4" />
          </Button>
        </form>

        {/* Quick Demo Logins Helper */}
        <div className="mt-6 pt-4 border-t border-slate-800/80">
          <p className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-2 flex items-center gap-1.5">
            <ShieldCheck className="h-3.5 w-3.5 text-amber-500" />
            Quick Test Credentials
          </p>
          <div className="grid grid-cols-2 gap-2 text-xs">
            <button
              type="button"
              onClick={() => handleQuickLogin("student@martialartsacademy.com", "password123")}
              className="px-2.5 py-1.5 rounded-lg bg-slate-800/60 hover:bg-slate-800 text-slate-300 hover:text-slate-100 text-left border border-slate-700/50 transition-colors"
            >
              <span className="font-semibold block text-slate-200">Student (Justin)</span>
              <span className="text-[10px] text-slate-400">student@martialartsacademy.com</span>
            </button>
            <button
              type="button"
              onClick={() => handleQuickLogin("admin@martialartsacademy.com", "password123")}
              className="px-2.5 py-1.5 rounded-lg bg-slate-800/60 hover:bg-slate-800 text-slate-300 hover:text-slate-100 text-left border border-slate-700/50 transition-colors"
            >
              <span className="font-semibold block text-amber-400">Admin Staff</span>
              <span className="text-[10px] text-slate-400">admin@martialartsacademy.com</span>
            </button>
          </div>
        </div>
      </CardContent>

      <CardFooter className="justify-center text-xs text-slate-400">
        Don&apos;t have an account yet?{" "}
        <Link href="/register" className="text-red-400 font-semibold hover:underline ml-1">
          Register as New Student
        </Link>
      </CardFooter>
    </Card>
  );
}

export default function LoginPage() {
  return (
    <div className="min-h-screen flex flex-col justify-center items-center px-4 py-12 bg-gradient-to-b from-[#090d16] via-slate-950 to-[#020617]">
      <div className="w-full max-w-md space-y-6">
        {/* Brand Header */}
        <div className="text-center space-y-2">
          <Link href="/" className="inline-flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-red-600 to-amber-600 text-white shadow-xl shadow-red-950/60">
              <Flame className="h-7 w-7" />
            </div>
          </Link>
          <h1 className="text-2xl font-black tracking-wider text-slate-100 uppercase">
            Martial Arts Academy Portal
          </h1>
          <p className="text-xs text-slate-400">
            Sign in to manage class bookings, grading progress & tokens
          </p>
        </div>

        <Suspense fallback={<div className="text-center text-xs text-slate-400 py-8">Loading sign-in...</div>}>
          <LoginForm />
        </Suspense>
      </div>
    </div>
  );
}
