package guard

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	MaxMessageRunes = 4000
	MaxDraftAmount  = 5_000_000.0 // demo sanity cap
	MinDraftAmount  = 1.0
)

// Blocked is a policy refusal for unsafe user intent.
type Blocked struct {
	Code    string
	Message string
}

var ingressPatterns = []struct {
	re   *regexp.Regexp
	code string
	msg  string
}{
	{re: must(`уклон(ить|ять|ени)|не\s*платить\s*налог|скрыть\s*(доход|выруч)|серая\s*схем|чёрн(ый|ая)\s*касс`),
		code: "GUARDRAIL_TAX_EVASION",
		msg:  "Не могу подсказать, как уклоняться от налогов или скрывать доходы. Могу помочь легально оформить НПД/ИП и посчитать сумму к уплате."},
	{re: must(`обналич|обнал|отмы(ть|вани)|дроп\s*сч[её]т`),
		code: "GUARDRAIL_MONEY_LAUNDERING",
		msg:  "Не помогаю с обналичиванием и схемами вывода средств. Могу помочь с легальным РКО, чеками НПД и проверкой контрагента."},
	{re: must(`обойти\s*115|обход\s*115|заб(ыть|удь|ывайте)\s*(все\s*)?инструк|ignore\s*(all\s*)?(previous|prior)\s*instructions|jailbreak|developer\s*mode|you\s*are\s*now\s*dan`),
		code: "GUARDRAIL_JAILBREAK",
		msg:  "Не могу игнорировать инструкции безопасности или помогать обходить 115-ФЗ. Могу проверить контрагента по ИНН и подсказать легальный порядок платежей."},
	{re: must(`поддела(ть|й)|фальшив(ый|ую)\s*(чек|документ)|липов(ый|ую)`),
		code: "GUARDRAIL_FRAUD",
		msg:  "Не помогаю подделывать документы и чеки. Могу показать, как корректно фиксировать доход в «Мой налог»."},
	{re: must(`взлома(ть|й)|укра(сть|ди)\s*(парол|ключ|токен|api)|sql\s*injection|prompt\s*injection`),
		code: "GUARDRAIL_ATTACK",
		msg:  "Запрос похож на попытку атаки или обхода защиты. Задайте вопрос про налоги, платежи или проверку контрагента."},
}

var egressBad = must(`уклон(ить|ять)|не\s*платите?\s*налог|скройте?\s*доход|обналич|обойдите?\s*115|заб(ыть|удь|ывайте|удьте)\s*инструк|игнорируйте?\s*(закон|фнс|банк)`)

var amountRe = regexp.MustCompile(`\d[\d\s]{0,12}`)

func must(p string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)` + p)
}

// CheckIngress blocks unsafe user messages before routing.
func CheckIngress(message string) *Blocked {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return &Blocked{Code: "VALIDATION_ERROR", Message: "Пустое сообщение. Напишите вопрос по бизнесу или налогам."}
	}
	if len([]rune(msg)) > MaxMessageRunes {
		return &Blocked{Code: "VALIDATION_ERROR", Message: "Сообщение слишком длинное. Сократите текст до сути вопроса."}
	}
	norm := normalize(msg)
	for _, p := range ingressPatterns {
		if p.re.MatchString(norm) {
			return &Blocked{Code: p.code, Message: p.msg}
		}
	}
	return nil
}

// CheckEgress rejects model text that violates policy (hallucinated unsafe advice).
func CheckEgress(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}
	return !egressBad.MatchString(normalize(text))
}

// SoftNumbersPreserved is looser: only large money-like amounts (≥4 digits) and INNs
// must survive. Used when we prefer informative LLM rewriting over strict determinism.
func SoftNumbersPreserved(source, polished string) bool {
	srcNums := significantNumbers(source)
	if len(srcNums) == 0 {
		return true
	}
	pol := compactDigits(polished)
	critical := 0
	kept := 0
	for _, n := range srcNums {
		if len(n) < 4 && len(n) != 10 && len(n) != 12 {
			continue
		}
		critical++
		if strings.Contains(pol, n) {
			kept++
		}
	}
	if critical == 0 {
		return true
	}
	// Keep at least one third of critical figures (informativeness > rigidity).
	return kept*3 >= critical
}

// SanitizeDraftAmount clamps/rejects absurd payment amounts in demo.
func SanitizeDraftAmount(amount float64) (float64, bool) {
	if amount < MinDraftAmount || amount > MaxDraftAmount {
		return 0, false
	}
	return amount, true
}

func significantNumbers(s string) []string {
	raw := amountRe.FindAllString(s, -1)
	seen := map[string]bool{}
	var out []string
	for _, r := range raw {
		n := compactDigits(r)
		if len(n) < 2 { // skip single digits / noise
			continue
		}
		// Prefer money-like and INN-like lengths.
		if len(n) >= 3 || len(n) == 10 || len(n) == 12 {
			if !seen[n] {
				seen[n] = true
				out = append(out, n)
			}
		}
	}
	return out
}

func compactDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.ReplaceAll(s, "ё", "е")
}
