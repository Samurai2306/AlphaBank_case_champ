"use client";

import { useEffect, useState } from "react";
import { Plus } from "lucide-react";
import { EmptyHint } from "@/components/EmptyHint";
import { PageHeader } from "@/components/PageHeader";
import { formatRub, getPayments } from "@/lib/api";

const statusRu: Record<string, string> = {
  draft: "черновик",
  confirmed: "принято",
  confirmed_mock: "принято",
  cancelled: "отменено",
};

export default function PaymentsPage() {
  const [items, setItems] = useState<
    Awaited<ReturnType<typeof getPayments>>["items"]
  >([]);
  const [err, setErr] = useState("");

  useEffect(() => {
    getPayments()
      .then((r) => setItems(r.items))
      .catch((e: Error) => setErr(e.message));
  }, []);

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <PageHeader
        title="Платежи"
        subtitle="История поручений в банк. Новое поручение удобнее собрать в чате."
        backHref="/app/profile"
        backLabel="Ещё"
        actions={[
          {
            href: "/app/chat?q=" + encodeURIComponent("Сформируй платёжку."),
            label: "Новое",
            icon: Plus,
            primary: true,
          },
        ]}
      />

      {err ? <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p> : null}

      <ul className="mt-6 space-y-3">
        {items.map((p) => (
          <li key={p.draft_id} className="rounded-card bg-white p-4 shadow-soft">
            <div className="flex items-start justify-between gap-3">
              <div>
                <p className="font-semibold">{p.purpose}</p>
                <p className="mt-1 text-sm text-ink/55">{p.payee_name}</p>
                <p className="mt-1 text-xs text-ink/40">
                  {p.updated_at
                    ? new Date(p.updated_at).toLocaleString("ru-RU")
                    : ""}
                </p>
              </div>
              <div className="text-right">
                <p className="font-bold tabular">{formatRub(p.amount)} ₽</p>
                <p className="mt-1 text-xs text-ink/50">
                  {statusRu[p.status] ?? p.status}
                </p>
              </div>
            </div>
          </li>
        ))}
        {!err && items.length === 0 ? (
          <li>
            <EmptyHint
              title="Поручений пока нет"
              text="Сначала посчитайте налог за месяц — затем соберите платёжку в чате. Средства не списываются, пока вы не подтвердите черновик."
              primary={{
                href:
                  "/app/chat?q=" +
                  encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
                label: "Рассчитать налог",
              }}
              secondary={{
                href: "/app/chat?q=" + encodeURIComponent("Сформируй платёжку."),
                label: "Сразу платёжка",
              }}
            />
          </li>
        ) : null}
      </ul>
    </main>
  );
}
