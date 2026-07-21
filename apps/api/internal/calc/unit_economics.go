package calc

import "math"

type UnitEconomicsInput struct {
	FixedCostsMonthly   float64
	VariableCostPerUnit float64
	PricePerUnit        float64
	WorkingDaysPerMonth int
	// TaxRate — доля налога с выручки на 1 клиента (для НПД beauty ~0.043).
	TaxRate float64
	// SafetyBuffer — запас на отмены/простои (0.15 = +15% к точке безубыточности).
	SafetyBuffer float64
}

type UnitEconomicsResult struct {
	BreakevenUnitsPerDay   float64       `json:"breakeven_units_per_day"`
	RecommendedUnitsPerDay float64       `json:"recommended_units_per_day"`
	MarginPerUnit          float64       `json:"margin_per_unit"`
	GrossMarginPerUnit     float64       `json:"gross_margin_per_unit"`
	TaxPerUnit             float64       `json:"tax_per_unit"`
	TaxRate                float64       `json:"tax_rate"`
	PricePerUnit           float64       `json:"price_per_unit"`
	VariableCostPerUnit    float64       `json:"variable_cost_per_unit"`
	FixedCostsMonthly      float64       `json:"fixed_costs_monthly"`
	WorkingDaysPerMonth    int           `json:"working_days_per_month"`
	ClientsPerMonth        float64       `json:"clients_per_month"`
	Series                 []ProfitPoint `json:"series"`
	Assumptions            []string      `json:"assumptions,omitempty"`
}

type ProfitPoint struct {
	Units  float64 `json:"units"`
	Profit float64 `json:"profit"`
}

func ComputeUnitEconomics(in UnitEconomicsInput) UnitEconomicsResult {
	days := in.WorkingDaysPerMonth
	if days <= 0 {
		days = 22
	}
	taxRate := in.TaxRate
	if taxRate < 0 {
		taxRate = 0
	}
	buffer := in.SafetyBuffer
	if buffer < 0 {
		buffer = 0
	}

	gross := in.PricePerUnit - in.VariableCostPerUnit
	taxPerUnit := round2(in.PricePerUnit * taxRate)
	margin := round2(gross - taxPerUnit)

	assumptions := []string{
		"22 рабочих дня в месяце (не календарные 30)",
		"Налог НПД оценён с чека (смешанная ставка ~4.3% при 85% B2C)",
		"Без учёта своего времени, рекламы, комиссии эквайринга и кассовых разрывов",
	}

	var beDay, recDay, clientsMonth float64
	if margin > 0 {
		clientsMonth = in.FixedCostsMonthly / margin
		beDay = math.Ceil(clientsMonth / float64(days))
		recDay = math.Ceil(clientsMonth * (1 + buffer) / float64(days))
		if recDay < beDay {
			recDay = beDay
		}
	}

	// Monthly profit = clients/day × days × net margin − fixed. At 0 clients = −fixed.
	profitAt := func(unitsPerDay float64) float64 {
		return round2(unitsPerDay*float64(days)*margin - in.FixedCostsMonthly)
	}

	maxU := recDay + 4
	if maxU < 6 {
		maxU = 6
	}
	if maxU > 20 {
		maxU = 20
	}
	series := make([]ProfitPoint, 0, int(maxU)+1)
	for u := 0.0; u <= maxU; u++ {
		series = append(series, ProfitPoint{Units: u, Profit: profitAt(u)})
	}

	return UnitEconomicsResult{
		BreakevenUnitsPerDay:   beDay,
		RecommendedUnitsPerDay: recDay,
		MarginPerUnit:          margin,
		GrossMarginPerUnit:     round2(gross),
		TaxPerUnit:             taxPerUnit,
		TaxRate:                taxRate,
		PricePerUnit:           in.PricePerUnit,
		VariableCostPerUnit:    in.VariableCostPerUnit,
		FixedCostsMonthly:      in.FixedCostsMonthly,
		WorkingDaysPerMonth:    days,
		ClientsPerMonth:        round2(clientsMonth),
		Series:                 series,
		Assumptions:            assumptions,
	}
}
