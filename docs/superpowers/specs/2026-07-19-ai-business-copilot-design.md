# AI Business Copilot — Design Spec

**Date:** 2026-07-19  
**Status:** Approved for documentation phase; implementation follows demo plan  
**Sources:** AlfaBank Developer PRD, Master Strategy, championship case, ЦА research

## Problem

Молодые предприниматели (ядро — beauty 17–25) сталкиваются с фрагментированным стартом бизнеса, страхом налогов/ФНС и слабой поддержкой. Средний клиент Платежного бизнеса ~42; сегмент Gen Z недообслужен. РКО — commodity; нужна уверенность и «CFO в кармане».

## Goals

1. Доказать на качественном **web demo**, что Copilot закрывает боли легализации/налогов/платежей через Generative UI.  
2. Зафиксировать **production blueprint** (RN, микросервисы, Kafka, on-prem LLM) с общими контрактами.  
3. Связать путь пользователя с продуктами Платежного бизнеса Альфа-Банка.  
4. Соблюсти юр./репутационные ограничения (RAG, HITL, дисклеймеры).

## Non-goals (now)

- Писать application code в фазе документации.  
- Реальные деньги, СМЭВ, on-prem GPU в demo.  
- Замена официального налогового консультанта.

## Users

- **Primary:** мастер beauty (НПД / «в серую»).  
- **Secondary:** жюри кейса / стейкхолдеры банка.  
- **Later:** HoReCa, e-com.

## Solution summary

Multi-agent RAG Copilot: Router → Tax / Legal / CFO / Onboarding. Ответы = текст + SDUI-карточки. Demo: Next.js + **Go API** + mocks (SSE). Prod: RN SuperApp + Go/Java/Python + Kafka + vLLM. VPS: Docker + Nginx/Caddy — `docs/engineering/07-server-deployment.md`.

## Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Delivery model | Hybrid demo + prod docs | Pitch speed + bank realism |
| Demo client | Next.js PWA | Быстрый GenUI, питч |
| Demo backend | **Go** monolith (SSE+agents) | Один бинарь на VPS; ближе к prod BFF |
| Prod client | React Native | SuperApp PRD |
| Contracts | Shared tools/SDUI/intents | Не выбрасывать демо |
| Tax UX | Calculator + RAG + disclaimer | Anti-hallucination |
| Delight colors | Landing only | Trust in money UI |
| Niche | Beauty first | Research + case priorities |

## Information architecture

См. `docs/design/03-screen-map-and-flows.md`, `docs/strategy/03-cjm-10-levels.md`.

## UX / UI

См. `docs/design/*`. Brand tokens PRD + championship delight на маркетинге. Voice: уверенный, лаконичный, эмпатичный.

Дополнено research (2026-07-19):

- Визуальный язык Альфы / Альфа-Бизнес → `docs/design/06-alfa-visual-language.md` (`#EF3124`, платина, Interface Sans, core-components).
- Decision-first Home dashboard → `docs/design/07-dashboard-patterns.md` (Role–Metric–Density–Action, hybrid ask+structure).
- Актуальные скрины alfabank.ru → `docs/design/08-alfa-marketing-ui-patterns.md` + `docs/design/references/alfa-web-2026/` (bento, pastel, 3D, pill segments, services rows).

## Architecture

См. `docs/architecture/*`. C4, agentic RAG, security guardrails, dual deploy.

## Risks

См. `docs/strategy/06-risk-compliance.md`.

## Success metrics

- Demo: сценарий «Маша» без оператора; ≥4 GenUI widgets.  
- Strategy KPIs: CAC↓, activation↑, survival↑ (prod horizon).  
- Case rubric: реальные боли, ИИ, платежный бизнес, long-term, constraints.

## Open points (non-blocking)

- Точная лицензия шрифта Alfa Slab One → fallback Manrope/Geist.  
- Выбор Qdrant vs Milvus — за ИБ банка.  
- Установка `parallel-cli` для углублённого research — опционально.

## Spec self-review

- [x] Нет TBD по стеку demo  
- [x] Нет противоречия demo vs prod contracts  
- [x] Scope документации отделён от кода  
- [x] Production описан достаточно для следующей команды  

## Next

Implement demo per `docs/superpowers/plans/2026-07-19-demo-mvp-implementation.md`.
