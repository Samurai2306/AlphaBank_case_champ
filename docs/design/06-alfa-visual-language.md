# Визуальный язык Альфа-Банка (референс для Copilot)

Документ фиксирует, **как выглядят и ощущаются** цифровые продукты Альфы, и как мы к ним прижимаем demo/production UI. Не заменяет официальный брендбук; опирается на открытые гайды и публикации дизайн-системы.

**Скриншоты актуального UI (2026):** [`references/alfa-web-2026/`](references/alfa-web-2026/) — разбор паттернов в [08-alfa-marketing-ui-patterns.md](08-alfa-marketing-ui-patterns.md).

## Официальные опоры бренда

| Элемент | Значение | Источник |
|---------|----------|----------|
| Альфа-красный | `#EF3124` (RGB 239, 49, 36) | [Корпоративные веб-стандарты / Art. Lebedev](https://www.artlebedev.com/alfabank/guides/Alfa-guides.pdf), брендбук 2020 |
| Альфа-платиновый | `#D1D5D8` | Тот же веб-стандарт (вторичный / плашки) |
| Чёрный / белый | Контрастные пары логотипа (позитив/негатив) | Веб-стандарт |
| Слоган / тон | Прямой, уверенный, без канцелярита | Бренд-коммуникации Альфы |

**Правило логотипа:** не вращать, не менять пропорции, не красить вне палитры ([веб-стандарты](https://www.artlebedev.com/alfabank/guides/Alfa-guides.pdf)).

## Дизайн-система в коде

История и актуальное состояние:

| Поколение | Репозиторий | Статус |
|-----------|-------------|--------|
| Feather / `arui-feather` | [alfa-laboratory/arui-feather](https://github.com/alfa-laboratory/arui-feather) | **Deprecated** — не развивать |
| Core Components | [alfa-laboratory/core-components](https://github.com/alfa-laboratory/core-components) | **Актуальная** React UI-библиотека |
| Tokens / primitives | [alfa-ui-primitives](https://github.com/alfa-laboratory/alfa-ui-primitives) | Цвета, типографика, gaps, иконки (часть color JSON deprecated) |
| Storybook | [alfa-laboratory.github.io/core-components](https://alfa-laboratory.github.io/core-components/) / [core-ds Storybook](https://core-ds.github.io/core-components/master/) | Песочница компонентов |

Для production RN/Web внутри банка: ориентир — **@alfalab/core-components** + токены, не самодельный «ещё один purple SaaS kit».

Для **demo MVP** (чемпионат): визуально **имитируем** язык Альфы (красный, воздух, типографика), без обязательной npm-зависимости от core-components (лицензии/сборка). При возможности — подключить токены/шрифт с CDN Альфы.

## Типографика продуктов

С 2024–2025 в продуктах Альфы (в т.ч. **Альфа-Бизнес**) закреплён фирменный UI-шрифт **Alfa Interface Sans** ([Habr: дизайн-система Альфа-Бизнес / иконки](https://habr.com/ru/companies/alfa/articles/1048062/), [коммит core-components](https://github.com/core-ds/core-components/commit/b832d8507355178a2b2d3fa24220be6e1f9d8ee3)).

Характер интерфейса после редизайна Альфа-Бизнес ([Habr](https://habr.com/ru/companies/alfa/articles/1048062/)):

- легче и аккуратнее;
- больше воздуха и ритма;
- выразительнее формы;
- иконки подружены с новым шрифтом в плотных таблицах, сайдбарах, формах.

**CDN (как в core-components docs):**  
`https://alfabank.servicecdn.ru/media/fonts/alfa-interface-sans_{regular|medium|bold}.woff2`

**Demo fallback:** Manrope / system-ui, если лицензия/оффлайн не позволяют CDN.

PRD упоминал Alfa Slab One для маркетинга — допустим на landing; **продуктовый чат и дашборд** — Interface Sans / UI sans, не slab.

## Как устроен UI Альфа-Бизнес (паттерны)

По разбору навигации B2B ([Habr: обновление бокового меню](https://habr.com/ru/companies/alfa/articles/809991/)):

1. **SharedUI shell** — боковое + верхнее меню, точка входа в чат.  
2. **Контент на белых карточках-подложках** на нейтральном фоне.  
3. **Сайдбар** совмещает навигацию и лёгкое промо; визуально отделён от рабочей области.  
4. **F-паттерн:** важные действия слева/сверху.  
5. Акценты витрин — из **банковской** палитры (красный / нейтрали), не произвольный неон.

### Что это значит для Copilot

| Зона | Поведение |
|------|-----------|
| Встраивание | FAB / Tab внутри shell Альфа-Бизнес (как точка входа рядом с чатом) |
| Рабочая область | Белые поверхности, soft gray canvas `#F2F3F5` / близкий нейтрал |
| Акцент CTA | Только `#EF3124` |
| Вторичный металл | `#D1D5D8` для dividers, chips, inactive |
| Copilot vs chrome | Чуть «воздушнее» и умнее основного банка, но **тот же бренд**, не отдельный стартап-скин |

## Чемпионат / «Альфа Будущее»

Маркетинг кейса использует purple / lime / 3D — это **кампанийный** слой.  
В продукте Платежного бизнеса и Copilot money UI этот слой **не доминирует** (см. [01-design-system.md](01-design-system.md)).

## Чеклист «похоже на Альфу»

- [ ] Красный `#EF3124` на primary CTA, не градиент purple→indigo  
- [ ] Платиновый / серый для вторички; pastel — на promo-карточках  
- [ ] UI-шрифт Interface Sans (или близкий grotesk), 16px/150% в чате  
- [ ] Крупный radius 24–32px; pill segment с тёмным active  
- [ ] Canvas светло-серый + белые/pastel панели, воздух  
- [ ] Витрина: glossy 3D; money UI: спокойнее  
- [ ] Service rows + опциональная чёрная CTA-карточка (легализация)  
- [ ] Суммы — tabular, крупная иерархия  
- [ ] Нет emoji как UI-иконок в money flows  
- [ ] Landing ярче (bento/carousel); `/app` — trust + паттерн «Сервисы»  

## Связанные файлы

- Токены и компоненты демо: [01-design-system.md](01-design-system.md)  
- Дашборд и KPI-паттерны: [07-dashboard-patterns.md](07-dashboard-patterns.md)  
- Источники: [../research/sources.md](../research/sources.md)
