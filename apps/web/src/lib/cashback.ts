/** Demo cashback rates for business categories — client-only, seeded per user/month. */

export type CashbackCategoryId =
  | "services"
  | "rent"
  | "supplies"
  | "ads"
  | "fees"
  | "ops"
  | "fuel"
  | "market"
  | "software";

export type CashbackCategory = {
  id: CashbackCategoryId;
  title: string;
  /** Tailwind bg class */
  tone: string;
  /** Inclusive percent range for monthly seed */
  minPct: number;
  maxPct: number;
  partnerHint: string;
  /** Show in home strip even if no txns */
  strip: boolean;
};

export const CASHBACK_CATALOG: CashbackCategory[] = [
  {
    id: "services",
    title: "Услуги",
    tone: "bg-mint",
    minPct: 2,
    maxPct: 5,
    partnerHint: "эквайринг · партнёры",
    strip: true,
  },
  {
    id: "rent",
    title: "Аренда",
    tone: "bg-lavender",
    minPct: 1,
    maxPct: 3,
    partnerHint: "бизнес-карта",
    strip: true,
  },
  {
    id: "supplies",
    title: "Расходники",
    tone: "bg-peach",
    minPct: 3,
    maxPct: 7,
    partnerHint: "маркет · до 7%",
    strip: true,
  },
  {
    id: "ads",
    title: "Реклама",
    tone: "bg-rose",
    minPct: 5,
    maxPct: 10,
    partnerHint: "партнёры рекламы",
    strip: true,
  },
  {
    id: "fees",
    title: "Комиссии",
    tone: "bg-cyanSoft",
    minPct: 1,
    maxPct: 2,
    partnerHint: "тариф РКО",
    strip: true,
  },
  {
    id: "ops",
    title: "Операционные",
    tone: "bg-[#F2F3F5]",
    minPct: 2,
    maxPct: 4,
    partnerHint: "ежедневные траты",
    strip: true,
  },
  // Partner niches kept for rates, but not offered in the wheel (no seed txns).
  {
    id: "fuel",
    title: "Топливо",
    tone: "bg-peach",
    minPct: 4,
    maxPct: 8,
    partnerHint: "АЗС-партнёры",
    strip: false,
  },
  {
    id: "market",
    title: "Маркет",
    tone: "bg-mint",
    minPct: 3,
    maxPct: 6,
    partnerHint: "закупки онлайн",
    strip: false,
  },
  {
    id: "software",
    title: "Сервисы",
    tone: "bg-lavender",
    minPct: 5,
    maxPct: 12,
    partnerHint: "SaaS · подписки",
    strip: false,
  },
];

export const CAT_RU: Record<string, string> = {
  services: "услуги",
  rent: "аренда",
  supplies: "расходники",
  ads: "реклама",
  fees: "комиссии",
  ops: "операционные",
  fuel: "топливо",
  market: "маркет",
  software: "сервисы",
};

function hashSeed(s: string): number {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return h >>> 0;
}

/** Stable seed for the current calendar month + user. */
export function cashbackSeed(userId?: string | null, when = new Date()): string {
  const ym = `${when.getFullYear()}-${String(when.getMonth() + 1).padStart(2, "0")}`;
  return `${ym}:${userId?.trim() || "demo"}`;
}

export function getCategoryMeta(id: string): CashbackCategory | undefined {
  return CASHBACK_CATALOG.find((c) => c.id === id);
}

/** Percent rate for category this month (one decimal). */
export function getCashbackRate(category: string, seed: string): number {
  const meta = getCategoryMeta(category);
  if (!meta) {
    const h = hashSeed(`${seed}:fallback:${category}`);
    return 1 + (h % 40) / 10;
  }
  const h = hashSeed(`${seed}:${meta.id}`);
  const span = meta.maxPct - meta.minPct;
  const steps = Math.max(1, Math.round(span * 10));
  const pick = (h % (steps + 1)) / 10;
  return Math.round((meta.minPct + pick) * 10) / 10;
}

export function cashbackForAmount(amount: number, ratePct: number): number {
  if (amount <= 0 || ratePct <= 0) return 0;
  return Math.round(amount * (ratePct / 100) * 100) / 100;
}

export type TxnLike = {
  amount: number;
  direction: string;
  category: string;
};

export function taxOnIncome(amount: number): number {
  return Math.round(amount * 0.04 * 100) / 100;
}

export function taxHintForTxn(txn: TxnLike): {
  kind: "income_tax" | "expense_note";
  amount?: number;
  label: string;
  detail: string;
} {
  if (txn.direction === "in") {
    const tax = taxOnIncome(txn.amount);
    return {
      kind: "income_tax",
      amount: tax,
      label: `~${tax.toLocaleString("ru-RU")} ₽ налог`,
      detail:
        "Ориентир НПД ~4% с физлиц (оценка для кабинета, не квитанция ФНС). Удобно откладывать в копилку ЕНП.",
    };
  }
  return {
    kind: "expense_note",
    label: "Не уменьшает НПД",
    detail:
      "На режиме НПД расходы не вычитаются из налога. Имеет смысл смотреть кэшбек и оптимизировать категорию.",
  };
}

const PICK_KEY = "alfa_start_cashback_pick";

export type CashbackPick = {
  id: string;
  rate: number;
  /** ISO timestamp */
  at: string;
  /** YYYY-MM — invalidates after month change */
  month: string;
};

function currentMonthKey(when = new Date()): string {
  return `${when.getFullYear()}-${String(when.getMonth() + 1).padStart(2, "0")}`;
}

export function loadCashbackPick(when = new Date()): CashbackPick | null {
  if (typeof window === "undefined") return null;
  try {
    const raw = localStorage.getItem(PICK_KEY);
    if (!raw) return null;
    const p = JSON.parse(raw) as CashbackPick;
    if (!p?.id || typeof p.rate !== "number") return null;
    const month = p.month || (p.at ? currentMonthKey(new Date(p.at)) : "");
    if (month !== currentMonthKey(when)) return null;
    return { ...p, month };
  } catch {
    return null;
  }
}

export function saveCashbackPick(p: Omit<CashbackPick, "month"> & { month?: string }) {
  const month = p.month ?? currentMonthKey();
  const full: CashbackPick = { id: p.id, rate: p.rate, at: p.at, month };
  localStorage.setItem(PICK_KEY, JSON.stringify(full));
  return full;
}

/**
 * Rate applied to an expense category.
 * After roulette pick — only the chosen category earns; others 0%.
 * Before pick — preview seeded rates for all (so выписка не пустая).
 */
export function rateForExpense(
  category: string,
  seed: string,
  pick: CashbackPick | null,
): number {
  if (pick) {
    return pick.id === category ? pick.rate : 0;
  }
  return getCashbackRate(category, seed);
}

export function monthCashbackTotal(
  items: TxnLike[],
  seed: string,
  pick: CashbackPick | null = null,
): number {
  let sum = 0;
  for (const t of items) {
    if (t.direction !== "out") continue;
    const rate = rateForExpense(t.category, seed, pick);
    sum += cashbackForAmount(t.amount, rate);
  }
  return Math.round(sum * 100) / 100;
}

export function categoryCashbackTotal(
  items: TxnLike[],
  category: string,
  seed: string,
  pick: CashbackPick | null = null,
): number {
  const rate = rateForExpense(category, seed, pick);
  if (rate <= 0) return 0;
  let sum = 0;
  for (const t of items) {
    if (t.direction !== "out" || t.category !== category) continue;
    sum += cashbackForAmount(t.amount, rate);
  }
  return Math.round(sum * 100) / 100;
}

/** Prefer categories that appear in expenses when spinning the wheel. */
export function spinCandidateCategories(items: TxnLike[]): CashbackCategory[] {
  const expenseCats = new Set(
    items.filter((t) => t.direction === "out").map((t) => t.category),
  );
  const withSpend = CASHBACK_CATALOG.filter((c) => expenseCats.has(c.id));
  const partners = CASHBACK_CATALOG.filter(
    (c) => c.strip && !expenseCats.has(c.id),
  );
  if (withSpend.length >= 3) return withSpend;
  return [...withSpend, ...partners];
}
