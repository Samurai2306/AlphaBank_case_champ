# API Contracts — Tools & HTTP

Контракты **общие** для demo и production. Менять только вместе с `packages/shared` и GenUI catalog.

## HTTP (Copilot BFF / Demo API)

| Method | Path | Описание |
|--------|------|----------|
| POST | `/api/v1/chat` | Streaming chat (AI SDK / SSE) |
| GET | `/api/v1/me` | Профиль + CJM level |
| PATCH | `/api/v1/me/onboarding` | Обновление онбординга |
| GET | `/api/v1/transactions` | Список транзакций (read) |
| POST | `/api/v1/payments/draft` | Создать черновик |
| POST | `/api/v1/payments/draft/{id}/confirm` | Confirm (demo mock / prod sign) |
| POST | `/api/v1/documents/legal` | Upload PDF |
| GET | `/api/v1/health` | Health |

Auth: `Authorization: Bearer <jwt>` (demo: mock JWT).

## LLM Tools (function calling)

### `get_client_transactions`

```
GET /api/v1/core/transactions
```

| Param | Type | Desc |
|-------|------|------|
| from | ISO date | начало |
| to | ISO date | конец |

**Returns:** `{ items: [{ id, amount, currency, direction, counterparty, booked_at, category }] }`  
**LLM use:** P&L, налоги.

### `check_counterparty_risk`

```
POST /api/v1/compliance/115fz
```

| Param | Type |
|-------|------|
| inn | string |

**Returns:** `{ inn, level: "green"|"yellow"|"red", reasons: string[] }`  
**SDUI:** `ComplianceTrafficLight`

### `create_payment_draft`

```
POST /api/v1/payments/draft
```

| Param | Type |
|-------|------|
| amount | number |
| currency | "RUB" |
| purpose | string |
| payee | object |

**Returns:** `{ draft_id, status: "draft", ui: PaymentDraftCardProps }`  
**SDUI:** `PaymentDraftCard`

### `fetch_fns_debt`

```
GET /api/v1/fns/enp_balance
```

**Returns:** `{ balance, updated_at, currency: "RUB" }`

### `calculate_tax` (shared helper)

| Param | Type |
|-------|------|
| regime | `NPD` \| `USN_6` \| `USN_15` \| `PSN` |
| income | number |
| expenses | number \| null |
| period | string |

**Returns:** `{ tax_amount, rate, breakdown, disclaimer }`  
**SDUI:** `TaxCard`

### `compute_unit_economics`

| Param | Type |
|-------|------|
| fixed_costs | number |
| variable_cost_per_unit | number |
| price_per_unit | number |

**Returns:** `{ breakeven_units_per_day, margin, chart }`  
**SDUI:** `UnitEconomicsChart`

### `scan_legal_document`

| Param | Type |
|-------|------|
| document_id | string |

**Returns:** `{ flags: [{ severity, title, quote, recommendation }] }`  
**SDUI:** `LegalFlagsList`

### `update_onboarding_profile`

| Param | Type |
|-------|------|
| business_sphere | string |
| city | string |
| monthly_revenue_estimate | number |
| tax_regime_preference | string \| null |

**SDUI:** `OnboardingSummary`

### `retrieve_kb`

| Param | Type |
|-------|------|
| query | string |
| top_k | number |

**Returns:** `{ chunks: [{ id, text, source, score }] }`

## Error model

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human readable",
    "retryable": false
  }
}
```

Коды: `UNAUTHORIZED`, `VALIDATION_ERROR`, `RAG_LOW_CONFIDENCE`, `GUARDRAIL_BLOCKED`, `UPSTREAM_UNAVAILABLE`, `CONFIRM_REQUIRED`.

## Versioning

- URL `/api/v1`  
- SDUI `schema_version: 1`  
- Breaking changes → `/v2` + migration note в этом файле.
