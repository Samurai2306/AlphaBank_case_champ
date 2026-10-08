# AGENTS.md — AI Business Copilot

## Что это

Гибридный проект: качественный **demo MVP** (Next.js + **Go API**) + **production blueprint** (React Native, микросервисы, Kafka, on-prem LLM). Документация — Single Source of Truth.

## Обязательный порядок работы

1. Прочитать `docs/00-README.md` и релевантный раздел `docs/`.
2. Перед новым функционалом — `superpowers:brainstorming` → spec в `docs/superpowers/specs/`.
3. Перед кодом — план в `docs/superpowers/plans/` (`writing-plans`).
4. Реализация — `subagent-driven-development` или `executing-plans` + TDD где уместно.
5. Перед «готово» — `verification-before-completion`. После правок UI — `deslop`.

## Жёсткие правила

- Demo и production используют **одни и те же** contracts: tools, SDUI types, intents (`TAX_CALC`, `LEGAL_REVIEW`, `ONBOARDING`, `GENERAL_QA`, `TRANSACTION`).
- Демо: можно мокать банк/ФНС и cloud LLM. Нельзя имитировать реальный перевод денег без UI confirm (draft only).
- Production docs: on-prem LLM, PII masking (NER), human-in-the-loop на платежах, без репутационных/юр. рисков.
- Налоговые/юр. ответы — только через RAG + цитирование + дисклеймер.
- Не коммитить секреты (`.env`, ключи API).

## Стек (locked)

| Demo | Production |
|------|------------|
| Next.js 15, TS, Tailwind, SSE→Go | React Native (TS), Zustand, React Query |
| **Go** API + agents + mocks | Go BFF, Python AI, Java/Spring core |
| Память процесса, лексический корпус | Kafka, Qdrant/Milvus, vLLM on-prem |

Детали: `docs/engineering/01-tech-stack-demo.md`, `docs/architecture/09-demo-backend-go.md`.  
Деплой на VPS (когда будет SSH/домен): `docs/engineering/07-server-deployment.md`.

## ЦА и ценность

Молодые предприниматели 17–25 (расш. 18–35), beauty. Ценность: легализация + налоги + платежи («CFO в кармане»), не «чат с GPT».

## Дизайн (обязательно прочитать перед UI)

- `docs/design/06-alfa-visual-language.md` — как выглядит Альфа / Альфа-Бизнес  
- `docs/design/07-dashboard-patterns.md` — decision-first Home, не KPI soup  
- `docs/design/08-alfa-marketing-ui-patterns.md` + `docs/design/references/alfa-web-2026/` — актуальные скрины  
- Landing = pastel/bento/3D; Home ≈ «Сервисы»; chat/money = спокойный trust + `#EF3124`  
- Radius карточек 24–32px; segment active = тёмная pill

## Правила Cursor

- `.cursor/rules/copilot-project.mdc`
- `.cursor/rules/demo-vs-production.mdc`
