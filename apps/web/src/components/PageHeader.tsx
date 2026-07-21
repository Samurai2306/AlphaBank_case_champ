"use client";

import Link from "next/link";
import type { LucideIcon } from "lucide-react";
import { ArrowLeft } from "lucide-react";

export type HeaderAction = {
  href: string;
  label: string;
  icon: LucideIcon;
  primary?: boolean;
};

export function PageHeader({
  title,
  subtitle,
  backHref = "/app",
  backLabel = "Главная",
  actions = [],
}: {
  title: string;
  subtitle?: string;
  backHref?: string;
  backLabel?: string;
  actions?: HeaderAction[];
}) {
  return (
    <header className="mb-1">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <Link
            href={backHref}
            className="inline-flex items-center gap-1.5 text-sm text-ink/50 transition hover:text-ink"
          >
            <ArrowLeft className="h-4 w-4" aria-hidden />
            {backLabel}
          </Link>
          <h1 className="mt-2 text-2xl font-bold tracking-tight">{title}</h1>
          {subtitle ? (
            <p className="mt-1 text-sm leading-snug text-ink/60">{subtitle}</p>
          ) : null}
        </div>
        {actions.length ? (
          <div className="flex shrink-0 flex-wrap justify-end gap-2">
            {actions.map((a) => {
              const Icon = a.icon;
              return (
                <Link
                  key={a.href + a.label}
                  href={a.href}
                  className={`inline-flex items-center gap-1.5 rounded-pill px-3 py-2 text-xs font-semibold transition sm:text-sm ${
                    a.primary
                      ? "bg-brand text-white"
                      : "bg-white text-ink/70 shadow-soft hover:text-ink"
                  }`}
                >
                  <Icon className="h-3.5 w-3.5" aria-hidden />
                  {a.label}
                </Link>
              );
            })}
          </div>
        ) : null}
      </div>
    </header>
  );
}
