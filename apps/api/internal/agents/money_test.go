package agents

import (
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/sdui"
)

func TestExtractMoneyKSuffix(t *testing.T) {
	cases := map[string]float64{
		"аренда кабинета 40к":        40000,
		"аренда 40к, расходники 200": 40000,
		"оборот 180 тысяч":           180000,
		"180 тыс":                    180000,
		"двести тысяч":               200000,
		"цена услуги 1500":           1500,
		"аренда 40 000":              40000,
	}
	for in, want := range cases {
		got := extractMoney(in)
		if got != want {
			t.Fatalf("%q: got %v want %v", in, got, want)
		}
	}
}

func TestExtractAllMoneyAmountsUnitEconPrompt(t *testing.T) {
	msg := "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500."
	nums := extractAllMoneyAmounts(msg)
	if len(nums) < 3 {
		t.Fatalf("want ≥3 amounts, got %v", nums)
	}
	if nums[0] != 40000 || nums[1] != 200 || nums[2] != 1500 {
		t.Fatalf("got %v want [40000 200 1500]", nums)
	}
}

func TestUnitEconUses40kNot40(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded(), Offline: true}
	evs := o.Handle("Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.")
	text := collectTokenText(evs)
	if containsAny(text, "фикс 40 ", "40 ₽/мес", "Фикс (аренда) 40 ") {
		t.Fatalf("still treating rent as 40 rubles: %s", text)
	}

	var fixed, be float64
	for _, e := range evs {
		if e.Event != "sdui" {
			continue
		}
		env, ok := e.Data.(sdui.Envelope)
		if !ok || env.Component != "UnitEconomicsChart" {
			continue
		}
		fixed, _ = env.Props["fixed_costs_monthly"].(float64)
		be, _ = env.Props["breakeven_units_per_day"].(float64)
	}
	if fixed != 40000 {
		t.Fatalf("fixed_costs_monthly=%v want 40000; text=%s", fixed, text)
	}
	if be < 2 {
		t.Fatalf("breakeven_units_per_day=%v want ≥2 with realistic model", be)
	}
	if collectTokenText(evs) == "" {
		t.Fatal("expected narrative tokens")
	}
}
