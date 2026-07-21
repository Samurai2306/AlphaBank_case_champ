"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "framer-motion";
import {
  ArrowDownLeft,
  ArrowUpRight,
  Bell,
  ChevronDown,
  FileText,
  Home,
  Megaphone,
  MessageCircle,
  Package,
  Percent,
  Scissors,
  Search,
  Smartphone,
  Sparkles,
  Wallet,
} from "lucide-react";
import { EmptyHint } from "@/components/EmptyHint";
import { PageHeader } from "@/components/PageHeader";
import {
  CAT_RU,
  cashbackForAmount,
  cashbackSeed,
  categoryCashbackTotal,
  getCategoryMeta,
  loadCashbackPick,
  rateForExpense,
  taxHintForTxn,
  taxOnIncome,
  type CashbackPick,
} from "@/lib/cashback";
import { formatRub, getMe, getTransactions } from "@/lib/api";
import {
  buildPaymentTips,
  chatHref,
  promptCategoryDeepDive,
  promptCutCategory,
  promptCutExpenses,
  promptGrowIncome,
  promptStatementOverview,
  promptTaxFromStatement,
  promptTxnDetail,
  promptTxnTax,
} from "@/lib/operations-prompts";
import { ensureUserId, getUserId } from "@/lib/session";

type Txn = Awaited<ReturnType<typeof getTransactions>>;
type Item = Txn["items"][number];

const catIcon: Record<string, typeof Scissors> = {
  services: Scissors,
  rent: Home,
  supplies: Package,
  ads: Megaphone,
  fees: Percent,
  ops: Package,
};

function dayKeyLocal(d: Date): string {
  return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;
}

function dayLabel(iso: string): string {
  const d = new Date(iso);
  const today = new Date();
  const yday = new Date();
  yday.setDate(today.getDate() - 1);
  const key = dayKeyLocal(d);
  if (key === dayKeyLocal(today)) return "Сегодня";
  if (key === dayKeyLocal(yday)) return "Вчера";
  return d.toLocaleDateString("ru-RU", {
    day: "numeric",
    month: "long",
  });
}

function TxnRow({
  t,
  seed,
  pick,
  expanded,
  onToggle,
  showCatBadge,
}: {
  t: Item;
  seed: string;
  pick: CashbackPick | null;
  expanded: boolean;
  onToggle: () => void;
  showCatBadge?: boolean;
}) {
  const Icon = catIcon[t.category] ?? Package;
  const rate = rateForExpense(t.category, seed, pick);
  const activePick = !!pick && pick.id === t.category;
  const cb =
    t.direction === "out" ? cashbackForAmount(t.amount, rate) : 0;
  const taxHint = taxHintForTxn(t);
  const taxAmt = t.direction === "in" ? taxOnIncome(t.amount) : 0;
  const reduce = useReducedMotion();

  return (
    <li className="overflow-hidden rounded-card bg-white shadow-soft">
      <div className="flex items-start gap-3 px-4 py-3">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-2xl bg-canvas text-brand">
          <Icon className="h-4 w-4" aria-hidden />
        </span>
        <div className="min-w-0 flex-1">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <p className="truncate font-medium">{t.counterparty_name}</p>
              <p className="text-xs text-ink/45">
                {new Date(t.booked_at).toLocaleDateString("ru-RU")}
                {showCatBadge ? (
                  <>
                    {" · "}
                    <span className="rounded-pill bg-canvas px-1.5 py-0.5 text-[11px] text-ink/60">
                      {CAT_RU[t.category] ?? t.category}
                    </span>
                  </>
                ) : (
                  <> · {CAT_RU[t.category] ?? t.category}</>
                )}
                {t.comment ? ` · ${t.comment}` : ""}
              </p>
            </div>
            <p
              className={`shrink-0 font-semibold tabular ${
                t.direction === "in" ? "text-emerald-700" : "text-brand"
              }`}
            >
              {t.direction === "in" ? "+" : "−"}
              {formatRub(t.amount)} ₽
            </p>
          </div>

          <div className="mt-2 flex flex-wrap items-center gap-2">
            {t.direction === "out" && cb > 0 ? (
              <span className="inline-flex items-center gap-1 rounded-pill bg-mint px-2 py-0.5 text-[11px] font-semibold text-emerald-800">
                <Sparkles className="h-3 w-3" aria-hidden />
                {rate}% · +{formatRub(cb)} ₽
                {activePick ? " · ваш выбор" : ""}
              </span>
            ) : null}
            {t.direction === "out" && pick && !activePick ? (
              <span className="inline-flex items-center gap-1 rounded-pill bg-canvas px-2 py-0.5 text-[11px] text-ink/45">
                кэшбек в другой категории
              </span>
            ) : null}
            {t.direction === "in" ? (
              <span className="inline-flex items-center gap-1 rounded-pill bg-peach px-2 py-0.5 text-[11px] font-medium text-ink/70">
                ~{formatRub(taxAmt)} ₽ налог
              </span>
            ) : null}
            <Link
              href={chatHref(promptTxnDetail(t))}
              className="btn-touch inline-flex items-center gap-1 rounded-pill bg-brand/10 px-2.5 py-1 text-[11px] font-semibold text-brand"
            >
              <MessageCircle className="h-3 w-3" aria-hidden />
              В чат
            </Link>
            <button
              type="button"
              onClick={onToggle}
              aria-expanded={expanded}
              className="btn-touch inline-flex items-center gap-0.5 rounded-pill px-2 py-1 text-[11px] font-medium text-ink/55 hover:bg-canvas"
            >
              Подробнее
              <ChevronDown
                className={`h-3.5 w-3.5 transition-transform ${expanded ? "rotate-180" : ""}`}
                aria-hidden
              />
            </button>
          </div>
        </div>
      </div>

      <AnimatePresence initial={false}>
        {expanded ? (
          <motion.div
            initial={reduce ? false : { height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={reduce ? undefined : { height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden border-t border-platinum/60"
          >
            <div className="space-y-2 bg-canvas/50 px-4 py-3 text-sm">
              <p className="text-ink/70">{taxHint.detail}</p>
              {t.direction === "out" ? (
                <p className="text-ink/70">
                  {cb > 0 ? (
                    <>
                      Кэшбек по категории «{CAT_RU[t.category] ?? t.category}»:{" "}
                      <span className="font-semibold text-emerald-800">
                        {rate}% → +{formatRub(cb)} ₽
                      </span>
                      {activePick ? " (категория из рулетки)" : ""}
                    </>
                  ) : pick ? (
                    <>
                      Кэшбек месяца выбран на другую категорию — по этой операции
                      начисление 0 ₽. Сменить можно в рулетке на главной.
                    </>
                  ) : (
                    <>
                      Превью ставки {rate}% → ~{formatRub(cashbackForAmount(t.amount, rate))} ₽.
                      Зафиксируйте категорию в рулетке на главной.
                    </>
                  )}
                </p>
              ) : (
                <p className="text-ink/70">
                  Ориентир налога:{" "}
                  <span className="font-semibold tabular">
                    ~{formatRub(taxAmt)} ₽
                  </span>{" "}
                  (4% НПД)
                </p>
              )}
              <div className="flex flex-wrap gap-2 pt-1">
                {t.direction === "in" ? (
                  <Link
                    href={chatHref(promptTxnTax(t))}
                    className="btn-touch rounded-pill bg-brand px-3 py-1.5 text-xs font-semibold text-white"
                  >
                    Налог с операции
                  </Link>
                ) : (
                  <Link
                    href={chatHref(promptCutCategory(t.category, t.amount))}
                    className="btn-touch rounded-pill bg-brand px-3 py-1.5 text-xs font-semibold text-white"
                  >
                    Как снизить такие расходы
                  </Link>
                )}
                <Link
                  href={chatHref(promptTxnDetail(t))}
                  className="btn-touch rounded-pill bg-white px-3 py-1.5 text-xs font-semibold text-ink shadow-soft ring-1 ring-platinum"
                >
                  Спросить в чате
                </Link>
              </div>
            </div>
          </motion.div>
        ) : null}
      </AnimatePresence>
    </li>
  );
}

export default function OperationsPage() {
  const params = useSearchParams();
  const catParam = params.get("cat");
  const [data, setData] = useState<Txn | null>(null);
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<"all" | "in" | "out">("all");
  const [mode, setMode] = useState<"category" | "date">("category");
  const [q, setQ] = useState("");
  const [openCat, setOpenCat] = useState<string | null>(catParam);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [pick, setPick] = useState<CashbackPick | null>(null);
  const reduce = useReducedMotion();

  useEffect(() => {
    setUserId(getUserId());
    setPick(loadCashbackPick());
    getTransactions()
      .then(setData)
      .catch((e: Error) => setErr(e.message))
      .finally(() => setLoading(false));
    if (!getUserId()) {
      getMe()
        .then((p) => {
          if (p.user_id) {
            ensureUserId(p.user_id);
            setUserId(p.user_id);
          }
        })
        .catch(() => undefined);
    }
    const syncPick = () => setPick(loadCashbackPick());
    window.addEventListener("focus", syncPick);
    window.addEventListener("storage", syncPick);
    return () => {
      window.removeEventListener("focus", syncPick);
      window.removeEventListener("storage", syncPick);
    };
  }, []);

  useEffect(() => {
    if (catParam) {
      setMode("category");
      setOpenCat(catParam);
    }
  }, [catParam]);

  const seed = cashbackSeed(userId);

  const items = useMemo(() => {
    let list = data?.items ?? [];
    if (filter !== "all") list = list.filter((t) => t.direction === filter);
    const needle = q.trim().toLowerCase();
    if (needle) {
      list = list.filter(
        (t) =>
          t.counterparty_name.toLowerCase().includes(needle) ||
          (t.comment ?? "").toLowerCase().includes(needle) ||
          (CAT_RU[t.category] ?? t.category).includes(needle),
      );
    }
    return list;
  }, [data, filter, q]);

  const categoryGroups = useMemo(() => {
    const map = new Map<string, Item[]>();
    for (const t of items) {
      const arr = map.get(t.category) ?? [];
      arr.push(t);
      map.set(t.category, arr);
    }
    const groups = [...map.entries()].map(([category, rows]) => {
      const income = rows
        .filter((r) => r.direction === "in")
        .reduce((s, r) => s + r.amount, 0);
      const expense = rows
        .filter((r) => r.direction === "out")
        .reduce((s, r) => s + r.amount, 0);
      const primary =
        filter === "in" ? income : filter === "out" ? expense : expense || income;
      return {
        category,
        rows: rows.sort(
          (a, b) =>
            new Date(b.booked_at).getTime() - new Date(a.booked_at).getTime(),
        ),
        income,
        expense,
        primary,
        cashback: categoryCashbackTotal(rows, category, seed, pick),
      };
    });
    groups.sort((a, b) => b.primary - a.primary);
    return groups;
  }, [items, filter, seed, pick]);

  useEffect(() => {
    if (openCat != null) return;
    if (categoryGroups.length === 0) return;
    setOpenCat(categoryGroups[0].category);
  }, [categoryGroups, openCat]);

  const dateGroups = useMemo(() => {
    const map = new Map<string, { label: string; rows: Item[]; sort: number }>();
    for (const t of items) {
      const key = dayKeyLocal(new Date(t.booked_at));
      const cur = map.get(key);
      if (cur) cur.rows.push(t);
      else
        map.set(key, {
          label: dayLabel(t.booked_at),
          rows: [t],
          sort: new Date(t.booked_at).getTime(),
        });
    }
    return [...map.values()]
      .map((g) => ({
        ...g,
        rows: g.rows.sort(
          (a, b) =>
            new Date(b.booked_at).getTime() - new Date(a.booked_at).getTime(),
        ),
      }))
      .sort((a, b) => b.sort - a.sort);
  }, [items]);

  const topExpense = useMemo(() => {
    const byCat = new Map<string, number>();
    for (const t of data?.items ?? []) {
      if (t.direction !== "out") continue;
      byCat.set(t.category, (byCat.get(t.category) ?? 0) + t.amount);
    }
    let top: { category: string; amount: number } | null = null;
    for (const [category, amount] of byCat) {
      if (!top || amount > top.amount) top = { category, amount };
    }
    return top;
  }, [data]);

  const summary = data?.summary ?? { income: 0, expense: 0, net: 0 };

  const paymentTips = useMemo(() => {
    let rentExpense = 0;
    let servicesIncome = 0;
    let rentInn = "";
    for (const t of data?.items ?? []) {
      if (t.direction === "out" && t.category === "rent") {
        rentExpense += t.amount;
        if (!rentInn && t.counterparty_inn) rentInn = t.counterparty_inn;
      }
      if (t.direction === "in" && t.category === "services") {
        servicesIncome += t.amount;
      }
    }
    return buildPaymentTips({
      income: summary.income,
      expense: summary.expense,
      rentExpense,
      servicesIncome,
      rentInn,
    });
  }, [data, summary.income, summary.expense]);

  const chips = [
    {
      label: "Уменьшить расходы",
      href: chatHref(
        promptCutExpenses({
          expense: summary.expense,
          topCategory: topExpense?.category,
          topAmount: topExpense?.amount,
        }),
      ),
    },
    {
      label: "Увеличить доходы",
      href: chatHref(promptGrowIncome({ income: summary.income })),
    },
    {
      label: "Разбор категории",
      href: chatHref(
        promptCategoryDeepDive({
          category: openCat ?? topExpense?.category ?? "rent",
          amount:
            categoryGroups.find((g) => g.category === (openCat ?? topExpense?.category))
              ?.primary ??
            topExpense?.amount ??
            0,
          sharePct: (() => {
            const g = categoryGroups.find(
              (x) => x.category === (openCat ?? topExpense?.category),
            );
            const base =
              filter === "in"
                ? summary.income
                : filter === "out"
                  ? summary.expense
                  : Math.max(summary.expense, summary.income) || 1;
            if (!g || !base) return 0;
            return (g.primary / base) * 100;
          })(),
          directionHint: filter === "in" ? "приход" : "оборот",
        }),
      ),
    },
    {
      label: "Налог с выписки",
      href: chatHref(
        promptTaxFromStatement({
          income: summary.income,
          expense: summary.expense,
        }),
      ),
    },
  ];

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <PageHeader
        title="Операции"
        subtitle="Выписка счёта ···4582. Здесь видно, откуда берутся цифры налога и отчёта."
        actions={[
          { href: "/app/report", label: "Отчёт", icon: FileText, primary: true },
          { href: "/app/alerts", label: "Алерты", icon: Bell },
        ]}
      />

      {err ? (
        <p className="mt-6 rounded-card bg-rose p-4 text-sm">{err}</p>
      ) : null}

      <section className="mt-5 grid grid-cols-3 gap-3">
        <article className="relative overflow-hidden rounded-card bg-mint p-4 shadow-soft">
          <ArrowDownLeft
            className="pointer-events-none absolute -bottom-2 -right-2 h-16 w-16 text-emerald-700/15"
            aria-hidden
          />
          <p className="relative z-10 text-xs text-ink/50">Приход</p>
          <p className="relative z-10 mt-1 text-lg font-bold tabular text-emerald-800">
            {data ? `${formatRub(data.summary.income)} ₽` : loading ? "…" : "—"}
          </p>
        </article>
        <article className="relative overflow-hidden rounded-card bg-rose p-4 shadow-soft">
          <ArrowUpRight
            className="pointer-events-none absolute -bottom-2 -right-2 h-16 w-16 text-brand/15"
            aria-hidden
          />
          <p className="relative z-10 text-xs text-ink/50">Расход</p>
          <p className="relative z-10 mt-1 text-lg font-bold tabular">
            {data ? `${formatRub(data.summary.expense)} ₽` : loading ? "…" : "—"}
          </p>
        </article>
        <article className="relative overflow-hidden rounded-card bg-white p-4 shadow-soft">
          <p className="text-xs text-ink/50">Нетто</p>
          <p
            className={`mt-1 text-lg font-bold tabular ${
              data && data.summary.net < 0 ? "text-brand" : "text-ink"
            }`}
          >
            {data ? `${formatRub(data.summary.net)} ₽` : loading ? "…" : "—"}
          </p>
        </article>
      </section>

      {!loading && data ? (
        <section className="mt-5 rounded-card bg-white p-4 shadow-soft">
          <div className="flex items-center gap-2">
            <Wallet className="h-4 w-4 text-brand" aria-hidden />
            <h2 className="text-sm font-semibold">Платежи по выписке</h2>
          </div>
          <p className="mt-1 text-xs text-ink/50">
            Следующие шаги по РКО: аренда → приём оплаты → налог ЕНП
          </p>
          <ul className="mt-3 space-y-2">
            {paymentTips.map((tip) => (
              <li key={tip.id}>
                <Link
                  href={tip.href}
                  className="flex items-start gap-3 rounded-2xl bg-canvas px-3 py-3 transition hover:bg-mint/50"
                >
                  <Smartphone
                    className="mt-0.5 h-4 w-4 shrink-0 text-brand"
                    aria-hidden
                  />
                  <span className="min-w-0">
                    <span className="block text-sm font-semibold">{tip.title}</span>
                    <span className="mt-0.5 block text-xs text-ink/55">
                      {tip.text}
                    </span>
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      {pick && getCategoryMeta(pick.id) ? (
        <div className="mt-5 flex items-center justify-between gap-3 rounded-[22px] bg-mint px-4 py-3 text-sm shadow-soft">
          <p className="text-emerald-900">
            <span className="font-semibold">Кэшбек месяца:</span>{" "}
            {getCategoryMeta(pick.id)!.title} · {pick.rate}%
          </p>
          <Link href="/app" className="shrink-0 font-semibold text-brand hover:underline">
            Рулетка →
          </Link>
        </div>
      ) : (
        <div className="mt-5 rounded-[22px] bg-canvas px-4 py-3 text-sm text-ink/60">
          Кэшбек ещё не выбран —{" "}
          <Link href="/app" className="font-semibold text-brand hover:underline">
            откройте рулетку на главной
          </Link>
          . Пока ставки в превью.
        </div>
      )}

      <section className="relative mt-4 overflow-hidden rounded-[28px] bg-white p-4 shadow-soft ring-1 ring-platinum/70">
        <div className="flex items-start gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-brand text-white">
            <MessageCircle className="h-5 w-5" aria-hidden />
          </span>
          <div className="min-w-0 flex-1">
            <p className="font-semibold">Разобрать выписку в чате</p>
            <p className="mt-0.5 text-sm text-ink/55">
              Расходы, доходы, налог и сервисы Альфа-Банка — по вашим цифрам.
            </p>
            <Link
              href={chatHref(promptStatementOverview(summary))}
              className="btn-touch mt-3 inline-flex rounded-pill bg-brand px-4 py-2 text-sm font-semibold text-white"
            >
              Спросить в чате
            </Link>
          </div>
        </div>
        <div className="hide-scrollbar mt-3 flex gap-2 overflow-x-auto pb-0.5">
          {chips.map((c) => (
            <Link
              key={c.label}
              href={c.href}
              className="btn-touch shrink-0 rounded-pill bg-canvas px-3 py-1.5 text-xs font-semibold text-ink/80 ring-1 ring-platinum/80"
            >
              {c.label}
            </Link>
          ))}
        </div>
      </section>

      <label className="relative mt-5 block">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ink/35" />
        <span className="sr-only">Поиск</span>
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Поиск по контрагенту или комментарию"
          className="w-full rounded-2xl border border-platinum bg-white py-3 pl-10 pr-4 text-sm outline-none ring-brand focus:ring-2"
        />
      </label>

      <div className="mt-4 flex flex-wrap items-center gap-2">
        <div className="inline-flex rounded-pill bg-platinum/40 p-1 text-sm">
          {(
            [
              ["all", "Все"],
              ["in", "Приход"],
              ["out", "Расход"],
            ] as const
          ).map(([k, label]) => (
            <button
              key={k}
              type="button"
              onClick={() => setFilter(k)}
              className={`btn-touch rounded-pill px-4 py-1.5 ${
                filter === k ? "bg-ink text-white" : "text-ink/70"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
        <div className="inline-flex rounded-pill bg-white p-1 text-sm shadow-soft ring-1 ring-platinum/80">
          {(
            [
              ["category", "Категории"],
              ["date", "Даты"],
            ] as const
          ).map(([k, label]) => (
            <button
              key={k}
              type="button"
              onClick={() => setMode(k)}
              className={`btn-touch rounded-pill px-3.5 py-1.5 ${
                mode === k ? "bg-brand text-white" : "text-ink/70"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </div>

      {loading ? (
        <ul className="mt-4 space-y-2">
          {[1, 2, 3].map((i) => (
            <li
              key={i}
              className="h-16 animate-pulse rounded-card bg-white/80 shadow-soft"
            />
          ))}
        </ul>
      ) : null}

      {!loading && !err && items.length === 0 ? (
        <div className="mt-4">
          {(data?.items.length ?? 0) === 0 ? (
            <EmptyHint
              title="Пока нет движений"
              text="Выписка пустая — так бывает у нового кабинета. Укажите оборот в профиле и посчитайте налог; операции появятся, когда подключим банк или добавите сценарии в чате."
              primary={{ href: "/app/profile", label: "Профиль" }}
              secondary={{
                href: chatHref(
                  "Сколько мне отложить на налог за этот месяц?",
                ),
                label: "Рассчитать налог",
              }}
            />
          ) : (
            <p className="rounded-card bg-white p-6 text-sm text-ink/55 shadow-soft">
              Ничего не найдено по фильтру.
            </p>
          )}
        </div>
      ) : null}

      {!loading && mode === "category" && categoryGroups.length > 0 ? (
        <div className="mt-4 space-y-3">
          {categoryGroups.map((g) => {
            const Icon = catIcon[g.category] ?? Package;
            const open = openCat === g.category;
            const base =
              filter === "in"
                ? summary.income
                : filter === "out"
                  ? summary.expense
                  : Math.max(summary.expense, summary.income) || 1;
            const share = base ? (g.primary / base) * 100 : 0;
            return (
              <section
                key={g.category}
                id={`cat-${g.category}`}
                className="overflow-hidden rounded-[24px] bg-white shadow-soft"
              >
                <button
                  type="button"
                  aria-expanded={open}
                  onClick={() => setOpenCat(open ? null : g.category)}
                  className="btn-touch flex w-full items-center gap-3 px-4 py-3.5 text-left"
                >
                  <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-canvas text-brand">
                    <Icon className="h-5 w-5" aria-hidden />
                  </span>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between gap-2">
                      <p className="font-semibold capitalize">
                        {CAT_RU[g.category] ?? g.category}
                      </p>
                      <p className="shrink-0 font-bold tabular">
                        {formatRub(g.primary)} ₽
                      </p>
                    </div>
                    <p className="mt-0.5 text-xs text-ink/50">
                      {g.rows.length} оп. · {share.toFixed(0)}%
                      {g.cashback > 0 ? (
                        <>
                          {" · "}
                          <span className="font-medium text-emerald-700">
                            кэшбек +{formatRub(g.cashback)} ₽
                          </span>
                        </>
                      ) : null}
                    </p>
                  </div>
                  <ChevronDown
                    className={`h-5 w-5 shrink-0 text-ink/35 transition-transform ${open ? "rotate-180" : ""}`}
                    aria-hidden
                  />
                </button>
                <AnimatePresence initial={false}>
                  {open ? (
                    <motion.div
                      initial={reduce ? false : { height: 0, opacity: 0 }}
                      animate={{ height: "auto", opacity: 1 }}
                      exit={reduce ? undefined : { height: 0, opacity: 0 }}
                      transition={{ duration: 0.22 }}
                      className="overflow-hidden"
                    >
                      <ul className="space-y-2 border-t border-platinum/50 px-3 pb-3 pt-2">
                        {g.rows.map((t) => (
                          <TxnRow
                            key={t.id}
                            t={t}
                            seed={seed}
                            pick={pick}
                            expanded={expandedId === t.id}
                            onToggle={() =>
                              setExpandedId((id) =>
                                id === t.id ? null : t.id,
                              )
                            }
                          />
                        ))}
                      </ul>
                    </motion.div>
                  ) : null}
                </AnimatePresence>
              </section>
            );
          })}
        </div>
      ) : null}

      {!loading && mode === "date" && dateGroups.length > 0 ? (
        <div className="mt-4 space-y-5">
          {dateGroups.map((g) => (
            <section key={g.label + g.sort}>
              <h2 className="mb-2 text-sm font-semibold text-ink/55">
                {g.label}
              </h2>
              <ul className="space-y-2">
                {g.rows.map((t) => (
                  <TxnRow
                    key={t.id}
                    t={t}
                    seed={seed}
                    pick={pick}
                    showCatBadge
                    expanded={expandedId === t.id}
                    onToggle={() =>
                      setExpandedId((id) => (id === t.id ? null : t.id))
                    }
                  />
                ))}
              </ul>
            </section>
          ))}
        </div>
      ) : null}
    </main>
  );
}
