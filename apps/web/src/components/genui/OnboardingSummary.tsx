"use client";

import { Check } from "lucide-react";
import { formatRub } from "@/lib/api";

export function OnboardingSummary({
  props,
  onAction,
}: {
  props: Record<string, unknown>;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const regimes = (props.suggested_regimes as string[]) ?? [];
  const confirmed = Boolean(props.confirmed);
  return (
    <div
      className="rounded-card bg-lavender/80 p-5 shadow-soft"
      role="region"
      aria-label="Резюме профиля"
    >
      <p className="text-sm text-ink/60">Профиль предпринимателя</p>
      <p className="mt-1 text-2xl font-bold">
        {String(props.display_name ?? "")}
      </p>
      <dl className="mt-4 grid gap-2 text-sm">
        <div className="flex justify-between gap-4">
          <dt className="text-ink/60">Сфера</dt>
          <dd className="font-medium">{String(props.business_sphere ?? "")}</dd>
        </div>
        <div className="flex justify-between gap-4">
          <dt className="text-ink/60">Город</dt>
          <dd className="font-medium">{String(props.city ?? "")}</dd>
        </div>
        <div className="flex justify-between gap-4">
          <dt className="text-ink/60">Оборот</dt>
          <dd className="font-medium tabular">
            {formatRub(Number(props.monthly_revenue_estimate ?? 0))} ₽/мес
          </dd>
        </div>
        <div className="flex justify-between gap-4">
          <dt className="text-ink/60">Режимы</dt>
          <dd className="font-medium">{regimes.join(" / ")}</dd>
        </div>
      </dl>
      {confirmed ? (
        <p className="mt-4 text-sm font-medium text-emerald-700" aria-live="polite">
          Профиль сохранён. Можно считать налог.
        </p>
      ) : (
        <button
          type="button"
          className="btn-touch mt-4 inline-flex items-center gap-2 rounded-pill bg-ink px-5 py-2.5 text-sm font-medium text-white"
          onClick={() =>
            onAction?.(String(props.cta_action ?? "confirm_onboarding"), {
              business_sphere: props.business_sphere,
              city: props.city,
              monthly_revenue_estimate: props.monthly_revenue_estimate,
            })
          }
        >
          <Check className="h-4 w-4" aria-hidden />
          {String(props.cta_label ?? "Всё верно")}
        </button>
      )}
    </div>
  );
}
