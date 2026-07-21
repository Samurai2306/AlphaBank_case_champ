# RAG MVP — корпус и проводка в чат

Дата: 2026-07-19 (обновлено: расширенный правовой корпус)

## Что сделано

- Корпус 80+ чанков в `apps/api/internal/rag/data/*.jsonl`
  - налоги, банк, FAQ, 115-ФЗ, legal
  - `laws.jsonl` — ГК (ст. 2, 23, гл. 4), НК (спецрежимы, ЕНП, УСН, календарные статьи), ТК/КоАП/УК (обзор), ФЗ 422/209/129/54/152/149/294/98/63/115/44 и др.
  - banking — публичные ориентиры РКО/эквайринга/налоговой копилки с alfabank.ru
- Пакет `internal/rag` с embed + lexical retrieve
- `retrieve_kb` через `tools.RetrieveKB` → контракт `{id,title,text,source,score}`
- Чат:
  - `GENERAL_QA` — grounded ответ + карточка `KnowledgeSources`
  - `TAX_CALC` / `LEGAL_REVIEW` / `COMPLIANCE` — цитаты из KB
- При `DEMO_OFFLINE=0` FAQ-ответы идут через LLM с CONTEXT из корпуса (цифры налогов по-прежнему из `calc`)

## Важно для питча

- Чанки — **дайджесты со ссылками**, не полные тексты Consultant
- Тарифы Альфа помечены как публичные ориентиры; всегда отсылать на актуальный сайт банка и календарь ФНС

## Модель на VPS

`LLM_BASE_URL=https://text.pollinations.ai/openai`, модель `openai-fast` (бесплатный OpenAI-compatible endpoint, доступен из РФ).
