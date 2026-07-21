# Design System

## Принцип двойного слоя

| Слой | Где | Характер |
|------|-----|----------|
| **Bank Trust** | Чат, суммы, платежи, налоги | Альфа-красная дисциплина, воздух, минимум декора |
| **Championship Delight** | Landing, первый onboarding splash | Энергия кейса (purple/lime), не переносить в money UI |

## Цвета (PRD + бренд Альфы)

Официальный **Альфа-красный** `#EF3124` и **Альфа-платиновый** `#D1D5D8` зафиксированы в [веб-стандартах alfabank.ru](https://www.artlebedev.com/alfabank/guides/Alfa-guides.pdf). Полный визуальный контекст: [06-alfa-visual-language.md](06-alfa-visual-language.md).

| Token | Hex | Use |
|-------|-----|-----|
| `--color-primary` | `#EF3124` | CTA, акценты, нотификации (Альфа-красный) |
| `--color-platinum` | `#D1D5D8` | Вторичные плашки, dividers, inactive chrome |
| `--color-ink` | `#0B1117` | Заголовки, dark bg |
| `--color-surface` | `#FFFFFF` | Карточки, light bg (паттерн Альфа-Бизнес) |
| `--color-muted` | `#F2F3F5` | Фон чата / canvas, disabled |
| `--color-success` | `#00C853` | Успех, green compliance |

### Delight (только маркетинг)

| Token | Approx | Use |
|-------|--------|-----|
| `--color-spark` | `#6C3CF0` | Hero background |
| `--color-volt` | `#C8F560` | Hero accents |
| `--color-smile` | `#FFB020` | Мелкие delight-детали |

**Запрет:** purple-on-white тема во всём приложении; glow-everywhere; «AI purple SaaS» клише в продуктовых экранах.

## Типографика

Продукты Альфы (Альфа-Бизнес) переходят на **Alfa Interface Sans** ([Habr](https://habr.com/ru/companies/alfa/articles/1048062/), [core-components](https://github.com/core-ds/core-components/commit/b832d8507355178a2b2d3fa24220be6e1f9d8ee3)).

| Role | Font | Notes |
|------|------|-------|
| UI / chat / dashboard | **Alfa Interface Sans** (CDN банка) или Manrope fallback | 16px body, line-height 150%; ритм «легче и воздушнее» как в редизайне Альфа-Бизнес |
| Display / marketing | Slab/display (Alfa Slab One если лицензия) или тот же Interface Bold | Только landing / чемпионат |
| Figures | Tabular nums | Суммы налогов и платежей — крупнейший акцент иерархии |

Не строить UI только на Inter. Не использовать slab в money UI.

**Production:** по возможности `@alfalab/core-components` Typography + tokens ([Storybook](https://alfa-laboratory.github.io/core-components/)).  
**Demo:** CSS-токены + CDN/fallback шрифт; визуальный паритет важнее npm-зависимости.

## Spacing & radius

По референсам сайта Альфы ([08-alfa-marketing-ui-patterns.md](08-alfa-marketing-ui-patterns.md)):

- База 4px; внутренний padding карточек 16–24.  
- **Витрина / landing / Home rows:** radius карточек **24–32px**; pills `999px`.  
- **GenUI money cards в чате:** radius 20–24px; только интерактивные поверхности.  
- Canvas страницы: `#F5F5F5` / `#F2F3F5`, карточки белые или pastel.  
- Макс. ширина колонки чата 720px desktop; mobile full-bleed.

## Elevation

- Мягкая рассеянная тень **или** разделение тоном (canvas vs surface) — как на alfabank.ru.  
- Без тяжёлого multi-layer glow.

## Icons & 3D

| Зона | Стиль |
|------|--------|
| Landing / promo | Glossy 3D / clay (карта, монета, %) — реф. `references/alfa-web-2026/` |
| Home service rows | Компактный 3D справа (как «Сервисы») |
| Chat / compliance | Линейные 24px или простой светофор; без emoji как UI |

## Dark mode

- Marketing dark bento (чёрный canvas + цветные карточки) — **landing only**.  
- Prod SuperApp: deep black `#0B1117` опционально.  
- Demo workspace (налоги/чат): light-first.

## Components (atomic)

`ButtonPrimary`, `ButtonGhost`, `ChatBubble`, `Composer`, `DisclaimerBar`, `CjmStepper`, `HeroMetric`, `InsightRow`, плюс GenUI registry.

Дашборд Home: [07-dashboard-patterns.md](07-dashboard-patterns.md) (Role–Metric–Density–Action).  
UI-паттерны Альфы: [08-alfa-marketing-ui-patterns.md](08-alfa-marketing-ui-patterns.md) + скрины в `references/alfa-web-2026/`.

Доп. компоненты: `SegmentedPill`, `ServiceRow`, `PastelPromoCard`, `InvertCtaCard` (чёрная row как «Оформить самозанятость»).

## Motion tokens

См. [05-motion-voice-tone.md](05-motion-voice-tone.md). Длительности: 150 / 250 / 400ms; easing standard.
