# Tech Stack — Demo Foundation

## Frontend

| Tech | Version target | Role |
|------|----------------|------|
| Next.js | 15 App Router | Web PWA / pitch |
| React | 19 | UI |
| TypeScript | 5.x | Strict |
| Tailwind CSS | 4.x / 3.4+ | Styling |
| Framer Motion | latest | Motion |
| Fonts | Alfa Interface Sans (CDN) or Manrope | Visual parity with Альфа-Бизнес |
| Zod | 3.x | Shared schemas with API |
| Recharts | latest | Unit economics |
| Playwright | latest | E2E smoke |

Chat client: native **SSE reader** к Go API (см. `docs/architecture/09-demo-backend-go.md`). Vercel AI SDK — опционально для UI helpers, не как единственный протокол к бэкенду.

## Backend (Go — единственный API)

| Tech | Role |
|------|------|
| Go 1.22+ | Runtime |
| chi (или echo) | HTTP router |
| pgx | Postgres driver |
| goose | SQL migrations |
| redis (go-redis) | Rate limit / session |
| slog | Structured JSON logs |
| OpenAI-compatible HTTP client | LLM (OpenAI / OpenRouter / later vLLM) |
| testcontainers / httptest | Tests |

Детали пакетов и агентов: [../architecture/09-demo-backend-go.md](../architecture/09-demo-backend-go.md).

## Data & AI

| Tech | Role |
|------|------|
| Postgres 16 + pgvector | Primary + embeddings |
| Redis 7 | Limits |
| Cloud LLM via env | Chat; `DEMO_OFFLINE=1` fixtures |
| Provider abstraction in Go | Swap base URL → on-prem later |

## Tooling

- pnpm/npm workspaces (web + shared)  
- Go modules (`apps/api`)  
- Docker Compose (local + `deploy/` prod)  
- Makefile: `up`, `seed`, `test`, `test-api`  
- ESLint / Vitest (web); `go test -race` (api)

## Explicit non-choices (demo)

- Не FastAPI / не Python в runtime демо  
- Не Kafka  
- Не React Native в Sprint A  
- Не обязательный GPU  
- Не публикация Postgres/Redis портов на VPS наружу  

## Server

Nginx или Caddy, TLS, DNS — [07-server-deployment.md](07-server-deployment.md).
