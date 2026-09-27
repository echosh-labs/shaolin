"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Input, Label } from "@/components/ui/Input";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/Card";
import { Alert } from "@/components/ui/Alert";
import { Flame, Lock, Mail, User, Calendar, ArrowRight } from "lucide-react";

export default function RegisterPage() {
  const router = useRouter();
  const { register } = useAuth();

  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    email: "",
    date_of_birth: "2000-01-01",
    password: "",
    confirm_password: "",
  });

  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.password !== formData.confirm_password) {
      setError("Passwords do not match");
      return;
    }

    if (formData.password.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }

    setLoading(true);

    try {
      await register({
        first_name: formData.first_name,
        last_name: formData.last_name,
        email: formData.email,
        date_of_birth: formData.date_of_birth,
        password: formData.password,
      });

      router.push("/portal");
    } catch (err: any) {
      setError(err.message || "Registration failed. Please check your information.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col justify-center items-center px-4 py-12 bg-gradient-to-b from-[#090d16] via-slate-950 to-[#020617]">
      <div className="w-full max-w-lg space-y-6">
        <div className="text-center space-y-2">
          <Link href="/" className="inline-flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-red-600 to-amber-600 text-white shadow-xl shadow-red-950/60">
              <Flame className="h-7 w-7" />
            </div>
          </Link>
          <h1 className="text-2xl font-black tracking-wider text-slate-100 uppercase">
            Join Martial Arts Academy
          </h1>
          <p className="text-xs text-slate-400">
            Create your martial arts student account to register for classes and track progress
          </p>
        </div>

        <Card glow="red">
          <CardHeader>
            <CardTitle className="text-lg">Student Onboarding Registration</CardTitle>
            <CardDescription>
              Complete the form below to establish your student training profile.
            </CardDescription>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <Alert variant="error" title="Registration Error">
                  {error}
                </Alert>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-1.5">
                  <Label htmlFor="first_name">First Name</Label>
                  <div className="relative">
                    <Input
                      id="first_name"
                      placeholder="e.g. Bruce"
                      value={formData.first_name}
                      onChange={(e) =>
                        setFormData({ ...formData, first_name: e.target.value })
                      }
                      required
                      className="pl-9"
                    />
                    <User className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="last_name">Last Name</Label>
                  <div className="relative">
                    <Input
                      id="last_name"
                      placeholder="e.g. Lee"
                      value={formData.last_name}
                      onChange={(e) =>
                        setFormData({ ...formData, last_name: e.target.value })
                      }
                      required
                      className="pl-9"
                    />
                    <User className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                  </div>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="email">Email Address</Label>
                <div className="relative">
                  <Input
                    id="email"
                    type="email"
                    placeholder="student@example.com"
                    value={formData.email}
                    onChange={(e) =>
                      setFormData({ ...formData, email: e.target.value })
                    }
                    required
                    className="pl-9"
                  />
                  <Mail className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="date_of_birth">Date of Birth (YYYY-MM-DD)</Label>
                <div className="relative">
                  <Input
                    id="date_of_birth"
                    type="date"
                    value={formData.date_of_birth}
                    onChange={(e) =>
                      setFormData({ ...formData, date_of_birth: e.target.value })
                    }
                    required
                    className="pl-9"
                  />
                  <Calendar className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-1.5">
                  <Label htmlFor="password">Password</Label>
                  <div className="relative">
                    <Input
                      id="password"
                      type="password"
                      placeholder="••••••••"
                      value={formData.password}
                      onChange={(e) =>
                        setFormData({ ...formData, password: e.target.value })
                      }
                      required
                      className="pl-9"
                    />
                    <Lock className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="confirm_password">Confirm Password</Label>
                  <div className="relative">
                    <Input
                      id="confirm_password"
                      type="password"
                      placeholder="••••••••"
                      value={formData.confirm_password}
                      onChange={(e) =>
                        setFormData({ ...formData, confirm_password: e.target.value })
                      }
                      required
                      className="pl-9"
                    />
                    <Lock className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                  </div>
                </div>
              </div>

              <Button type="submit" variant="primary" className="w-full mt-3" isLoading={loading}>
                Create Student Account & Sign In
                <ArrowRight className="h-4 w-4" />
              </Button>
            </form>
          </CardContent>

          <CardFooter className="justify-center text-xs text-slate-400">
            Already have an account?{" "}
            <Link href="/login" className="text-red-400 font-semibold hover:underline ml-1">
              Sign In
            </Link>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
