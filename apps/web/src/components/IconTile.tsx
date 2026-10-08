"use client";

import Image from "next/image";
import Link from "next/link";

type Size = "sm" | "md" | "lg";

const sizeCls: Record<
  Size,
  { box: string; title: string; sub: string; textPad: string; artTop: string }
> = {
  sm: {
    box: "h-[168px]",
    title: "text-sm",
    sub: "text-[11px]",
    textPad: "px-3 pt-4",
    artTop: "top-[44%]",
  },
  md: {
    box: "h-[232px]",
    title: "text-[15px] sm:text-base",
    sub: "text-xs sm:text-sm",
    textPad: "px-3 pt-5 sm:px-4",
    artTop: "top-[36%]",
  },
  lg: {
    box: "h-[268px]",
    title: "text-base sm:text-lg",
    sub: "text-sm",
    textPad: "px-4 pt-6",
    artTop: "top-[34%]",
  },
};

/** Product tile: title on top, 3D art in the lower field without hard crop. */
export function IconTile({
  href,
  title,
  sub,
  icon,
  tone = "bg-[#F2F3F5]",
  size = "md",
  className = "",
  blend = true,
}: {
  href: string;
  title: string;
  sub?: string;
  icon: string;
  tone?: string;
  size?: Size;
  className?: string;
  /** Soften residual light fringe on pastel cards */
  blend?: boolean;
}) {
  const s = sizeCls[size];
  return (
    <Link
      href={href}
      className={`group relative block overflow-hidden rounded-[32px] ${tone} ${s.box} shadow-soft transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_16px_40px_rgba(11,17,23,0.12)] active:translate-y-0 active:scale-[0.985] ${className}`}
    >
      <div className={`relative z-20 text-center ${s.textPad}`}>
        <p className={`font-bold leading-snug tracking-tight text-ink ${s.title}`}>
          {title}
        </p>
        {sub ? (
          <p className={`mt-1 leading-snug text-ink/65 ${s.sub}`}>{sub}</p>
        ) : null}
      </div>
      <div
        className={`pointer-events-none absolute inset-x-[-4%] bottom-[-4%] ${s.artTop} transition-transform duration-500 ease-out group-hover:scale-[1.06]`}
      >
        <Image
          src={icon}
          alt=""
          fill
          sizes="(max-width: 640px) 50vw, 240px"
          className={`object-contain object-bottom ${blend ? "mix-blend-multiply" : ""}`}
          priority={false}
        />
      </div>
    </Link>
  );
}
