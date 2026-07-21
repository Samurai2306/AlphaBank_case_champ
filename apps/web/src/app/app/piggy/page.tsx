"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import { EmptyHint } from "@/components/EmptyHint";
import {
  formatRub,
  getHome,
  getPiggy,
  getPiggyContributions,
  patchPiggy,
} from "@/lib/api";

export default function PiggyPage() {
  const [enabled, setEnabled] = useState(true);
  const [rate, setRate] = useState(6);
  const [balance, setBalance] = useState(0);
  const [last, setLast] = useState(0);
  const [target, setTarget] = useState(0);
  const [gap, setGap] = useState<number | null>(null);
  const [contribs, setContribs] = useState<
    Awaited<ReturnType<typeof getPiggyContributions>>["items"]
  >([]);
  const [saving, setSaving] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  useEffect(() => {
    Promise.all([getPiggy(), getHome(), getPiggyContributions()])
      .then(([p, h, c]) => {
        setEnabled(p.enabled);
        setRate(p.rate_percent);
        setBalance(p.balance);
        setLast(p.last_contribution);
        setTarget(p.target_amount ?? h.tax_due);
        setGap(h.piggy_gap ?? null);
        setContribs(c.items.slice(0, 8));
      })
      .catch((e: Error) => setErr(e.message));
  }, []);

  async function save() {
    setSaving(true);
    setMsg("");
    setErr("");
    try {
      const p = await patchPiggy({ enabled, rate_percent: rate });
      setEnabled(p.enabled);
      setRate(p.rate_percent);
      setBalance(p.balance);
      setLast(p.last_contribution);
      if (p.target_amount) setTarget(p.target_amount);
      setMsg("Настройки копилки сохранены. Отчисления с новых поступлений.");
    } catch (e) {
      setErr((e as Error).message);
    } finally {
      setSaving(false);
    }
  }

  const progress =
    target > 0 ? Math.min(100, Math.round((balance / target) * 100)) : 0;

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <Link href="/app" className="text-sm text-ink/50">
        ← Главная
      </Link>
      <h1 className="mt-3 text-2xl font-bold">Копилка ЕНП</h1>
      <p className="mt-1 text-sm text-ink/60">
        Откладывает % с поступлений к сроку уплаты налога.
      </p>

      {!enabled && balance === 0 ? (
        <div className="mt-6">
          <EmptyHint
            title="Копилка выключена"
            text="Включите автоотчисления после расчёта налога — копилка будет копить % с поступлений к сроку ЕНП. Пока операций нет, ориентир — оборот из профиля."
            primary={{
              href:
                "/app/chat?q=" + encodeURIComponent("Включи копилку 6%."),
              label: "Включить в чате",
            }}
            secondary={{
              href:
                "/app/chat?q=" +
                encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
              label: "Сначала налог",
            }}
          />
        </div>
      ) : null}

      <section className="relative mt-6 min-h-[200px] overflow-hidden rounded-[32px] bg-mint p-6 shadow-soft">
        <div className="relative z-10 max-w-[60%]">
          <p className="text-sm text-ink/60">Баланс</p>
          <p className="mt-1 text-4xl font-bold tabular">{formatRub(balance)} ₽</p>
          <p className="mt-2 text-sm text-ink/55">
            Цель {formatRub(target)} ₽ · последний вклад ~{formatRub(last)} ₽
            {gap != null ? ` · ещё ${formatRub(gap)} ₽` : ""}
          </p>
          <div className="mt-4 h-2 overflow-hidden rounded-pill bg-white/70">
            <div className="h-full rounded-pill bg-brand" style={{ width: `${progress}%` }} />
          </div>
          <p className="mt-2 text-xs text-ink/45">{progress}% от суммы «с запасом»</p>
        </div>
        <div className="pointer-events-none absolute -bottom-4 -right-2 h-44 w-44 sm:h-52 sm:w-52">
          <Image
            src="/icons/3d/09-enp-piggy-bank.png"
            alt=""
            fill
            className="object-contain object-bottom drop-shadow-md"
            sizes="208px"
          />
        </div>
      </section>

      <section className="mt-4 space-y-4 rounded-card bg-white p-5 shadow-soft">
        <label className="flex items-center justify-between gap-3">
          <span className="font-medium">Автоотчисления</span>
          <button
            type="button"
            role="switch"
            aria-checked={enabled}
            onClick={() => setEnabled((v) => !v)}
            className={`relative h-8 w-14 rounded-pill transition ${
              enabled ? "bg-brand" : "bg-platinum"
            }`}
          >
            <span
              className={`absolute top-1 h-6 w-6 rounded-full bg-white transition ${
                enabled ? "left-7" : "left-1"
              }`}
            />
          </button>
        </label>

        <div>
          <div className="flex justify-between text-sm">
            <span>Ставка</span>
            <span className="font-semibold tabular">{rate}%</span>
          </div>
          <input
            type="range"
            min={1}
            max={15}
            value={rate}
            onChange={(e) => setRate(Number(e.target.value))}
            className="mt-2 w-full accent-[var(--brand-red)]"
            disabled={!enabled}
          />
          <p className="mt-1 text-xs text-ink/45">Допустимый диапазон: 1–15%</p>
        </div>

        <button
          type="button"
          disabled={saving}
          onClick={() => void save()}
          className="w-full rounded-pill bg-ink px-5 py-3 text-sm font-semibold text-white disabled:opacity-50"
        >
          {saving ? "Сохраняю…" : "Сохранить"}
        </button>
        {msg ? <p className="text-sm text-emerald-700">{msg}</p> : null}
        {err ? <p className="text-sm text-brand">{err}</p> : null}
      </section>

      <section className="mt-6">
        <h2 className="text-lg font-bold">Последние взносы</h2>
        <ul className="mt-3 space-y-2">
          {contribs.map((c) => (
            <li
              key={c.id}
              className="flex justify-between rounded-card bg-white px-4 py-3 text-sm shadow-soft"
            >
              <span className="text-ink/55">
                {new Date(c.booked_at).toLocaleDateString("ru-RU")}
              </span>
              <span className="font-semibold tabular">+{formatRub(c.amount)} ₽</span>
            </li>
          ))}
          {contribs.length === 0 && !err ? (
            <li className="text-sm text-ink/50">Взносов пока нет</li>
          ) : null}
        </ul>
      </section>

      <Link
        href={"/app/chat?q=" + encodeURIComponent("Включи копилку 6%.")}
        className="mt-4 inline-flex text-sm font-medium text-brand"
      >
        Спросить в чате про копилку →
      </Link>
    </main>
  );
}
