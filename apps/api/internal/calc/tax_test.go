package calc

import "testing"

func TestCalculateTaxNPDMixed(t *testing.T) {
	// 85% B2C default: 180k * (0.85*0.04 + 0.15*0.06) = 180k * 0.043 = 7740
	r := CalculateTax(TaxInput{Regime: RegimeNPD, Income: 180000, Period: "Июль 2026"})
	if r.TaxAmount != 7740 {
		t.Fatalf("want 7740 got %v", r.TaxAmount)
	}
	if r.CeilingAmount != 10800 {
		t.Fatalf("ceiling want 10800 got %v", r.CeilingAmount)
	}
}

func TestCalculateTaxNPDCeiling(t *testing.T) {
	r := CalculateTax(TaxInput{Regime: RegimeNPD, Income: 180000, ConservativeCeiling: true})
	if r.TaxAmount != 10800 {
		t.Fatalf("want 10800 got %v", r.TaxAmount)
	}
}

func TestCalculateTaxNPDPureB2C(t *testing.T) {
	r := CalculateTax(TaxInput{Regime: RegimeNPD, Income: 100000, B2CShare: 1})
	if r.TaxAmount != 4000 {
		t.Fatalf("want 4000 got %v", r.TaxAmount)
	}
}

func TestCalculateTaxUSN6(t *testing.T) {
	r := CalculateTax(TaxInput{Regime: RegimeUSN6, Income: 100000, Period: "Q3"})
	if r.TaxAmount != 6000 {
		t.Fatalf("want 6000 got %v", r.TaxAmount)
	}
}

func TestCalculateTaxUSN15(t *testing.T) {
	r := CalculateTax(TaxInput{Regime: RegimeUSN15, Income: 200000, Expenses: 80000, Period: "Q3"})
	if r.TaxAmount != 18000 {
		t.Fatalf("want 18000 got %v", r.TaxAmount)
	}
}

func TestUnitEconomics(t *testing.T) {
	taxRate := 0.85*0.04 + 0.15*0.06
	r := ComputeUnitEconomics(UnitEconomicsInput{
		FixedCostsMonthly:   40000,
		VariableCostPerUnit: 200,
		PricePerUnit:        1500,
		WorkingDaysPerMonth: 22,
		TaxRate:             taxRate,
		SafetyBuffer:        0.15,
	})
	if r.BreakevenUnitsPerDay < 2 {
		t.Fatalf("breakeven too low (want ≥2 with tax+22 days): %v", r.BreakevenUnitsPerDay)
	}
	if r.Series[0].Profit != -40000 {
		t.Fatalf("at 0 clients want −40000, got %v", r.Series[0].Profit)
	}
	if r.TaxPerUnit <= 0 {
		t.Fatal("expected tax per unit")
	}
	if r.MarginPerUnit >= 1300 {
		t.Fatalf("net margin should be below gross 1300 after tax, got %v", r.MarginPerUnit)
	}
	if r.RecommendedUnitsPerDay < r.BreakevenUnitsPerDay {
		t.Fatalf("recommended %v < breakeven %v", r.RecommendedUnitsPerDay, r.BreakevenUnitsPerDay)
	}
	for i := 1; i < len(r.Series); i++ {
		if r.Series[i].Profit < r.Series[i-1].Profit {
			t.Fatalf("series not increasing at %d: %v -> %v", i, r.Series[i-1], r.Series[i])
		}
	}
}
