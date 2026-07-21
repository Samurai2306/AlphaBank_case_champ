"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  Bell,
  FileText,
  Landmark,
  ListChecks,
  LogOut,
  Mail,
  MapPin,
  Phone,
  PiggyBank,
  Route,
  Save,
  Wallet,
} from "lucide-react";
import { getMe, patchMe, type Profile } from "@/lib/api";
import { clearSession } from "@/lib/session";
import { getCategoryMeta, loadCashbackPick } from "@/lib/cashback";

const hub = [
  {
    href: "/app/journey",
    label: "Мой путь",
    hint: "Трекборд 5 уровней",
    icon: Route,
    progress: true,
  },
  {
    href: "/app/piggy",
    label: "Копилка ЕНП",
    hint: "Автоотчисления",
    icon: PiggyBank,
  },
  {
    href: "/app/payments",
    label: "Платежи",
    hint: "Поручения в банк",
    icon: Wallet,
  },
  {
    href: "/app/report",
    label: "Отчёт",
    hint: "Доходы и расходы",
    icon: FileText,
  },
  {
    href: "/app/alerts",
    label: "Алерты",
    hint: "Что требует внимания",
    icon: Bell,
  },
  {
    href: "/app/tasks",
    label: "Задачи",
    hint: "Рекомендуемые шаги",
    icon: ListChecks,
  },
];

export default function ProfilePage() {
  const router = useRouter();
  const [form, setForm] = useState<Profile | null>(null);
  const [saving, setSaving] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(true);
  const [cashPick, setCashPick] = useState<ReturnType<typeof loadCashbackPick>>(null);
  const cashMeta = cashPick ? getCategoryMeta(cashPick.id) : null;

  useEffect(() => {
    setCashPick(loadCashbackPick());
    getMe()
      .then(setForm)
      .catch((e: Error) => setErr(e.message))
      .finally(() => setLoading(false));
  }, []);

  async function onSave(e: React.FormEvent) {
    e.preventDefault();
    if (!form) return;
    setSaving(true);
    setMsg("");
    setErr("");
    try {
      const next = await patchMe({
        display_name: form.display_name,
        business_sphere: form.business_sphere,
        city: form.city,
        monthly_revenue_estimate: form.monthly_revenue_estimate,
        tax_regime: form.tax_regime || "NPD",
      });
      setForm(next);
      setMsg("Профиль обновлён.");
    } catch (ex) {
      setErr((ex as Error).message);
    } finally {
      setSaving(false);
    }
  }

  const level = form?.cjm_level ?? 0;

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <p className="text-sm text-ink/50">Кабинет</p>
      <h1 className="mt-1 text-2xl font-bold">Ещё</h1>
      <p className="mt-1 text-sm text-ink/60">
        Путь, копилка, платежи, отчёты и профиль — всё в одном месте.
      </p>

      <nav className="mt-5 grid grid-cols-2 gap-3" aria-label="Разделы кабинета">
        {hub.map((item) => {
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              className="card-touch relative overflow-hidden rounded-[24px] bg-white p-4 shadow-soft"
            >
              <Icon
                className="pointer-events-none absolute bottom-2 right-2 h-11 w-11 text-brand/[0.08]"
                strokeWidth={1.25}
                aria-hidden
              />
              <span className="relative z-10 flex h-10 w-10 items-center justify-center rounded-2xl bg-canvas text-brand">
                <Icon className="h-5 w-5" strokeWidth={2} aria-hidden />
              </span>
              <span className="relative z-10 mt-3 block text-sm font-semibold">
                {item.label}
              </span>
              <span className="relative z-10 mt-0.5 block text-xs text-ink/50">
                {item.progress && level
                  ? `Уровень ${level}/5`
                  : item.hint}
              </span>
              {item.progress && level ? (
                <span
                  className="relative z-10 mt-2.5 block h-1 overflow-hidden rounded-full bg-platinum/80"
                  aria-hidden
                >
                  <span
                    className="block h-full rounded-full bg-brand/85"
                    style={{ width: `${Math.min(100, (level / 5) * 100)}%` }}
                  />
                </span>
              ) : null}
            </Link>
          );
        })}
      </nav>

      {loading ? (
        <div className="mt-6 h-40 animate-pulse rounded-card bg-white shadow-soft" />
      ) : null}

      {err && !form ? (
        <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p>
      ) : null}

      {form ? (
        <>
          <section className="mt-6 rounded-card bg-white p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <Landmark className="h-5 w-5 text-brand" aria-hidden />
              <p className="text-lg font-bold">{form.full_name || form.display_name}</p>
            </div>
            <dl className="mt-4 grid gap-3 text-sm sm:grid-cols-2">
              <div className="flex items-center gap-2.5">
                <Phone className="h-4 w-4 shrink-0 text-brand/70" strokeWidth={2} aria-hidden />
                <div>
                  <dt className="text-ink/45">Телефон</dt>
                  <dd className="mt-0.5 font-medium">{form.phone || "—"}</dd>
                </div>
              </div>
              <div className="flex items-center gap-2.5">
                <Mail className="h-4 w-4 shrink-0 text-brand/70" strokeWidth={2} aria-hidden />
                <div>
                  <dt className="text-ink/45">Email</dt>
                  <dd className="mt-0.5 break-all font-medium">{form.email || "—"}</dd>
                </div>
              </div>
              <div className="flex items-center gap-2.5">
                <FileText className="h-4 w-4 shrink-0 text-brand/70" strokeWidth={2} aria-hidden />
                <div>
                  <dt className="text-ink/45">ИНН</dt>
                  <dd className="mt-0.5 font-medium tabular">{form.inn || "—"}</dd>
                </div>
              </div>
              <div className="flex items-center gap-2.5">
                <Wallet className="h-4 w-4 shrink-0 text-brand/70" strokeWidth={2} aria-hidden />
                <div>
                  <dt className="text-ink/45">Счёт</dt>
                  <dd className="mt-0.5 font-medium">
                    {form.bank_name || "Альфа-Банк"} {form.account_masked || ""}
                  </dd>
                </div>
              </div>
              <div className="flex items-center gap-2.5 sm:col-span-2">
                <MapPin className="h-4 w-4 shrink-0 text-brand/70" strokeWidth={2} aria-hidden />
                <div>
                  <dt className="text-ink/45">Адрес</dt>
                  <dd className="mt-0.5 font-medium">{form.address || "—"}</dd>
                </div>
              </div>
            </dl>
          </section>

          <form onSubmit={onSave} className="mt-4 space-y-4 rounded-card bg-white p-5 shadow-soft">
            <p className="font-semibold">Параметры помощника</p>
            <p className="mt-1 text-xs text-ink/50">
              От оборота считается налог на главной, даже если операций ещё нет.
            </p>
            <p className="text-xs text-ink/50">
              Эти данные используются в расчётах налога и копилке.
            </p>
            {(
              [
                ["display_name", "Имя в интерфейсе", "text"],
                ["business_sphere", "Сфера", "text"],
                ["city", "Город", "text"],
              ] as const
            ).map(([key, label, type]) => (
              <label key={key} className="block text-sm">
                <span className="text-ink/60">{label}</span>
                <input
                  type={type}
                  value={form[key]}
                  onChange={(e) => setForm({ ...form, [key]: e.target.value })}
                  className="mt-1 w-full rounded-2xl border border-platinum bg-canvas px-4 py-3 outline-none ring-brand focus:ring-2"
                />
              </label>
            ))}
            <label className="block text-sm">
              <span className="text-ink/60">Оборот, ₽/мес</span>
              <input
                type="number"
                min={0}
                value={form.monthly_revenue_estimate}
                onChange={(e) =>
                  setForm({
                    ...form,
                    monthly_revenue_estimate: Number(e.target.value),
                  })
                }
                className="mt-1 w-full rounded-2xl border border-platinum bg-canvas px-4 py-3 outline-none ring-brand focus:ring-2"
              />
            </label>
            <label className="block text-sm">
              <span className="text-ink/60">Режим</span>
              <select
                value={form.tax_regime || "NPD"}
                onChange={(e) => setForm({ ...form, tax_regime: e.target.value })}
                className="mt-1 w-full rounded-2xl border border-platinum bg-canvas px-4 py-3 outline-none ring-brand focus:ring-2"
              >
                <option value="NPD">НПД</option>
                <option value="USN_6">УСН 6%</option>
                <option value="IP">ИП (упрощённо)</option>
                <option value="UNDECIDED">Пока думаю</option>
              </select>
            </label>
            <p className="text-xs text-ink/45">Уровень пути: {form.cjm_level}/5</p>
            {cashPick && cashMeta ? (
              <p className="rounded-2xl bg-mint/80 px-3 py-2 text-xs text-emerald-900">
                Кэшбек месяца: {cashMeta.title} · {cashPick.rate}%{" "}
                <Link href="/app" className="font-semibold underline">
                  изменить
                </Link>
              </p>
            ) : null}
            <button
              type="submit"
              disabled={saving}
              className="btn-touch inline-flex w-full items-center justify-center gap-2 rounded-pill bg-brand px-5 py-3 text-sm font-semibold text-white disabled:opacity-50"
            >
              <Save className="h-4 w-4" aria-hidden />
              {saving ? "Сохраняю…" : "Сохранить"}
            </button>
            {msg ? <p className="text-sm text-emerald-700">{msg}</p> : null}
            {err ? <p className="text-sm text-brand">{err}</p> : null}
          </form>

          <button
            type="button"
            onClick={() => {
              clearSession();
              router.replace("/login");
            }}
            className="btn-touch mt-4 inline-flex w-full items-center justify-center gap-2 rounded-pill bg-white py-3 text-sm font-semibold text-ink shadow-soft ring-1 ring-platinum"
          >
            <LogOut className="h-4 w-4" aria-hidden />
            Выйти из кабинета
          </button>
        </>
      ) : null}
    </main>
  );
}
