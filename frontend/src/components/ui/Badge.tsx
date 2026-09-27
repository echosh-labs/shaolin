import * as React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?:
    | "default"
    | "primary"
    | "gold"
    | "success"
    | "warning"
    | "danger"
    | "outline"
    | "beltWhite"
    | "beltYellow"
    | "beltOrange"
    | "beltGreen"
    | "beltBlue"
    | "beltRed"
    | "beltBlack";
  size?: "sm" | "md";
}

export function Badge({
  className,
  variant = "default",
  size = "md",
  ...props
}: BadgeProps) {
  const variantStyles = {
    default: "bg-slate-800 text-slate-200 border-slate-700",
    primary: "bg-red-950/80 text-red-300 border-red-800/60 shadow-sm shadow-red-950/50",
    gold: "bg-amber-950/80 text-amber-300 border-amber-800/60 shadow-sm shadow-amber-950/50",
    success: "bg-emerald-950/80 text-emerald-300 border-emerald-800/60",
    warning: "bg-yellow-950/80 text-yellow-300 border-yellow-800/60",
    danger: "bg-rose-950/80 text-rose-300 border-rose-800/60",
    outline: "bg-transparent text-slate-300 border-slate-700",
    // Martial Arts Belt Variations
    beltWhite: "bg-slate-100 text-slate-900 border-slate-300 font-bold",
    beltYellow: "bg-yellow-400 text-yellow-950 border-yellow-500 font-bold",
    beltOrange: "bg-orange-500 text-orange-950 border-orange-600 font-bold",
    beltGreen: "bg-emerald-600 text-emerald-50 border-emerald-700 font-bold",
    beltBlue: "bg-blue-600 text-blue-50 border-blue-700 font-bold",
    beltRed: "bg-red-600 text-red-50 border-red-700 font-bold",
    beltBlack: "bg-slate-950 text-amber-400 border-amber-500/80 shadow-md shadow-amber-500/20 font-bold",
  };

  const sizeStyles = {
    sm: "px-2 py-0.5 text-xs",
    md: "px-2.5 py-1 text-xs tracking-wide",
  };

  return (
    <div
      className={cn(
        "inline-flex items-center gap-1.5 font-medium rounded-full border transition-colors",
        variantStyles[variant],
        sizeStyles[size],
        className
      )}
      {...props}
    />
  );
}
