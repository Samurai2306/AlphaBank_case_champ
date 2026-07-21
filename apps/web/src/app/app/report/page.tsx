"use client";

import { useEffect, useState } from "react";
import { ArrowLeftRight, MessageCircle } from "lucide-react";
import { PageHeader } from "@/components/PageHeader";
import { formatRub, getMonthReport } from "@/lib/api";

const catRu: Record<string, string> = {
  rent: "Аренда",
  supplies: "Расходники",
  ads: "Реклама",
  fees: "Комиссии",
  ops: "Операционные",
  services: "Услуги",
};

export default function ReportPage() {
  const [data, setData] = useState<Awaited<ReturnType<typeof getMonthReport>> | null>(
    null,
  );
  const [err, setErr] = useState("");

  useEffect(() => {
    getMonthReport()
      .then(setData)
      .catch((e: Error) => setErr(e.message));
  }, []);

  const cats = Object.entries(data?.expense_by_category ?? {}).sort(
    (a, b) => b[1] - a[1],
  );

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <PageHeader
        title="Отчёт за месяц"
        subtitle="Сводка по операциям счёта — чтобы понимать, откуда налог и куда уходят расходы."
        actions={[
          {
            href: "/app/operations",
            label: "Операции",
            icon: ArrowLeftRight,
          },
          {
            href:
              "/app/chat?q=" +
              encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
            label: "В чат",
            icon: MessageCircle,
            primary: true,
          },
        ]}
      />
      <p className="mt-2 text-sm text-ink/50">Период · {data?.period ?? "…"}</p>

      {err ? <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p> : null}

      <section className="mt-6 grid grid-cols-2 gap-3">
        <article className="rounded-card bg-mint p-5 shadow-soft">
          <p className="text-xs text-ink/50">Доходы</p>
          <p className="mt-1 text-2xl font-bold tabular text-emerald-800">
            {data ? `${formatRub(data.income)} ₽` : "—"}
          </p>
        </article>
        <article className="rounded-card bg-rose p-5 shadow-soft">
          <p className="text-xs text-ink/50">Расходы</p>
          <p className="mt-1 text-2xl font-bold tabular">
            {data ? `${formatRub(data.expense)} ₽` : "—"}
          </p>
        </article>
        <article className="rounded-card bg-white p-5 shadow-soft">
          <p className="text-xs text-ink/50">Нетто</p>
          <p
            className={`mt-1 text-2xl font-bold tabular ${
              data && data.net < 0 ? "text-brand" : ""
            }`}
          >
            {data ? `${formatRub(data.net)} ₽` : "—"}
          </p>
        </article>
        <article className="rounded-card bg-lavender p-5 shadow-soft">
          <p className="text-xs text-ink/50">Услуг проведено</p>
          <p className="mt-1 text-2xl font-bold tabular">
            {data?.service_count ?? "—"}
          </p>
        </article>
      </section>

      <section className="mt-6 rounded-card bg-white p-5 shadow-soft">
        <p className="font-semibold">Оценка налога НПД</p>
        <p className="mt-2 text-3xl font-bold tabular">
          {data ? `${formatRub(data.tax_estimate)} ₽` : "—"}
        </p>
      </section>

      <section className="mt-6">
        <h2 className="text-lg font-bold">Расходы по категориям</h2>
        <ul className="mt-3 space-y-2">
          {cats.map(([k, v]) => (
            <li
              key={k}
              className="flex justify-between rounded-card bg-white px-4 py-3 text-sm shadow-soft"
            >
              <span>{catRu[k] ?? k}</span>
              <span className="font-semibold tabular">{formatRub(v)} ₽</span>
            </li>
          ))}
          {cats.length === 0 && !err ? (
            <li className="text-sm text-ink/50">Нет расходов за период</li>
          ) : null}
        </ul>
      </section>
    </main>
  );
}
