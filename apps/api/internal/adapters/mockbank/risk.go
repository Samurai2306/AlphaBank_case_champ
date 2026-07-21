package mockbank

import "github.com/alphabank-case-champ/copilot-api/internal/domain"

func CheckCounterpartyRisk(inn string) domain.RiskResult {
	switch inn {
	case "7707083893":
		return domain.RiskResult{
			INN: inn, Level: "yellow", Title: "Есть замечания — лучше уточнить",
			Reasons: []string{
				"Недавно сменён директор",
				"Короткий срок регистрации компании",
			},
			Recommendation: "Перед оплатой запросите счёт, акт/УПД и реквизиты. Для крупной суммы — разбейте платёж или проверьте договор.",
		}
	case "0000000000":
		return domain.RiskResult{
			INN: inn, Level: "red", Title: "Высокий риск — платёж лучше отложить",
			Reasons:        []string{"Контрагент в стоп-листе проверки"},
			Recommendation: "Не отправляйте деньги без ручной проверки. В Альфа-Банке можно уточнить статус у менеджера или через поддержку РКО.",
		}
	default:
		return domain.RiskResult{
			INN: inn, Level: "green", Title: "Рисков не найдено",
			Reasons: []string{
				"ИНН распознан, экспресс-проверка пройдена",
				"Признаков блокировки по 115-ФЗ не найдено",
			},
			Recommendation: "Можно готовить платёжку. Перед крупной арендой или закупкой всё равно сохраните закрывающие документы.",
		}
	}
}
