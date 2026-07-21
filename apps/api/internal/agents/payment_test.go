package agents

import (
	"strings"
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/sdui"
)

func TestPaymentTaxDraftUsesUFK(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded(), Offline: true}
	evs := o.Handle("Сформируй платёжку.")
	var payee string
	for _, e := range evs {
		if e.Event != "sdui" {
			continue
		}
		env, ok := e.Data.(sdui.Envelope)
		if !ok || env.Component != "PaymentDraftCard" {
			continue
		}
		payee, _ = env.Props["payee_name"].(string)
	}
	if !strings.Contains(payee, "УФК") {
		t.Fatalf("tax draft payee=%q, want УФК", payee)
	}
}

func TestPaymentRentDraftUsesCounterparty(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded(), Offline: true}
	evs := o.Handle("Сформируй платёжку аренды 40000 на ИНН 1650987654")
	var payee, inn, purpose, risk string
	var hasCompliance bool
	for _, e := range evs {
		if e.Event != "sdui" {
			continue
		}
		env, ok := e.Data.(sdui.Envelope)
		if !ok {
			continue
		}
		switch env.Component {
		case "PaymentDraftCard":
			payee, _ = env.Props["payee_name"].(string)
			inn, _ = env.Props["payee_inn"].(string)
			purpose, _ = env.Props["purpose"].(string)
			risk, _ = env.Props["risk_level"].(string)
		case "ComplianceTrafficLight":
			hasCompliance = true
		}
	}
	if strings.Contains(payee, "УФК") {
		t.Fatalf("rent draft must not use УФК, payee=%q", payee)
	}
	if inn != "1650987654" {
		t.Fatalf("payee_inn=%q", inn)
	}
	if !strings.Contains(strings.ToLower(purpose), "аренд") {
		t.Fatalf("purpose=%q", purpose)
	}
	if risk != "green" {
		t.Fatalf("risk=%q", risk)
	}
	if !hasCompliance {
		t.Fatal("expected ComplianceTrafficLight alongside rent draft")
	}
}

func TestComplianceStoresRiskForRentFollowUp(t *testing.T) {
	store := memory.NewSeeded()
	o := &Orchestrator{Store: store, Offline: true}
	_ = o.Handle("Проверь ИНН 0000000000")
	evs := o.Handle("Сформируй платёжку аренды 25000")
	var risk, inn string
	for _, e := range evs {
		if e.Event != "sdui" {
			continue
		}
		env, ok := e.Data.(sdui.Envelope)
		if !ok || env.Component != "PaymentDraftCard" {
			continue
		}
		risk, _ = env.Props["risk_level"].(string)
		inn, _ = env.Props["payee_inn"].(string)
	}
	if inn != "0000000000" {
		t.Fatalf("expected last risk INN, got %q", inn)
	}
	if risk != "red" {
		t.Fatalf("expected red risk on draft, got %q", risk)
	}
}

func TestUnitEconEmitsProductOffers(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded(), Offline: true}
	evs := o.Handle("Точка безубыточности при аренде 40к")
	var hasChart, hasOffers bool
	for _, e := range evs {
		if e.Event != "sdui" {
			continue
		}
		env, ok := e.Data.(sdui.Envelope)
		if !ok {
			continue
		}
		if env.Component == "UnitEconomicsChart" {
			hasChart = true
		}
		if env.Component == "ProductOffers" {
			hasOffers = true
		}
	}
	if !hasChart || !hasOffers {
		t.Fatalf("chart=%v offers=%v", hasChart, hasOffers)
	}
}

func TestWantsCounterpartyPayment(t *testing.T) {
	cases := map[string]bool{
		"сформируй платежку аренды 40000": true,
		"сформируй платежку":              false,
		"платежка на налог нпд":           false,
		"оплата арендодателю":             true,
	}
	for msg, want := range cases {
		if got := wantsCounterpartyPayment(normalizeRU(msg)); got != want {
			t.Fatalf("%q: got %v want %v", msg, got, want)
		}
	}
}
