# YA registration & empty cabinet (18–25)

Date: 2026-07-19  
Status: approved

## Goals

- Primary path: register a new young entrepreneur (18–25), enter data manually, open an empty cabinet that still works (tax from revenue estimate).
- Secondary: “Посмотреть пример” → seeded Маша.
- Tone: short, mobile-first, no bank jargon; empty-states guide what to fill in.

## Decisions

| Topic | Choice |
|-------|--------|
| Auth | Memory multi-user; `POST /register`, `POST /login` by phone; Bearer token per user |
| Wizard | 3 steps: who → business → regime |
| Cabinet data | Empty (0 txns, piggy off); tax from `monthly_revenue_estimate` |
| Persistence | In-memory only (no Postgres/SMS) |
| Маша | Seed user + example CTA |

## API

- `POST /api/v1/register` — create user, return `{ token, profile }`
- `POST /api/v1/login` — `{ phone }` → `{ token, profile }`
- Auth middleware resolves `token → user_id`; handlers use user-scoped store
- New user: `cjm_level=1`, `persona_key=custom`, empty drafts/txns

## Web

- `/register` wizard; `/login` phone + example; landing primary → register
- `localStorage.copilot_token`; app layout redirects if missing
- EmptyHint on home / ops / piggy / payments / journey / tasks

## Out of scope

Real SMS, KYC, Postgres, cloning Masha’s full ledger onto new users.
