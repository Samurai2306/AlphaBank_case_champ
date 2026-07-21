"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { motion, useReducedMotion } from "framer-motion";
import {
  ArrowLeft,
  CheckCircle2,
  CircleDollarSign,
  Percent,
  Target,
  Zap,
} from "lucide-react";
import { getMe } from "@/lib/api";

type Stage = {
  id: number;
  name: string;
  short: string;
  blurb: string;
  reward: string;
  xp: number;
  icon: string;
  tone: string;
  href: string;
  cta: string;
  checklist: string[];
};

const stages: Stage[] = [
  {
    id: 1,
    name: "Идея",
    short: "Профиль",
    blurb: "Собрать картину бизнеса: сфера, город, оборот — чтобы расчёты шли по вашим цифрам.",
    reward: "Карточка профиля",
    xp: 20,
    icon: "/icons/3d/07-ai-onboarding.png",
    tone: "bg-cyanSoft",
    href:
      "/app/chat?q=" +
      encodeURIComponent(
        "Привет! Я делаю маникюр на дому в Казани, где-то 180 тысяч в месяц, пока без ИП.",
      ),
    cta: "Заполнить профиль",
    checklist: ["Назвать сферу и город", "Оценить оборот", "Подтвердить резюме в чате"],
  },
  {
    id: 2,
    name: "Безубыточность",
    short: "Цифры",
    blurb: "Понять, сколько клиентов в день нужно, чтобы аренда и расходники окупались.",
    reward: "Точка безубыточности",
    xp: 20,
    icon: "/icons/3d/05-unit-economics.png",
    tone: "bg-peach",
    href:
      "/app/chat?q=" +
      encodeURIComponent(
        "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
      ),
    cta: "Посчитать безубыточность",
    checklist: ["Ввести аренду и цену", "Получить клиентов/день", "Сверить с фактом"],
  },
  {
    id: 3,
    name: "Легализация",
    short: "НПД",
    blurb: "Выйти из серой зоны: режим НПД и понятная сумма к отложению на налог.",
    reward: "Налоговая карта",
    xp: 20,
    icon: "/icons/3d/15-legalization-npd.png",
    tone: "bg-mint",
    href: "/app/chat?q=" + encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
    cta: "Рассчитать налог",
    checklist: ["Выбрать режим НПД", "Посчитать смешанную ставку", "Увидеть сумму с запасом"],
  },
  {
    id: 4,
    name: "РКО / платёж",
    short: "Платёжка",
    blurb: "Сформировать поручение на налог и принять его к исполнению в банке.",
    reward: "Платёжное поручение",
    xp: 20,
    icon: "/icons/3d/10-payment-draft.png",
    tone: "bg-rose",
    href: "/app/chat?q=" + encodeURIComponent("Сформируй платёжку."),
    cta: "Собрать платёжку",
    checklist: ["Создать черновик", "Проверить назначение", "Подтвердить поручение"],
  },
  {
    id: 5,
    name: "ЕНП-копилка",
    short: "Копилка",
    blurb: "Включить автоотчисления, чтобы к сроку налог уже лежал отдельно от оборота.",
    reward: "Автокопилка ЕНП",
    xp: 20,
    icon: "/icons/3d/09-enp-piggy-bank.png",
    tone: "bg-lavender",
    href: "/app/piggy",
    cta: "Настроить копилку",
    checklist: ["Включить автоотчисления", "Выбрать ставку 1–15%", "Дойти до цели месяца"],
  },
];

const checkIcons = [Zap, Percent, Target];

const STORAGE_KEY = "alfa_start_journey_checks";

function checklistHits(stage: Stage, level: number, checks: Record<string, boolean>) {
  if (stage.id < level) return stage.checklist.length;
  return stage.checklist.filter((_, i) => checks[`${stage.id}:${i}`]).length;
}

export default function JourneyPage() {
  const [level, setLevel] = useState(4);
  const [checks, setChecks] = useState<Record<string, boolean>>({});
  const [celebrate, setCelebrate] = useState(false);
  const reduce = useReducedMotion();

  useEffect(() => {
    getMe()
      .then((p) => setLevel(p.cjm_level || 4))
      .catch(() => setLevel(4));
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) setChecks(JSON.parse(raw) as Record<string, boolean>);
    } catch {
      /* ignore */
    }
  }, []);

  const progress = useMemo(() => {
    const doneLevels = Math.max(0, level - 1);
    const current = stages.find((s) => s.id === level);
    let micro = 0;
    if (current) {
      micro = checklistHits(current, level, checks) / current.checklist.length;
      if (current.id < level) micro = 1;
    }
    return Math.min(100, Math.round(((doneLevels + micro) / stages.length) * 100));
  }, [level, checks]);

  const xp = useMemo(() => {
    let total = 0;
    for (const s of stages) {
      if (s.id < level) total += s.xp;
      else if (s.id === level) {
        const hit = checklistHits(s, level, checks);
        total += Math.round((hit / s.checklist.length) * s.xp);
      }
    }
    return total;
  }, [level, checks]);

  function toggleCheck(stageId: number, idx: number) {
    const key = `${stageId}:${idx}`;
    setChecks((prev) => {
      const next = { ...prev, [key]: !prev[key] };
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
      } catch {
        /* ignore */
      }
      const stage = stages.find((s) => s.id === stageId);
      if (stage && stage.checklist.every((_, i) => next[`${stageId}:${i}`])) {
        setCelebrate(true);
        window.setTimeout(() => setCelebrate(false), 1800);
      }
      return next;
    });
  }

  const current = stages.find((s) => s.id === level) ?? stages[3];

  return (
    <main className="mx-auto max-w-3xl px-4 py-6 pb-10">
      <Link
        href="/app"
        className="inline-flex items-center gap-1.5 text-sm text-ink/50 hover:text-ink"
      >
        <ArrowLeft className="h-4 w-4" aria-hidden />
        Главная
      </Link>

      <section className="relative mt-3 overflow-hidden rounded-[32px] bg-ink p-6 text-white shadow-soft">
        <div className="pointer-events-none absolute -right-8 -top-8 h-40 w-40 rounded-full bg-brand/40 blur-3xl" />
        <div className="relative z-10 flex items-start justify-between gap-4">
          <div className="min-w-0">
            <p className="text-sm text-white/60">Трекборд · путь предпринимателя</p>
            <h1 className="mt-1 text-2xl font-bold sm:text-3xl">Мой путь предпринимателя</h1>
            <p className="mt-2 max-w-md text-sm text-white/70">
              Пять уровней от идеи до копилки ЕНП. Сейчас — этап {level}: {current.name}.
            </p>
          </div>
          <div className="relative h-24 w-24 shrink-0">
            <Image
              src="/icons/3d/13-cjm-journey.png"
              alt=""
              fill
              className="object-contain drop-shadow-lg"
              sizes="96px"
            />
          </div>
        </div>

        <div className="relative z-10 mt-6">
          <div className="flex items-end justify-between text-sm">
            <span className="font-semibold tabular">{progress}% пути</span>
            <span className="text-white/60 tabular">{xp} / 100 XP</span>
          </div>
          <div className="mt-2 h-3 overflow-hidden rounded-pill bg-white/15">
            <motion.div
              className="h-full rounded-pill bg-gradient-to-r from-brand to-[#ff7a6e]"
              initial={false}
              animate={{ width: `${progress}%` }}
              transition={reduce ? { duration: 0 } : { type: "spring", stiffness: 120, damping: 20 }}
            />
          </div>
        </div>

        <ol className="relative z-10 mt-6 flex justify-between gap-1">
          {stages.map((s) => {
            const done = s.id < level;
            const cur = s.id === level;
            return (
              <li key={s.id} className="flex flex-1 flex-col items-center gap-1.5">
                <span
                  className={`flex h-9 w-9 items-center justify-center rounded-full text-xs font-bold transition ${
                    done
                      ? "bg-emerald-400 text-ink"
                      : cur
                        ? "bg-brand text-white ring-4 ring-brand/30"
                        : "bg-white/15 text-white/50"
                  }`}
                >
                  {done ? "✓" : s.id}
                </span>
                <span
                  className={`text-[10px] font-medium sm:text-xs ${
                    cur ? "text-white" : "text-white/45"
                  }`}
                >
                  {s.short}
                </span>
              </li>
            );
          })}
        </ol>
      </section>

      {celebrate ? (
        <p
          className="mt-3 rounded-card bg-mint px-4 py-3 text-sm font-medium text-emerald-800"
          aria-live="polite"
        >
          Чеклист закрыт — награда «{current.reward}» почти ваша. Завершите шаг в чате.
        </p>
      ) : null}

      <motion.section
        key={current.id}
        initial={reduce ? false : { opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        className={`relative mt-5 min-h-[200px] overflow-hidden rounded-[32px] ${current.tone} p-5 shadow-soft`}
      >
        <div className="relative z-10 max-w-[68%]">
          <p className="text-xs font-semibold uppercase tracking-wide text-ink/45">
            Вы здесь · уровень {current.id}
          </p>
          <h2 className="mt-1 text-2xl font-bold">{current.name}</h2>
          <p className="mt-2 text-sm text-ink/70">{current.blurb}</p>
          <p className="mt-3 inline-flex items-center gap-1.5 rounded-pill bg-white/80 px-3 py-1 text-xs font-medium text-ink/70">
            <CircleDollarSign className="h-3.5 w-3.5 text-brand" aria-hidden />
            Награда: {current.reward}
          </p>
          <div className="mt-4">
            <Link
              href={current.href}
              className="inline-flex rounded-pill bg-brand px-5 py-2.5 text-sm font-semibold text-white"
            >
              {current.cta}
            </Link>
          </div>
        </div>
        <div className="pointer-events-none absolute -bottom-4 -right-2 h-[120%] w-[52%]">
          <Image
            src={current.icon}
            alt=""
            fill
            className="object-contain object-bottom opacity-95"
            sizes="220px"
          />
        </div>
      </motion.section>

      <section className="mt-5 rounded-[28px] bg-white p-5 shadow-soft">
        <h3 className="font-bold">Чеклист текущего этапа</h3>
        <p className="mt-1 text-sm text-ink/55">
          Отмечайте шаги — прогресс и XP обновятся сразу.
        </p>
        <ul className="mt-4 space-y-2">
          {current.checklist.map((item, idx) => {
            const key = `${current.id}:${idx}`;
            const on = !!checks[key];
            const Icon = checkIcons[idx % checkIcons.length];
            return (
              <li key={key}>
                <button
                  type="button"
                  onClick={() => toggleCheck(current.id, idx)}
                  className={`btn-touch flex w-full items-center gap-3 rounded-2xl px-3 py-3 text-left text-sm ${
                    on ? "bg-mint" : "bg-canvas"
                  }`}
                >
                  <span
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full ${
                      on ? "bg-emerald-600 text-white" : "bg-white text-brand ring-1 ring-platinum"
                    }`}
                  >
                    {on ? (
                      <CheckCircle2 className="h-4 w-4" aria-hidden />
                    ) : (
                      <Icon className="h-4 w-4" aria-hidden />
                    )}
                  </span>
                  <span className={on ? "text-ink line-through decoration-ink/30" : ""}>
                    {item}
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      </section>

      <section className="mt-6">
        <h3 className="text-lg font-bold">Все уровни</h3>
        <ol className="mt-3 space-y-3">
          {stages.map((s) => {
            const done = s.id < level;
            const cur = s.id === level;
            const locked = s.id > level;
            const hit = checklistHits(s, level, checks);
            return (
              <li
                key={s.id}
                className={`relative overflow-hidden rounded-[28px] p-4 shadow-soft transition ${
                  cur
                    ? `ring-2 ring-brand/35 ${s.tone}`
                    : done
                      ? s.tone
                      : "bg-platinum/25"
                }`}
              >
                <div className="relative z-10 max-w-[72%]">
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <p className="text-xs text-ink/45">Уровень {s.id}</p>
                      <p className="font-bold">{s.name}</p>
                    </div>
                    <span
                      className={`shrink-0 rounded-pill px-2.5 py-1 text-[11px] font-semibold ${
                        done
                          ? "bg-emerald-600/15 text-emerald-800"
                          : cur
                            ? "bg-brand text-white"
                            : "bg-white/80 text-ink/45"
                      }`}
                    >
                      {done ? "Пройдено" : cur ? "Сейчас" : "Закрыто"}
                    </span>
                  </div>
                  <p className="mt-1 text-sm text-ink/60">{s.blurb}</p>
                  <p className="mt-2 text-xs text-ink/45">
                    Награда: {s.reward}
                    {cur || done
                      ? ` · чеклист ${hit}/${s.checklist.length}`
                      : ""}
                  </p>
                  {cur ? (
                    <Link
                      href={s.href}
                      className="mt-3 inline-flex rounded-pill bg-brand px-4 py-2 text-xs font-semibold text-white"
                    >
                      {s.cta}
                    </Link>
                  ) : done ? (
                    <Link
                      href={s.href}
                      className="mt-3 inline-flex rounded-pill bg-ink/90 px-4 py-2 text-xs font-semibold text-white"
                    >
                      Открыть снова
                    </Link>
                  ) : (
                    <p className="mt-2 text-xs text-ink/40">
                      Откроется после уровня {s.id - 1}
                    </p>
                  )}
                </div>
                <div
                  className={`pointer-events-none absolute -bottom-4 -right-2 h-32 w-32 ${
                    locked ? "opacity-30 grayscale" : "opacity-90"
                  }`}
                >
                  <Image
                    src={s.icon}
                    alt=""
                    fill
                    className="object-contain object-bottom"
                    sizes="128px"
                  />
                </div>
              </li>
            );
          })}
        </ol>
      </section>

      <section className="mt-6 grid grid-cols-2 gap-3">
        <Link href="/app/tasks" className="card-touch rounded-card bg-white p-4 shadow-soft">
          <p className="font-semibold">Рекомендуемые задачи</p>
          <p className="mt-1 text-sm text-ink/55">Быстрый список действий</p>
        </Link>
        <Link
          href="/app/chat"
          className="card-touch rounded-card bg-brand p-4 text-white shadow-soft"
        >
          <p className="font-semibold">Спросить в чате</p>
          <p className="mt-1 text-sm text-white/80">Помощь на текущем этапе</p>
        </Link>
      </section>
    </main>
  );
}
