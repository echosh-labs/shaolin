import * as React from "react";
import { cn } from "@/lib/utils";
import { AlertCircle, CheckCircle2, Info, AlertTriangle } from "lucide-react";

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "info" | "success" | "warning" | "error";
  title?: string;
}

export function Alert({
  className,
  variant = "info",
  title,
  children,
  ...props
}: AlertProps) {
  const variantIcons = {
    info: <Info className="h-5 w-5 text-sky-400 shrink-0 mt-0.5" />,
    success: <CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" />,
    warning: <AlertTriangle className="h-5 w-5 text-amber-400 shrink-0 mt-0.5" />,
    error: <AlertCircle className="h-5 w-5 text-rose-400 shrink-0 mt-0.5" />,
  };

  const variantStyles = {
    info: "bg-sky-950/40 text-sky-200 border-sky-800/50",
    success: "bg-emerald-950/40 text-emerald-200 border-emerald-800/50",
    warning: "bg-amber-950/40 text-amber-200 border-amber-800/50",
    error: "bg-rose-950/40 text-rose-200 border-rose-800/50",
  };

  return (
    <div
      className={cn(
        "flex gap-3 rounded-xl border p-4 text-sm leading-relaxed",
        variantStyles[variant],
        className
      )}
      {...props}
    >
      {variantIcons[variant]}
      <div className="space-y-0.5">
        {title && <h5 className="font-semibold text-slate-100">{title}</h5>}
        <div className="text-slate-300">{children}</div>
      </div>
    </div>
  );
}
