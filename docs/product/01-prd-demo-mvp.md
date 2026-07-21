# PRD — Demo MVP (Foundation)

**Статус:** Ready for implementation  
**Цель:** Питч-/чемпионат-ready web-приложение, доказывающее ценность Copilot на сценарии beauty-мастера.  
**Стек:** Next.js 15 + **Go API** (agents/SSE) + Postgres/pgvector + Redis.  
Спека бэкенда: [../architecture/09-demo-backend-go.md](../architecture/09-demo-backend-go.md).  
Сервер: [../engineering/07-server-deployment.md](../engineering/07-server-deployment.md).

## Problem / Outcome

Молодой мастер beauty боится налогов и оформления. За 5–8 минут демо пользователь проходит: идея → расчёт → выбор режима → «открытие» счёта (мок) → налоговая копилка / расчёт налога с GenUI-карточками.

**Success metrics демо:**

- Сценарий «Маша, маникюр» проходит без оператора.
- ≥4 GenUI-виджета реально рендерятся из tool results.
- Налоговый ответ содержит сумму + дисклеймер + (желательно) ссылку на источник RAG.
- Нет side-effect реальных платежей.

## Personas

1. **Маша, 22** — мастер маникюра, частично в серую, доход ~150–250к/мес, смотрит Т-Банк.
2. **Жюри / стейкхолдер банка** — смотрит логику стратегии и связь с Платежным бизнесом.

## Must Have (P0)

| ID | Фича | Описание |
|----|------|----------|
| D1 | Chat shell | Стриминг чата, состояния думает/отвечает, дисклеймер |
| D2 | Router agent | Классификация intent → агент |
| D3 | Onboarding chat | Извлечение сферы, оборота, города; `OnboardingSummary` |
| D4 | Tax Calculator | НПД / УСН 6% на мок-данных; `TaxCard` |
| D5 | Unit economics | Ввод параметров → точка безубыточности; `UnitEconomicsChart` |
| D6 | Payment draft | Черновик платежа налога; `PaymentDraftCard` + confirm без перевода |
| D7 | Compliance light | Мок проверки ИНН; `ComplianceTrafficLight` |
| D8 | ENP piggy bank | % откладывается с поступлений (мок); `EnpPiggyBank` |
| D9 | Seed + script | Демо-аккаунт и scripted prompts |
| D10 | Landing | Экран «Банк, с которым начинается бизнес» |
| D10b | Home dashboard | Decision-first Home: hero-налог + копилка + инсайты + CTA в чат ([design/07-dashboard-patterns.md](../design/07-dashboard-patterns.md)) |

## Should Have (P1)

| ID | Фича |
|----|------|
| D11 | Legal Scanner: upload PDF → `LegalFlagsList` |
| D12 | Голосовой ввод (Web Speech API) — nice-to-have |
| D13 | CJM trackboard (уровни 1–5 визуально) |
| D14 | Eval harness: 20 golden tax questions |

## Out of Scope (demo)

- Реальные API банка/ФНС/Госуслуг
- React Native / Kafka / on-prem LLM
- УКЭП/УНЭП
- Реальные деньги и SMS-подпись

## Functional requirements

### Chat & GenUI

- Ответы не только Markdown: tool result → React component по типу SDUI.
- Каталог: см. [../design/04-generative-ui-catalog.md](../design/04-generative-ui-catalog.md).

### Tools (имена = production)

- `get_client_transactions`
- `check_counterparty_risk`
- `create_payment_draft`
- `fetch_fns_debt`
- `calculate_tax` (demo helper, остаётся в shared)
- `compute_unit_economics`
- `scan_legal_document`
- `update_onboarding_profile`

### Non-functional

- RU UI; mobile-first web (390px+).
- TTI чата < 2s на warm; первый токен стрима < 3s при доступном LLM.
- Доступность: контраст WCAG AA на финансовых суммах.
- `.env.example` без секретов.

## Acceptance (сводка)

Полные AC — в [03-user-stories-acceptance.md](03-user-stories-acceptance.md). Демо принимается, если P0 D1–D10 зелёные и сценарий жюри воспроизводим из README.

## Roadmap demo

1. Scaffold monorepo + chat streaming  
2. Router + tax/onboarding tools + GenUI  
3. Seed persona + landing  
4. Legal + unit econ + compliance polish  
5. E2E smoke + eval  
