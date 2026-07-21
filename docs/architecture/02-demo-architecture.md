# Demo Architecture

## Обзор

Монорепо: **Next.js** (web) + **Go API** (единственный backend) + Postgres/pgvector + Redis.  
Банк и ФНС — in-process mocks с теми же tool-именами, что в production.  
Детальная спецификация Go: [09-demo-backend-go.md](09-demo-backend-go.md).  
Деплой на VPS: [../engineering/07-server-deployment.md](../engineering/07-server-deployment.md).

```mermaid
flowchart LR
  Browser[NextjsPWA]
  API[Go_copilot_api]
  Agents[Go_Agents_Router]
  PG[(Postgres_pgvector)]
  Redis[(Redis)]
  LLM[CloudLLM_OpenAI_compatible]
  Mocks[MockAdapters]

  Browser -->|SSE_plus_REST| API
  API --> Agents
  Agents --> LLM
  Agents --> PG
  Agents --> Mocks
  API --> Redis
```

## Компоненты

| Компонент | Ответственность |
|-----------|-----------------|
| `apps/web` | Landing, Home dashboard, chat GenUI, demo script |
| `apps/api` | **Go**: HTTP, SSE chat, agents, tools, mocks, RAG |
| `packages/shared` | Zod/JSON schemas: intents, tools, SDUI (контракт с фронтом) |
| `deploy/` | prod compose, Nginx/Caddy, remote scripts |
| pgvector | Чанки НК (урезанный корпус), FAQ тарифов |
| MockAdapters | transactions, 115-ФЗ, FNS debt, payment draft |

## Chat flow

1. Client → `POST /api/v1/chat` (SSE).  
2. Guard (regex + policy).  
3. Router agent → JSON intent.  
4. Specialist ± tools; цифры налогов из `internal/calc` (детерминированно).  
5. Stream: `token` / `sdui` / `done` events.  
6. Client maps SDUI → React components.

## Data stores (demo)

- `users`, `profiles`, `conversations`, `messages`
- `transactions` (seed)
- `payment_drafts`
- `documents` (legal uploads metadata)
- `kb_chunks` (embeddings)

## Config

См. полный список в [09-demo-backend-go.md](09-demo-backend-go.md):  
`LLM_*`, `DATABASE_URL`, `REDIS_URL`, `DEMO_TOKEN`, `CORS_ORIGINS`, `DEMO_OFFLINE`, …

## Deploy demo

| Env | How |
|-----|-----|
| Local | `docker compose up` — api + web + postgres + redis |
| Server | Nginx/Caddy + TLS + DNS — [07-server-deployment.md](../engineering/07-server-deployment.md) |

GPU не нужен. Один бинарь Go + Node web.

## Соответствие production

| Demo | Production |
|------|------------|
| Go monolith (API+agents) | Go BFF + Python AI swarm |
| mock adapters | Java/Go core adapters |
| Cloud LLM via OpenAI-compatible URL | on-prem vLLM (тот же client interface) |
| SSE | сохранить или эволюционировать |

Контракты tools/SDUI/intents — общие ([05-api-contracts.md](05-api-contracts.md)).
