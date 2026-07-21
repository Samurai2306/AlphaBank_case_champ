"use client";

const levelMeta: Record<
  string,
  { dot: string; badge: string; label: string }
> = {
  green: {
    dot: "bg-emerald-500",
    badge: "bg-mint text-emerald-800",
    label: "Зелёный",
  },
  yellow: {
    dot: "bg-amber-400",
    badge: "bg-peach text-amber-900",
    label: "Жёлтый",
  },
  red: {
    dot: "bg-brand",
    badge: "bg-rose text-brand",
    label: "Красный",
  },
};

export function ComplianceTrafficLight({
  props,
  onAction,
}: {
  props: Record<string, unknown>;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const level = String(props.level ?? "green");
  const meta = levelMeta[level] ?? levelMeta.green;
  const reasons = (props.reasons as string[]) ?? [];
  const cta = props.cta as
    | { action?: string; label?: string; text?: string }
    | undefined;
  return (
    <div className="rounded-card bg-white p-5 shadow-soft">
      <div className="flex items-start gap-3">
        <span
          className={`mt-1 h-4 w-4 shrink-0 rounded-full ${meta.dot}`}
          aria-hidden
        />
        <div className="min-w-0">
          <p className="font-semibold">{String(props.title ?? "Проверка")}</p>
          <p className="mt-1 text-sm text-ink/60">
            ИНН {String(props.inn ?? "")}
            <span
              className={`ml-2 inline-flex rounded-pill px-2 py-0.5 text-[11px] font-semibold ${meta.badge}`}
            >
              {meta.label}
            </span>
          </p>
        </div>
      </div>
      {reasons.length > 0 ? (
        <ul className="mt-3 space-y-1.5 text-sm text-ink/70">
          {reasons.map((r) => (
            <li key={r} className="flex gap-2">
              <span className="text-ink/35" aria-hidden>
                ·
              </span>
              <span>{r}</span>
            </li>
          ))}
        </ul>
      ) : null}
      <p className="mt-3 rounded-2xl bg-canvas px-3 py-2.5 text-sm text-ink/80">
        {String(props.recommendation ?? "")}
      </p>
      {cta?.label && cta.text ? (
        <button
          type="button"
          className="btn-touch mt-3 inline-flex rounded-pill bg-brand px-5 py-2.5 text-sm font-medium text-white"
          onClick={() =>
            onAction?.("send_message", { text: String(cta.text) })
          }
        >
          {String(cta.label)}
        </button>
      ) : null}
    </div>
  );
}
