# Feature Backlog (приоритизированный)

Легенда: **D** = Demo, **P** = Production. Score = Impact × (1/Effort) качественно.

## Sprint A — Demo Foundation

| ID | Фича | Layer | Priority | Depends |
|----|------|-------|----------|---------|
| A1 | Monorepo scaffold `apps/web`, `apps/api` (**Go**), `packages/shared`, `deploy/` stubs | D | P0 | — |
| A2 | Chat SSE API (Go) + UI client | D | P0 | A1 |
| A3 | Router agent + intent enum (Go) | D | P0 | A2 |
| A4 | Tax tools + `TaxCard` | D | P0 | A3 |
| A5 | Onboarding + `OnboardingSummary` | D | P0 | A3 |
| A6 | Payment draft + confirm mock | D | P0 | A4 |
| A7 | Seed persona «Маша» + demo script | D | P0 | A5 |
| A8 | Landing / pitch screen | D | P0 | A1 |

## Sprint B — Demo Polish

| ID | Фича | Layer | Priority | Depends |
|----|------|-------|----------|---------|
| B1 | Unit economics + chart | D | P0 | A3 |
| B2 | ENP piggy bank | D | P0 | A7 |
| B3 | Compliance traffic light | D | P0 | A3 |
| B4 | Legal Scanner PDF | D | P1 | A2 |
| B5 | CJM trackboard 1–5 | D | P1 | A5 |
| B6 | Golden set eval (tax) | D | P1 | A4 |
| B7 | E2E smoke Playwright | D | P1 | A8 |
| B8 | Deslop / a11y pass | D | P0 | B1–B3 |

## Production Phase 1

| ID | Фича | Priority |
|----|------|----------|
| P1.1 | RN shell + Copilot Tab/FAB | P0 |
| P1.2 | On-prem LLM gateway + provider abstraction | P0 |
| P1.3 | Qdrant/Milvus + корпус НК/тарифов | P0 |
| P1.4 | PII NER masking middleware | P0 |
| P1.5 | Read-only transactions adapter | P0 |
| P1.6 | Anti-jailbreak + audit log | P0 |

## Production Phase 2

| ID | Фича | Priority |
|----|------|----------|
| P2.1 | `create_payment_draft` → core + sign UX | P0 |
| P2.2 | `check_counterparty_risk` 115-ФЗ | P0 |
| P2.3 | `fetch_fns_debt` СМЭВ | P0 |
| P2.4 | SDUI native component pack parity with web | P0 |
| P2.5 | OCR pipeline (Vision) | P1 |
| P2.6 | Госключ / УКЭП в чате | P1 |

## Production Phase 3

| ID | Фича | Priority |
|----|------|----------|
| P3.1 | Kafka txn stream → alert service | P0 |
| P3.2 | Cash-gap ML + овердрафт offer | P0 |
| P3.3 | CRM/Yclients growth agent | P1 |
| P3.4 | Copilot Pro BI | P2 |

## Icebox

- Мультиязычность  
- Полная замена бухгалтера  
- Автоплатежи без confirm  
- Обучение LLM на клиентских данных  
