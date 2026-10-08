# Demo elevation notes (2026-07-19)

## What changed after pitch MVP

### Critical UX
- PaymentDraftCard updates to `confirmed_mock` / `cancelled` in-place after API calls.
- OnboardingSummary «Всё верно» → `PATCH /me/onboarding` + saved state.

### Tax model
- NPD mixed rate: default 85% B2C → **7 740 ₽** on 180k (4%/6%).
- Home hero uses **conservative ceiling 10 800 ₽** («с запасом»).
- TaxCard shows assumptions + `source_ref` (`kb:npd-rates-demo#2026`).

### Backend contracts
- `internal/tools` registry with production tool names.
- Legal upload: PDF magic bytes, ≤10MB, fixture flags by filename.
- Cancel draft endpoint; `/ready` documents memory-store demo mode.
- Eval harness: 20 golden offline prompts (`internal/eval`).

### UI elevation
- Landing: dark bento + phone mock + carousel + motion.
- Home segments Today / Month / Path actually change hero content.
- CJM journey reads `/me` and links next actions.
- Framer Motion on GenUI; chart axes/tooltip/breakeven marker.

## Verification evidence
- `go test ./...` green (calc, agents, eval×20, sdui golden).
- `npm run build` green.
- Playwright `e2e/masha.spec.ts` — 6 passed (landing, home, demo, tax→confirm, onboarding, guardrail).
- `scripts/debug-demo-session.ps1` for API smoke.

## Remaining (honest, not blockers for pitch)
- Сессия и корпус живут в памяти процесса.
- Cloud LLM is narrative polish, not token-streaming upstream.
- Legal scan is fixture/heuristic, not NLP extraction of full PDF text.
