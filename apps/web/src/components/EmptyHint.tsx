import Link from "next/link";

export function EmptyHint({
  title,
  text,
  primary,
  secondary,
}: {
  title: string;
  text: string;
  primary?: { href: string; label: string };
  secondary?: { href: string; label: string };
}) {
  return (
    <div className="rounded-[28px] bg-white p-5 shadow-soft ring-1 ring-platinum/70">
      <p className="font-semibold text-ink">{title}</p>
      <p className="mt-2 text-sm leading-relaxed text-ink/60">{text}</p>
      <div className="mt-4 flex flex-wrap gap-2">
        {primary ? (
          <Link
            href={primary.href}
            className="btn-touch inline-flex rounded-pill bg-brand px-4 py-2.5 text-sm font-semibold text-white"
          >
            {primary.label}
          </Link>
        ) : null}
        {secondary ? (
          <Link
            href={secondary.href}
            className="btn-touch inline-flex rounded-pill bg-canvas px-4 py-2.5 text-sm font-semibold text-ink"
          >
            {secondary.label}
          </Link>
        ) : null}
      </div>
    </div>
  );
}
