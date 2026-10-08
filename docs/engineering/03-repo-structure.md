# Repo Structure (target)

Создаётся при старте implementation plan. Сейчас в репозитории — документация.

```text
AlphaBank_case_champ/
├── AGENTS.md
├── README.md
├── docs/
├── .cursor/rules/
├── apps/
│   ├── web/                       # Next.js demo
│   │   ├── app/
│   │   ├── components/
│   │   │   ├── chat/
│   │   │   ├── home/              # HeroMetric, insights
│   │   │   └── genui/
│   │   ├── lib/api/               # SSE client → Go
│   │   └── package.json
│   ├── api/                       # Go backend (module)
│   │   ├── cmd/api/
│   │   ├── internal/
│   │   │   ├── agents/
│   │   │   ├── adapters/
│   │   │   ├── calc/
│   │   │   ├── httpserver/
│   │   │   ├── tools/
│   │   │   └── ...
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── mobile/                    # RN later (prod track)
├── packages/
│   └── shared/                    # Zod/JSON: intents, SDUI, tools
├── deploy/
│   ├── docker-compose.prod.yml
│   ├── nginx/copilot.conf
│   └── scripts/
│       ├── remote-bootstrap.sh
│       └── remote-deploy.sh
├── docker-compose.yml             # local
├── Makefile
└── .env.example
```

Подробности Go: `docs/architecture/09-demo-backend-go.md`.

## Ownership

| Path | Owner concern |
|------|----------------|
| `packages/shared` | Contract changes → sync Go structs + docs |
| `apps/api/internal/adapters` | Mocks → real ports later |
| `apps/api/internal/calc` | Deterministic money math |
| `deploy/` | VPS-only compose & proxy |
| `docs/` | Product/architecture truth |

## Branching (suggested)

- `main` — stable docs + demo  
- `feat/*` — features  
