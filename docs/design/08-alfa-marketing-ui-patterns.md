# Паттерны UI Альфы (референсы сайта / приложения 2026)

Источник: скриншоты актуального маркетингового и сервисного UI Альфа-Банка (июль 2026).  
Файлы: [`references/alfa-web-2026/`](references/alfa-web-2026/).

Этот документ **дополняет** [06-alfa-visual-language.md](06-alfa-visual-language.md) конкретными layout/component-паттернами, которые жюри узнает как «Альфа».

## Карта референсов

| # | Файл | Паттерн |
|---|------|---------|
| 01 | `01-product-grid-choose-best.png` | Product grid 2×4, pill segment control, 3D-иконки в карточках |
| 02 | `02-dark-bento-osago.png` | Dark canvas + color bento + красный CTA |
| 03 | `03-carousel-get-more.png` | Horizontal promo carousel, pastel cards, circular arrows |
| 04 | `04-cashback-phone-mock.png` | Split banner + phone mock + pastel app tiles |
| 05 | `05-services-dashboard.png` | Services: currency table + stacked service rows + inverted CTA |
| 06 | `06-offers-podeli.png` | Pill icon filter bar + pastel offer cards |
| 07 | `07-about-bank-achievements.png` | Bento «о банке»: black feature + light stats |

---

## Визуальный язык (что повторять)

### Canvas и поверхности

| Слой | Значение |
|------|----------|
| Page canvas | Светло-серый / off-white (`#F5F5F5`–`#F7F7F7`), не чисто «плоский» `#FFF` на весь экран |
| Card surface | Белый или **мягкий pastel** (мята, лаванда, персик, роза, циан) |
| Emphasis card | **Чёрный / charcoal** для ключевого оффера (самозанятость, достижения) |
| Marketing dark | Полностью чёрный фон + цветные bento-карточки (реф. 02) |

### Радиусы и плотность

- Карточки: **очень крупный radius** ≈ **24–32px** (мягкий «soap»-вид).  
- Pills / segment control / CTA: **full pill** (капсула).  
- Внутри карточек — щедрый padding (16–24px), воздух между блоками.  
- Тени: мягкие, рассеянные **или** почти отсутствуют (слой за счёт тона canvas/card).

### Цвет

| Роль | Как у Альфы |
|------|-------------|
| Brand / CTA | Красный `#EF3124` — кнопка «Войти», «Стать клиентом», логотип A |
| Pastel fills | Фон промо-карточек (не весь chrome) |
| Text | Ink/black заголовки; secondary — тёмно-серый |
| 3D assets | Насыщенные глянцевые объекты (карты, монеты, % ) — главный «wow» |
| Gold / premium | Металл/золото на achievement-карточках |

**Не путать:** pastel на **маркетинге и витрине**; money/Copilot workspace — светлый canvas + белые панели + красный CTA (см. ниже «где что»).

### Типографика

- Крупный bold sans заголовок секции слева.  
- В карточке продукта: **title centered** + короткий subtext серым.  
- Цифры офферов огромные («30%») — иерархия как у hero-метрик.

### 3D / иллюстрации

- Стиль: glossy / clay / pearlescent 3D, не flat outline icons.  
- Объекты часто «сидят» в центре или справа карточки, иногда выходят за край.  
- В приложении: мелкие 3D в правом нижнем углу плитки.

Для демо Copilot:  
- Готовый пак: [`assets/3d-icons/`](assets/3d-icons/) (18 шт.) — Router, Tax, Legal, Transactions, Unit Econ, Compliance, Onboarding, Alerts, Piggy, Payment, Beauty, Card, CJM, Chat, НПД, Cash gap, pastel badge, hero cluster.  
- Landing / pitch — смело использовать 3D из пака.  
- Home service rows — компактные иконки из пака справа.  
- Chat/confirm — сдержаннее, не каждые 3D на каждом сообщении.

---

## Компонентные паттерны

### 1. Pill segmented control

Реф. 01, 07.

- Трек: светло-серый pill.  
- Active: **тёмный charcoal** + белый текст.  
- Inactive: прозрачный/серый текст.  
- Применение в Copilot: фильтр Home «Сегодня / Месяц / Путь»; CJM tabs.

### 2. Product / feature grid

Реф. 01.

- Равные карточки, large radius, soft fill.  
- Иконка 3D по центру, текст сверху/снизу.  
- Применение: выбор сценария онбординга (маникюр / салон / retail).

### 3. Dark bento hero

Реф. 02.

- Чёрный фон страницы.  
- Крупная цветная hero-карточка + меньшие pastel.  
- Красная pill-кнопка внутри цветной карточки.  
- Применение: **только landing/pitch**, не workspace налогов.

### 4. Promo carousel

Реф. 03.

- Заголовок слева, круглые ← → справа.  
- Горизонтальный скролл pastel-карточек.  
- Применение: «Что умеет Copilot» на landing; инсайты на Home (опционально).

### 5. Phone mock + pastel app tiles

Реф. 04.

- Слева текст + красная pill CTA.  
- Справа red panel + mock телефона с сеткой pastel tiles.  
- Применение: hero лендинга «Банк, с которым начинается бизнес».

### 6. Services dashboard

Реф. 05 — **ближайший к рабочему UI**.

- Заголовок «Сервисы».  
- Слева: data card (курсы) — таблица, флаги, disclaimer.  
- Справа: **stacked rows** — title + subtext + 3D icon справа.  
- Одна row **inverted black** = primary action (у Альфы «Оформить самозанятость» → у нас «Оформить НПД / открыть путь»).  
- Применение: **Home Copilot** и список действий.

### 7. Icon filter bar

Реф. 06.

- Светло-серый horizontal pill с иконками+лейблами.  
- Active — тёмная капсула.  
- Применение: категории GenUI / инструменты Copilot.

### 8. About / trust bento

Реф. 07.

- Black feature card + light stat cards + 3D premium objects.  
- Применение: блок доверия на landing («почему Альфа»).

---

## Где применять в Copilot

| Поверхность | Паттерны | Плотность 3D |
|-------------|----------|--------------|
| Landing `/` | 02, 03, 04, 07 | Высокая |
| Demo script `/demo` | 01 grid кнопок сценариев | Средняя |
| Home `/app` | 05 services + hero metric + soft insights | Низкая–средняя |
| Chat `/app/chat` | Белые GenUI cards на gray canvas; red CTA | Низкая |
| Payment confirm | Intentional friction; без «игрушечных» pastel | Нет |

Правило: **чем ближе к деньгам и налогам — тем спокойнее палитра**; pastel/3D — для привлечения и навигации по возможностям.

---

## Токены (дополнение к design-system)

```text
--radius-card: 28px;
--radius-pill: 999px;
--canvas: #F5F5F5;
--surface: #FFFFFF;
--pastel-mint: #E8F6EE;
--pastel-lavender: #EDE7F6;
--pastel-peach: #FFF0E6;
--pastel-rose: #FEE1E1;
--pastel-cyan: #E3F4F8;
--ink: #0B1117;
--ink-secondary: #6B7280;
--emphasis-invert: #0B1117; /* black CTA cards */
--brand-red: #EF3124;
```

Segment active: `background: #1A1A1A; color: #FFF`.

---

## Anti-patterns (относительно референсов)

- Мелкий radius 8px «bootstrap»  
- Фиолетовый градиент SaaS вместо red+pastel  
- Плоские серые иконки Material на витрине (ок в dense tables, не на promo)  
- KPI soup без bento-иерархии  
- Pastel фон на экране подтверждения платежа  

---

## Чеклист для UI-агента

- [ ] Landing узнаваем как Альфа (red CTA, large radius, 3D/pastel bento)  
- [ ] Home ближе к реф. 05 (rows + optional black CTA «легализация»)  
- [ ] Segment/filter = dark pill active  
- [ ] Chat/money = white cards, `#EF3124` actions, дисклеймер  
- [ ] Референсы просмотрены из `references/alfa-web-2026/` перед вёрсткой  
