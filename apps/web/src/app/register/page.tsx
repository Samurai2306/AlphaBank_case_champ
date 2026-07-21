"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { registerAccount } from "@/lib/api";
import { enterExampleSession, setSession } from "@/lib/session";

const spheres = ["Бьюти", "Контент / креатив", "Услуги", "Другое"];
const revenues = [
  { label: "30к", value: 30000 },
  { label: "80к", value: 80000 },
  { label: "150к", value: 150000 },
];

export default function RegisterPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState("");
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("+7 ");
  const [sphere, setSphere] = useState("Бьюти");
  const [sphereCustom, setSphereCustom] = useState("");
  const [city, setCity] = useState("");
  const [revenue, setRevenue] = useState(80000);
  const [customRev, setCustomRev] = useState("");
  const [regime, setRegime] = useState("NPD");
  const [email, setEmail] = useState("");

  async function submit() {
    setBusy(true);
    setErr("");
    try {
      const sphereFinal =
        sphere === "Другое" && sphereCustom.trim()
          ? sphereCustom.trim()
          : sphere;
      const rev = customRev ? Number(customRev.replace(/\s/g, "")) : revenue;
      const res = await registerAccount({
        display_name: name.trim(),
        phone: phone.trim(),
        email: email.trim() || undefined,
        business_sphere: sphereFinal,
        city: city.trim() || "Город",
        monthly_revenue_estimate: rev > 0 ? rev : 50000,
        tax_regime: regime,
      });
      setSession(res.token, res.profile.user_id);
      router.push("/app");
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось создать кабинет");
      setBusy(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen max-w-md flex-col px-4 py-8">
      <Link href="/" className="text-sm text-ink/50 hover:text-ink">
        ← На сайт
      </Link>
      <p className="mt-6 text-sm font-bold text-ink">Альфа-Бизнес: Старт</p>
      <p className="mt-4 text-xs font-semibold uppercase tracking-wide text-brand">
        Шаг {step} из 3
      </p>
      <h1 className="mt-2 text-3xl font-bold tracking-tight">Открыть кабинет</h1>
      <p className="mt-2 text-sm text-ink/60">
        Для тех, кто начинает бизнес в 18–25. Укажите данные — кабинет откроется
        сразу.
      </p>

      <div className="mt-4 flex gap-1.5">
        {[1, 2, 3].map((n) => (
          <span
            key={n}
            className={`h-1.5 flex-1 rounded-full ${n <= step ? "bg-brand" : "bg-platinum"}`}
          />
        ))}
      </div>

      {step === 1 ? (
        <div className="mt-8 space-y-4">
          <label className="block text-sm">
            <span className="text-ink/60">Как вас зовут</span>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Например, Алина"
              className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 outline-none focus:ring-2 focus:ring-brand/30"
            />
          </label>
          <label className="block text-sm">
            <span className="text-ink/60">Телефон — ваш логин</span>
            <input
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="+7 900 000-00-00"
              className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 outline-none focus:ring-2 focus:ring-brand/30"
            />
          </label>
          <button
            type="button"
            disabled={!name.trim() || phone.replace(/\D/g, "").length < 11}
            onClick={() => setStep(2)}
            className="btn-touch mt-4 w-full rounded-pill bg-brand py-3.5 text-sm font-semibold text-white disabled:opacity-40"
          >
            Дальше
          </button>
        </div>
      ) : null}

      {step === 2 ? (
        <div className="mt-8 space-y-4">
          <p className="text-sm font-medium text-ink/70">Чем занимаетесь</p>
          <div className="flex flex-wrap gap-2">
            {spheres.map((s) => (
              <button
                key={s}
                type="button"
                onClick={() => setSphere(s)}
                className={`btn-touch rounded-pill px-3 py-2 text-sm font-semibold ${
                  sphere === s ? "bg-ink text-white" : "bg-white text-ink/70 shadow-soft"
                }`}
              >
                {s}
              </button>
            ))}
          </div>
          {sphere === "Другое" ? (
            <input
              value={sphereCustom}
              onChange={(e) => setSphereCustom(e.target.value)}
              placeholder="Своя сфера"
              className="w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-sm outline-none"
            />
          ) : null}
          <label className="block text-sm">
            <span className="text-ink/60">Город</span>
            <input
              value={city}
              onChange={(e) => setCity(e.target.value)}
              placeholder="Казань, Москва…"
              className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 outline-none focus:ring-2 focus:ring-brand/30"
            />
          </label>
          <p className="text-sm font-medium text-ink/70">Оборот в месяц ≈</p>
          <div className="flex flex-wrap gap-2">
            {revenues.map((r) => (
              <button
                key={r.label}
                type="button"
                onClick={() => {
                  setRevenue(r.value);
                  setCustomRev("");
                }}
                className={`btn-touch rounded-pill px-3 py-2 text-sm font-semibold ${
                  !customRev && revenue === r.value
                    ? "bg-ink text-white"
                    : "bg-white text-ink/70 shadow-soft"
                }`}
              >
                {r.label}
              </button>
            ))}
          </div>
          <input
            value={customRev}
            onChange={(e) => setCustomRev(e.target.value)}
            placeholder="Или своя сумма, ₽"
            className="w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-sm outline-none"
          />
          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={() => setStep(1)}
              className="btn-touch flex-1 rounded-pill bg-canvas py-3 text-sm font-semibold"
            >
              Назад
            </button>
            <button
              type="button"
              disabled={!city.trim()}
              onClick={() => setStep(3)}
              className="btn-touch flex-1 rounded-pill bg-brand py-3 text-sm font-semibold text-white disabled:opacity-40"
            >
              Дальше
            </button>
          </div>
        </div>
      ) : null}

      {step === 3 ? (
        <div className="mt-8 space-y-4">
          <p className="text-sm font-medium text-ink/70">Налоговый режим</p>
          {(
            [
              ["NPD", "НПД — самозанятость"],
              ["UNDECIDED", "Пока думаю"],
              ["IP", "ИП"],
            ] as const
          ).map(([k, label]) => (
            <button
              key={k}
              type="button"
              onClick={() => setRegime(k)}
              className={`btn-touch flex w-full items-center rounded-2xl px-4 py-3 text-left text-sm font-semibold ${
                regime === k ? "bg-mint ring-2 ring-emerald-600/30" : "bg-white shadow-soft"
              }`}
            >
              {label}
            </button>
          ))}
          <label className="block text-sm">
            <span className="text-ink/60">Email (по желанию)</span>
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              type="email"
              placeholder="you@mail.ru"
              className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 outline-none"
            />
          </label>
          {err ? <p className="text-sm text-brand">{err}</p> : null}
          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={() => setStep(2)}
              className="btn-touch flex-1 rounded-pill bg-canvas py-3 text-sm font-semibold"
            >
              Назад
            </button>
            <button
              type="button"
              disabled={busy}
              onClick={() => void submit()}
              className="btn-touch flex-1 rounded-pill bg-brand py-3 text-sm font-semibold text-white disabled:opacity-50"
            >
              {busy ? "Создаём…" : "Создать кабинет"}
            </button>
          </div>
        </div>
      ) : null}

      <p className="mt-8 text-center text-sm text-ink/45">
        Уже есть аккаунт?{" "}
        <Link href="/login" className="font-semibold text-brand hover:underline">
          Войти
        </Link>
      </p>
      <button
        type="button"
        onClick={() => {
          enterExampleSession();
          router.push("/app");
        }}
        className="mt-3 text-center text-sm text-ink/40 hover:text-ink/70"
      >
        Посмотреть пример кабинета →
      </button>
    </main>
  );
}
