export type Scenario = {
  id: string;
  label: string;
  hint: string;
  text: string;
  group: "Старт" | "Деньги" | "Риски" | "Рост";
};

/** Deterministic demo paths — same prompts as API golden / tasks. */
export const SCENARIOS: Scenario[] = [
  {
    id: "onboarding",
    label: "Онбординг",
    hint: "Собрать профиль Маши",
    group: "Старт",
    text: "Привет! Я делаю маникюр на дому в Казани, где-то 180 тысяч в месяц, пока без ИП.",
  },
  {
    id: "unit",
    label: "Безубыточность",
    hint: "Сколько клиентов в день",
    group: "Рост",
    text: "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
  },
  {
    id: "tax",
    label: "Налог НПД",
    hint: "Сколько отложить",
    group: "Деньги",
    text: "Сколько мне отложить на налог за этот месяц?",
  },
  {
    id: "payment",
    label: "Платёжка",
    hint: "Черновик поручения",
    group: "Деньги",
    text: "Сформируй платёжку.",
  },
  {
    id: "piggy",
    label: "Копилка ЕНП",
    hint: "Автоотчисления 6%",
    group: "Деньги",
    text: "Включи копилку 6%.",
  },
  {
    id: "inn",
    label: "Проверка ИНН",
    hint: "Светофор 115-ФЗ",
    group: "Риски",
    text: "Проверь ИНН 1650987654",
  },
  {
    id: "legal",
    label: "Договор",
    hint: "Или загрузите PDF ниже",
    group: "Риски",
    text: "Разбери договор аренды PDF",
  },
];

export const QUICK_CHIPS = SCENARIOS.filter((s) =>
  ["tax", "payment", "unit", "inn", "piggy"].includes(s.id),
);
