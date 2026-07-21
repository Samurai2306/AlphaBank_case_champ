"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import { getAlerts } from "@/lib/api";

const tone: Record<string, string> = {
  high: "bg-rose",
  medium: "bg-peach",
  low: "bg-mint",
};

export default function AlertsPage() {
  const [items, setItems] = useState<
    Awaited<ReturnType<typeof getAlerts>>["items"]
  >([]);
  const [err, setErr] = useState("");

  useEffect(() => {
    getAlerts()
      .then((r) => setItems(r.items))
      .catch((e: Error) => setErr(e.message));
  }, []);

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <Link href="/app" className="text-sm text-ink/50 hover:text-ink">
        ← Главная
      </Link>

      <section className="relative mt-3 min-h-[160px] overflow-hidden rounded-[32px] bg-lavender p-5 shadow-soft">
        <div className="relative z-10 max-w-[65%]">
          <h1 className="text-2xl font-bold">Алерты</h1>
          <p className="mt-2 text-sm text-ink/65">
            Проактивные сигналы — факт, зачем важно, что сделать.
          </p>
        </div>
        <div className="pointer-events-none absolute -bottom-3 -right-2 h-40 w-40">
          <Image
            src="/icons/3d/08-proactive-alerts.png"
            alt=""
            fill
            className="object-contain object-bottom"
            sizes="160px"
          />
        </div>
      </section>

      {err ? (
        <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p>
      ) : null}

      <ul className="mt-6 space-y-3">
        {items.map((a) => (
          <li
            key={a.id}
            className={`card-touch rounded-card p-5 shadow-soft ${tone[a.severity] ?? "bg-white"}`}
          >
            <p className="text-xs uppercase tracking-wide text-ink/45">
              {a.severity === "high"
                ? "Важно"
                : a.severity === "medium"
                  ? "Внимание"
                  : "Подсказка"}
            </p>
            <p className="mt-1 font-semibold">{a.title}</p>
            <p className="mt-2 text-sm text-ink/70">{a.text}</p>
            <Link
              href={a.href}
              className="mt-3 inline-flex rounded-pill bg-ink px-4 py-2 text-sm font-medium text-white"
            >
              {a.cta}
            </Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
