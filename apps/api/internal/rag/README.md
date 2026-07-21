# Demo knowledge corpus (RAG)

Embedded JSONL under `data/`:

| File | Topics |
|------|--------|
| `tax.jsonl` | НПД, УСН, ЕНП/ЕНС, календарь, ПСН, дисклеймеры |
| `banking.jsonl` | РКО/эквайринг Альфа (публичные ориентиры), копилка, кэшбек, демо-персона |
| `faq.jsonl` | FAQ питча, карта законов, продукты банка |
| `compliance.jsonl` | 115-ФЗ, 152-ФЗ, светофор, guard |
| `legal.jsonl` | Риски договора аренды |
| `laws.jsonl` | Дайджесты ГК/НК/ТК/КоАП/УК и ключевых ФЗ со ссылками Consultant |

Retrieval: in-memory lexical ranking (`Index.Retrieve`). Wired via `tools.RetrieveKB` into `GENERAL_QA`, tax, legal, compliance.

All chunks are **demo-labeled digests** — not official Alfa Bank, FNS, or Consultant full-text republication. Prefer links + сжатый смысл for the pitch.
