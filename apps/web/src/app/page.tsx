"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion, useReducedMotion } from "framer-motion";
import { IconTile } from "@/components/IconTile";
import { enterExampleSession } from "@/lib/session";

const features = [
  {
    title: "Налоги без страха",
    sub: "НПД 4%/6% с запасом",
    icon: "/icons/3d/02-tax-calculator.png",
    tone: "bg-mint",
    href: "/register",
  },
  {
    title: "Сканер договора",
    sub: "Флаги в договоре",
    icon: "/icons/3d/03-legal-scanner.png",
    tone: "bg-rose",
    href: "/register",
  },
  {
    title: "Точка безубыточности",
    sub: "Сколько клиентов в день",
    icon: "/icons/3d/05-unit-economics.png",
    tone: "bg-peach",
    href: "/register",
  },
  {
    title: "Проверка контрагента",
    sub: "ИНН и риски 115-ФЗ",
    icon: "/icons/3d/06-compliance-risk.png",
    tone: "bg-lavender",
    href: "/register",
  },
  {
    title: "Старт в чате",
    sub: "От идеи до НПД",
    icon: "/icons/3d/07-ai-onboarding.png",
    tone: "bg-cyanSoft",
    href: "/register",
  },
  {
    title: "Копилка на налог",
    sub: "Копит к сроку уплаты",
    icon: "/icons/3d/09-enp-piggy-bank.png",
    tone: "bg-mint",
    href: "/register",
  },
];

export default function LandingPage() {
  const reduce = useReducedMotion();
  const router = useRouter();

  function openExample() {
    enterExampleSession();
    router.push("/app");
  }

  return (
    <main className="min-h-screen overflow-hidden bg-ink text-white">
      <header className="relative z-20 mx-auto flex max-w-6xl items-center justify-between gap-3 px-4 py-5">
        <div className="flex min-w-0 items-center gap-2 text-base font-bold tracking-tight sm:text-lg">
          <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand text-sm text-white shadow-[0_0_24px_rgba(239,49,36,0.45)]">
            A
          </span>
          <span className="sm:hidden">Старт</span>
          <span className="hidden truncate sm:inline">Альфа-Бизнес: Старт</span>
        </div>
        <div className="flex shrink-0 gap-2">
          <Link
            href="/register"
            className="rounded-pill bg-brand px-3 py-2 text-xs font-semibold text-white sm:px-4 sm:text-sm"
          >
            Открыть кабинет
          </Link>
          <Link
            href="/login"
            className="rounded-pill bg-white/10 px-3 py-2 text-xs text-white ring-1 ring-white/15 hover:bg-white/20 sm:px-4 sm:text-sm"
          >
            Войти
          </Link>
        </div>
      </header>

      <section className="relative mx-auto max-w-6xl px-4 pb-14 pt-4">
        <div className="pointer-events-none absolute -left-20 top-10 h-64 w-64 rounded-full bg-brand/20 blur-3xl" />
        <div className="pointer-events-none absolute right-0 top-32 h-56 w-56 rounded-full bg-cyanSoft/15 blur-3xl" />

        <div className="relative z-10 grid items-stretch gap-5 lg:grid-cols-[1.1fr_0.9fr]">
          <motion.div
            initial={reduce ? false : { opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            className="relative overflow-hidden rounded-[36px] bg-gradient-to-br from-[#2a3038] via-[#171c22] to-[#3a1612] p-7 shadow-[0_30px_80px_rgba(0,0,0,0.35)] sm:p-10"
          >
            <div className="pointer-events-none absolute -right-16 -top-16 h-64 w-64 rounded-full bg-brand/30 blur-3xl" />
            <p className="text-sm font-medium text-white/70">
              Альфа-Бизнес: Старт · для тех, кто начинает в 18–25
            </p>
            <h1 className="mt-4 max-w-xl text-4xl font-bold leading-[1.05] tracking-tight text-white md:text-5xl lg:text-[3.4rem]">
              Банк, с которым начинается бизнес
            </h1>
            <p className="mt-5 max-w-lg text-base leading-relaxed text-white/85">
              Свой налог и копилка ЕНП за пару минут. Считаем НПД, проверяем
              договоры и контрагентов — в чате, без канцелярита.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link
                href="/register"
                className="btn-touch inline-flex rounded-pill bg-brand px-7 py-3.5 text-sm font-semibold text-white shadow-[0_8px_28px_rgba(239,49,36,0.45)]"
              >
                Открыть кабинет
              </Link>
              <Link
                href="/login"
                className="btn-touch inline-flex rounded-pill bg-white/15 px-7 py-3.5 text-sm font-semibold text-white ring-1 ring-white/35 hover:bg-white/25"
              >
                Войти
              </Link>
              <button
                type="button"
                onClick={openExample}
                className="btn-touch inline-flex rounded-pill px-5 py-3.5 text-sm font-medium text-white/70 underline-offset-4 hover:text-white hover:underline"
              >
                Посмотреть пример
              </button>
            </div>
          </motion.div>

          <motion.div
            initial={reduce ? false : { opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="relative overflow-hidden rounded-[36px] bg-gradient-to-b from-peach to-[#ffe8d6] p-5 text-ink sm:p-6"
          >
            <div className="mb-3 flex items-center justify-between px-1">
              <p className="text-sm font-bold">Кабинет Маши</p>
              <span className="h-2 w-2 animate-pulse rounded-full bg-brand" />
            </div>
            <div className="rounded-[28px] bg-white/70 p-4 shadow-soft backdrop-blur">
              <p className="text-xs text-ink/50">К отложению с запасом</p>
              <p className="mt-1 text-4xl font-bold tabular tracking-tight">10 800 ₽</p>
              <p className="mt-1 text-xs text-ink/50">оценка ~7 740 ₽ · срок 28 числа</p>
            </div>
            <div className="mt-3 grid grid-cols-2 gap-2.5">
              <IconTile
                href="/app?example=1"
                title="Копилка"
                sub="5 400 ₽"
                icon="/icons/3d/09-enp-piggy-bank.png"
                tone="bg-mint"
                size="sm"
              />
              <IconTile
                href="/app?example=1"
                title="Мой путь"
                sub="Уровень 4/5"
                icon="/icons/3d/13-cjm-journey.png"
                tone="bg-lavender"
                size="sm"
              />
            </div>
            <Image
              src="/icons/3d/18-copilot-hero-cluster.png"
              alt=""
              width={120}
              height={120}
              className="pointer-events-none absolute -bottom-2 -right-1 opacity-80 sm:opacity-90"
            />
          </motion.div>
        </div>
      </section>

      <section id="features" className="border-t border-white/10 bg-[#0e1217] py-14">
        <div className="mx-auto max-w-6xl px-4">
          <div className="flex flex-wrap items-end justify-between gap-3">
            <div>
              <h2 className="text-2xl font-bold text-white sm:text-3xl">
                Что умеет Альфа-Бизнес: Старт
              </h2>
              <p className="mt-2 text-sm text-white/55">
                Инструменты для первого бизнеса — в одном кабинете.
              </p>
            </div>
            <Link href="/register" className="text-sm font-medium text-brand hover:underline">
              Открыть кабинет →
            </Link>
          </div>
          <div className="mt-8 grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-3">
            {features.map((f, i) => (
              <motion.div
                key={f.title}
                initial={reduce ? false : { opacity: 0, y: 12 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: "-40px" }}
                transition={{ delay: i * 0.04 }}
              >
                <IconTile
                  href={f.href}
                  title={f.title}
                  sub={f.sub}
                  icon={f.icon}
                  tone={f.tone}
                  size="lg"
                />
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      <section className="mx-auto grid max-w-6xl gap-4 px-4 py-14 md:grid-cols-3">
        <article className="relative overflow-hidden rounded-[32px] bg-white p-7 text-ink md:col-span-2">
          <div className="relative z-10 max-w-[min(100%,28rem)] pr-2 sm:max-w-md sm:pr-0 md:max-w-lg">
            <p className="text-sm text-ink/50">Кейс в кабинете</p>
            <h3 className="mt-2 text-2xl font-bold sm:text-3xl">
              Мария Соколова · маникюр, Казань, НПД
            </h3>
            <p className="mt-3 text-ink/70">
              Пример заполненного кабинета: выписка, копилка, налог и путь. Свой
              кабинет лучше открыть с нуля — за пару минут.
            </p>
            <div className="mt-6 flex flex-wrap gap-2">
              <Link
                href="/register"
                className="btn-touch inline-flex rounded-pill bg-ink px-6 py-3 text-sm font-semibold text-white"
              >
                Открыть свой
              </Link>
              <button
                type="button"
                onClick={openExample}
                className="btn-touch inline-flex rounded-pill bg-canvas px-6 py-3 text-sm font-semibold text-ink ring-1 ring-platinum"
              >
                Пример Маши
              </button>
            </div>
          </div>
          <div className="pointer-events-none absolute -bottom-8 -right-8 h-40 w-40 opacity-90 sm:-bottom-6 sm:-right-2 sm:h-48 sm:w-48 md:h-56 md:w-56">
            <Image
              src="/icons/3d/11-beauty-niche.png"
              alt=""
              fill
              className="object-contain object-bottom"
              sizes="224px"
            />
          </div>
        </article>
        <article className="relative overflow-hidden rounded-[32px] bg-brand p-7 text-white">
          <p className="text-5xl font-bold tabular">180к</p>
          <p className="mt-2 text-sm text-white/90">оборот Маши в месяц</p>
          <p className="mt-4 text-sm text-white/75">
            10 800 ₽ к отложению · копилка уже 5 400 ₽
          </p>
          <Link
            href="/register"
            className="mt-8 inline-flex rounded-pill bg-ink px-5 py-2.5 text-sm font-semibold text-white"
          >
            Начать сейчас
          </Link>
        </article>
      </section>
    </main>
  );
}
