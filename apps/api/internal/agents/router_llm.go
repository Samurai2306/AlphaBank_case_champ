package agents

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/domain"
)

var validIntents = map[string]bool{
	"TAX_CALC": true, "TRANSACTION": true, "LEGAL_REVIEW": true,
	"ONBOARDING": true, "UNIT_ECON": true, "PIGGY": true,
	"COMPLIANCE": true, "GENERAL_QA": true,
}

type llmRoute struct {
	Intent     string         `json:"intent"`
	Confidence float64        `json:"confidence"`
	Entities   map[string]any `json:"extracted_entities"`
}

// resolveIntent prefers LLM understanding when online; keywords are a fallback.
func (o *Orchestrator) resolveIntent(ctx context.Context, msg, lower string, history []string) string {
	kw := route(lower)
	if o.Offline || o.LLM == nil {
		return kw
	}

	// Hard override only for ultra-explicit pitch commands.
	if explicitCommand(lower) {
		return kw
	}
	// Jury/demo phrasing is already unambiguous — skip an extra LLM hop.
	if kw != "GENERAL_QA" && keywordConfident(lower, kw) {
		return kw
	}

	intent, conf := o.classifyLLM(ctx, msg, history)
	if !validIntents[intent] {
		return kw
	}

	// Prefer model understanding for freer phrasing.
	if conf >= 0.42 {
		// If keyword says TAX but user asked «что такое / отличие» — trust LLM/GENERAL.
		if kw == "TAX_CALC" && intent == "GENERAL_QA" {
			return intent
		}
		return intent
	}

	// Low LLM confidence → keyword, then GENERAL_QA as soft default for chatty asks.
	if kw != "GENERAL_QA" {
		return kw
	}
	return "GENERAL_QA"
}

func keywordConfident(lower, kw string) bool {
	switch kw {
	case "TAX_CALC":
		return containsAny(lower, "сколько отложить", "налог за этот месяц", "налог за месяц", "посчитай налог")
	case "TRANSACTION":
		return containsAny(lower, "сформируй плат")
	case "PIGGY":
		return containsAny(lower, "включи копил", "выключи копил")
	case "COMPLIANCE":
		return containsAny(lower, "проверь инн")
	case "UNIT_ECON":
		return containsAny(lower, "клиентов в день", "безубыточ")
	case "LEGAL_REVIEW":
		return containsAny(lower, "разбери договор")
	default:
		return false
	}
}

func explicitCommand(lower string) bool {
	return containsAny(lower,
		"сформируй платёжк", "сформируй платежк", "сформируй плат",
		"включи копил", "выключи копил",
		"проверь инн",
		"разбери договор",
	)
}

func (o *Orchestrator) classifyLLM(ctx context.Context, msg string, history []string) (string, float64) {
	system := `Ты маршрутизатор продукта «Альфа-Бизнес: Старт».
Понимай смысл свободной русской речи. Верни ТОЛЬКО JSON без markdown:
{"intent":"TAX_CALC|TRANSACTION|LEGAL_REVIEW|ONBOARDING|UNIT_ECON|PIGGY|COMPLIANCE|GENERAL_QA","confidence":0.0,"extracted_entities":{}}

Правила выбора intent:
- TAX_CALC: пользователь хочет ИМЕННО посчитать/оценить сумму налога к уплате или «сколько отложить»
- TRANSACTION: сформировать платёжку/поручение сейчас (налог ЕНП ИЛИ аренда/контрагент)
- LEGAL_REVIEW: разобрать договор/PDF/условия аренды/штрафы в договоре (не сама платёжка)
- UNIT_ECON: сколько клиентов в день, безубыточность, маржа услуги при аренде/цене
- PIGGY: копилка ЕНП (включить/ставку/накопить)
- COMPLIANCE: проверить ИНН/контрагента/риски 115-ФЗ (без «сформируй платёжку»)
- ONBOARDING: новый бизнес, собрать/уточнить профиль
- GENERAL_QA: объяснения, сравнения, тарифы, СБП/эквайринг/РКО, FAQ, советы, «что такое»

«Сформируй платёжку аренды…» → TRANSACTION. «Как принимать оплату по СБП» → GENERAL_QA.

Если вопрос объяснительный («что такое», «чем отличается», «расскажи») — GENERAL_QA, даже если есть слова «налог/НПД».
Если сомневаешься между калькулятором и объяснением — GENERAL_QA (информативность важнее).
В extracted_entities клади числа если есть: rent, variable_cost, price, income, inn, piggy_rate.`
	var b strings.Builder
	if len(history) > 0 {
		b.WriteString("Recent dialogue:\n")
		for _, h := range history {
			b.WriteString(h)
			b.WriteByte('\n')
		}
	}
	b.WriteString("User message:\n")
	b.WriteString(msg)
	raw, err := o.LLM.CompleteBudget(ctx, system, b.String(), 0.2, 220)
	if err != nil || raw == "" {
		return "", 0
	}
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			raw = raw[i : j+1]
		}
	}
	var out llmRoute
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return "", 0
	}
	out.Intent = strings.ToUpper(strings.TrimSpace(out.Intent))
	if !validIntents[out.Intent] {
		return "", 0
	}
	if out.Confidence <= 0 {
		out.Confidence = 0.55
	}
	return out.Intent, out.Confidence
}

func historyLines(turns []domain.ChatMessage) []string {
	out := make([]string, 0, len(turns))
	for _, t := range turns {
		role := t.Role
		if role == "" {
			role = "user"
		}
		out = append(out, role+": "+t.Content)
	}
	return out
}
