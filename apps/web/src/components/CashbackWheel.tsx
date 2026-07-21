"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "framer-motion";
import {
  Fuel,
  Home,
  Megaphone,
  Monitor,
  Package,
  Percent,
  Scissors,
  ShoppingBag,
  Sparkles,
  X,
  type LucideIcon,
} from "lucide-react";
import {
  cashbackSeed,
  getCashbackRate,
  getCategoryMeta,
  loadCashbackPick,
  monthCashbackTotal,
  saveCashbackPick,
  spinCandidateCategories,
  type CashbackCategory,
  type CashbackPick,
  type TxnLike,
} from "@/lib/cashback";
import { formatRub } from "@/lib/api";

const icons: Record<string, LucideIcon> = {
  services: Scissors,
  rent: Home,
  supplies: Package,
  ads: Megaphone,
  fees: Percent,
  ops: Package,
  fuel: Fuel,
  market: ShoppingBag,
  software: Monitor,
};

function shufflePick3(cats: CashbackCategory[], salt: string): CashbackCategory[] {
  const pool = cats.length >= 3 ? cats : spinCandidateCategories([]);
  const scored = pool.map((c, i) => ({
    c,
    s: (hash(`${salt}:spin:${c.id}`) + i * 17) % 997,
  }));
  scored.sort((a, b) => a.s - b.s);
  const out: CashbackCategory[] = [];
  for (const row of scored) {
    if (out.length >= 3) break;
    if (!out.some((x) => x.id === row.c.id)) out.push(row.c);
  }
  return out;
}

function hash(s: string): number {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return h >>> 0;
}

export function CashbackWheel({
  userId,
  transactions,
  emptyCabinet,
}: {
  userId?: string | null;
  transactions: TxnLike[];
  emptyCabinet?: boolean;
}) {
  const reduce = useReducedMotion();
  const seed = cashbackSeed(userId);
  const candidates = useMemo(
    () => spinCandidateCategories(transactions),
    [transactions],
  );

  const [open, setOpen] = useState(false);
  const [phase, setPhase] = useState<"idle" | "spin" | "pick">("idle");
  const [spinKey, setSpinKey] = useState(0);
  const [options, setOptions] = useState<CashbackCategory[]>([]);
  const [picked, setPicked] = useState<CashbackPick | null>(null);

  useEffect(() => {
    setPicked(loadCashbackPick());
  }, []);

  const earned = monthCashbackTotal(transactions, seed, picked);

  const pickedMeta = useMemo(
    () => (picked ? getCategoryMeta(picked.id) : null),
    [picked],
  );

  function startSpin() {
    if (
      picked &&
      !window.confirm(
        "Новая рулетка заменит текущий кэшбек месяца. Продолжить?",
      )
    ) {
      return;
    }
    setOpen(true);
    setPhase("spin");
    setSpinKey((k) => k + 1);
    const next = shufflePick3(candidates, `${seed}:${Date.now()}`);
    setOptions(next);
    window.setTimeout(() => setPhase("pick"), reduce ? 200 : 2200);
  }

  function choose(c: CashbackCategory) {
    const rate = getCashbackRate(c.id, seed);
    const p = saveCashbackPick({
      id: c.id,
      rate,
      at: new Date().toISOString(),
    });
    setPicked(p);
    setPhase("idle");
    setOpen(false);
  }

  return (
    <>
      <motion.section
        initial={reduce ? false : { opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        className="mt-4 overflow-hidden rounded-[28px] bg-white p-4 shadow-soft ring-1 ring-platinum/70"
      >
        <div className="flex items-start justify-between gap-3">
          <div>
            <h2 className="text-lg font-bold">Кэшбек на бизнес</h2>
            <p className="mt-0.5 text-xs text-ink/50">
              {emptyCabinet || transactions.length === 0
                ? "Крутите рулетку — выберите одну категорию на месяц"
                : picked
                  ? `По выбранной категории ~${formatRub(earned)} ₽ за месяц`
                  : "Пока превью по всем категориям · выберите одну в рулетке"}
            </p>
          </div>
          <Link
            href={
              picked
                ? `/app/operations?cat=${picked.id}`
                : "/app/operations"
            }
            className="shrink-0 text-sm font-medium text-brand hover:underline"
          >
            Выписка →
          </Link>
        </div>

        {picked && pickedMeta ? (
          <div
            className={`mt-3 flex items-center gap-3 rounded-[22px] ${pickedMeta.tone} px-4 py-3`}
          >
            <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-white/80 text-brand">
              {(() => {
                const Icon = icons[pickedMeta.id] ?? Sparkles;
                return <Icon className="h-5 w-5" aria-hidden />;
              })()}
            </span>
            <div className="min-w-0 flex-1">
              <p className="text-sm font-semibold">{pickedMeta.title}</p>
              <p className="text-xs text-ink/55">
                Активный кэшбек месяца · {pickedMeta.partnerHint}
              </p>
            </div>
            <p className="text-2xl font-bold tabular">{picked.rate}%</p>
          </div>
        ) : (
          <p className="mt-3 text-sm text-ink/55">
            Откройте рулетку, выпадут 3 категории — заберите одну. В выписке
            кэшбек будет только по ней.
          </p>
        )}

        <button
          type="button"
          onClick={startSpin}
          className="btn-touch mt-3 inline-flex w-full items-center justify-center gap-2 rounded-pill bg-brand py-3 text-sm font-semibold text-white"
        >
          <Sparkles className="h-4 w-4" aria-hidden />
          {picked ? "Крутить ещё раз" : "Открыть рулетку кэшбека"}
        </button>
      </motion.section>

      <AnimatePresence>
        {open ? (
          <motion.div
            className="fixed inset-0 z-50 flex items-end justify-center bg-ink/45 p-4 sm:items-center"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            role="dialog"
            aria-modal
            aria-label="Рулетка кэшбека"
            onClick={(e) => {
              if (e.target === e.currentTarget && phase !== "spin") {
                setOpen(false);
                setPhase("idle");
              }
            }}
          >
            <motion.div
              initial={reduce ? false : { y: 40, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              exit={reduce ? undefined : { y: 24, opacity: 0 }}
              className="w-full max-w-md overflow-hidden rounded-[28px] bg-white p-5 shadow-[0_24px_80px_rgba(0,0,0,0.25)]"
            >
              <div className="flex items-center justify-between gap-2">
                <div>
                  <p className="font-bold">Рулетка кэшбека</p>
                  <p className="text-xs text-ink/50">
                    Три варианта — выберите один на этот месяц
                  </p>
                </div>
                <button
                  type="button"
                  disabled={phase === "spin"}
                  onClick={() => {
                    setOpen(false);
                    setPhase("idle");
                  }}
                  className="btn-touch rounded-full p-2 text-ink/45 hover:bg-canvas disabled:opacity-40"
                  aria-label="Закрыть"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              {phase === "spin" ? (
                <div className="relative mt-6 flex h-44 items-center justify-center">
                  <motion.div
                    key={spinKey}
                    className="flex h-36 w-36 items-center justify-center rounded-full bg-gradient-to-br from-brand via-[#ff6b5a] to-peach text-white shadow-[0_0_40px_rgba(239,49,36,0.35)]"
                    animate={reduce ? {} : { rotate: 360 * 4 + 40 }}
                    transition={{ duration: 2.1, ease: [0.15, 0.85, 0.25, 1] }}
                  >
                    <Sparkles className="h-12 w-12" aria-hidden />
                  </motion.div>
                  <p className="absolute bottom-0 text-sm font-medium text-ink/55">
                    Крутим…
                  </p>
                </div>
              ) : null}

              {phase === "pick" ? (
                <div className="mt-5 space-y-2.5">
                  <p className="text-sm text-ink/60">
                    Выпало три варианта — нажмите на самый интересный:
                  </p>
                  {options.map((c, i) => {
                    const rate = getCashbackRate(c.id, seed);
                    const Icon = icons[c.id] ?? Sparkles;
                    return (
                      <motion.button
                        key={c.id}
                        type="button"
                        initial={reduce ? false : { opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: i * 0.06 }}
                        onClick={() => choose(c)}
                        className={`card-touch flex w-full items-center gap-3 rounded-[22px] ${c.tone} px-4 py-3.5 text-left`}
                      >
                        <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-white/80 text-brand">
                          <Icon className="h-5 w-5" aria-hidden />
                        </span>
                        <div className="min-w-0 flex-1">
                          <p className="font-semibold">{c.title}</p>
                          <p className="text-xs text-ink/55">{c.partnerHint}</p>
                        </div>
                        <span className="text-2xl font-bold tabular">{rate}%</span>
                      </motion.button>
                    );
                  })}
                </div>
              ) : null}
            </motion.div>
          </motion.div>
        ) : null}
      </AnimatePresence>
    </>
  );
}
