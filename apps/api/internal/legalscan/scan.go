package legalscan

import (
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/domain"
)

// Result of heuristic document review.
type Result struct {
	DocumentName string
	Mode         string
	Flags        []domain.LegalFlag
	Excerpt      string
	Disclaimer   string
}

// ScanText finds risk patterns in extracted contract text.
func ScanText(name, text string) Result {
	if name == "" {
		name = "dogovor.pdf"
	}
	lower := strings.ToLower(strings.ReplaceAll(text, "ё", "е"))
	var flags []domain.LegalFlag

	type rule struct {
		keys   []string
		flag   domain.LegalFlag
		quoteK string
	}
	rules := []rule{
		{
			keys: []string{"односторонн", "изменить стоимость", "изменить цену", "индексац"},
			flag: domain.LegalFlag{
				Severity: "high", Title: "Одностороннее изменение цены",
				Recommendation: "Пропишите срок уведомления (например 30 дней) и право расторгнуть договор без штрафа, если цена выросла.",
			},
			quoteK: "односторон",
		},
		{
			keys: []string{"штраф", "неустойк", "расторжен", "досрочн"},
			flag: domain.LegalFlag{
				Severity: "medium", Title: "Штраф / неустойка при выходе",
				Recommendation: "Согласуйте соразмерную неустойку или период предупреждения (1–2 месяца), а не полный остаток срока.",
			},
			quoteK: "штраф",
		},
		{
			keys: []string{"подсудност", "арбитраж", "по месту нахождения арендодател"},
			flag: domain.LegalFlag{
				Severity: "low", Title: "Подсудность",
				Recommendation: "Можно предложить суд по месту вашей регистрации или по месту кабинета — так проще при конфликте.",
			},
			quoteK: "подсудн",
		},
		{
			keys: []string{"без уведомлен", "не уведомля", "в любое время"},
			flag: domain.LegalFlag{
				Severity: "medium", Title: "Слабые гарантии уведомления",
				Recommendation: "Зафиксируйте письменное уведомление и срок, иначе условия могут меняться внезапно.",
			},
			quoteK: "уведомл",
		},
		{
			keys: []string{"обеспечительн", "залог", "депозит"},
			flag: domain.LegalFlag{
				Severity: "low", Title: "Депозит / обеспечительный платёж",
				Recommendation: "Уточните условия возврата депозита и зачёт в последний месяц аренды.",
			},
			quoteK: "депозит",
		},
	}

	for _, r := range rules {
		hit := false
		for _, k := range r.keys {
			if strings.Contains(lower, k) {
				hit = true
				break
			}
		}
		if !hit {
			continue
		}
		f := r.flag
		f.Quote = quoteAround(text, r.quoteK)
		if f.Quote == "" {
			f.Quote = "…фрагмент по теме «" + f.Title + "»…"
		}
		flags = append(flags, f)
		if len(flags) >= 4 {
			break
		}
	}

	mode := "heuristic_text"
	if strings.TrimSpace(text) == "" {
		mode = "heuristic_fallback"
		flags = defaultFlags()
	} else if len(flags) == 0 {
		flags = []domain.LegalFlag{{
			Severity: "low", Title: "Явных красных флагов по шаблону не найдено",
			Quote:          truncate(text, 160),
			Recommendation: "Всё равно проверьте сроки, оплату и ответственность сторон с юристом перед подписью.",
		}}
	}

	return Result{
		DocumentName: name,
		Mode:         mode,
		Flags:        flags,
		Excerpt:      truncate(text, 12000),
		Disclaimer:   "Предварительный разбор документа, не юридическое заключение",
	}
}

func defaultFlags() []domain.LegalFlag {
	return []domain.LegalFlag{
		{
			Severity: "high", Title: "Одностороннее изменение цены",
			Quote: "Арендодатель вправе в одностороннем порядке изменить стоимость аренды…",
			Recommendation: "Пропишите срок уведомления (например 30 дней) и право расторгнуть договор без штрафа, если цена выросла.",
		},
		{
			Severity: "medium", Title: "Штраф за досрочный выход",
			Quote: "При расторжении по инициативе арендатора взыскивается штраф…",
			Recommendation: "Согласуйте соразмерную неустойку или период предупреждения (1–2 месяца).",
		},
		{
			Severity: "low", Title: "Подсудность",
			Quote: "Споры рассматриваются по месту нахождения арендодателя.",
			Recommendation: "Можно предложить суд по месту вашей регистрации или кабинета.",
		},
	}
}

func quoteAround(text, key string) string {
	lower := strings.ToLower(strings.ReplaceAll(text, "ё", "е"))
	key = strings.ToLower(strings.ReplaceAll(key, "ё", "е"))
	i := strings.Index(lower, key)
	if i < 0 {
		return ""
	}
	start := i - 40
	if start < 0 {
		start = 0
	}
	end := i + 120
	if end > len(text) {
		end = len(text)
	}
	q := strings.TrimSpace(text[start:end])
	if start > 0 {
		q = "…" + q
	}
	if end < len(text) {
		q += "…"
	}
	return q
}

func truncate(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n]) + "…"
}
