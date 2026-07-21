# Security & Guardrails

## Threat model (кратко)

| Угроза | Контроль |
|--------|----------|
| Утечка PII в LLM | NER masking → токены `[ACCOUNT_1]`, `[PHONE_1]` |
| Jailbreak / уклонение от налогов | Ingress classifier (Llama Guard) + policy replies |
| Ложный налоговый совет | RAG + citation + low-confidence refusal |
| Несанкционированный платёж | Strict output parser ≠ LLM; confirm SMS/bio |
| Prompt injection via PDF | Sanitize OCR text; tool allowlist |
| Model supply chain | On-prem images, signed artifacts |

## PII Masking

**До** вызова LLM (даже on-prem):

1. NER: ФИО, паспорт, телефон, счёт, карта, адрес.  
2. Замена на стабильные токены в рамках сессии.  
3. После ответа — detokenize только в безопасных полях UI (не в логах промптов).

Demo: можно упростить regex-masker на seed-данных.

## Strict Output Parser (money)

Любой tool/`create_payment_draft` / transfer intent:

1. JSON schema validation (Zod/Pydantic) — **не** «доверься модели».  
2. Amount/payee whitelist checks.  
3. `CONFIRM_REQUIRED` пока нет второго фактора.  
4. Audit event: who/when/what draft_id.

## Anti-Jailbreak

- Блокировать: «забудь инструкции», «как не платить налог», обход 115-ФЗ.  
- Ответ: отказ + легальные альтернативы (НПД, вычет, рассрочка пени — только из KB).  

## Human-in-the-loop

| Действие | HITL |
|----------|------|
| Q&A налоги | Дисклеймер |
| Payment draft | Confirm |
| Sign payment | SMS/Push/FaceID |
| Credit offer accept | Credit policy + confirm |
| Legal conclusions | Дисклеймер «не юр. консультация» |

## Diskлеймер (обязательный copy)

> ИИ-помощник Альфа-Банка. Рекомендации носят справочный характер и не заменяют консультацию специалиста или официальные разъяснения ФНС. Финальное решение — за вами.

## Audit

Логировать: user_id hash, intent, tool name, args redacted, latency, model id, corpus version, confirm id.  
Не логировать: сырые промпты с PII, полные тексты документов без необходимости.

## Demo specifics

- Cloud LLM допустим с синтетическими данными.  
- Не использовать реальные паспорта/счета в seed.  
- Confirm не вызывает внешний money API.
