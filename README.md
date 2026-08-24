# Альфа-Бизнес: Старт

Демо для кейса Альфа-Банка: кабинет молодого предпринимателя (beauty / НПД) с чатом, налогом, платёжками и Generative UI.

**Стек:** Next.js 15 + один Go API (агенты, RAG, SSE). Python в runtime демо нет.

Практический гайд: **[docs/PROJECT.md](docs/PROJECT.md)**  
Вся документация: **[docs/00-README.md](docs/00-README.md)**

---

## Быстрый старт

По умолчанию `DEMO_OFFLINE=1` — LLM-ключ не нужен. Цифры налогов всегда из `apps/api/internal/calc`.

```bash
cp .env.example .env
docker compose up -d --build
```

| Что | URL |
|-----|-----|
| Приложение | http://localhost:3000 |
| API health | http://localhost:8080/api/v1/health |
| Сценарий жюри | http://localhost:3000/demo |
| Брифинг | http://localhost:3000/jury |

Токен Маши: `demo-masha-token`.

Без Docker:

```bash
cd apps/api && set DEMO_OFFLINE=1 && go run ./cmd/api
cd apps/web && npm install && npm run dev
```

---

## Что умеет демо

- Налог НПД/УСН, копилка ЕНП, черновик платёжки (без реального списания)
- Светофор 115-ФЗ по ИНН → платёжка аренды/контрагенту
- Безубыточность + подсказки СБП/эквайринга
- Выписка в «Операциях», разбор договора PDF, RAG по НК/банковским материалам

Это **не** официальная консультация банка или ФНС.

---

## Тесты

```bash
cd apps/api && go test ./...
cd apps/web && npx tsc --noEmit
```

E2E: API на `:8080`, затем `cd apps/web && npm run test:e2e`.

---

## LLM (опционально)

`DEMO_OFFLINE=0` + OpenAI-compatible endpoint. С RU VPS удобен OpenCode Zen с ротацией моделей (`LLM_MODELS`). Второй провайдер — `LLM_FALLBACK_*`. См. `.env.example` и [docs/PROJECT.md](docs/PROJECT.md).

---

## VPS / откат домена

Заготовки: [deploy/](deploy/), overlay `docker-compose.vps.yml`.  
Runbook: [docs/engineering/07-server-deployment.md](docs/engineering/07-server-deployment.md).  
Соседство с прежним **B.O.T.-Project** на `bot-project.ru`: [docs/ops/vps-bot-project-swap-reversible.md](docs/ops/vps-bot-project-swap-reversible.md).

Секреты в git не коммитить.
