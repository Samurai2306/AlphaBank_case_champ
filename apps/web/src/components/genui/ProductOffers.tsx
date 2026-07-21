"use client";

import { ArrowUpRight, CreditCard, Home, Smartphone } from "lucide-react";

const icons: Record<string, typeof Smartphone> = {
  sbp: Smartphone,
  acquiring: CreditCard,
  rent: Home,
};

export function ProductOffers({
  props,
  onAction,
}: {
  props: Record<string, unknown>;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const title = String(props.title ?? "Сервисы оплаты");
  const subtitle = String(props.subtitle ?? "");
  const items = (props.items as Array<Record<string, unknown>>) ?? [];

  return (
    <div className="rounded-card bg-white p-5 shadow-soft">
      <p className="font-semibold">{title}</p>
      {subtitle ? (
        <p className="mt-1 text-sm text-ink/55">{subtitle}</p>
      ) : null}
      <ul className="mt-4 space-y-2">
        {items.map((item) => {
          const id = String(item.id ?? item.title ?? "");
          const Icon = icons[id] ?? ArrowUpRight;
          const chat = String(item.chat ?? "");
          const href = String(item.href ?? "");
          return (
            <li key={id}>
              <button
                type="button"
                className="btn-touch flex w-full items-start gap-3 rounded-2xl bg-canvas px-3 py-3 text-left transition hover:bg-mint/60"
                onClick={() => {
                  if (chat) {
                    onAction?.("send_message", { text: chat });
                    return;
                  }
                  if (href) {
                    onAction?.("navigate", { href });
                  }
                }}
              >
                <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-white text-brand shadow-soft">
                  <Icon className="h-4 w-4" aria-hidden />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="block text-sm font-semibold text-ink">
                    {String(item.title ?? "")}
                  </span>
                  <span className="mt-0.5 block text-xs text-ink/55">
                    {String(item.subtitle ?? "")}
                  </span>
                </span>
                <ArrowUpRight
                  className="mt-1 h-4 w-4 shrink-0 text-ink/30"
                  aria-hidden
                />
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
