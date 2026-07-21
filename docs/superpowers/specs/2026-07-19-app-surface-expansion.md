# App surface expansion (quality-only)

## Missing → added (no decorative junk)

| Surface | Why |
|---------|-----|
| Bottom nav (`AppNav`) | Alfa-like app shell; Chat hides nav for composer UX |
| `/app/operations` | Real seed transactions + income/expense/net |
| `/app/services` | Product grid → working deep-links only |
| `/app/piggy` | PATCH piggy rate/enable; gap from home |
| `/app/profile` | Editable persona feeding calc |
| `/app/alerts` | Proactive feed with CTA |
| `/jury` | Stakeholder briefing, not marketing fluff |
| Chat mic | Web Speech API (D12), fallback message |
| 404 | Clear recovery |

## API

- `PATCH /me`, `PATCH /piggy`, `GET /alerts`
- `GET /transactions` + `summary`

## Out of this pass (intentionally)

- Fake KPI dashboards, social feed, settings with 40 toggles
- Real bank/FNS connectors
