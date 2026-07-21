# AI Business Copilot — Demo MVP

Питч-ready демо для кейса Альфа-Банка: **Next.js** + **Go API**, offline agents, Generative UI, персона «Маша».

## Быстрый старт (<15 минут)

По умолчанию `DEMO_OFFLINE=1` — LLM-ключ не нужен. Цифры налогов всегда из `internal/calc`.

### 1. API (Go)

```bash
cd apps/api
set DEMO_OFFLINE=1
go run ./cmd/api
```

Health: http://localhost:8080/api/v1/health  
Token: `demo-masha-token`

### 2. Web

```bash
cd apps/web
npm install
npm run dev
```

Открыть: http://localhost:3000

### Сценарий жюри

- http://localhost:3000/demo — 6 промптов по порядку  
- http://localhost:3000/jury — краткий брифинг для жюри  

Приложение: Home · Операции · Сервисы · Chat · Профиль (+ копилка, алерты, путь).

## Переменные

См. [.env.example](.env.example).

| Режим | Как |
|-------|-----|
| Offline (питч) | `DEMO_OFFLINE=1` (default) — fixture/orchestrator |
| Cloud narrative | `DEMO_OFFLINE=0` + `LLM_API_KEY` — OpenAI-compatible; SDUI/calc без изменений |

## Docker (удобный режим)

```bash
docker compose up -d --build
```

| Сервис | URL |
|--------|-----|
| Web | http://localhost:3000 |
| API health | http://localhost:8080/api/v1/health |
| Сценарий | http://localhost:3000/demo |
| Жюри | http://localhost:3000/jury |

Остановить: `docker compose down`  
Логи: `docker compose logs -f`

## Тесты

```bash
cd apps/api && go test ./...
cd apps/web && npx playwright install chromium && npm run test:e2e
```

E2E ожидает API на `:8080`. Home показывает **10 800 ₽** (запас 6%); в чате TaxCard — **~7 740 ₽** (смешанная НПД 4%/6%).

Debug API: `powershell -File scripts/debug-demo-session.ps1`

## Deploy (VPS later)

Заготовки: [deploy/](deploy/) — `docker-compose.prod.yml`, nginx, scripts.  
Runbook: [docs/engineering/07-server-deployment.md](docs/engineering/07-server-deployment.md).  
Нужны: SSH, домен, `DEMO_TOKEN`, опционально `LLM_API_KEY`. Секреты не коммитить.

## Документация

Стартовая точка: [docs/00-README.md](docs/00-README.md)

- Backend Go: [docs/architecture/09-demo-backend-go.md](docs/architecture/09-demo-backend-go.md)
- PRD демо: [docs/product/01-prd-demo-mvp.md](docs/product/01-prd-demo-mvp.md)
- 3D icons: [docs/design/assets/3d-icons/](docs/design/assets/3d-icons/)
