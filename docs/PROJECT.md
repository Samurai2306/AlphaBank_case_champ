# Альфа-Бизнес: Старт — рабочая документация

Практический гайд по **демо-продукту** в этом репозитории. Стратегия, PRD и production-blueprint лежат отдельно: [00-README.md](00-README.md).

Продуктовое имя в UI: **«Альфа-Бизнес: Старт»**. Это не официальная консультация банка или ФНС.

---

## 1. Что это

Проактивный кабинет для молодых предпринимателей (фокус — beauty / НПД). Не «чат с GPT»: чат связан с банковскими действиями через карточки (SDUI).

Ценность в демо:

- налог НПД/УСН → копилка ЕНП → черновик платёжки;
- проверка ИНН (светофор 115-ФЗ) → платёжка аренды/контрагенту;
- безубыточность → офферы СБП/эквайринга;
- выписка в «Операциях» → следующие шаги по РКО.

Персона для жюри: **Маша**, Казань, маникюр, НПД. Токен: `demo-masha-token`.

---

## 2. Стек демо

| Слой | Технология |
|------|------------|
| Web | Next.js 15, TypeScript, Tailwind, PWA |
| API | один сервис на **Go** (агенты, tools, SSE, моки) |
| Данные | in-memory store (профиль, операции, черновики). Поиск по корпусу — лексический индекс в процессе |
| LLM | OpenAI-compatible, с failover моделей; без ключа — `DEMO_OFFLINE=1` |
| RAG | ~90 чанков в `apps/api/internal/rag/data/` |

Python runtime в демо **не поднимать**.

---

## 3. Локальный запуск

Скопируйте `.env.example` → `.env`. По умолчанию `DEMO_OFFLINE=1` — нейросеть не нужна, калькулятор налогов работает.

### Docker

```bash
docker compose up -d --build
```

| Сервис | URL |
|--------|-----|
| Web | http://localhost:3000 |
| API health | http://localhost:8080/api/v1/health |
| Сценарий жюри | http://localhost:3000/demo |
| Брифинг | http://localhost:3000/jury |

Войти как Маша: логин на `/login` или токен `demo-masha-token`.

### Без Docker

```bash
# API
cd apps/api
set DEMO_OFFLINE=1
go run ./cmd/api

# Web (другой терминал)
cd apps/web
npm install
npm run dev
```

---

## 4. Ключевые сценарии (чат)

| Запрос | Intent | Карточка |
|--------|--------|----------|
| Сколько отложить на налог… | `TAX_CALC` | TaxCard |
| Сформируй платёжку | `TRANSACTION` | PaymentDraftCard (УФК / ЕНП) |
| Сформируй платёжку аренды 40000 на ИНН … | `TRANSACTION` | платёжка контрагента + светофор |
| Проверь ИНН … | `COMPLIANCE` | ComplianceTrafficLight |
| Аренда 40к, расходники 200, цена 1500 | `UNIT_ECON` | UnitEconomicsChart + ProductOffers |
| Включи копилку 6% | `PIGGY` | EnpPiggyBank |
| Разбери договор PDF | `LEGAL_REVIEW` | LegalFlagsList |

Реальные деньги не уходят: confirm в демо только меняет статус черновика.

---

## 5. LLM

Налоги и суммы **всегда** из `internal/calc`. LLM пишет текст и классифицирует intent.

| Переменная | Назначение |
|------------|------------|
| `DEMO_OFFLINE` | `1` — без облака; `0` — живой LLM |
| `LLM_BASE_URL` / `LLM_API_KEY` / `LLM_MODEL` | основной OpenAI-compatible endpoint |
| `LLM_PROVIDER` | `opencode-zen` или пусто (обычный OpenAI-стиль) |
| `LLM_MODELS` | список моделей через запятую — ротация при сбое |
| `LLM_FALLBACK_*` | второй провайдер (например DeepSeek) |

Проверка: `GET /api/v1/ready` и `GET /api/v1/ready?probe=1`.

Рабочий контур: OpenCode Zen, модели `big-pickle` → `mimo-v2.5-free` → `nemotron-3.5-lightning-free` → `ling-3.0-flash-fin-free`. На Vercel этот контур включается сам (`DEMO_OFFLINE=0`), локально по умолчанию офлайн. Старые id (`deepseek-v4-flash-free`, `minimax-m2.5-free`) больше не принимаются Zen.

---

## 6. Тесты

```bash
cd apps/api && go test ./...
cd apps/web && npx tsc --noEmit
# e2e (нужен API на :8080)
cd apps/web && npx playwright install chromium && npm run test:e2e
```

Golden-диалоги: `apps/api/internal/eval/golden_test.go` (offline).

---

## 7. Структура репо

```text
apps/web          Next.js UI
apps/api          Go API + agents + RAG + calc
docs/             стратегия, PRD, архитектура, дизайн
docker-compose.yml
docker-compose.vps.yml   bind 127.0.0.1:3100/8100
deploy/           nginx/compose заготовки
```

Контракты demo↔prod (intents, tools, SDUI) не ломать без обновления:

- [architecture/05-api-contracts.md](architecture/05-api-contracts.md)
- [design/04-generative-ui-catalog.md](design/04-generative-ui-catalog.md)

---

## 8. Деплой и соседство с B.O.T.-Project

Этот демо **временно** жил на `https://bot-project.ru` (VPS `155.212.170.159`), заняв порты прежнего продукта B.O.T.-Project (`127.0.0.1:3100` web, `:8100` api).

Полный обратимый план, бэкапы и откат:

- [ops/vps-bot-project-swap-reversible.md](ops/vps-bot-project-swap-reversible.md)
- [engineering/07-server-deployment.md](engineering/07-server-deployment.md)

Каталог AlphaBank на сервере: `/opt/AlphaBank_case_champ` (дерево можно оставить).  
Прежний стек: `/opt/B.O.T.-Project`, systemd `bot-project.service`, volumes `BOT_project_pgdata`.

Секреты (`.env`, пароли SSH) в git не класть.

---

## 9. Ограничения демо

- Не выдавать ответы за официальную консультацию банка/ФНС.
- Не хранить реальные паспортные/банковские данные.
- Не обходить confirm на финансовых действиях.
- Jailbreak / уклонение от налогов — блок на входе (`internal/guard`).
