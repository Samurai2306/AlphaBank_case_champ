"use client";

import { Check, FileText, X } from "lucide-react";
import { formatRub } from "@/lib/api";

export function PaymentDraftCard({
  props,
  onAction,
}: {
  props: Record<string, unknown>;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const id = String(props.draft_id ?? "");
  const status = String(props.status ?? "draft");
  const risk = String(props.risk_level ?? "green");
  const payeeInn = String(props.payee_inn ?? "");
  const confirmed = status === "confirmed" || status === "confirmed_mock";
  const riskBadge =
    risk === "red"
      ? "bg-rose text-brand"
      : risk === "yellow"
        ? "bg-peach text-amber-900"
        : "bg-mint text-emerald-800";
  return (
    <div className="relative mt-1 overflow-hidden rounded-card border border-platinum/80 bg-white p-5 shadow-soft">
      <div className="pointer-events-none absolute -bottom-6 -right-4 h-28 w-28 text-brand/10">
        <FileText className="h-full w-full" strokeWidth={1} aria-hidden />
      </div>
      <div className="relative z-10 flex items-center justify-between gap-3">
        <p className="font-semibold">Платёжное поручение</p>
        <span
          className={`rounded-pill px-3 py-1 text-xs font-semibold ${
            confirmed
              ? "bg-mint text-emerald-800"
              : status === "cancelled"
                ? "bg-canvas text-ink/55"
                : "bg-amber-100 text-amber-900"
          }`}
        >
          {confirmed ? "принято" : status === "cancelled" ? "отменено" : "черновик"}
        </span>
      </div>
      <p className="relative z-10 mt-3 text-3xl font-bold tabular">
        {formatRub(Number(props.amount ?? 0))} ₽
      </p>
      <p className="relative z-10 mt-2 text-sm text-ink/70">{String(props.purpose ?? "")}</p>
      <p className="relative z-10 mt-1 text-sm text-ink/50">{String(props.payee_name ?? "")}</p>
      {payeeInn ? (
        <p className="relative z-10 mt-0.5 text-xs text-ink/45">ИНН {payeeInn}</p>
      ) : null}
      {risk && risk !== "green" ? (
        <p
          className={`relative z-10 mt-2 inline-flex rounded-pill px-2.5 py-1 text-[11px] font-semibold ${riskBadge}`}
        >
          риск: {risk === "red" ? "красный" : "жёлтый"}
        </p>
      ) : null}
      {status === "draft" ? (
        <div className="relative z-10 mt-4 flex flex-wrap gap-2">
          <button
            type="button"
            className="btn-touch inline-flex items-center gap-2 rounded-pill bg-brand px-5 py-2.5 text-sm font-medium text-white"
            onClick={() => onAction?.("confirm_draft", { draft_id: id })}
          >
            <Check className="h-4 w-4" aria-hidden />
            {String(props.cta_confirm_label ?? "Подтвердить")}
          </button>
          <button
            type="button"
            className="btn-touch inline-flex items-center gap-2 rounded-pill bg-canvas px-5 py-2.5 text-sm font-medium text-ink"
            onClick={() => onAction?.("cancel_draft", { draft_id: id })}
          >
            <X className="h-4 w-4" aria-hidden />
            {String(props.cta_cancel_label ?? "Отменить")}
          </button>
        </div>
      ) : status === "cancelled" ? (
        <p className="relative z-10 mt-4 text-sm text-ink/55" aria-live="polite">
          Черновик отменён. Можно сформировать новый.
        </p>
      ) : (
        <p className="relative z-10 mt-4 text-sm text-emerald-700" aria-live="polite">
          Поручение принято к исполнению. Статус — в разделе «Платежи».
        </p>
      )}
    </div>
  );
}
