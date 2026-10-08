package agents

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/calc"
	"github.com/alphabank-case-champ/copilot-api/internal/domain"
	"github.com/alphabank-case-champ/copilot-api/internal/guard"
	"github.com/alphabank-case-champ/copilot-api/internal/llm"
	"github.com/alphabank-case-champ/copilot-api/internal/rag"
	"github.com/alphabank-case-champ/copilot-api/internal/sdui"
	"github.com/alphabank-case-champ/copilot-api/internal/tools"
	"github.com/google/uuid"
)

type StreamEvent struct {
	Event string
	Data  any
}

type Orchestrator struct {
	Store   *memory.Store
	Tools   *tools.Registry
	LLM     *llm.Client
	Offline bool
}

func (o *Orchestrator) scope(ctx context.Context) *memory.UserScope {
	return o.Store.ForContext(ctx)
}

func (o *Orchestrator) tools(ctx context.Context) *tools.Registry {
	return &tools.Registry{Store: o.scope(ctx)}
}

func (o *Orchestrator) Handle(message string) []StreamEvent {
	return o.HandleCtx(memory.WithUserID(context.Background(), memory.MashaUserID), message, nil)
}

func (o *Orchestrator) HandleCtx(ctx context.Context, message string, history []domain.ChatMessage) []StreamEvent {
	msg := strings.TrimSpace(message)
	lower := normalizeRU(strings.ToLower(msg))

	if blocked := guard.CheckIngress(msg); blocked != nil {
		evs := []StreamEvent{
			{Event: "meta", Data: map[string]any{"intent": "GUARD"}},
			{Event: "error", Data: map[string]any{"code": blocked.Code, "message": blocked.Message}},
		}
		evs = append(evs, tokens(blocked.Message)...)
		evs = append(evs, suggestions(
			chip("Посчитать налог легально", "Сколько мне отложить на налог за этот месяц?"),
			chip("Оформить НПД", "Привет! Я начинаю бизнес, помоги оформить НПД."),
			chip("Проверить ИНН", "Проверь ИНН 7707083893"),
		)...)
		evs = append(evs, StreamEvent{Event: "done", Data: map[string]any{"finish_reason": "stop"}})
		return evs
	}

	sess := o.scope(ctx)
	if len(history) == 0 {
		history = sess.RecentChat(8)
	}
	// Follow-ups about an uploaded contract.
	if doc, ok := sess.LastDocument(); ok && mentionsDocument(lower) {
		lower = lower + " договор"
		_ = doc
	}

	intent := o.resolveIntent(ctx, msg, lower, historyLines(history))
	events := []StreamEvent{
		{Event: "meta", Data: map[string]any{"intent": intent, "conversation_id": uuid.NewString()}},
	}

	switch intent {
	case "ONBOARDING":
		events = append(events, o.onboarding(ctx, msg)...)
	case "TAX_CALC":
		events = append(events, o.tax(ctx, msg)...)
	case "TRANSACTION":
		events = append(events, o.payment(ctx, msg)...)
	case "LEGAL_REVIEW":
		events = append(events, o.legal(ctx, msg)...)
	case "UNIT_ECON":
		events = append(events, o.unitEcon(ctx, msg)...)
	case "PIGGY":
		events = append(events, o.piggy(ctx, msg)...)
	case "COMPLIANCE":
		events = append(events, o.compliance(ctx, msg)...)
	default:
		events = append(events, o.general(ctx, msg, history)...)
	}

	// GENERAL_QA / LEGAL already call the LLM when online.
	// Other intents: expand deterministic FACTS into a practical answer (numbers locked).
	if !o.Offline && o.LLM != nil && intent != "GENERAL_QA" && intent != "LEGAL_REVIEW" {
		events = o.enrichFacts(ctx, msg, intent, events)
	}

	// Persist short memory for follow-ups.
	asst := collectTokenText(events)
	sess.AppendChat("user", msg)
	if asst != "" {
		sess.AppendChat("assistant", asst)
	}

	events = append(events, StreamEvent{Event: "done", Data: map[string]any{"finish_reason": "stop"}})
	return events
}

func mentionsDocument(lower string) bool {
	return containsAny(lower, "договор", "документ", "pdf", "файл", "пункт", "штраф", "аренд", "что там", "разбери")
}

func collectTokenText(events []StreamEvent) string {
	var b strings.Builder
	for _, e := range events {
		if e.Event != "token" {
			continue
		}
		if m, ok := e.Data.(map[string]any); ok {
			b.WriteString(fmt.Sprint(m["text"]))
		}
	}
	return strings.TrimSpace(b.String())
}

func route(lower string) string {
	if containsAny(lower, "платежк", "платежн", "поручен", "сформируй плат", "нужна платеж") {
		return "TRANSACTION"
	}
	// FAQ / definitional / comparison questions before tax keywords.
	if containsAny(lower,
		"что такое", "чем отличается", "расскажи про", "как работает",
		"какие тариф", "тариф рко", "комиссия за", "сколько стоит счет",
		"кэшбек", "кешбек", "что умеет", "как считает",
		"отлича", "разниц", "сравни", "простыми словами", "в чём разница", "чем нпд", "или усн",
	) {
		return "GENERAL_QA"
	}
	if containsAny(lower, "налог", "нпд", "усн", "сколько отложить", "фнс", "отложить на налог") {
		return "TAX_CALC"
	}
	// Legal before unit-econ: "договор аренды" must not fall into UNIT_ECON via «аренд».
	if containsAny(lower, "договор", "legal", "оферт") {
		return "LEGAL_REVIEW"
	}
	if containsAny(lower, "безубыт", "unit", "юнит", "расходник") ||
		(containsAny(lower, "клиент") && containsAny(lower, "день", "сколько", "нужно", "цена")) ||
		(containsAny(lower, "аренд") && containsAny(lower, "клиент", "цена", "сколько", "посчита")) {
		return "UNIT_ECON"
	}
	if containsAny(lower, "копил") {
		return "PIGGY"
	}
	if containsAny(lower, "инн", "контрагент", "115", "светофор") {
		return "COMPLIANCE"
	}
	if containsAny(lower, "маникюр", "ресниц", "без ип", "начинаю бизнес", "зовут", "на дому", "открыть бизнес") ||
		(containsAny(lower, "привет") && containsAny(lower, "казан", "тысяч", "месяц")) {
		return "ONBOARDING"
	}
	return "GENERAL_QA"
}

func chip(label, text string) map[string]string {
	return map[string]string{"label": label, "text": text}
}

func suggestions(items ...map[string]string) []StreamEvent {
	return []StreamEvent{{Event: "suggestions", Data: map[string]any{"items": items}}}
}

func (o *Orchestrator) onboarding(ctx context.Context, msg string) []StreamEvent {
	cur := o.scope(ctx).Profile()
	sphere, city, rev := cur.BusinessSphere, cur.City, 0.0
	if sphere == "" {
		sphere = "Услуги"
	}
	if city == "" {
		city = "Город"
	}
	if v := extractMoney(msg); v > 0 {
		rev = v
	}
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "казан") {
		city = "Казань"
	}
	if strings.Contains(lower, "маникюр") {
		sphere = "Маникюр"
	}
	if strings.Contains(lower, "ресниц") {
		sphere = "Наращивание ресниц"
	}
	if strings.Contains(lower, "контент") || strings.Contains(lower, "блог") {
		sphere = "Контент / креатив"
	}
	p := o.tools(ctx).UpdateOnboardingProfile(sphere, city, rev)
	if p.CJMLevel < 2 {
		p.CJMLevel = 2
	}
	o.scope(ctx).UpdateProfile(p)

	text := fmt.Sprintf(
		"Черновик профиля: %s · %s · %s · оборот ~%s ₽/мес.\n\n"+
			"Что это значит на практике:\n"+
			"• На старте без сотрудников обычно берут НПД: 4%% с физлиц / 6%% с юрлиц, чеки в «Мой налог», лимит 2,4 млн ₽/год\n"+
			"• ИП на УСН 6%% имеет смысл, если появятся сотрудники, крупные B2B-контракты или выход за лимит НПД\n"+
			"• До подтверждения профиля цифры в налоге и копилке считаются по этому черновику\n\n"+
			"Следующие шаги: подтвердите карточку → посчитайте налог за месяц → при аренде кабинета проверьте безубыточность и ИНН арендодателя.\n\n%s",
		p.DisplayName, p.BusinessSphere, p.City, formatRub(p.MonthlyRevenueEstimate), domain.Disclaimer,
	)
	env := sdui.New("OnboardingSummary", map[string]any{
		"display_name":             p.DisplayName,
		"business_sphere":          p.BusinessSphere,
		"city":                     p.City,
		"monthly_revenue_estimate": p.MonthlyRevenueEstimate,
		"suggested_regimes":        []string{"NPD", "USN_6"},
		"suggested_okved":          p.OKVED,
		"cta_label":                "Всё верно",
		"cta_action":               "confirm_onboarding",
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	out = append(out, suggestions(
		chip("Посчитать налог", "Сколько мне отложить на налог за этот месяц?"),
		chip("Безубыточность", "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500."),
		chip("Копилка 6%", "Включи копилку 6%."),
	)...)
	return out
}

func (o *Orchestrator) tax(ctx context.Context, msg string) []StreamEvent {
	p := o.scope(ctx).Profile()
	income := p.MonthlyRevenueEstimate
	if income <= 0 {
		income = o.scope(ctx).IncomeThisMonth()
	}
	if v := extractMoney(msg); v >= 10000 {
		income = v
	}
	regime := calc.RegimeNPD
	if strings.Contains(strings.ToLower(msg), "усн") {
		regime = calc.RegimeUSN6
	}
	period := monthLabel()
	hits := o.tools(ctx).RetrieveKB(msg+" НПД УСН ЕНП срок уплаты 28", 5)
	res := o.tools(ctx).CalculateTax(regime, income, period, false)
	debt := o.tools(ctx).FetchFNSDebt()
	sourceRef := res.SourceRef
	if len(hits) > 0 {
		sourceRef = hits[0].ID
	}
	b2cRub := round2(res.Income * res.B2CShare)
	b2bRub := round2(res.Income - b2cRub)
	text := fmt.Sprintf(
		"Ориентир налога по режиму %s за %s: %s ₽ (%s) с оборота %s ₽.\n\n"+
			"Как сложилась сумма:\n"+
			"• от физлиц ~%s ₽ (доля %.0f%%) — ставка ближе к 4%%\n"+
			"• от юрлиц/ИП ~%s ₽ — ставка ближе к 6%%\n"+
			"• с запасом к сроку удобнее держать %s ₽ (оценка по верхней ставке 6%%)\n\n"+
			"Практика уплаты:\n"+
			"1) Начисление смотрите в «Мой налог» / ЛК; через ЕНП/ЕНС обычно платят до 28-го числа следующего месяца\n"+
			"2) Включите копилку — откладывайте %% с поступлений, чтобы не вынимать налог из оборота в последний день\n"+
			"3) Когда сумма ясна — сформируйте черновик платёжки в чате и подтвердите только после проверки\n\n"+
			"Не учтено: точные чеки периода, региональные льготы, вычеты. Это не квитанция ФНС.",
		res.Regime, period, formatRub(res.TaxAmount), res.RateLabel, formatRub(res.Income),
		formatRub(b2cRub), res.B2CShare*100, formatRub(b2bRub), formatRub(res.CeilingAmount),
	)
	if regime == calc.RegimeUSN6 {
		text = fmt.Sprintf(
			"Ориентир УСН «доходы» 6%% за %s: %s ₽ с оборота %s ₽.\n\n"+
				"Важно:\n"+
				"• Нужен статус ИП/ООО; авансы по кварталам, налог по году — через ЕНП/ЕНС\n"+
				"• Страховые взносы и региональная ставка в этой оценке не вычтены\n"+
				"• Сравнивайте с НПД на том же обороте — при высокой доле B2C самозанятость часто дешевле\n\n"+
				"Дальше: включите копилку под ЕНП и при необходимости сформируйте черновик платёжки.\n\n"+
				"Справочный расчёт, не декларация.",
			period, formatRub(res.TaxAmount), formatRub(res.Income),
		)
	}
	if cite := pickCitation(hits, "енп", "срок", "28"); cite != "" {
		text += "\n\n" + cite
	} else if cite := pickCitation(hits, "нпд", "усн", "ставк"); cite != "" {
		text += "\n\n" + cite
	}
	if hasDebt, _ := debt["has_debt"].(bool); !hasDebt {
		text += "\n\nПо доступной проверке текущих задолженностей не видно."
	}
	text += "\n\n" + domain.Disclaimer
	env := sdui.New("TaxCard", map[string]any{
		"regime":         res.Regime,
		"period_label":   period,
		"tax_amount":     res.TaxAmount,
		"ceiling_amount": res.CeilingAmount,
		"currency":       "RUB",
		"rate_label":     res.RateLabel,
		"income":         res.Income,
		"b2c_share":      res.B2CShare,
		"assumptions":    res.Assumptions,
		"disclaimer":     domain.Disclaimer,
		"source_ref":     sourceRef,
		"cta":            map[string]any{"action": "create_payment_draft", "label": "Сформировать платёжку"},
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	if len(hits) > 0 {
		out = append(out, knowledgeSourcesEvent(hits)...)
	}
	out = append(out, suggestions(
		chip("Сформировать ЕНП", "Сформируй платёжку."),
		chip("Включить копилку", "Включи копилку 6%."),
		chip("Сравнить с УСН 6%", "Посчитай УСН 6% с оборота 180 тысяч"),
		chip("Приём оплаты", "Как принимать оплату по СБП и эквайрингу?"),
	)...)
	return out
}

func wantsCounterpartyPayment(lower string) bool {
	taxExplicit := containsAny(lower, "налог", "нпд", "енп", "фнс", "усн", "уфк", "бюджет")
	counterparty := containsAny(lower,
		"аренд", "арендодател", "контрагент", "поставщик", "закуп", "оплат аренд",
		"платежк аренд", "платёжк аренд", "платежку аренд",
	)
	if counterparty && !taxExplicit {
		return true
	}
	// «платёжка аренды на налог» is odd; prefer counterparty when both rent + draft verbs.
	return counterparty && containsAny(lower, "платежк", "платежн", "поручен", "сформируй плат")
}

func (o *Orchestrator) payment(ctx context.Context, msg string) []StreamEvent {
	lower := normalizeRU(strings.ToLower(msg))
	if wantsCounterpartyPayment(lower) {
		return o.paymentCounterparty(ctx, msg, lower)
	}
	return o.paymentTax(ctx, msg)
}

func (o *Orchestrator) paymentTax(ctx context.Context, msg string) []StreamEvent {
	p := o.scope(ctx).Profile()
	income := p.MonthlyRevenueEstimate
	if income <= 0 {
		income = o.scope(ctx).IncomeThisMonth()
	}
	amount := 0.0
	if v := extractMoney(msg); v >= 100 {
		if safe, ok := guard.SanitizeDraftAmount(v); ok {
			amount = safe
		}
	}
	purpose := "Налог НПД за " + monthLabel()
	if amount <= 0 {
		res := o.tools(ctx).CalculateTax(calc.RegimeNPD, income, monthLabel(), true)
		amount = res.TaxAmount
	}
	d := o.tools(ctx).CreatePaymentDraft(amount, purpose)
	text := fmt.Sprintf(
		"Черновик платёжки на %s ₽ — %s.\n\n"+
			"Реквизиты:\n"+
			"• получатель: %s (бюджет / ЕНП)\n"+
			"• статус: draft — деньги ещё не списаны\n"+
			"• риск-метка: %s\n\n"+
			"Перед подтверждением:\n"+
			"1) Сверьте сумму с начислением в «Мой налог» / расчётом налога\n"+
			"2) Убедитесь, что в копилке ЕНП достаточно средств или пополните её\n"+
			"3) Нажмите «Подтвердить» только когда сумма и назначение верны\n\n"+
			"Для аренды/поставщика — отдельная платёжка после проверки ИНН. История — в «Платежи». %s",
		formatRub(d.Amount), d.Purpose, d.PayeeName, d.RiskLevel, domain.Disclaimer,
	)
	env := paymentDraftEnvelope(d)
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	out = append(out, suggestions(
		chip("К платежам", "__href__:/app/payments"),
		chip("Платёжка аренды", "Сформируй платёжку аренды 40000 на ИНН 1650987654"),
		chip("Открыть копилку", "Включи копилку 6%."),
		chip("Пересчитать налог", "Сколько мне отложить на налог за этот месяц?"),
	)...)
	return out
}

func (o *Orchestrator) paymentCounterparty(ctx context.Context, msg, lower string) []StreamEvent {
	sess := o.scope(ctx)
	inn := extractINN(msg)
	risk, hasRisk := sess.LastRisk()
	if inn == "" && hasRisk {
		inn = risk.INN
	}
	if inn == "" {
		inn = "1650987654"
	}
	r := o.tools(ctx).CheckCounterpartyRisk(inn)
	sess.SetLastRisk(r)

	amount := 40000.0
	if v := extractMoney(msg); v >= 100 {
		if safe, ok := guard.SanitizeDraftAmount(v); ok {
			amount = safe
		}
	} else if v := extractMoneyNear(msg, "аренд"); v >= 100 {
		if safe, ok := guard.SanitizeDraftAmount(v); ok {
			amount = safe
		}
	}

	purpose := "Аренда кабинета за " + monthLabel()
	if containsAny(lower, "поставщик", "закуп", "материал", "расходник") && !containsAny(lower, "аренд") {
		purpose = "Оплата поставщику за " + monthLabel()
	}
	payee := "Контрагент ИНН " + inn
	if strings.HasPrefix(purpose, "Аренда") {
		payee = "Арендодатель ИНН " + inn
	}

	d := o.tools(ctx).CreatePaymentDraftSpec(tools.PaymentDraftSpec{
		Amount: amount, Purpose: purpose, PayeeName: payee, PayeeINN: inn, RiskLevel: r.Level,
	})

	levelHint := map[string]string{
		"green":  "зелёный — можно подтверждать после сверки договора",
		"yellow": "жёлтый — уточните документы до подтверждения",
		"red":    "красный — лучше не подтверждать без ручной проверки",
	}[r.Level]
	if levelHint == "" {
		levelHint = r.Level
	}

	text := fmt.Sprintf(
		"Черновик платежа контрагенту на %s ₽ — %s.\n\n"+
			"Реквизиты:\n"+
			"• получатель: %s\n"+
			"• ИНН: %s\n"+
			"• светофор 115-ФЗ: %s\n"+
			"• статус: draft — деньги ещё не списаны\n\n"+
			"%s\n\n"+
			"Перед подтверждением:\n"+
			"1) Сверьте сумму с договором/счётом\n"+
			"2) Проверьте назначение платежа (аренда / поставка)\n"+
			"3) При красном/жёлтом светофоре — не подтверждайте «вслепую»\n\n"+
			"Налог в бюджет — отдельной платёжкой на УФК. %s",
		formatRub(d.Amount), d.Purpose, d.PayeeName, inn, levelHint, r.Recommendation, domain.Disclaimer,
	)
	env := paymentDraftEnvelope(d)
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	out = append(out, StreamEvent{Event: "sdui", Data: sdui.New("ComplianceTrafficLight", map[string]any{
		"inn": r.INN, "level": r.Level, "title": r.Title, "reasons": r.Reasons, "recommendation": r.Recommendation,
	})})
	out = append(out, suggestions(
		chip("К платежам", "__href__:/app/payments"),
		chip("Платёжка налога", "Сформируй платёжку."),
		chip("Другой ИНН", "Проверь ИНН 7707083893"),
		chip("Разбор договора", "Разбери договор аренды pdf"),
	)...)
	return out
}

func paymentDraftEnvelope(d domain.PaymentDraft) sdui.Envelope {
	return sdui.New("PaymentDraftCard", map[string]any{
		"draft_id":          d.ID,
		"amount":            d.Amount,
		"currency":          "RUB",
		"purpose":           d.Purpose,
		"payee_name":        d.PayeeName,
		"payee_inn":         d.PayeeINN,
		"status":            d.Status,
		"risk_level":        d.RiskLevel,
		"cta_confirm_label": "Подтвердить",
		"cta_cancel_label":  "Отменить",
	})
}

func (o *Orchestrator) unitEcon(ctx context.Context, msg string) []StreamEvent {
	fixed, variable, price := 40000.0, 200.0, 1500.0
	if v := extractMoneyNear(msg, "аренд"); v > 0 {
		fixed = v
	}
	if v := extractMoneyNear(msg, "расход"); v > 0 {
		variable = v
	}
	if v := extractMoneyNear(msg, "цена"); v > 0 {
		price = v
	}
	// Informal "аренда 40" without «к» still means thousands in this pitch domain.
	if fixed > 0 && fixed < 1000 {
		if v := extractMoneyNear(msg, "аренд"); v > 0 && v < 1000 && !hasExplicitRublesNear(msg, "аренд") {
			fixed = v * 1000
		}
	}
	nums := extractAllMoneyAmounts(msg)
	if len(nums) >= 3 {
		fixed, variable, price = nums[0], nums[1], nums[2]
		if fixed > 0 && fixed < 1000 && !hasExplicitRublesNear(msg, "аренд") {
			fixed *= 1000
		}
	}
	res := o.tools(ctx).ComputeUnitEconomics(fixed, variable, price)
	days := res.WorkingDaysPerMonth
	if days <= 0 {
		days = 22
	}
	ratePct := res.TaxRate * 100
	revAtBE := round2(res.BreakevenUnitsPerDay * float64(days) * res.PricePerUnit)
	text := fmt.Sprintf(
		"Точка безубыточности ≈ %.0f %s в день (реалистичная модель).\n\n"+
			"Разбор:\n"+
			"• Цена %s ₽ − расходники %s ₽ = валовая маржа %s ₽\n"+
			"• Минус НПД ~%.1f%% с чека (−%s ₽) → чистый вклад в фикс %s ₽/клиент\n"+
			"• Аренда %s ₽/мес при %d рабочих днях → ≈ %.0f клиентов/мес\n"+
			"• С запасом 15%% на отмены/простои — ориентир %.0f %s/день\n"+
			"• Выручка «на ноль» при такой загрузке ≈ %s ₽/мес\n\n"+
			"На графике при 0 клиентах убыток −%s ₽. "+
			"Не учтены: ваше время, реклама, эквайринг, комиссии банка и кассовые разрывы.\n\n"+
			"Что сделать:\n"+
			"1) Проверьте ИНН арендодателя и сформируйте платёжку аренды\n"+
			"2) Посчитайте налог с плановой выручки и включите копилку ЕНП\n"+
			"3) Ускорьте приход: СБП/эквайринг снижают кассовые разрывы до безубыточности\n\n%s",
		res.BreakevenUnitsPerDay, clientsWord(res.BreakevenUnitsPerDay),
		formatRub(res.PricePerUnit), formatRub(res.VariableCostPerUnit), formatRub(res.GrossMarginPerUnit),
		ratePct, formatRub(res.TaxPerUnit), formatRub(res.MarginPerUnit),
		formatRub(res.FixedCostsMonthly), days, res.ClientsPerMonth,
		res.RecommendedUnitsPerDay, clientsWord(res.RecommendedUnitsPerDay),
		formatRub(revAtBE), formatRub(res.FixedCostsMonthly), domain.Disclaimer,
	)
	series := make([]map[string]any, 0, len(res.Series))
	for _, p := range res.Series {
		series = append(series, map[string]any{"units": p.Units, "profit": p.Profit})
	}
	env := sdui.New("UnitEconomicsChart", map[string]any{
		"breakeven_units_per_day":   res.BreakevenUnitsPerDay,
		"recommended_units_per_day": res.RecommendedUnitsPerDay,
		"price_per_unit":            res.PricePerUnit,
		"variable_cost_per_unit":    res.VariableCostPerUnit,
		"fixed_costs_monthly":       res.FixedCostsMonthly,
		"margin_per_unit":           res.MarginPerUnit,
		"gross_margin_per_unit":     res.GrossMarginPerUnit,
		"tax_per_unit":              res.TaxPerUnit,
		"tax_rate":                  res.TaxRate,
		"working_days_per_month":    res.WorkingDaysPerMonth,
		"series":                    series,
	})
	rentAmt := int(res.FixedCostsMonthly)
	if rentAmt <= 0 {
		rentAmt = 40000
	}
	offers := sdui.New("ProductOffers", map[string]any{
		"title":    "Оплаты вокруг безубыточности",
		"subtitle": "Связка: принять деньги от клиентов → оплатить аренду → отложить налог",
		"items": []map[string]any{
			{
				"id": "sbp", "title": "СБП для клиентов",
				"subtitle": "Мгновенный приход на РКО без терминала",
				"chat":     "Как принимать оплату по СБП на счёт ИП/самозанятого?",
			},
			{
				"id": "acquiring", "title": "Эквайринг",
				"subtitle": "Карты в кабинете — выше средний чек",
				"chat":     "Какой эквайринг Альфа-Банка подойдёт салону красоты?",
			},
			{
				"id": "rent", "title": "Платёжка аренды",
				"subtitle": fmt.Sprintf("Черновик на %s ₽ после проверки ИНН", formatRub(float64(rentAmt))),
				"chat":     fmt.Sprintf("Сформируй платёжку аренды %d на ИНН 1650987654", rentAmt),
			},
		},
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	out = append(out, StreamEvent{Event: "sdui", Data: offers})
	out = append(out, suggestions(
		chip("Посчитать налог", "Сколько мне отложить на налог за этот месяц?"),
		chip("Проверить аренду", "Проверь ИНН 1650987654"),
		chip("Платёжка аренды", fmt.Sprintf("Сформируй платёжку аренды %d на ИНН 1650987654", rentAmt)),
		chip("СБП / эквайринг", "Как принимать оплату по СБП и эквайрингу?"),
	)...)
	return out
}

func (o *Orchestrator) piggy(ctx context.Context, msg string) []StreamEvent {
	sess := o.scope(ctx)
	p := sess.Piggy()
	wasEnabled := p.Enabled
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "выкл") {
		p.Enabled = false
	} else {
		p.Enabled = true
	}
	if r := extractPercent(msg); r > 0 {
		p.RatePercent = r
	}
	prof := sess.Profile()
	p.LastContribution = round2(prof.MonthlyRevenueEstimate * p.RatePercent / 100 / 30)
	ceiling := calc.CalculateTax(calc.TaxInput{
		Regime: calc.RegimeNPD, Income: prof.MonthlyRevenueEstimate, ConservativeCeiling: true,
	})
	p.TargetAmount = ceiling.TaxAmount
	// Only credit once when turning on (or explicit top-up) — avoid stacking on every chat hit.
	if p.Enabled && (!wasEnabled || containsAny(lower, "пополн", "вклад", "отлож")) {
		p.Balance = round2(p.Balance + p.LastContribution)
	}
	sess.SetPiggy(p)
	gap := round2(p.TargetAmount - p.Balance)
	if gap < 0 {
		gap = 0
	}
	daysToFill := 0.0
	if p.LastContribution > 0 && gap > 0 {
		daysToFill = round2(gap / p.LastContribution)
	}
	text := fmt.Sprintf(
		"Копилка ЕНП %s · ставка %.0f%% с поступлений.\n\n"+
			"Состояние:\n"+
			"• баланс %s ₽\n"+
			"• цель к сроку уплаты (ориентир налога с запасом) %s ₽\n"+
			"• осталось накопить %s ₽\n"+
			"• дневной ориентир вклада ~%s ₽",
		map[bool]string{true: "включена", false: "выключена"}[p.Enabled],
		p.RatePercent, formatRub(p.Balance), formatRub(p.TargetAmount), formatRub(gap), formatRub(p.LastContribution),
	)
	if daysToFill > 0 {
		text += fmt.Sprintf("\n• в текущем темпе до цели ≈ %.0f дн.", daysToFill)
	}
	text += "\n\nЗачем: уплата НПД/налога через ЕНС обычно к 28-му числу следующего месяца — "+
		"копилка не даёт снять налог одним платежом из оборота.\n\n"+
		"Дальше: когда баланс ≈ цели, сформируйте черновик платёжки. Ставку можно менять (например 6–8%%) в чате или в разделе «Копилка».\n\n"+
		domain.Disclaimer
	env := sdui.New("EnpPiggyBank", map[string]any{
		"enabled": p.Enabled, "rate_percent": p.RatePercent, "balance": p.Balance,
		"currency": "RUB", "last_contribution": p.LastContribution, "target_amount": p.TargetAmount,
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	out = append(out, suggestions(
		chip("Сформировать ЕНП", "Сформируй платёжку."),
		chip("Пересчитать налог", "Сколько мне отложить на налог за этот месяц?"),
		chip("Ставка 8%", "Включи копилку 8%."),
	)...)
	return out
}

func (o *Orchestrator) compliance(ctx context.Context, msg string) []StreamEvent {
	sess := o.scope(ctx)
	inn := extractINN(msg)
	if inn == "" {
		inn = "1650987654"
	}
	r := o.tools(ctx).CheckCounterpartyRisk(inn)
	sess.SetLastRisk(r)
	levelRu := map[string]string{
		"green": "зелёный — можно платить", "yellow": "жёлтый — есть замечания", "red": "красный — риск",
	}[r.Level]
	if levelRu == "" {
		levelRu = r.Level
	}
	text := fmt.Sprintf(
		"Светофор 115-ФЗ по ИНН %s: %s.\n\n%s\n\n%s\n\n"+
			"Что обычно снижает риск задержки платежа:\n"+
			"• договор / акт с понятным назначением платежа\n"+
			"• совпадение суммы с договором и реальной деятельностью\n"+
			"• не дробить платежи «ради обхода» — только по делу\n"+
			"• проверять арендодателя и поставщиков до перевода крупной суммы\n\n"+
			"Дальше: сформируйте черновик платёжки аренды/поставщику на этот ИНН — "+
			"светофор подтянется в поручение. Налог в бюджет — отдельной платёжкой.\n\n"+
			"Это экспресс-проверка по справочнику, не полноценный KYC банка. %s",
		inn, levelRu, r.Title, r.Recommendation, domain.Disclaimer,
	)
	hits := o.tools(ctx).RetrieveKB(msg+" 115-ФЗ светофор контрагент проверка ИНН", 4)
	if cite := pickCitation(hits, "115", "светофор", "инн"); cite != "" {
		text += "\n\n" + cite
	}
	rentChat := fmt.Sprintf("Сформируй платёжку аренды 40000 на ИНН %s", inn)
	env := sdui.New("ComplianceTrafficLight", map[string]any{
		"inn": r.INN, "level": r.Level, "title": r.Title, "reasons": r.Reasons, "recommendation": r.Recommendation,
		"cta": map[string]any{"action": "chat", "label": "Платёжка аренды 40к", "text": rentChat},
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	if len(hits) > 0 {
		out = append(out, knowledgeSourcesEvent(hits)...)
	}
	out = append(out, suggestions(
		chip("Платёжка аренды 40к", rentChat),
		chip("Платёжка налога", "Сформируй платёжку."),
		chip("Разобрать договор", "Разбери договор аренды PDF"),
		chip("Другой ИНН", "Проверь ИНН 7707083893"),
	)...)
	return out
}

func (o *Orchestrator) legal(ctx context.Context, msg string) []StreamEvent {
	sess := o.scope(ctx)
	name := "dogovor-arenda-kabinet.pdf"
	excerpt := ""
	var flags any
	mode := "heuristic_fallback"
	docBody := ""
	if doc, ok := sess.LastDocument(); ok {
		name = doc.Name
		excerpt = doc.Excerpt
		flags = doc.Flags
		mode = doc.Mode
		docBody = doc.Markdown
		if docBody == "" {
			docBody = doc.FullText
		}
		if docBody == "" {
			docBody = doc.Excerpt
		}
	}
	if flags == nil {
		scan := o.tools(ctx).ScanLegalDocument(name, firstNonEmpty(docBody, excerpt))
		flags = scan["flags"]
		mode = fmt.Sprint(scan["mode"])
		if excerpt == "" {
			if ex, ok := scan["excerpt"].(string); ok {
				excerpt = ex
			}
		}
		if docBody == "" {
			docBody = excerpt
		}
	}
	hits := o.tools(ctx).RetrieveKB("договор аренды одностороннее изменение цены штраф подсудность ГК", 4)
	text := o.legalNarrative(ctx, msg, name, firstNonEmpty(docBody, excerpt), flags)
	env := sdui.New("LegalFlagsList", map[string]any{
		"document_name": name,
		"flags":         flags,
		"mode":          mode,
	})
	out := append(tokens(text), StreamEvent{Event: "sdui", Data: env})
	if len(hits) > 0 {
		out = append(out, knowledgeSourcesEvent(hits)...)
	}
	out = append(out, suggestions(
		chip("Проверить ИНН арендодателя", "Проверь ИНН 1650987654"),
		chip("Посчитать аренду в юнит-экономике", "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500."),
		chip("Загрузить свой PDF", "Разбери договор аренды pdf"),
	)...)
	return out
}

func (o *Orchestrator) legalNarrative(ctx context.Context, msg, name, excerpt string, flags any) string {
	fallback := fmt.Sprintf(
		"Разбор договора «%s» — риски на карточке ниже (по убыванию важности).\n\n"+
			"Как пользоваться:\n"+
			"• high — просите правку до подписания (срок уведомления, потолок роста цены, выход без штрафа)\n"+
			"• medium — торгуйтесь по неустойке и срокам предупреждения\n"+
			"• low — зафиксируйте удобную подсудность, если можете\n\n"+
			"Свой PDF загрузите кнопкой «Документ». Перед крупной арендой проверьте ИНН контрагента.\n\n%s",
		name, domain.Disclaimer,
	)
	if o.Offline || o.LLM == nil {
		return fallback
	}
	body := excerpt
	if len([]rune(body)) > 14000 {
		body = string([]rune(body)[:14000]) + "\n…[обрезано]"
	}
	user := fmt.Sprintf(
		"User: %s\nDocument: %s\nFLAGS:\n%v\n\nDOCUMENT (markdown/text extracted for you):\n%s\n\n"+
			"Сделай прикладной разбор для переговоров. Цитируй документ. Не выдумывай отсутствующие пункты.",
		msg, name, flags, body,
	)
	if text, err := o.LLM.CompleteBudget(ctx, legalSystemPrompt(), user, 0.4, 480); err == nil {
		text = strings.TrimSpace(text)
		if text != "" && guard.CheckEgress(text) && !looksLikeWidgetGarbage(text) {
			if !strings.Contains(strings.ToLower(text), "справоч") && !strings.Contains(strings.ToLower(text), "не заменяет") {
				text += "\n\n" + domain.Disclaimer
			}
			return text
		}
	}
	return fallback
}

func (o *Orchestrator) general(ctx context.Context, msg string, history []domain.ChatMessage) []StreamEvent {
	hits := o.tools(ctx).RetrieveKB(msg, 6)
	lower := normalizeRU(strings.ToLower(msg))
	greetingOnly := containsAny(lower, "привет", "здравств", "хай", "hello") &&
		len([]rune(msg)) < 40 &&
		!containsAny(lower, "налог", "нпд", "тариф", "договор", "копил", "инн")

	if greetingOnly && rag.LowConfidence(hits) {
		return o.generalMenu(ctx)
	}

	text := o.answerFromKB(ctx, msg, hits, history)
	out := tokens(text)
	if !rag.LowConfidence(hits) {
		out = append(out, knowledgeSourcesEvent(hits)...)
	}
	out = append(out, suggestions(
		chip("Налог за месяц", "Сколько мне отложить на налог за этот месяц?"),
		chip("Что такое НПД?", "Что такое НПД и какая ставка?"),
		chip("НПД или УСН?", "Чем отличается НПД от УСН 6%?"),
		chip("Тариф РКО", "Какие тарифы РКО для микробизнеса?"),
		chip("Копилка ЕНП", "Как работает копилка ЕНП?"),
		chip("Безубыточность", "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500."),
	)...)
	return out
}

func (o *Orchestrator) generalMenu(ctx context.Context) []StreamEvent {
	out := tokens(o.generalMenuText(ctx))
	out = append(out, suggestions(
		chip("Налог за месяц", "Сколько мне отложить на налог за этот месяц?"),
		chip("Что такое НПД?", "Что такое НПД и какая ставка?"),
		chip("Платёжка", "Сформируй платёжку."),
		chip("Безубыточность", "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500."),
		chip("Проверить ИНН", "Проверь ИНН 1650987654"),
		chip("Копилка 6%", "Включи копилку 6%."),
	)...)
	return out
}

func (o *Orchestrator) answerFromKB(ctx context.Context, msg string, hits []rag.Hit, history []domain.ChatMessage) string {
	var ctxBlock strings.Builder
	for i, h := range hits {
		if i >= 6 {
			break
		}
		src := h.Source
		if src == "" {
			src = "kb"
		}
		ctxBlock.WriteString(fmt.Sprintf("[%d] id=%s | %s | source=%s\n%s\n\n", i+1, h.ID, h.Title, src, h.Text))
	}
	p := o.scope(ctx).Profile()
	profileLine := fmt.Sprintf(
		"Клиент: %s; сфера: %s; город: %s; оборот ≈ %.0f ₽/мес; режим: %s; CJM: %d.",
		p.DisplayName, p.BusinessSphere, p.City, p.MonthlyRevenueEstimate, p.TaxRegime, p.CJMLevel,
	)
	if doc, ok := o.scope(ctx).LastDocument(); ok {
		body := doc.Markdown
		if body == "" {
			body = doc.FullText
		}
		if body == "" {
			body = doc.Excerpt
		}
		ctxBlock.WriteString(fmt.Sprintf("[last_document name=%s quality=%s pages=%d]\n%s\n\n",
			doc.Name, doc.Quality, doc.Pages, truncateRunes(body, 4000)))
	}

	if !o.Offline && o.LLM != nil {
		var hist strings.Builder
		for _, h := range history {
			hist.WriteString(h.Role + ": " + truncateRunes(h.Content, 400) + "\n")
		}
		user := fmt.Sprintf(
			"PROFILE:\n%s\n\nDIALOGUE (recent):\n%s\nQUESTION:\n%s\n\nCONTEXT from knowledge base:\n%s\n"+
				"Ответь информативно и по делу. Если для точных цифр нужен калькулятор — предложи конкретную фразу запроса.",
			profileLine, hist.String(), msg, ctxBlock.String(),
		)
		for attempt := 0; attempt < 2; attempt++ {
			text, err := o.LLM.CompleteBudget(ctx, generalSystemPrompt(), user, 0.5, 420)
			if err != nil || strings.TrimSpace(text) == "" {
				continue
			}
			text = strings.TrimSpace(text)
			if !guard.CheckEgress(text) || looksLikeWidgetGarbage(text) {
				continue
			}
			if !strings.Contains(strings.ToLower(text), "справоч") &&
				!strings.Contains(strings.ToLower(text), "не заменяет") {
				text += "\n\n" + domain.Disclaimer
			}
			return text
		}
	}
	if len(hits) == 0 {
		return o.generalMenuText(ctx)
	}
	return groundedFallbackAnswer(hits)
}

func groundedFallbackAnswer(hits []rag.Hit) string {
	var b strings.Builder
	b.WriteString(hits[0].Title)
	b.WriteString("\n\n")
	b.WriteString(hits[0].Text)
	if len(hits) > 1 {
		b.WriteString("\n\nСвязанные материалы:\n")
		for i, h := range hits {
			if i == 0 {
				continue
			}
			if i >= 3 {
				break
			}
			b.WriteString("• ")
			b.WriteString(h.Title)
			b.WriteString(" — ")
			b.WriteString(h.Text)
			b.WriteString("\n")
		}
	}
	b.WriteString("\nЧто можно сделать дальше в чате: посчитать налог, сравнить НПД и УСН, включить копилку ЕНП, проверить ИНН или разобрать договор.\n\n")
	b.WriteString(domain.Disclaimer)
	return b.String()
}

func (o *Orchestrator) generalMenuText(ctx context.Context) string {
	p := o.scope(ctx).Profile()
	sphere := p.BusinessSphere
	if sphere == "" {
		sphere = "ваш бизнес"
	}
	city := p.City
	if city == "" {
		city = "ваш город"
	}
	rev := ""
	if p.MonthlyRevenueEstimate > 0 {
		rev = fmt.Sprintf(", оборот ~%s ₽/мес", formatRub(p.MonthlyRevenueEstimate))
	}
	return fmt.Sprintf(
		"Я «Альфа-Бизнес: Старт» — помощник для %s (%s, %s%s).\n\n"+
			"Могу помочь по делу:\n"+
			"• налог НПД/УСН и срок через ЕНП\n"+
			"• черновик платёжки и копилку под уплату\n"+
			"• точку безубыточности кабинета\n"+
			"• проверку ИНН (115-ФЗ) и риски в договоре аренды\n"+
			"• тарифы РКО, кэшбек и частые вопросы микробизнеса\n\n"+
			"Спросите своими словами — отвечу с опорой на базу знаний и расчёты. %s",
		p.DisplayName, sphere, city, rev, domain.Disclaimer,
	)
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func knowledgeSourcesEvent(hits []rag.Hit) []StreamEvent {
	items := make([]map[string]any, 0, len(hits))
	for i, h := range hits {
		if i >= 4 {
			break
		}
		snippet := h.Text
		if r := []rune(snippet); len(r) > 160 {
			snippet = string(r[:160]) + "…"
		}
		items = append(items, map[string]any{
			"id": h.ID, "title": h.Title, "source": h.Source,
			"snippet": snippet, "score": h.Score,
		})
	}
	return []StreamEvent{{
		Event: "sdui",
		Data: sdui.New("KnowledgeSources", map[string]any{
			"title": "Источники из базы знаний",
			"items": items,
		}),
	}}
}

func pickCitation(hits []rag.Hit, keys ...string) string {
	for _, h := range hits {
		blob := normalizeRU(strings.ToLower(h.Title + " " + h.Text + " " + h.ID))
		for _, k := range keys {
			if strings.Contains(blob, normalizeRU(k)) {
				snip := h.Text
				if r := []rune(snip); len(r) > 180 {
					snip = string(r[:180]) + "…"
				}
				return fmt.Sprintf("Из базы знаний (%s): %s", h.ID, snip)
			}
		}
	}
	return ""
}

func tokens(text string) []StreamEvent {
	parts := chunkRunes(text, 28)
	out := make([]StreamEvent, 0, len(parts))
	for _, p := range parts {
		out = append(out, StreamEvent{Event: "token", Data: map[string]any{"text": p}})
	}
	return out
}

func textOnly(code, msg string) []StreamEvent {
	return []StreamEvent{
		{Event: "meta", Data: map[string]any{"intent": "GUARD"}},
		{Event: "error", Data: map[string]any{"code": code, "message": msg}},
		{Event: "token", Data: map[string]any{"text": msg}},
		{Event: "done", Data: map[string]any{"finish_reason": "stop"}},
	}
}

func containsAny(s string, keys ...string) bool {
	for _, k := range keys {
		if strings.Contains(s, normalizeRU(k)) {
			return true
		}
	}
	return false
}

func normalizeRU(s string) string {
	return strings.ReplaceAll(s, "ё", "е")
}

func chunkRunes(s string, n int) []string {
	r := []rune(s)
	if len(r) == 0 {
		return []string{""}
	}
	var out []string
	for i := 0; i < len(r); i += n {
		j := i + n
		if j > len(r) {
			j = len(r)
		}
		out = append(out, string(r[i:j]))
	}
	return out
}

func monthLabel() string {
	months := []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	now := time.Now()
	return fmt.Sprintf("%s %d", months[int(now.Month())-1], now.Year())
}

func formatRub(v float64) string {
	n := int(v + 0.5)
	s := strconv.Itoa(n)
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteString("\u00a0")
		}
		b.WriteRune(c)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// moneyNumRe finds numeric cores like "40", "40 000", "1.5". Suffix (к/тыс) is read manually —
// Go RE2 has no lookahead, and \b is ASCII-only (breaks Cyrillic «к»).
var moneyNumRe = regexp.MustCompile(`\d[\d\s]*(?:[.,]\d+)?`)

func parseMoneyCore(num string) float64 {
	raw := strings.ReplaceAll(strings.ReplaceAll(num, " ", ""), ",", ".")
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

func thousandSuffixLen(rest string) int {
	trimmed := strings.TrimLeft(rest, " \t")
	if trimmed == "" {
		return 0
	}
	spaces := len(rest) - len(trimmed)
	for _, p := range []string{"тысячи", "тысячу", "тысяча", "тысяч", "тыс.", "тыс"} {
		if strings.HasPrefix(trimmed, p) {
			return spaces + len(p)
		}
	}
	r := []rune(trimmed)
	if r[0] == 'к' || r[0] == 'k' || r[0] == 'т' {
		if len(r) == 1 || !unicode.IsLetter(r[1]) {
			return spaces + len(string(r[0]))
		}
	}
	return 0
}

func extractMoney(s string) float64 {
	lower := normalizeRU(strings.ToLower(s))
	if amounts := extractAllMoneyAmounts(lower); len(amounts) > 0 {
		return amounts[0]
	}
	words := map[string]float64{
		"сто": 100, "двести": 200, "триста": 300, "четыреста": 400,
		"пятьсот": 500, "шестьсот": 600, "семьсот": 700, "восемьсот": 800, "девятьсот": 900,
		"пятьдесят": 50, "шестьдесят": 60, "семьдесят": 70, "восемьдесят": 80, "девяносто": 90,
	}
	for w, base := range words {
		if strings.Contains(lower, w) {
			v := base
			if strings.Contains(lower, "тыс") {
				v *= 1000
			}
			return v
		}
	}
	return 0
}

func extractMoneyNear(s, key string) float64 {
	lower := normalizeRU(strings.ToLower(s))
	idx := strings.Index(lower, normalizeRU(key))
	if idx < 0 {
		return 0
	}
	return extractMoney(lower[idx:])
}

func extractAllMoneyAmounts(s string) []float64 {
	lower := normalizeRU(strings.ToLower(s))
	idxs := moneyNumRe.FindAllStringIndex(lower, -1)
	var out []float64
	for _, idx := range idxs {
		core := lower[idx[0]:idx[1]]
		v := parseMoneyCore(core)
		if v <= 0 {
			continue
		}
		rest := lower[idx[1]:]
		if thousandSuffixLen(rest) > 0 {
			v *= 1000
		}
		out = append(out, v)
	}
	return out
}

func extractAllNumbers(s string) []float64 {
	return extractAllMoneyAmounts(s)
}

func hasExplicitRublesNear(s, key string) bool {
	lower := normalizeRU(strings.ToLower(s))
	idx := strings.Index(lower, normalizeRU(key))
	if idx < 0 {
		return false
	}
	window := lower[idx:]
	if len(window) > 48 {
		window = window[:48]
	}
	return strings.Contains(window, "руб") || strings.Contains(window, "₽")
}

func clientsWord(n float64) string {
	abs := int(math.Abs(n)+0.5) % 100
	d := abs % 10
	if abs > 10 && abs < 20 {
		return "клиентов"
	}
	if d == 1 {
		return "клиент"
	}
	if d >= 2 && d <= 4 {
		return "клиента"
	}
	return "клиентов"
}

func extractPercent(s string) float64 {
	re := regexp.MustCompile(`(\d+)\s*%`)
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0
	}
	v, _ := strconv.ParseFloat(m[1], 64)
	return v
}

func extractINN(s string) string {
	re := regexp.MustCompile(`\b\d{10}\b|\b\d{12}\b`)
	return re.FindString(s)
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
