# Demo MVP Implementation Plan

> **Status:** Implemented (Phase 0–4). Cursor plan file left untouched; this doc is the living checklist.

**Goal:** Ship a pitch-ready AI Business Copilot web demo (beauty persona «Маша») with streaming chat, router agents, and Generative UI cards on mocked banking/FNS data.

**Architecture:** Next.js 15 ↔ **Go API** (chi, agents, tools, SSE); shared Zod/JSON contracts; memory store (+ Postgres/Redis optional); OpenAI-compatible LLM client in Go. Spec: `docs/architecture/09-demo-backend-go.md`.

**LLM:** `DEMO_OFFLINE=1` by default; with `LLM_API_KEY` + `DEMO_OFFLINE=0` narrative polish via `internal/llm`. Tax numbers always from `internal/calc`.

---

## Global Constraints

- RU UI; disclaimer on all tax surfaces.  
- Same tool names/SDUI types/intents as production docs.  
- **Demo backend is Go only** — no FastAPI/Python runtime.  
- No real money transfers — draft + `confirmed_mock` only.  
- No secrets in git; `.env.example` only.  

---

### Task 1: Monorepo scaffold (Next + Go)

- [x] Initialize web (Next.js TS App Router) and Go API with `GET /api/v1/health` + `/api/v1/ready`
- [x] Add Postgres+pgvector and Redis to Compose; api depends_on healthy db
- [x] Wire `.env.example`: `DATABASE_URL`, `REDIS_URL`, `LLM_*`, `DEMO_TOKEN`, `CORS_ORIGINS`
- [x] Deploy stubs under `deploy/`
- [x] 3D icons in `apps/web/public/icons/3d/`

### Task 2: Shared contracts

- [x] Intent enum and SDUI component union (`packages/shared`)
- [x] Golden JSON `packages/shared/fixtures/tax_card.json` + `go test` unmarshal

### Task 3: DB migrations + seed persona

- [x] Migration stub + memory seed `masha_nails` (profile, txn, piggy 6%)
- [x] `cmd/seed` present; README documents offline memory default

### Task 4: Calc + mock banking ports

- [x] TDD `CalculateTax(NPD/USN_6)` + unit economics
- [x] Mocks: memory transactions, 115fz risk, payment draft

### Task 5: Chat SSE endpoint

- [x] Go `POST /api/v1/chat` emits SSE (`token` / `sdui` / `done`)
- [x] Web SSE client + disclaimer + thinking state
- [x] Offline path works without API key

### Task 6: Router + Tax agent + TaxCard

- [x] Router intents + Tax path → calc + TaxCard SDUI
- [x] Web renders `TaxCard`

### Task 7: Onboarding + PaymentDraft + Piggy

- [x] OnboardingSummary, PaymentDraftCard (`confirmed_mock`), EnpPiggyBank

### Task 8: Unit economics + Compliance light

- [x] UnitEconomicsChart + ComplianceTrafficLight

### Task 9: Landing + Home + demo script

- [x] Landing Alfa bento; Home decision-first tax hero; `/demo` script buttons
- [x] CJM `/app/journey`

### Task 10: Legal scanner (P1)

- [x] `POST /api/v1/documents/legal` PDF ≤10MB + LegalFlagsList UI upload

### Task 11: E2E smoke + README polish

- [x] Playwright `e2e/masha.spec.ts`
- [x] README <15 min + `DEMO_OFFLINE=1`
- [x] `go test ./...` green; web production build green

### Task 12: Deploy stubs for future VPS

- [x] `deploy/docker-compose.prod.yml`, nginx, remote scripts
- [x] Documented in README «Deploy»

---

## Demo script (wired in `/demo` + offline fixtures)

1. Онбординг → `OnboardingSummary`
2. Unit-экономика → chart
3. Налог → `TaxCard`
4. Платёжка → `PaymentDraftCard` → confirm mock
5. Копилка 6% → `EnpPiggyBank`
6. ИНН → `ComplianceTrafficLight`

## Done when

See checked pitch checklist in `docs/engineering/06-definition-of-done.md`.
