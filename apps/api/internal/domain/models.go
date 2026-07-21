package domain

import "time"

type Profile struct {
	UserID                 string   `json:"user_id"`
	DisplayName            string   `json:"display_name"`
	FullName               string   `json:"full_name"`
	Phone                  string   `json:"phone"`
	Email                  string   `json:"email"`
	INN                    string   `json:"inn"`
	BusinessSphere         string   `json:"business_sphere"`
	City                   string   `json:"city"`
	Address                string   `json:"address"`
	MonthlyRevenueEstimate float64  `json:"monthly_revenue_estimate"`
	TaxRegime              string   `json:"tax_regime"`
	OKVED                  []string `json:"okved_codes"`
	CJMLevel               int      `json:"cjm_level"`
	PersonaKey             string   `json:"persona_key"`
	AccountMasked          string   `json:"account_masked"`
	BankName               string   `json:"bank_name"`
	RegisteredAt           string   `json:"registered_at"`
	NPDStatus              string   `json:"npd_status"`
}

type Transaction struct {
	ID               string    `json:"id"`
	Amount           float64   `json:"amount"`
	Direction        string    `json:"direction"`
	CounterpartyName string    `json:"counterparty_name"`
	CounterpartyINN  string    `json:"counterparty_inn,omitempty"`
	BookedAt         time.Time `json:"booked_at"`
	Category         string    `json:"category"`
	Comment          string    `json:"comment,omitempty"`
}

type PaymentDraft struct {
	ID        string    `json:"draft_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Purpose   string    `json:"purpose"`
	PayeeName string    `json:"payee_name"`
	PayeeINN  string    `json:"payee_inn,omitempty"`
	Status    string    `json:"status"`
	RiskLevel string    `json:"risk_level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EnpPiggy struct {
	Enabled          bool    `json:"enabled"`
	RatePercent      float64 `json:"rate_percent"`
	Balance          float64 `json:"balance"`
	Currency         string  `json:"currency"`
	LastContribution float64 `json:"last_contribution"`
	TargetAmount     float64 `json:"target_amount"`
}

type PiggyContribution struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	BookedAt  time.Time `json:"booked_at"`
	SourceRef string    `json:"source_ref"`
}

type RiskResult struct {
	INN            string   `json:"inn"`
	Level          string   `json:"level"`
	Title          string   `json:"title"`
	Reasons        []string `json:"reasons"`
	Recommendation string   `json:"recommendation"`
}

type LegalFlag struct {
	Severity       string `json:"severity"`
	Title          string `json:"title"`
	Quote          string `json:"quote"`
	Recommendation string `json:"recommendation"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LegalDocument is the last uploaded contract kept in session for follow-ups.
type LegalDocument struct {
	ID       string      `json:"document_id"`
	Name     string      `json:"document_name"`
	Excerpt  string      `json:"excerpt"`
	FullText string      `json:"full_text,omitempty"`
	Markdown string      `json:"markdown,omitempty"`
	Flags    []LegalFlag `json:"flags"`
	Mode     string      `json:"mode"`
	Pages    int         `json:"pages,omitempty"`
	Quality  string      `json:"quality,omitempty"`
}

type TaxCalendar struct {
	PeriodLabel   string  `json:"period_label"`
	DueDate       string  `json:"due_date"`
	DaysLeft      int     `json:"days_left"`
	AmountDue     float64 `json:"amount_due"`
	AmountSaved   float64 `json:"amount_saved"`
	AmountGap     float64 `json:"amount_gap"`
	Status        string  `json:"status"`
}

const Disclaimer = "Альфа-Бизнес: Старт. Рекомендации носят справочный характер и не заменяют консультацию специалиста или официальные разъяснения ФНС. Финальное решение — за вами."
