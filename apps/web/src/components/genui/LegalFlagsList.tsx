"use client";

type Flag = {
  severity: string;
  title: string;
  quote: string;
  recommendation: string;
};

const sev: Record<string, { label: string; cls: string }> = {
  high: { label: "Важно", cls: "text-brand" },
  medium: { label: "Средне", cls: "text-amber-700" },
  low: { label: "На заметку", cls: "text-ink/55" },
};

export function LegalFlagsList({ props }: { props: Record<string, unknown> }) {
  const flags = (props.flags as Flag[]) ?? [];
  const doc = String(props.document_name ?? "документ");
  return (
    <div className="rounded-card bg-white p-5 shadow-soft">
      <p className="font-semibold">Разбор договора</p>
      <p className="mt-0.5 text-xs text-ink/45">{doc}</p>
      {flags.length === 0 ? (
        <p className="mt-4 rounded-2xl bg-canvas px-3 py-3 text-sm text-ink/60">
          Явных риск-флагов не найдено. Всё равно перечитайте срок, цену и
          штрафы перед подписью.
        </p>
      ) : (
        <ul className="mt-4 space-y-3">
          {flags.map((f) => {
            const s = sev[f.severity.toLowerCase()] ?? sev.medium;
            return (
              <li key={f.title} className="rounded-2xl bg-canvas p-3.5">
                <p className={`text-xs font-semibold uppercase tracking-wide ${s.cls}`}>
                  {s.label}
                </p>
                <p className="mt-1 font-medium">{f.title}</p>
                <p className="mt-1.5 text-sm italic text-ink/55">«{f.quote}»</p>
                <p className="mt-2 text-sm text-ink/80">
                  <span className="font-medium text-ink">Что сделать: </span>
                  {f.recommendation}
                </p>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
