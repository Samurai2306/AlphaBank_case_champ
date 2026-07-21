# Sources & Research Notes

## Primary documents (in-repo)

| File | Description |
|------|-------------|
| [source-pdfs/AlfaBank_Copilot_Developer_PRD.pdf](source-pdfs/AlfaBank_Copilot_Developer_PRD.pdf) | Master Blueprint & PRD: functions, UI, stack, prompts, tools, guardrails, roadmap |
| [source-pdfs/AlfaBank_Business_Copilot_Master_Strategy.pdf](source-pdfs/AlfaBank_Business_Copilot_Master_Strategy.pdf) | Strategy 2026+: TAM, CJM×10, agents, SWOT, GTM, monetization, risks |

## Championship case (screenshots / brief)

Ключевые тезисы, использованные в `docs/strategy/*` и `docs/product/*`:

- Контекст: Платежный бизнес Альфы; ЦА 17–25; средний клиент ~42; фрагментированный путь новичка.
- Образ результата: фокус 17–25; 2–3 отрасли (HoReCa, beauty, retail…); боли; осмысленный ИИ; продукты ПБ в пути; пилот РФ.
- Критерии успеха: реальное понимание ЦА; польза банка; long-term стратегия; ИИ на болях; интеграция ПБ; новые продукты; ограничения банка.
- Ограничения: внешний AI на идее ок → финал совместим с Альфой; без репутационных/юр. рисков.
- Постановка: 8 шагов (ЦА → отрасли → конкуренты → инсайты → стратегия+AI → промо → продукты → эффект).
- Критерии оценки: качество исследования ЦА/рынка; структура и логика; роль ИИ в стратегии.
- Исследование ЦА (Telegram): beauty masters; Value Proposition Canvas; боли налогов/легализации; ожидание all-in-one + сильный UX; конкурент Т-Банк.
- Команда: AI-консультант по налогам/документам; направления beauty-медицина, retail, beauty.

## Web research (patterns, July 2026 session)

Использовано для архитектуры GenUI/RAG и UX-принципов (не для юридических норм РФ):

- [AI SDK UI: Generative User Interfaces](https://ai-sdk.dev/v7/docs/ai-sdk-ui/generative-user-interfaces) — tools → React components.
- [Building production RAG with Next.js, LangChain, Vercel AI SDK](https://www.reactify-solutions.com/articles/rag-nextjs-langchain-vercel-ai-sdk) — разделение ingestion vs streaming UI.
- [Tool calling / multi-step (AI SDK)](https://github.com/vercel/ai/blob/08cdf6ae/content/docs/03-ai-sdk-core/15-tools-and-tool-calling.mdx)
- [GenUI in banking account opening](https://verygood.ventures/blog/genui-last-mile-gap-banking-account-openings/) — intent → assembled UI under bank allowlist.
- [Fintech UI/UX 2026 — The Skins Factory](https://www.theskinsfactory.com/uiux-design-blog/fintech-uiux-design) — trust, intentional friction.
- [Fintech UX Guide 2026 — Fuselab](https://fuselabcreative.com/fintech-ux-design-guide-2026-user-experience/) — proactive insights.
- Multi-agent RAG references: LangGraph supervisor patterns (e.g. production RAG + FastAPI + Next.js open templates) — отражено в `docs/architecture/04-agentic-rag.md`.

## Alfa design system & brand (July 2026 session)

Отражено в `docs/design/06-alfa-visual-language.md` и обновлениях `01-design-system.md`:

- [Корпоративные веб-стандарты alfabank.ru (Art. Lebedev PDF)](https://www.artlebedev.com/alfabank/guides/Alfa-guides.pdf) — Альфа-красный `#EF3124`, платиновый `#D1D5D8`, правила логотипа.
- [alfa-laboratory/core-components](https://github.com/alfa-laboratory/core-components) — актуальная UI-библиотека; [Storybook](https://alfa-laboratory.github.io/core-components/).
- [arui-feather](https://github.com/alfa-laboratory/arui-feather) — deprecated Feather; не использовать для нового кода.
- [alfa-ui-primitives](https://github.com/alfa-laboratory/alfa-ui-primitives) — design tokens (web/iOS/Android).
- [Habr: иконки и Alfa Interface Sans в Альфа-Бизнес](https://habr.com/ru/companies/alfa/articles/1048062/) — лёгкость, воздух, фирменный UI-шрифт.
- [Habr: боковое меню Альфа-Бизнес](https://habr.com/ru/companies/alfa/articles/809991/) — SharedUI, белые карточки, F-паттерн.
- [core-components: Alfa Interface Sans](https://github.com/core-ds/core-components/commit/b832d8507355178a2b2d3fa24220be6e1f9d8ee3) — подключение шрифта / Typography.
- [Feather by Alfa-Bank (catalog)](https://www.designsystemcookbooks.com/system/alfa-bank) — историческое имя DS «Feather».

## Modern fintech dashboards (July 2026 session)

Отражено в `docs/design/07-dashboard-patterns.md`:

- [Fintech Dashboard Design: Patterns & Real Examples (2026) — Masterly](https://www.themasterly.com/blog/fintech-dashboard-design-guide) — Role–Metric–Density–Action; intentional friction; trust structural.
- [How to Build a Financial Dashboard in 2026 — Outcrowd](https://medium.com/outcrowd/how-to-decide-what-your-financial-dashboard-actually-needs-in-2026-b4bddf5bc0ba) — AI в workflow, explainability, hybrid ask+structure.
- [UX Design for Dashboards best practices](https://webdesignerindia.medium.com/ux-design-for-dashboards-best-practices-663767ff19c1) — decision-first, progressive disclosure, AI citations, &lt;1.5s FMC.
- [Future of Financial UX: Proactive Copilots](https://craftinnovations.global/future-of-financial-ux-webinar/) — intent-driven, features find the user, time-to-task.
- [Fintech Dashboard Design Patterns 2025 — SaaS Tree](https://thesaastree.com/fintech-dashboard-design-patterns-2025-the-new-rules/) — trust, adaptive, AI-native.

**Note:** `parallel-cli` в среде разработки не был установлен; поиск выполнен через доступный web search. При необходимости углубления — установить Parallel CLI и расширить этот файл.

## How sources map to docs

| Source theme | Doc |
|--------------|-----|
| PRD features / stack / prompts | `product/02-prd-production.md`, `engineering/*`, `architecture/*` |
| Strategy CJM / GTM / KPI | `strategy/*` |
| ЦА beauty / VPC | `strategy/02-target-audience-value-canvas.md` |
| GenUI industry practice | `design/04-generative-ui-catalog.md`, `design/02-ux-principles.md` |
| Alfa brand / DS | `design/06-alfa-visual-language.md`, `design/01-design-system.md` |
| Alfa UI screenshots 2026 | `design/references/alfa-web-2026/`, `design/08-alfa-marketing-ui-patterns.md` |
| Dashboard UX 2025–26 | `design/07-dashboard-patterns.md` |
| Demo cut | `product/01-prd-demo-mvp.md`, `superpowers/plans/*` |
| Demo Go backend | `architecture/09-demo-backend-go.md` (decision: Go, not FastAPI) |
| VPS deploy prep | `engineering/07-server-deployment.md` |
