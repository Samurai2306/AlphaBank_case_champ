package tools

import (
	"time"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/mockbank"
	"github.com/alphabank-case-champ/copilot-api/internal/calc"
	"github.com/alphabank-case-champ/copilot-api/internal/domain"
	"github.com/alphabank-case-champ/copilot-api/internal/legalscan"
	"github.com/alphabank-case-champ/copilot-api/internal/rag"
)

// Backend is the per-user data surface used by tool adapters.
type Backend interface {
	Profile() domain.Profile
	UpdateProfile(domain.Profile)
	Transactions() []domain.Transaction
	CreateDraft(domain.PaymentDraft) domain.PaymentDraft
	IncomeThisMonth() float64
	Piggy() domain.EnpPiggy
	SetPiggy(domain.EnpPiggy)
}

type Registry struct {
	Store Backend
}

func (r *Registry) GetClientTransactions() []domain.Transaction {
	return r.Store.Transactions()
}

func (r *Registry) CheckCounterpartyRisk(inn string) domain.RiskResult {
	return mockbank.CheckCounterpartyRisk(inn)
}

// PaymentDraftSpec describes a tax (УФК) or counterparty payment draft.
type PaymentDraftSpec struct {
	Amount    float64
	Purpose   string
	PayeeName string
	PayeeINN  string
	RiskLevel string
}

func (r *Registry) CreatePaymentDraft(amount float64, purpose string) domain.PaymentDraft {
	return r.CreatePaymentDraftSpec(PaymentDraftSpec{
		Amount: amount, Purpose: purpose,
		PayeeName: "УФК по Республике Татарстан",
		RiskLevel: "green",
	})
}

func (r *Registry) CreatePaymentDraftSpec(spec PaymentDraftSpec) domain.PaymentDraft {
	if spec.PayeeName == "" {
		spec.PayeeName = "УФК по Республике Татарстан"
	}
	if spec.RiskLevel == "" {
		spec.RiskLevel = "green"
	}
	if spec.Purpose == "" {
		spec.Purpose = "Платёж"
	}
	return r.Store.CreateDraft(domain.PaymentDraft{
		Amount:    spec.Amount,
		Purpose:   spec.Purpose,
		PayeeName: spec.PayeeName,
		PayeeINN:  spec.PayeeINN,
		RiskLevel: spec.RiskLevel,
	})
}

func (r *Registry) FetchFNSDebt() map[string]any {
	p := r.Store.Profile()
	tax := calc.CalculateTax(calc.TaxInput{
		Regime: calc.RegimeNPD, Income: p.MonthlyRevenueEstimate, Period: time.Now().Format("2006-01"),
	})
	return map[string]any{
		"status": "mock", "has_debt": false, "estimated_current_period": tax.TaxAmount,
		"currency": "RUB", "as_of": time.Now().Format(time.RFC3339),
		"note": "Задолженностей не найдено; сумма текущего периода — оценка сервиса",
	}
}

func (r *Registry) CalculateTax(regime calc.Regime, income float64, period string, ceiling bool) calc.TaxResult {
	return calc.CalculateTax(calc.TaxInput{
		Regime: regime, Income: income, Period: period, ConservativeCeiling: ceiling,
	})
}

func (r *Registry) ComputeUnitEconomics(fixed, variable, price float64) calc.UnitEconomicsResult {
	// Beauty-default NPD mix: 85% B2C @4% + 15% B2B @6% ≈ 4.3%.
	taxRate := 0.85*0.04 + 0.15*0.06
	return calc.ComputeUnitEconomics(calc.UnitEconomicsInput{
		FixedCostsMonthly:   fixed,
		VariableCostPerUnit: variable,
		PricePerUnit:        price,
		WorkingDaysPerMonth: 22,
		TaxRate:             taxRate,
		SafetyBuffer:        0.15,
	})
}

func (r *Registry) ScanLegalDocument(name, text string) map[string]any {
	res := legalscan.ScanText(name, text)
	return map[string]any{
		"document_name": res.DocumentName,
		"mode":          res.Mode,
		"flags":         res.Flags,
		"excerpt":       res.Excerpt,
		"disclaimer":    res.Disclaimer,
	}
}

func (r *Registry) UpdateOnboardingProfile(sphere, city string, revenue float64) domain.Profile {
	p := r.Store.Profile()
	if sphere != "" {
		p.BusinessSphere = sphere
	}
	if city != "" {
		p.City = city
	}
	if revenue > 0 {
		p.MonthlyRevenueEstimate = revenue
	}
	if p.CJMLevel < 2 {
		p.CJMLevel = 2
	}
	r.Store.UpdateProfile(p)
	return p
}

// RetrieveKB returns ranked knowledge chunks for the query (demo in-memory RAG).
func (r *Registry) RetrieveKB(query string, topK int) []rag.Hit {
	if topK <= 0 {
		topK = 4
	}
	return rag.Default().Retrieve(query, topK)
}
