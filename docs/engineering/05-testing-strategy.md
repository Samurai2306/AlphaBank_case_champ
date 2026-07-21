# Testing Strategy

## Пирамида

| Layer | Demo | Production |
|-------|------|------------|
| Unit | `go test` `internal/calc`; Vitest formatters | + Java unit |
| Contract | Zod (TS) ↔ Go structs / golden JSON SDUI | Pact/OpenAPI vs adapters |
| Agent eval | Golden 20 tax Qs (`calc` + offline SSE) | Expanded + legal suite |
| Integration | Go API + postgres testcontainers | NFT with mocks of core |
| E2E | Playwright: landing → home → tax card | Detox/RN + bank UAT |
| Security | Jailbreak suite smoke | Full red team + ИБ |

## Golden tax eval (минимум)

Хранить `apps/api/testdata/tax_golden.jsonl` (+ тесты в `internal/calc`):

- вопрос, режим, входные цифры, ожидаемая сумма ±1₽, must_include disclaimer  

CI gate: accuracy ≥ 95% через **детерминированный** Go `internal/calc` (без LLM). LLM/SSE — offline fixtures отдельно.

## GenUI contract tests

Для каждого component в catalog — fixture JSON → parse props → render smoke (React Testing Library).

## Manual pitch checklist

См. demo script в [../design/03-screen-map-and-flows.md](../design/03-screen-map-and-flows.md).

## TDD guidance

- Чистые функции (`calculate_tax`, breakeven) — тесты сначала.  
- Agents — eval, не хрупкие snapshot всего промпта.  
- UI — smoke на критический путь, не screenshot-всё.
