# Современные fintech-дашборды и Home Copilot

Как сделать **топовый** финансовый UI 2025–2026 и как это ложится на AI Business Copilot (не «ещё одна сетка KPI»).

## Главный сдвиг

| Было | Стало |
|------|-------|
| Дашборд = склад метрик | Дашборд = **ответ на один главный вопрос** + следующий шаг |
| Пользователь ищет инсайт | Система **проактивно** подсвечивает аномалии |
| Чат отдельно от цифр | **Hybrid:** structured analysis + ask-first Copilot |
| AI как баннер «powered by AI» | AI **тише**, внутри workflow, с объяснимостью |

Опора: [Fintech Dashboard Design Guide (Masterly, 2026)](https://www.themasterly.com/blog/fintech-dashboard-design-guide), [Outcrowd / Financial Dashboard 2026](https://medium.com/outcrowd/how-to-decide-what-your-financial-dashboard-actually-needs-in-2026-b4bddf5bc0ba), [Dashboard UX 2026 practices](https://webdesignerindia.medium.com/ux-design-for-dashboards-best-practices-663767ff19c1), [Proactive copilots / intent-driven UX](https://craftinnovations.global/future-of-financial-ux-webinar/).

---

## Framework: Role → Metric → Density → Action

Из [Masterly](https://www.themasterly.com/blog/fintech-dashboard-design-guide):

| Шаг | Вопрос | Copilot (beauty / Маша) |
|-----|--------|-------------------------|
| **Role** | Кто смотрит? | Мастер 17–25, не CFO корпорации |
| **Metric** | Одно число сверху | «Налог к отложению в этом месяце» или «Баланс копилки ЕНП» |
| **Density** | Сколько информации | **Minimal / founder-like** (как Mercury), не dense Ramp/Brex |
| **Action** | Что сделать дальше | CTA: «Сформировать платёжку» / «Включить копилку» / «Открыть Copilot» |

Правило: если убрать все виджеты кроме hero-метрики и одного CTA — экран всё ещё должен иметь смысл.

---

## Паттерны топового дашборда (чеклист)

### 1. Decision-first layout

- Иерархия: **1 hero metric → 2–3 secondary → feed инсайтов → список действий**.  
- Не начинать с таблицы из 12 KPI ([Dashboard UX 2026](https://webdesignerindia.medium.com/ux-design-for-dashboards-best-practices-663767ff19c1)).

### 2. Progressive disclosure

- Новичок (CJM 1–5): налог, копилка, путь легализации.  
- Позже: P&L, 115-ФЗ, овердрафт.  
- Сложность прячется, пока не созрел уровень ([Outcrowd](https://medium.com/outcrowd/how-to-decide-what-your-financial-dashboard-actually-needs-in-2026-b4bddf5bc0ba)).

### 3. Proactive insights (не спам)

Примеры строк для Copilot Home:

- «До срока ЕНП 9 дней — в копилке не хватает 2 400 ₽»  
- «Выручка за неделю −18% к прошлой — риск аренды»  
- «Контрагент в жёлтой зоне — проверьте перед оплатой»

Каждый инсайт: **факт → почему важно → CTA**.  
AI-инсайты визуально отделены от «сырых» чисел ([Dashboard UX 2026](https://webdesignerindia.medium.com/ux-design-for-dashboards-best-practices-663767ff19c1)).

### 4. Explainable AI

Обязательно:

1. Ссылка на источник данных / RAG chunk  
2. Уровень уверенности или режим («расчётная модель демо»)  
3. Override / «уточнить»  
4. Дисклеймер  

Без этого — не банковский продукт.

### 5. Hybrid: Ask + Structure

- **Ask-first** (чат): «Сколько налогов?», разведка, онбординг.  
- **Structured**: карточки налога, график безубыточности, светофор — для точности и повторяемости ([Outcrowd](https://medium.com/outcrowd/how-to-decide-what-your-financial-dashboard-actually-needs-in-2026-b4bddf5bc0ba)).  
- Generative UI = мост: ответ чата **собирает** structured-карточку.

### 6. Intentional friction

Высокое трение только на деньгах/подписи; низкое — на просмотре баланса и FAQ ([Masterly](https://www.themasterly.com/blog/fintech-dashboard-design-guide)).

### 7. Trust = структура, не бейджи

Прозрачные статусы («черновик», «ожидает подтверждения»), видимые комиссии/суммы, processing states. «Secure»-наклейки вторичны ([Masterly](https://www.themasterly.com/blog/fintech-dashboard-design-guide)).

### 8. Performance UX

Цель: first meaningful content **&lt; 1.5s** на Home ([Dashboard UX 2026](https://webdesignerindia.medium.com/ux-design-for-dashboards-best-practices-663767ff19c1)). Скелетоны метрик, не пустой белый экран.

### 9. Mobile-first one-hand

Нижний composer / FAB; hero metric читается на 390px без pinch. Thumb-zone для confirm.

### 10. Proactive copilots / features find the user

Не заставлять искать «налоги» в меню — сигнал приходит в ленту/чат, когда пора ([Craft Innovations webinar](https://craftinnovations.global/future-of-financial-ux-webinar/)). Метрика успеха: **time-to-task**, не time-in-app.

---

## Home Copilot (demo + prod) — wireframe логика

Ориентир по композиции: реф. «Сервисы» Альфы (`references/alfa-web-2026/05-services-dashboard.png`) + decision-first metric. Паттерны: [08-alfa-marketing-ui-patterns.md](08-alfa-marketing-ui-patterns.md).

```text
┌─────────────────────────────────────────┐
│  Copilot                    [Сегодня|…] │  ← pill segment
│  Привет, Маша · уровень: Налоги         │
├─────────────────────────────────────────┤
│  HERO (white / soft card, r=28)         │
│  К отложению в июле                     │
│  10 800 ₽              [Платёжка red]   │
├───────────────┬─────────────────────────┤
│ Копилка ЕНП   │ Светофор                │
│ pastel/white  │                         │
├───────────────┴─────────────────────────┤
│ Service rows (как Альфа «Сервисы»)      │
│ Рассчитать налог              [3D/%]    │
│ Проверить контрагента         […]       │
│ ■ Оформить НПД / легализация  [black]   │  ← invert CTA
├─────────────────────────────────────────┤
│ Инсайты · горизонтальный peek            │
├─────────────────────────────────────────┤
│ [Спросить Copilot………………]        [🎤]   │
└─────────────────────────────────────────┘
```

Компоненты = существующий SDUI-каталог ([04-generative-ui-catalog.md](04-generative-ui-catalog.md)):  
`TaxCard` / hero, `EnpPiggyBank`, `ComplianceTrafficLight`, CJM stepper, composer → чат.

---

## Визуальные правила дашборда (под Альфу)

См. также [06-alfa-visual-language.md](06-alfa-visual-language.md).

| Элемент | Do | Don't |
|---------|-----|-------|
| Hero metric | Огромная цифра, tabular, ink/black | Мелкий KPI в сетке 2×4 |
| Chart | Один простой line/area (безубыточность) | 3D pie, радуга серий |
| Cards | Белая поверхность, radius 16, hairline | Glassmorphism + glow |
| Accent | Красный CTA | Rainbow status chips |
| AI block | Тихая плашка «Copilot заметил» | «✨ AI MAGIC» |
| Empty | Спокойный next step | Иллюстрация-мем |

Плотность: **founder-minimal**. Beauty-ЦА = Mercury-класс ясности, не Bloomberg terminal.

---

## Data viz guidelines

1. Одна идея на график.  
2. Ось и единицы на русском (`₽`, «клиентов/день»).  
3. Forecast — отдельный стиль (пунктир + band), не смешивать с fact без подписи.  
4. Accessibility: не только цвет (иконка уровня светофора + текст).  
5. Налог/деньги — всегда явное «на какую дату».

---

## Scope по поставкам

| Поставка | Dashboard |
|----------|-----------|
| **Demo MVP** | Home с hero-налогом + копилка + 1–2 инсайта + CTA в чат; полный чат GenUI |
| **Prod Phase 1** | Read-only метрики из транзакций + RAG FAQ |
| **Prod Phase 2** | Actionable cards (платёж, контрагент) |
| **Prod Phase 3** | Proactive feed из Kafka (касса, ЕНП) |

---

## Anti-patterns (запрет)

- KPI soup без главного числа  
- Dark neon cyberpunk «AI bank»  
- Автоплей видео на Home  
- Скрытые комиссии  
- AI-инсайт без CTA и без источника  
- Дашборд, который нельзя объяснить за 10 секунд жюри  

---

## Acceptance (дизайн)

- [ ] Жюри за 5 секунд понимает: «ей надо отложить X ₽»  
- [ ] Один primary CTA на первом экране  
- [ ] Экран узнаваем как Альфа (красный + белые панели + UI sans)  
- [ ] Инсайт Copilot ведёт в чат или GenUI-действие  
- [ ] Mobile 390px без горизонтального скролла контента  
