package calc

import "math"

type Regime string

const (
	RegimeNPD   Regime = "NPD"
	RegimeUSN6  Regime = "USN_6"
	RegimeUSN15 Regime = "USN_15"
	RegimePSN   Regime = "PSN"
)

type TaxInput struct {
	Regime Regime
	Income float64
	// B2CShare — доля дохода от физлиц для НПД (0..1). 0 = default 0.85 (beauty).
	B2CShare float64
	// ConservativeCeiling: для Home «отложить с запасом» — всегда 6% от дохода.
	ConservativeCeiling bool
	Expenses            float64
	Period              string
}

type TaxResult struct {
	TaxAmount      float64  `json:"tax_amount"`
	Rate           float64  `json:"rate"`
	RateLabel      string   `json:"rate_label"`
	Income         float64  `json:"income"`
	Expenses       float64  `json:"expenses"`
	Period         string   `json:"period"`
	Regime         string   `json:"regime"`
	B2CShare       float64  `json:"b2c_share,omitempty"`
	B2BShare       float64  `json:"b2b_share,omitempty"`
	Assumptions    []string `json:"assumptions,omitempty"`
	CeilingAmount  float64  `json:"ceiling_amount,omitempty"`
	SourceRef      string   `json:"source_ref"`
}

func CalculateTax(in TaxInput) TaxResult {
	income := math.Max(0, in.Income)
	expenses := math.Max(0, in.Expenses)
	assumptions := []string{}

	switch in.Regime {
	case RegimeNPD:
		return calcNPD(income, in.B2CShare, in.ConservativeCeiling, in.Period, assumptions)
	case RegimeUSN6:
		amount := round2(income * 0.06)
		return TaxResult{
			TaxAmount: amount, Rate: 0.06, RateLabel: "6%", Income: income, Expenses: expenses,
			Period: in.Period, Regime: string(in.Regime), SourceRef: "calc:usn6-demo",
			Assumptions: []string{"УСН «доходы» 6% без учёта страховых взносов и вычетов"},
		}
	case RegimeUSN15:
		base := income - expenses
		if base < 0 {
			base = 0
		}
		amount := round2(base * 0.15)
		return TaxResult{
			TaxAmount: amount, Rate: 0.15, RateLabel: "15%", Income: income, Expenses: expenses,
			Period: in.Period, Regime: string(in.Regime), SourceRef: "calc:usn15-demo",
			Assumptions: []string{"УСН «доходы минус расходы» 15% по введённым данным"},
		}
	case RegimePSN:
		amount := round2(income * 0.06)
		return TaxResult{
			TaxAmount: amount, Rate: 0.06, RateLabel: "патент (оценка)", Income: income, Expenses: expenses,
			Period: in.Period, Regime: string(in.Regime), SourceRef: "calc:psn-approx",
			Assumptions: []string{"Патент упрощённо оценён как 6% от оборота — ориентир, не патентный расчёт"},
		}
	default:
		return calcNPD(income, in.B2CShare, in.ConservativeCeiling, in.Period, assumptions)
	}
}

func calcNPD(income, b2cShare float64, ceiling bool, period string, assumptions []string) TaxResult {
	share := b2cShare
	if share <= 0 {
		share = 0.85
		assumptions = append(assumptions, "Доля физлиц по умолчанию 85% (типично для beauty)")
	}
	if share > 1 {
		share = 1
	}
	b2b := 1 - share
	mixedRate := share*0.04 + b2b*0.06
	mixed := round2(income * mixedRate)
	ceilingAmt := round2(income * 0.06)

	assumptions = append(assumptions,
		"НПД: 4% с физлиц, 6% с юрлиц/ИП (упрощённая модель без региональных льгот)",
		"Без учёта вычета страховых и региональных особенностей",
	)

	amount := mixed
	rateLabel := formatPct(mixedRate) + " смешанная"
	rate := mixedRate
	if ceiling {
		amount = ceilingAmt
		rate = 0.06
		rateLabel = "6% с запасом"
		assumptions = append(assumptions, "Показана консервативная оценка «отложить с запасом» по верхней ставке 6%")
	}

	return TaxResult{
		TaxAmount:     amount,
		Rate:          rate,
		RateLabel:     rateLabel,
		Income:        income,
		Period:        period,
		Regime:        string(RegimeNPD),
		B2CShare:      share,
		B2BShare:      b2b,
		Assumptions:   assumptions,
		CeilingAmount: ceilingAmt,
		SourceRef:     "kb:npd-rates-demo#2026",
	}
}

func formatPct(rate float64) string {
	return formatFloat(rate*100) + "%"
}

func formatFloat(v float64) string {
	if math.Abs(v-float64(int(math.Round(v)))) < 1e-9 {
		return itoa(int(math.Round(v)))
	}
	return formatNumber(v, 1)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func formatNumber(v float64, prec int) string {
	neg := v < 0
	if neg {
		v = -v
	}
	pow := math.Pow10(prec)
	iv := int(math.Round(v * pow))
	whole := iv / int(pow)
	frac := iv % int(pow)
	s := itoa(whole) + "."
	fs := itoa(frac)
	for len(fs) < prec {
		fs = "0" + fs
	}
	if neg {
		return "-" + s + fs
	}
	return s + fs
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
