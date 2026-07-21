"use client";

import { Calculator, FilePlus2 } from "lucide-react";
import { formatRub } from "@/lib/api";

export function TaxCard({
  props,
  onAction,
}: {
  props: Record<string, unknown>;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const amount = Number(props.tax_amount ?? 0);
  const ceiling = Number(props.ceiling_amount ?? 0);
  const cta = props.cta as { action?: string; label?: string } | undefined;
  const assumptions = (props.assumptions as string[]) ?? [];
  return (
    <div
      className="relative mt-1 overflow-hidden rounded-card bg-white p-5 shadow-soft"
      role="region"
      aria-label="Карточка налога"
    >
      <Calculator
        className="pointer-events-none absolute -bottom-4 -right-4 h-28 w-28 text-brand/10"
        strokeWidth={1}
        aria-hidden
      />
      <p className="text-sm text-ink/60">{String(props.period_label ?? "Период")}</p>
      <p className="mt-1 text-sm font-medium">
        Режим {String(props.regime)} · {String(props.rate_label)}
      </p>
      <p className="mt-3 text-4xl font-bold tabular tracking-tight">
        {formatRub(amount)} ₽
      </p>
      {ceiling > 0 && ceiling !== amount ? (
        <p className="mt-1 text-sm text-ink/55">
          С запасом (6%): {formatRub(ceiling)} ₽
        </p>
      ) : null}
      <p className="mt-2 text-sm text-ink/60">
        С дохода {formatRub(Number(props.income ?? 0))} ₽
        {props.b2c_share != null
          ? ` · физлица ~${Math.round(Number(props.b2c_share) * 100)}%`
          : ""}
      </p>
      {assumptions.length ? (
        <ul className="mt-3 list-disc space-y-1 pl-4 text-xs text-ink/50">
          {assumptions.slice(0, 3).map((a) => (
            <li key={a}>{a}</li>
          ))}
        </ul>
      ) : null}
      {props.source_ref ? (
        <p className="mt-2 text-[11px] text-ink/40">
          Источник: {String(props.source_ref)}
        </p>
      ) : null}
      {cta?.label ? (
        <button
          type="button"
          onClick={() => onAction?.(cta.action ?? "create_payment_draft")}
          className="btn-touch mt-4 inline-flex items-center gap-2 rounded-pill bg-brand px-5 py-2.5 text-sm font-medium text-white"
        >
          <FilePlus2 className="h-4 w-4" aria-hidden />
          {cta.label}
        </button>
      ) : null}
    </div>
  );
}
