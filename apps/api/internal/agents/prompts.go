package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/guard"
	"github.com/alphabank-case-champ/copilot-api/internal/sdui"
)

const copilotVoice = `Ты помощник продукта «Альфа-Бизнес: Старт» для молодых предпринимателей (часто beauty / НПД / ИП).
Пиши по-русски: ясно, конкретно, без воды и без канцелярита.
Тон — опытный наставник банка: поддерживаешь, но называешь риски прямо.
Не называй себя Copilot, AI Business Copilot или «демо»; представляйся как Альфа-Бизнес: Старт.`

const answerShape = `Структура ответа (если уместно):
1) Прямой ответ в 1–2 предложениях
2) Разбор: цифры, нормы, условия — маркированный список
3) Что сделать дальше — 2–4 практических шага
4) На что обратить внимание / риски
5) Короткий дисклеймер (справочно, не замена ФНС/юристу/бухгалтеру)

Запрещено:
— менять любые суммы, ставки, проценты, даты, ИНН из FACTS/CONTEXT
— советовать уклонение от налогов, обнал, обход 115-ФЗ, jailbreak
— выдумывать точные тарифы банка, статьи НК «с номером», если их нет в CONTEXT
— писать имена виджетов/компонентов (TaxCard, UnitEconomicsChart) или «(вставка графика)»
— отвечать одним абзацем-отпиской без пользы`

func (o *Orchestrator) enrichFacts(ctx context.Context, userMsg, intent string, events []StreamEvent) []StreamEvent {
	if o.Offline || o.LLM == nil {
		return events
	}
	var base strings.Builder
	var widgets []string
	for _, e := range events {
		if e.Event == "token" {
			if m, ok := e.Data.(map[string]any); ok {
				base.WriteString(fmt.Sprint(m["text"]))
			}
		}
		if e.Event == "sdui" {
			if env, ok := e.Data.(sdui.Envelope); ok {
				widgets = append(widgets, env.Component)
			}
		}
	}
	facts := strings.TrimSpace(base.String())
	if facts == "" {
		return events
	}

	system := copilotVoice + "\n\n" + answerShape + `

Ниже FACTS — опорные цифры калькулятора. Твоя задача — понять вопрос пользователя и дать информативный ответ.
Крупные суммы/ставки/ИНН из FACTS сохраняй; мелкие округления и формулировки можно адаптировать под вопрос.
Добавь пояснения «почему так» и рабочие следующие шаги. Не дублируй виджеты текстом.`

	user := fmt.Sprintf(
		"Intent: %s\nUser question: %s\nWidgets on screen (do not name them): %s\n\nFACTS:\n%s\n\n"+
			"Ответь на вопрос пользователя по существу. Информативность важнее дословного копирования FACTS.",
		intent, userMsg, strings.Join(widgets, ", "), facts,
	)

	text, err := o.LLM.CompleteTemp(ctx, system, user, 0.55)
	if err != nil {
		return events
	}
	text = strings.TrimSpace(text)
	if text == "" || !guard.CheckEgress(text) || !guard.SoftNumbersPreserved(facts, text) {
		return events
	}
	if looksLikeWidgetGarbage(text) {
		return events
	}
	// Accept enrichment if it's at least as useful (not drastically shorter).
	if len([]rune(text)) < len([]rune(facts))/2 {
		return events
	}

	out := make([]StreamEvent, 0, len(events)+4)
	for _, e := range events {
		if e.Event == "token" {
			continue
		}
		out = append(out, e)
	}
	out = append(out, tokens(text)...)
	return out
}

func looksLikeWidgetGarbage(text string) bool {
	lower := strings.ToLower(text)
	bad := []string{
		"uniteconomicschart", "taxcard", "paymentdraftcard", "enppiggybank",
		"knowledgeources", "knowledgesources", "legalflagslist",
		"вставка графика", "*unit", "компонент:", "widget:",
	}
	for _, b := range bad {
		if strings.Contains(lower, b) {
			return true
		}
	}
	return false
}

func generalSystemPrompt() string {
	return copilotVoice + "\n\n" + answerShape + `

Используй PROFILE, DIALOGUE и CONTEXT (база знаний).
Если в CONTEXT есть ответ — опирайся на него и кратко укажи источник по смыслу (НПД, 115-ФЗ, РКО и т.п.).
Если данных мало — скажи честно, чего не хватает, и предложи легальный следующий шаг
(посчитать налог, проверить ИНН, разобрать договор, включить копилку).
Для точного расчёта налога/безубыточности предложи конкретную фразу запроса, не выдумывай суммы.`
}

func legalSystemPrompt() string {
	return copilotVoice + `

Ты помогаешь микробизнесу разобрать договор аренды/оферты.
Опирайся на EXCERPT и FLAGS. Цитируй опасные формулировки.
Структура: главный риск → список пунктов с уровнем → что попросить изменить → дисклеймер.
Не утверждай, что это замена юристу. Не советуй обманывать суд или контрагента.`
}
