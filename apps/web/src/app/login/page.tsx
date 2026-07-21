"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useState } from "react";
import { loginUser } from "@/lib/api";
import { enterExampleSession, setSession } from "@/lib/session";

function LoginInner() {
  const router = useRouter();
  const params = useSearchParams();
  const next = params.get("next") || "/app";
  const [busy, setBusy] = useState(false);
  const [phone, setPhone] = useState("+7 ");
  const [err, setErr] = useState("");

  async function enter() {
    setBusy(true);
    setErr("");
    try {
      const res = await loginUser(phone.trim());
      setSession(res.token, res.profile.user_id);
      router.push(next.startsWith("/") ? next : "/app");
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось войти");
      setBusy(false);
    }
  }

  function openExample() {
    enterExampleSession();
    router.push("/app");
  }

  return (
    <main className="relative mx-auto flex min-h-screen max-w-md flex-col justify-center overflow-hidden px-4 py-10">
      <div className="pointer-events-none absolute -right-4 top-16 h-40 w-40 opacity-80">
        <Image
          src="/icons/3d/17-alfa-badge-pastel.png"
          alt=""
          fill
          className="object-contain drop-shadow-sm"
          sizes="160px"
          priority
        />
      </div>

      <Link href="/" className="relative z-10 text-sm text-ink/50 hover:text-ink">
        ← На сайт
      </Link>
      <div className="relative z-10 mt-8 flex items-center gap-2">
        <span className="flex h-10 w-10 items-center justify-center rounded-full bg-brand text-sm font-bold text-white shadow-[0_0_20px_rgba(239,49,36,0.35)]">
          A
        </span>
        <div>
          <p className="font-bold">Альфа-Бизнес: Старт</p>
          <p className="text-sm text-ink/55">Кабинет для первого бизнеса</p>
        </div>
      </div>
      <h1 className="relative z-10 mt-10 text-3xl font-bold">Вход</h1>
      <p className="relative z-10 mt-2 text-sm text-ink/60">
        Войдите по телефону, который указали при регистрации. Без SMS на стенде.
      </p>

      <label className="relative z-10 mt-6 block text-sm">
        <span className="text-ink/60">Телефон — ваш логин</span>
        <input
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-ink outline-none focus:ring-2 focus:ring-brand/30"
        />
      </label>

      {err ? (
        <p className="relative z-10 mt-3 text-sm text-brand">
          {err}{" "}
          <Link href="/register" className="font-semibold underline">
            Зарегистрироваться
          </Link>
        </p>
      ) : null}

      <button
        type="button"
        disabled={busy || phone.replace(/\D/g, "").length < 11}
        onClick={() => void enter()}
        className="btn-touch relative z-10 mt-6 w-full rounded-pill bg-brand py-3.5 text-sm font-semibold text-white disabled:opacity-50"
      >
        {busy ? "Входим…" : "Войти в кабинет"}
      </button>

      <button
        type="button"
        onClick={openExample}
        className="btn-touch relative z-10 mt-3 w-full rounded-pill bg-white py-3.5 text-sm font-semibold text-ink shadow-soft ring-1 ring-platinum"
      >
        Посмотреть пример
      </button>

      <p className="relative z-10 mt-6 text-center text-sm text-ink/45">
        Нет аккаунта?{" "}
        <Link href="/register" className="font-semibold text-brand hover:underline">
          Открыть кабинет
        </Link>
      </p>
    </main>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center text-sm text-ink/50">
          Загрузка…
        </main>
      }
    >
      <LoginInner />
    </Suspense>
  );
}
