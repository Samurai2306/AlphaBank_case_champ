# PRD — Production (Bank Blueprint)

**Статус:** Specification complete (implementation after demo validation)  
**Цель:** AI Business Copilot внутри экосистемы Альфа-Бизнес / Платежный бизнес, совместимый с контуром банка.  
**Клиент:** React Native SuperApp; **AI:** on-prem multi-agent RAG; **Events:** Kafka.

## Vision

Проактивный интерфейс: чат — связующее звено между сложными банковскими микросервисами. ИИ генерирует UI-компоненты и инициирует **черновики** действий; исполнение — только после human confirm.

## Core Must Have (из Developer PRD)

1. **Onboarding Chatbot** — Voice/Text; извлечение ФИО, сферы, оборота, ОКВЭД.  
2. **Мультимодальный ввод** — фото счетов/чеков/договоров, PDF, голос; OCR (Yandex Vision / Tesseract) → LLM.  
3. **Interactive Widgets (Generative UI / SDUI)** — JSON тип+props → нативные компоненты.  
4. **Proactive Alerts** — фон: кассовый разрыв (ML), 115-ФЗ, срок ЕНП.  
5. **Госуслуги/ФНС** — авторизация и подписание УКЭП/УНЭП (Госключ) в чате.

## AI Skills

| Skill | Поведение |
|-------|-----------|
| Legal Scanner | PDF договор → красные флаги (штрафы, разрыв, скрытые комиссии) |
| Tax Calculator | ПСН, УСН 6/15%; API ФНС + БД приходов/расходов → драфт декларации |
| Unit Economics | CAC/LTV/Margin → P&L и график безубыточности |

## Integration в Альфа-Бизнес

- Точка входа: **FAB** или отдельный **Tab**.  
- UI Copilot — более воздушный и «умный», чем основной banking chrome, но на бренд-токенах Альфы.  
- Progressive disclosure по CJM-уровню клиента.

## Phased delivery

### Phase 1 — MVP (мес. 1–3)

- Базовый чат + LLM on-prem.  
- Навыки: тарифы банка, налоги по базе знаний (RAG).  
- Транзакции **read-only**.  
- Без Kafka proactive.

### Phase 2 — Actionable AI (мес. 4–6)

- Function calling tools.  
- Платежные поручения (draft→sign), проверка контрагентов по ИНН.  
- Generative UI карточки.  

### Phase 3 — Proactive CFO (мес. 7–9)

- Kafka consumers.  
- Предиктивный кассовый разрыв.  
- Микро-кредиты/овердрафт в чате (с credit policy gates).

## Non-functional (production)

| Область | Требование |
|---------|------------|
| Privacy | On-prem inference; NER masking; no training on client data |
| Latency | p95 router < 400ms; first token agent < 2s (GPU pool) |
| Availability | AI path degraded mode: FAQ + human handoff |
| Audit | Полный лог tool calls, prompts hash, user confirm IDs |
| Localization | RU; актуальная база НК (версияция корпуса) |
| Compliance | 152-ФЗ, банковская тайна, 115-ФЗ, требования ИБ банка |

## New / improved Payment Business products (обоснование)

| Продукт | Зачем ЦА | Связь с ИИ |
|---------|----------|------------|
| **ЕНП-копилка** | Убрать налоговый шок | Авто-% с поступлений + алерты |
| **Легализация one-tap** | Серый → НПД/ИП | Онбординг-агент + Госключ |
| **Светофор 115-ФЗ до платежа** | Страх блокировок | Compliance tool в чате |
| **Beauty Starter Pack** | Счёт + СБП/терминал + шаблон налогов | Пакет в уровне 4 CJM |
| **Овердрафт по сигналу Copilot** | Аренда/закуп | Proactive CFO + credit API |

## Explicit non-goals (production v1)

- Замена официального налогового консультанта.  
- Полная автоматизация платежей без biometrics.  
- Обучение foundation model на сырых клиентских данных.

## Dependencies

- Core banking APIs, СМЭВ/ФНС, Kafka platform, GPU inference cluster, Alfa design system RN, legal sign-off на дисклеймеры.
