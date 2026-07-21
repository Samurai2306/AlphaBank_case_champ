import { CAT_RU, type TxnLike } from "@/lib/cashback";
import { formatRub } from "@/lib/api";

export type PromptTxn = TxnLike & {
  counterparty_name: string;
  booked_at: string;
  comment?: string;
};

const ALFA_TAIL =
  " Какие сервисы Альфа-Банка и партнёров подойдут (РКО, бизнес-карта, копилка ЕНП, эквайринг, партнёрский кэшбек)?";

function catLabel(category: string): string {
  return CAT_RU[category] ?? category;
}

export function chatHref(q: string): string {
  return "/app/chat?q=" + encodeURIComponent(q);
}

export function promptTxnDetail(t: PromptTxn): string {
  const dir = t.direction === "in" ? "приход" : "расход";
  const date = new Date(t.booked_at).toLocaleDateString("ru-RU");
  const comment = t.comment?.trim() ? t.comment.trim() : "—";
  return (
    `Разбери операцию: ${dir} ${formatRub(t.amount)} ₽, контрагент «${t.counterparty_name}», ` +
    `категория «${catLabel(t.category)}», дата ${date}, комментарий «${comment}». ` +
    `Что это значит для бизнеса? Как оптимизировать?` +
    ALFA_TAIL
  );
}

export function promptTxnTax(t: PromptTxn): string {
  return (
    `Сколько налога НПД уйдёт с этой операции: приход ${formatRub(t.amount)} ₽ ` +
    `от «${t.counterparty_name}» (${catLabel(t.category)})? Свяжи с копилкой ЕНП.` +
    ALFA_TAIL
  );
}

export function promptCutExpenses(summary: {
  expense: number;
  topCategory?: string;
  topAmount?: number;
}): string {
  const top =
    summary.topCategory && summary.topAmount != null
      ? ` Топ расхода: «${catLabel(summary.topCategory)}» — ${formatRub(summary.topAmount)} ₽.`
      : "";
  return (
    `Как уменьшить расходы по выписке? Расход за месяц ${formatRub(summary.expense)} ₽.${top} ` +
    `Что срезать в первую очередь и где взять кэшбек?` +
    ALFA_TAIL
  );
}

export function promptGrowIncome(summary: { income: number }): string {
  return (
    `Как увеличить доходы? Приход за месяц ${formatRub(summary.income)} ₽. ` +
    `Дай идеи роста для малого бизнеса и партнёрские офферы.` +
    ALFA_TAIL
  );
}

export function promptCategoryDeepDive(input: {
  category: string;
  amount: number;
  sharePct: number;
  directionHint: string;
}): string {
  return (
    `Разбери категорию «${catLabel(input.category)}»: ${input.directionHint} ` +
    `${formatRub(input.amount)} ₽ (${input.sharePct.toFixed(0)}% от группы). ` +
    `Что делать дальше?` +
    ALFA_TAIL
  );
}

export function promptTaxFromStatement(summary: {
  income: number;
  expense: number;
}): string {
  return (
    `Свяжи выписку с налогом НПД: приход ${formatRub(summary.income)} ₽, ` +
    `расход ${formatRub(summary.expense)} ₽. Сколько отложить в копилку ЕНП?` +
    ALFA_TAIL
  );
}

export function promptStatementOverview(summary: {
  income: number;
  expense: number;
  net: number;
}): string {
  return (
    `Разбери выписку за месяц: приход ${formatRub(summary.income)} ₽, ` +
    `расход ${formatRub(summary.expense)} ₽, нетто ${formatRub(summary.net)} ₽. ` +
    `Дай выводы и следующий шаг.` +
    ALFA_TAIL
  );
}

export function promptCutCategory(category: string, amount: number): string {
  return (
    `Как снизить расходы в категории «${catLabel(category)}» ` +
    `(сейчас ${formatRub(amount)} ₽ за месяц)?` +
    ALFA_TAIL
  );
}

export type PaymentTip = {
  id: string;
  title: string;
  text: string;
  href: string;
};

/** RKO / payment next steps derived from statement categories. */
export function buildPaymentTips(input: {
  income: number;
  expense: number;
  rentExpense: number;
  servicesIncome: number;
  rentInn?: string;
}): PaymentTip[] {
  const tips: PaymentTip[] = [];
  const inn = input.rentInn?.trim() || "1650987654";
  if (input.rentExpense > 0) {
    tips.push({
      id: "rent",
      title: "Аренда",
      text: `В выписке аренда ≈ ${formatRub(input.rentExpense)} ₽ — проверьте ИНН и соберите платёжку.`,
      href: chatHref(`Сформируй платёжку аренды ${Math.round(input.rentExpense)} на ИНН ${inn}`),
    });
  } else {
    tips.push({
      id: "rent-check",
      title: "Перед арендой",
      text: "Перед крупным платежом арендодателю — светофор 115-ФЗ и черновик поручения.",
      href: chatHref(`Проверь ИНН ${inn}`),
    });
  }
  if (input.servicesIncome >= 20000) {
    tips.push({
      id: "sbp",
      title: "Приём оплаты",
      text: `Приход услуг ≈ ${formatRub(input.servicesIncome)} ₽ — СБП/эквайринг ускоряют деньги на РКО.`,
      href: chatHref("Как принимать оплату по СБП и эквайрингу?"),
    });
  } else {
    tips.push({
      id: "sbp",
      title: "СБП / эквайринг",
      text: "Чтобы быстрее выходить на безубыточность, принимайте оплату на счёт, а не только наличными.",
      href: chatHref("Как принимать оплату по СБП и эквайрингу?"),
    });
  }
  if (input.income > 0) {
    tips.push({
      id: "enp",
      title: "Налог ЕНП",
      text: `С прихода ${formatRub(input.income)} ₽ отложите налог в копилку и сформируйте ЕНП.`,
      href: chatHref("Сколько мне отложить на налог за этот месяц?"),
    });
  }
  return tips.slice(0, 3);
}
