"use client";

import { MessageCircle } from "lucide-react";
import { IconTile } from "@/components/IconTile";
import { PageHeader } from "@/components/PageHeader";

const catalog = [
  {
    title: "Налог НПД",
    sub: "Смешанная ставка 4%/6%",
    icon: "/icons/3d/02-tax-calculator.png",
    tone: "bg-mint",
    href: "/app/chat?q=" + encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
  },
  {
    title: "Сканер договора",
    sub: "PDF до 10 МБ в чате",
    icon: "/icons/3d/03-legal-scanner.png",
    tone: "bg-rose",
    href: "/app/chat?q=" + encodeURIComponent("Разбери договор аренды PDF"),
  },
  {
    title: "Точка безубыточности",
    sub: "Клиенты в день",
    icon: "/icons/3d/05-unit-economics.png",
    tone: "bg-peach",
    href:
      "/app/chat?q=" +
      encodeURIComponent(
        "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
      ),
  },
  {
    title: "Проверка контрагента",
    sub: "ИНН и 115-ФЗ",
    icon: "/icons/3d/06-compliance-risk.png",
    tone: "bg-lavender",
    href: "/app/chat?q=" + encodeURIComponent("Проверь ИНН 1650987654"),
  },
  {
    title: "Копилка ЕНП",
    sub: "Автоотчисления",
    icon: "/icons/3d/09-enp-piggy-bank.png",
    tone: "bg-cyanSoft",
    href: "/app/piggy",
  },
  {
    title: "Платёжка",
    sub: "Поручение в банк",
    icon: "/icons/3d/10-payment-draft.png",
    tone: "bg-mint",
    href: "/app/chat?q=" + encodeURIComponent("Сформируй платёжку."),
  },
  {
    title: "Операции",
    sub: "Выписка счёта",
    icon: "/icons/3d/04-transaction-manager.png",
    tone: "bg-[#F2F3F5]",
    href: "/app/operations",
  },
  {
    title: "Отчёт за месяц",
    sub: "Доходы и расходы",
    icon: "/icons/3d/05-unit-economics.png",
    tone: "bg-peach",
    href: "/app/report",
  },
  {
    title: "Мой путь",
    sub: "Трекборд предпринимателя",
    icon: "/icons/3d/13-cjm-journey.png",
    tone: "bg-lavender",
    href: "/app/journey",
  },
  {
    title: "Алерты",
    sub: "Что требует внимания",
    icon: "/icons/3d/08-proactive-alerts.png",
    tone: "bg-rose",
    href: "/app/alerts",
  },
];

export default function ServicesPage() {
  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <PageHeader
        title="Сервисы"
        subtitle="Каталог рабочих сценариев. Большинство открывают чат с готовым запросом."
        actions={[
          {
            href: "/app/chat",
            label: "Чат",
            icon: MessageCircle,
            primary: true,
          },
        ]}
      />
      <div className="mt-6 grid grid-cols-2 gap-3 sm:gap-4">
        {catalog.map((item) => (
          <IconTile
            key={item.title}
            href={item.href}
            title={item.title}
            sub={item.sub}
            icon={item.icon}
            tone={item.tone}
            size="md"
          />
        ))}
      </div>
    </main>
  );
}
