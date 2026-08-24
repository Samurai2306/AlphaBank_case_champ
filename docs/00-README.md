# Документация — Альфа-Бизнес: Старт

Навигация по Single Source of Truth для стратегии, ТЗ, архитектуры, дизайна и разработки.

**С чего начать разработчику / оператору:** [PROJECT.md](PROJECT.md) (запуск, сценарии, LLM, тесты, VPS).

## Как читать

1. **Стратегия** — зачем продукт и для кого.
2. **Product** — что строим (demo MVP vs production).
3. **Architecture** — как устроено (два контура, общие контракты).
4. **Design** — UI/UX и Generative UI.
5. **Engineering** — стек, репо, промпты, тесты, DoD.
6. **Superpowers** — утверждённый design spec и план реализации демо.

## Содержание

### Strategy

| Файл | Содержание |
|------|------------|
| [strategy/01-executive-summary.md](strategy/01-executive-summary.md) | Executive summary |
| [strategy/02-target-audience-value-canvas.md](strategy/02-target-audience-value-canvas.md) | ЦА и Value Proposition Canvas |
| [strategy/03-cjm-10-levels.md](strategy/03-cjm-10-levels.md) | CJM: 10 уровней |
| [strategy/04-competitive-swot.md](strategy/04-competitive-swot.md) | Конкуренты и SWOT |
| [strategy/05-gtm-and-monetization.md](strategy/05-gtm-and-monetization.md) | GTM и монетизация |
| [strategy/06-risk-compliance.md](strategy/06-risk-compliance.md) | Риски и compliance |

### Product

| Файл | Содержание |
|------|------------|
| [PROJECT.md](PROJECT.md) | **Рабочая дока демо: запуск, сценарии, LLM, деплой** |
| [product/01-prd-demo-mvp.md](product/01-prd-demo-mvp.md) | ТЗ демо (приоритет кода) |
| [product/02-prd-production.md](product/02-prd-production.md) | ТЗ production |
| [product/03-user-stories-acceptance.md](product/03-user-stories-acceptance.md) | User stories + AC |
| [product/04-feature-backlog-prioritized.md](product/04-feature-backlog-prioritized.md) | Бэклог |

### Architecture

| Файл | Содержание |
|------|------------|
| [architecture/01-system-context.md](architecture/01-system-context.md) | C4 L1 |
| [architecture/02-demo-architecture.md](architecture/02-demo-architecture.md) | Demo |
| [architecture/03-production-architecture.md](architecture/03-production-architecture.md) | Production |
| [architecture/04-agentic-rag.md](architecture/04-agentic-rag.md) | Multi-agent RAG |
| [architecture/05-api-contracts.md](architecture/05-api-contracts.md) | Tools / API |
| [architecture/06-data-model.md](architecture/06-data-model.md) | Данные |
| [architecture/07-security-guardrails.md](architecture/07-security-guardrails.md) | Security |
| [architecture/08-deployment-and-ops.md](architecture/08-deployment-and-ops.md) | Deploy & ops |
| [architecture/09-demo-backend-go.md](architecture/09-demo-backend-go.md) | **Demo backend на Go (детально)** |

### Design

| Файл | Содержание |
|------|------------|
| [design/01-design-system.md](design/01-design-system.md) | Design system |
| [design/02-ux-principles.md](design/02-ux-principles.md) | UX principles |
| [design/03-screen-map-and-flows.md](design/03-screen-map-and-flows.md) | Экраны и флоу |
| [design/04-generative-ui-catalog.md](design/04-generative-ui-catalog.md) | SDUI catalog |
| [design/05-motion-voice-tone.md](design/05-motion-voice-tone.md) | Motion, voice & tone |
| [design/06-alfa-visual-language.md](design/06-alfa-visual-language.md) | Как выглядит дизайн Альфы / Альфа-Бизнес |
| [design/07-dashboard-patterns.md](design/07-dashboard-patterns.md) | Топовые fintech-дашборды → Home Copilot |
| [design/08-alfa-marketing-ui-patterns.md](design/08-alfa-marketing-ui-patterns.md) | Паттерны UI с референсов alfabank.ru |
| `design/references/alfa-web-2026/` | Скриншоты-референсы (7 шт.) |
| [`design/assets/3d-icons/`](design/assets/3d-icons/) | Пак 3D-иконок Copilot (18 шт.) |

### Engineering

| Файл | Содержание |
|------|------------|
| [engineering/01-tech-stack-demo.md](engineering/01-tech-stack-demo.md) | Demo stack |
| [engineering/02-tech-stack-production.md](engineering/02-tech-stack-production.md) | Production stack |
| [engineering/03-repo-structure.md](engineering/03-repo-structure.md) | Структура репо |
| [engineering/04-prompt-library.md](engineering/04-prompt-library.md) | Промпты |
| [engineering/05-testing-strategy.md](engineering/05-testing-strategy.md) | Тесты |
| [engineering/06-definition-of-done.md](engineering/06-definition-of-done.md) | DoD |
| [engineering/07-server-deployment.md](engineering/07-server-deployment.md) | **VPS / DNS / TLS — runbook** |
| [ops/vps-bot-project-swap-reversible.md](ops/vps-bot-project-swap-reversible.md) | Обмен с B.O.T.-Project на `bot-project.ru` и откат |

### Superpowers & research

| Файл | Содержание |
|------|------------|
| [superpowers/specs/2026-07-19-ai-business-copilot-design.md](superpowers/specs/2026-07-19-ai-business-copilot-design.md) | Design spec |
| [superpowers/plans/2026-07-19-demo-mvp-implementation.md](superpowers/plans/2026-07-19-demo-mvp-implementation.md) | План кода демо |
| [research/sources.md](research/sources.md) | Источники |
| `research/source-pdfs/` | Исходные PDF |

## Критерии успеха кейса (чеклист)

Решение успешно, если оно:

1. Основано на реальном понимании молодых предпринимателей.
2. Показывает пользу Альфа-Банка.
3. Даёт долгосрочную стратегию привлечения / сопровождения / удержания.
4. Использует ИИ для конкретных болей.
5. Встраивает продукты Платежного бизнеса в путь.
6. Предлагает обоснованные новые/доработанные продукты.
7. Учитывает тех., юр., регуляторные и репутационные ограничения.
