"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Home,
  ArrowLeftRight,
  LayoutGrid,
  MessageCircle,
  UserRound,
} from "lucide-react";

const items = [
  {
    href: "/app",
    label: "Главная",
    icon: Home,
    match: (p: string) => p === "/app",
  },
  {
    href: "/app/operations",
    label: "Операции",
    icon: ArrowLeftRight,
    match: (p: string) => p.startsWith("/app/operations"),
  },
  {
    href: "/app/services",
    label: "Сервисы",
    icon: LayoutGrid,
    match: (p: string) => p.startsWith("/app/services"),
  },
  {
    href: "/app/chat",
    label: "Чат",
    icon: MessageCircle,
    match: (p: string) => p.startsWith("/app/chat"),
  },
  {
    href: "/app/profile",
    label: "Ещё",
    icon: UserRound,
    match: (p: string) =>
      p.startsWith("/app/profile") ||
      p.startsWith("/app/piggy") ||
      p.startsWith("/app/alerts") ||
      p.startsWith("/app/journey") ||
      p.startsWith("/app/payments") ||
      p.startsWith("/app/report") ||
      p.startsWith("/app/tasks"),
  },
];

export function AppNav() {
  const path = usePathname();
  if (path.startsWith("/app/chat")) return null;

  return (
    <nav
      className="fixed inset-x-0 bottom-0 z-40 border-t border-platinum/80 bg-white/95 pb-[env(safe-area-inset-bottom)] backdrop-blur"
      aria-label="Основная навигация"
    >
      <ul className="mx-auto flex max-w-3xl items-stretch justify-between px-1 py-1.5">
        {items.map((item) => {
          const active = item.match(path);
          const Icon = item.icon;
          return (
            <li key={item.href} className="flex-1">
              <Link
                href={item.href}
                className={`flex flex-col items-center gap-0.5 rounded-xl px-1 py-2 text-[11px] font-medium outline-none transition-colors focus-visible:ring-2 focus-visible:ring-brand/40 active:scale-95 ${
                  active ? "text-brand" : "text-ink/45 hover:text-ink/75"
                }`}
              >
                <Icon
                  className={`h-5 w-5 ${active ? "stroke-[2.25]" : "stroke-[1.75]"}`}
                  aria-hidden
                />
                {item.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
