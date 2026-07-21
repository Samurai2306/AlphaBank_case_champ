"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { motion, useReducedMotion } from "framer-motion";
import {
  Bell,
  FileText,
  ListChecks,
  PiggyBank,
  Route,
  Wallet,
} from "lucide-react";
import { CashbackWheel } from "@/components/CashbackWheel";
import { EmptyHint } from "@/components/EmptyHint";
import { IconTile } from "@/components/IconTile";
import { formatRub, getHome, getTasks, getTransactions } from "@/lib/api";

type HomeData = Awaited<ReturnType<typeof getHome>>;

const services = [
  {
    title: "Рассчитать налог",
    sub: "НПД 4%/6%",
    icon: "/icons/3d/02-tax-calculator.png",
    tone: "bg-mint",
    href:
      "/app/chat?q=" +
      encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
  },
  {
    title: "Безубыточность",
    sub: "Клиенты в день",
    icon: "/icons/3d/05-unit-economics.png",
    tone: "bg-peach",
    href:
      "/app/chat?q=" +
      encodeURIComponent(
        "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
      ),
  },
  {
    title: "Светофор 115-ФЗ",
    sub: "Проверка ИНН",
    icon: "/icons/3d/06-compliance-risk.png",
    tone: "bg-lavender",
    href: "/app/chat?q=" + encodeURIComponent("Проверь ИНН 1650987654"),
  },
  {
    title: "Мой путь",
    sub: "Трекборд",
    icon: "/icons/3d/13-cjm-journey.png",
    tone: "bg-rose",
    href: "/app/journey",
  },
];

export default function HomePage() {
  const [data, setData] = useState<HomeData | null>(null);
  const [tasks, setTasks] = useState<
    Awaited<ReturnType<typeof getTasks>>["items"]
  >([]);
  const [txns, setTxns] = useState<
    Awaited<ReturnType<typeof getTransactions>>["items"]
  >([]);
  const [err, setErr] = useState("");
  const [seg, setSeg] = useState<"today" | "month" | "path">("month");
  const reduce = useReducedMotion();

  useEffect(() => {
    getHome()
      .then(setData)
      .catch((e: Error) => setErr(e.message));
    getTasks()
      .then((r) => setTasks(r.items.slice(0, 3)))
      .catch(() => undefined);
    getTransactions()
      .then((r) => setTxns(r.items))
      .catch(() => undefined);
  }, []);

  const hero = useMemo(() => {
    if (!data) {
      return { title: "К отложению в этом месяце", value: null as number | null, hint: "Загрузка…" };
    }
    if (seg === "today") {
      return {
        title: "В копилке сейчас",
        value: data.piggy.balance,
        hint: `Последний вклад ~${formatRub(data.piggy.last_contribution)} ₽ · ставка ${data.piggy.rate_percent}%`,
      };
    }
    if (seg === "path") {
      return {
        title: "Уровень пути предпринимателя",
        value: data.profile.cjm_level,
        hint: `${data.profile.business_sphere} · ${data.profile.city} · следующий шаг на трекборде`,
        path: true as const,
      };
    }
    return {
      title: "К отложению в этом месяце",
      value: data.tax_due,
      hint: data.tax_estimate
        ? `Оценка ~${formatRub(data.tax_estimate)} ₽ (${data.tax_rate_label ?? "смеш."}) · с запасом ${formatRub(data.tax_due)} ₽`
        : `${data.profile.business_sphere} · ${data.profile.city}`,
    };
  }, [data, seg]);

  return (
    <main className="mx-auto max-w-3xl px-4 py-6 pb-8">
      <header className="flex items-center justify-between gap-3">
        <div>
          <p className="text-sm text-ink/50">Альфа-Бизнес: Старт · кабинет</p>
          <h1 className="text-2xl font-bold">
            Привет{data ? `, ${data.profile.display_name}` : ""}
          </h1>
          <p className="mt-1 text-xs text-ink/45">
            Сводка · налог · копилка · быстрые сервисы
          </p>
        </div>
        <div className="flex gap-2">
          <Link
            href="/app/alerts"
            className="inline-flex items-center gap-1 rounded-pill bg-white px-3 py-2 text-xs font-semibold text-brand shadow-soft"
            title="Алерты"
          >
            <Bell className="h-3.5 w-3.5" aria-hidden />
            Алерты
          </Link>
          <Link
            href="/app/report"
            className="inline-flex items-center gap-1 rounded-pill bg-white px-3 py-2 text-xs font-semibold text-ink/70 shadow-soft"
            title="Отчёт за месяц"
          >
            <FileText className="h-3.5 w-3.5" aria-hidden />
            Отчёт
          </Link>
        </div>
      </header>

      <nav
        className="mt-4 flex gap-2 overflow-x-auto pb-1"
        aria-label="Быстрые разделы"
      >
        {(
          [
            ["/app/journey", "Путь", Route],
            ["/app/piggy", "Копилка", PiggyBank],
            ["/app/payments", "Платежи", Wallet],
            ["/app/tasks", "Задачи", ListChecks],
          ] as const
        ).map(([href, label, Icon]) => (
          <Link
            key={href}
            href={href}
            className="btn-touch inline-flex shrink-0 items-center gap-1.5 rounded-pill bg-white px-3 py-2 text-xs font-semibold text-ink/70 shadow-soft"
          >
            <Icon className="h-3.5 w-3.5 text-brand" aria-hidden />
            {label}
          </Link>
        ))}
      </nav>

      <div className="mt-4 inline-flex rounded-pill bg-platinum/40 p-1 text-sm">
        {(
          [
            ["today", "Сегодня"],
            ["month", "Месяц"],
            ["path", "Путь"],
          ] as const
        ).map(([k, label]) => (
          <button
            key={k}
            type="button"
            onClick={() => setSeg(k)}
            className={`btn-touch rounded-pill px-4 py-1.5 ${
              seg === k ? "bg-ink text-white" : "text-ink/70 hover:bg-white/70"
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      {err ? (
        <p className="mt-6 rounded-card bg-rose p-4 text-sm">
          API недоступен. Запустите сервис на :8080. {err}
        </p>
      ) : null}

      {data?.empty_cabinet ? (
        <div className="mt-5">
          <EmptyHint
            title="Кабинет только открыт"
            text="Операций пока нет — это нормально. Налог ниже считается от оборота в профиле. Уточните цифры и посчитайте налог в чате, чтобы появились сценарии и копилка."
            primary={{ href: "/app/profile", label: "Профиль" }}
            secondary={{
              href:
                "/app/chat?q=" +
                encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
              label: "Рассчитать налог",
            }}
          />
        </div>
      ) : null}

      <motion.section
        key={seg}
        initial={reduce ? false : { opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative mt-5 overflow-hidden rounded-[32px] bg-white p-6 shadow-soft"
      >
        <p className="text-sm text-ink/60">{hero.title}</p>
        <p className="mt-2 text-5xl font-bold tabular tracking-tight">
          {hero.value == null
            ? "—"
            : "path" in hero && hero.path
              ? `${hero.value} / 5`
              : `${formatRub(hero.value)} ₽`}
        </p>
        <p className="mt-2 text-sm text-ink/55">
          {hero.hint}
          {data?.empty_cabinet && seg === "month"
            ? " · от оборота в профиле"
            : ""}
        </p>
        {seg === "path" ? (
          <Link
            href="/app/journey"
            className="mt-5 inline-flex rounded-pill bg-ink px-5 py-2.5 text-sm font-semibold text-white"
          >
            Открыть трекборд
          </Link>
        ) : (
          <Link
            href={
              "/app/chat?q=" +
              encodeURIComponent(
                seg === "today"
                  ? "Включи копилку 6%."
                  : data?.empty_cabinet
                    ? "Сколько мне отложить на налог за этот месяц?"
                    : "Сформируй платёжку.",
              )
            }
            className="mt-5 inline-flex rounded-pill bg-brand px-5 py-2.5 text-sm font-semibold text-white"
          >
            {seg === "today"
              ? "Управлять копилкой"
              : data?.empty_cabinet
                ? "Рассчитать налог"
                : "Сформировать платёжку"}
          </Link>
        )}
        {seg === "month" ? (
          <div className="pointer-events-none absolute -bottom-8 -right-6 h-40 w-40 opacity-95">
            <Image
              src="/icons/3d/02-tax-calculator.png"
              alt=""
              fill
              className="object-contain object-bottom"
              sizes="160px"
            />
          </div>
        ) : null}
      </motion.section>

      {data?.tax_due_date != null ? (
        <Link
          href="/app/payments"
          className="card-touch mt-4 flex items-center justify-between rounded-card bg-peach px-5 py-4 shadow-soft ring-1 ring-black/[0.04]"
        >
          <div>
            <p className="text-sm font-semibold">Календарь налога</p>
            <p className="mt-1 text-sm text-ink/65">
              Срок {new Date(data.tax_due_date).toLocaleDateString("ru-RU")}
              {data.tax_days_left != null
                ? ` · осталось ${data.tax_days_left} дн.`
                : ""}
            </p>
          </div>
          <p className="text-xl font-bold tabular">{formatRub(data.tax_due)} ₽</p>
        </Link>
      ) : null}

      <CashbackWheel
        userId={data?.profile.user_id}
        transactions={txns}
        emptyCabinet={data?.empty_cabinet}
      />

      <div className="mt-4 grid grid-cols-2 gap-3">
        <IconTile
          href="/app/piggy"
          title="Копилка ЕНП"
          sub={
            data
              ? `${formatRub(data.piggy.balance)} ₽ · ${data.piggy.rate_percent}%`
              : "Загрузка…"
          }
          icon="/icons/3d/09-enp-piggy-bank.png"
          tone="bg-mint"
          size="md"
        />
        <IconTile
          href="/app/report"
          title="Отчёт"
          sub="Доходы и расходы"
          icon="/icons/3d/05-unit-economics.png"
          tone="bg-lavender"
          size="md"
        />
      </div>

      <Link
        href="/app/chat"
        className="card-touch relative mt-4 flex min-h-[120px] items-center overflow-hidden rounded-[32px] bg-white px-5 py-4 shadow-soft ring-1 ring-platinum/80"
      >
        <div className="relative z-10 max-w-[70%]">
          <p className="font-semibold">Спросить в чате</p>
          <p className="text-sm text-ink/55">Налог, безубыточность, платёжка, ИНН</p>
          <span className="mt-3 inline-flex rounded-pill bg-brand px-4 py-2 text-sm font-semibold text-white">
            Чат
          </span>
        </div>
        <div className="pointer-events-none absolute -bottom-6 -right-4 h-40 w-40">
          <Image
            src="/icons/3d/14-copilot-chat.png"
            alt=""
            fill
            className="object-contain object-bottom"
            sizes="160px"
          />
        </div>
      </Link>

      {tasks.length > 0 ? (
        <section className="mt-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-bold">С чего начать</h2>
            <Link href="/app/tasks" className="text-sm font-medium text-brand hover:underline">
              Все →
            </Link>
          </div>
          <ul className="mt-3 space-y-2">
            {tasks.map((t) => (
              <li key={t.id}>
                <Link
                  href={t.href}
                  className="card-touch block rounded-card bg-white px-4 py-3 text-sm shadow-soft"
                >
                  <span className="font-medium">{t.label}</span>
                  <span className="mt-0.5 block text-ink/50">{t.text}</span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      <section className="mt-6">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-bold">Сервисы</h2>
          <Link href="/app/services" className="text-sm font-medium text-brand hover:underline">
            Все →
          </Link>
        </div>
        <div className="mt-3 grid grid-cols-2 gap-3">
          {services.map((s) => (
            <IconTile
              key={s.title}
              href={s.href}
              title={s.title}
              sub={s.sub}
              icon={s.icon}
              tone={s.tone}
              size="md"
            />
          ))}
        </div>
      </section>

      {(data?.insights?.length ?? 0) > 0 ? (
        <section className="mt-6">
          <h2 className="text-lg font-bold">Подсказки для бизнеса</h2>
          <div className="mt-3 space-y-3">
            {data!.insights.map((i) => (
              <Link
                key={i.id}
                href={i.href}
                className="card-touch block rounded-card bg-white p-4 shadow-soft"
              >
                <p className="text-sm leading-snug">{i.text}</p>
                {i.why ? (
                  <p className="mt-1 text-xs leading-snug text-ink/50">{i.why}</p>
                ) : null}
                <p className="mt-2 text-sm font-medium text-brand">{i.cta} →</p>
              </Link>
            ))}
          </div>
        </section>
      ) : null}
    </main>
  );
}
