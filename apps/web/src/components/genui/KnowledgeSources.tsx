"use client";

type Item = {
  id?: string;
  title?: string;
  source?: string;
  snippet?: string;
  score?: number;
};

export function KnowledgeSources({
  props,
}: {
  props: Record<string, unknown>;
}) {
  const title = String(props.title ?? "Источники из базы знаний");
  const items = (Array.isArray(props.items) ? props.items : []) as Item[];

  if (items.length === 0) {
    return null;
  }

  return (
    <div className="rounded-card border border-ink/10 bg-white/90 p-4 shadow-soft">
      <p className="text-xs font-medium uppercase tracking-wide text-ink/50">
        {title}
      </p>
      <ul className="mt-3 space-y-3">
        {items.map((item, i) => (
          <li key={item.id ?? i} className="border-t border-ink/8 pt-3 first:border-0 first:pt-0">
            <div className="flex items-baseline justify-between gap-2">
              <p className="text-sm font-semibold text-ink">
                {item.title ?? "Фрагмент"}
              </p>
              {item.id ? (
                <code className="shrink-0 text-[10px] text-ink/40">{item.id}</code>
              ) : null}
            </div>
            {item.snippet ? (
              <p className="mt-1 text-sm leading-snug text-ink/70">{item.snippet}</p>
            ) : null}
          </li>
        ))}
      </ul>
    </div>
  );
}
