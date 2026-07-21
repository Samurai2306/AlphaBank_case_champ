"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { MessageCircle } from "lucide-react";
import { PageHeader } from "@/components/PageHeader";
import { getTasks } from "@/lib/api";

export default function TasksPage() {
  const [title, setTitle] = useState("Рекомендуемые действия");
  const [items, setItems] = useState<
    Awaited<ReturnType<typeof getTasks>>["items"]
  >([]);
  const [err, setErr] = useState("");

  useEffect(() => {
    getTasks()
      .then((r) => {
        setTitle(r.title);
        setItems(r.items);
      })
      .catch((e: Error) => setErr(e.message));
  }, []);

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <PageHeader
        title={title}
        subtitle="Готовые шаги по кабинету Марии — каждый ведёт в чат или нужный раздел."
        actions={[
          {
            href: "/app/chat",
            label: "Чат",
            icon: MessageCircle,
            primary: true,
          },
        ]}
      />

      {err ? <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p> : null}

      <ol className="mt-6 space-y-3">
        {items.map((item, idx) => (
          <li key={item.id}>
            <Link
              href={item.href}
              className="card-touch flex items-start gap-3 rounded-card bg-white p-4 shadow-soft"
            >
              <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-canvas text-sm font-bold text-ink/60">
                {idx + 1}
              </span>
              <div>
                <p className="font-semibold">{item.label}</p>
                <p className="mt-1 text-sm text-ink/55">{item.text}</p>
              </div>
            </Link>
          </li>
        ))}
      </ol>
    </main>
  );
}
