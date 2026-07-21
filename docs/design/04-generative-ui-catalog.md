# Generative UI Catalog (SDUI)

`schema_version: 1`  
Backend возвращает tool result; клиент рендерит **только** типы из реестра.

## Envelope

```json
{
  "schema_version": 1,
  "component": "TaxCard",
  "props": {}
}
```

Неизвестный `component` → fallback Markdown + log.

---

## TaxCard

```json
{
  "component": "TaxCard",
  "props": {
    "regime": "NPD",
    "period_label": "Июль 2026",
    "tax_amount": 10800,
    "currency": "RUB",
    "rate_label": "6%",
    "income": 180000,
    "disclaimer": "ИИ-помощник…",
    "source_ref": "NK_RF:demo-model",
    "cta": { "action": "create_payment_draft", "label": "Сформировать платёжку" }
  }
}
```

## PaymentDraftCard

```json
{
  "component": "PaymentDraftCard",
  "props": {
    "draft_id": "uuid",
    "amount": 10800,
    "currency": "RUB",
    "purpose": "Налог НПД за июль 2026",
    "payee_name": "УФК … (мок)",
    "status": "draft",
    "risk_level": "green",
    "cta_confirm_label": "Подтвердить",
    "cta_cancel_label": "Отменить"
  }
}
```

## UnitEconomicsChart

```json
{
  "component": "UnitEconomicsChart",
  "props": {
    "breakeven_units_per_day": 12,
    "price_per_unit": 1500,
    "variable_cost_per_unit": 200,
    "fixed_costs_monthly": 40000,
    "margin_per_unit": 1300,
    "series": [{ "units": 0, "profit": -40000 }, { "units": 12, "profit": 0 }]
  }
}
```

## ComplianceTrafficLight

```json
{
  "component": "ComplianceTrafficLight",
  "props": {
    "inn": "7707083893",
    "level": "yellow",
    "title": "Есть замечания",
    "reasons": ["Недавно сменён директор (мок)"],
    "recommendation": "Запросите закрывающие документы до оплаты"
  }
}
```

## OnboardingSummary

```json
{
  "component": "OnboardingSummary",
  "props": {
    "display_name": "Маша",
    "business_sphere": "Маникюр",
    "city": "Казань",
    "monthly_revenue_estimate": 180000,
    "suggested_regimes": ["NPD", "USN_6"],
    "suggested_okved": ["96.02"],
    "cta_label": "Всё верно"
  }
}
```

## EnpPiggyBank

```json
{
  "component": "EnpPiggyBank",
  "props": {
    "enabled": true,
    "rate_percent": 6,
    "balance": 5400,
    "currency": "RUB",
    "last_contribution": 600
  }
}
```

## LegalFlagsList

```json
{
  "component": "LegalFlagsList",
  "props": {
    "document_name": "dogovor-arenda.pdf",
    "flags": [
      {
        "severity": "high",
        "title": "Штраф за досрочный выход",
        "quote": "…штраф в размере трёх месячных платежей…",
        "recommendation": "Согласуйте лимит штрафа до 1 месяца"
      }
    ]
  }
}
```

## Client rendering rules

1. Map `toolName` → preferred component (server may also send envelope).  
2. Money amounts: tabular nums, `ru-RU` formatting.  
3. CTA вызывает client action → API, не «выполняет» LLM.  
4. Accessibility: card = region with aria-label.
