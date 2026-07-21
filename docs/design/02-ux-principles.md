# UX Principles

## 1. Trust in first 3 seconds

Плотность, типографика, точность сумм отвечают на вопрос «можно ли доверить деньги» раньше, чем фичи. Никакого визуального шума в money flows.

## 2. Not a chatbot — a cockpit

Чат связывает сервисы. Ценность — **виджеты и действия**, не стены текста. Буллиты > лонгриды.

## 3. Generative UI with bank control

ИИ собирает **одобренные** компоненты (SDUI allowlist). Не исполняет произвольный код. Банк контролирует components + compliance; ИИ решает *какую* карточку показать.

## 4. Intentional friction

| Низкое трение | Высокое трение |
|---------------|----------------|
| Баланс, FAQ, расчёт налога | Подтверждение платежа, смена режима, кредит |
| Онбординг-чат | Подписание УКЭП |

## 5. Progressive disclosure (CJM)

Новичок видит уровни 1–5. Кредиты, найм, франшиза — когда дозрел уровень. ИИ не вываливает всё сразу.

## 6. Conversational bar + navigation

Гибрид: классические вкладки банка + contextual AI bar (подсказки на экране терминала/платежа) без обязательного открытия полного чата.

## 7. Plain language

«Вам нужно отложить 10 800 ₽ до 28 числа», не «согласно п. … рекомендуется сформировать платёжное поручение…» (ссылка на норму — вторична, в expand).

## 8. Empathy without cringe

При кассовом разрыве: признать риск → конкретный next step → оффер. Без сюсюканья и без высокомерия.

## 9. Fail gracefully

Низкий confidence RAG → «не уверена, давайте уточним / специалист». Никогда не выдумывать статью НК.

## 10. Accessibility

- Контраст AA на суммах и CTA.  
- Focus rings.  
- Не полагаться только на цвет светофора (текст уровня).  
- Touch targets ≥44px.

## 11. Decision-first dashboard (не KPI soup)

Home отвечает на **один** вопрос ЦА («сколько отложить на налог?») за секунды, затем даёт один primary CTA. Плотность — founder-minimal, не terminal. Framework Role→Metric→Density→Action: [07-dashboard-patterns.md](07-dashboard-patterns.md) · [Masterly 2026](https://www.themasterly.com/blog/fintech-dashboard-design-guide).

## 12. Proactive, explainable AI

ИИ тихий: summary / anomaly / next step. Каждый инсайт — источник + контроль пользователя ([Outcrowd 2026](https://medium.com/outcrowd/how-to-decide-what-your-financial-dashboard-actually-needs-in-2026-b4bddf5bc0ba)). Features find the user ([intent-driven financial UX](https://craftinnovations.global/future-of-financial-ux-webinar/)).

## 13. Выглядеть как Альфа, не как «AI-стартап»

Красный `#EF3124`, платина `#D1D5D8`, белые панели, Interface Sans, воздух Альфа-Бизнес — [06-alfa-visual-language.md](06-alfa-visual-language.md). Delight purple/lime только на маркетинге.
