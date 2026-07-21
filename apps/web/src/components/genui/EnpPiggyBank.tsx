"use client";

import Link from "next/link";
import { formatRub } from "@/lib/api";

export function EnpPiggyBank({ props }: { props: Record<string, unknown> }) {
  return (
    <div
      className="rounded-card bg-mint p-5 shadow-soft"
      role="region"
      aria-label="Копилка ЕНП"
    >
      <p className="font-semibold">Налоговая копилка ЕНП</p>
      <p className="mt-2 text-3xl font-bold tabular">
        {formatRub(Number(props.balance ?? 0))} ₽
      </p>
      <p className="mt-2 text-sm text-ink/70">
        {props.enabled ? "Включена" : "Выключена"} · ставка{" "}
        {Number(props.rate_percent ?? 0)}% · последнее отчисление{" "}
        {formatRub(Number(props.last_contribution ?? 0))} ₽
      </p>
      <Link
        href="/app/piggy"
        className="mt-4 inline-flex rounded-pill bg-ink px-4 py-2 text-sm font-medium text-white"
      >
        Настроить копилку
      </Link>
    </div>
  );
}
