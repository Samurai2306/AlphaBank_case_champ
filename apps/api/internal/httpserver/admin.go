package httpserver

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (s *Server) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		token = strings.TrimSpace(token)
		if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(s.cfg.AdminToken)) != 1 {
			writeErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "нужен вход в админку")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) adminLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "bad json")
		return
	}
	pass := strings.TrimSpace(body.Password)
	if pass == "" || subtle.ConstantTimeCompare([]byte(pass), []byte(s.cfg.AdminPassword)) != 1 {
		writeErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "неверный пароль")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": s.cfg.AdminToken})
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	items := s.store.AdminDirectory()
	writeJSON(w, http.StatusOK, map[string]any{
		"count": len(items),
		"items": items,
	})
}

func (s *Server) adminUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	detail, ok := s.store.AdminUser(id)
	if !ok {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	writeJSON(w, http.StatusOK, clipDetail(detail))
}

func (s *Server) adminCreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DisplayName            string  `json:"display_name"`
		Phone                  string  `json:"phone"`
		Email                  string  `json:"email"`
		BusinessSphere         string  `json:"business_sphere"`
		City                   string  `json:"city"`
		MonthlyRevenueEstimate float64 `json:"monthly_revenue_estimate"`
		TaxRegime              string  `json:"tax_regime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "bad json")
		return
	}
	_, profile, errMsg := s.store.Register(memory.RegisterInput{
		DisplayName:            body.DisplayName,
		Phone:                  body.Phone,
		Email:                  body.Email,
		BusinessSphere:         body.BusinessSphere,
		City:                   body.City,
		MonthlyRevenueEstimate: body.MonthlyRevenueEstimate,
		TaxRegime:              body.TaxRegime,
	})
	if errMsg != "" {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", adminErr(errMsg))
		return
	}
	detail, ok := s.store.AdminUser(profile.UserID)
	if !ok {
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "кабинет не сохранился")
		return
	}
	writeJSON(w, http.StatusCreated, clipDetail(detail))
}

func (s *Server) adminPatchUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		DisplayName            *string  `json:"display_name"`
		FullName               *string  `json:"full_name"`
		Phone                  *string  `json:"phone"`
		Email                  *string  `json:"email"`
		City                   *string  `json:"city"`
		BusinessSphere         *string  `json:"business_sphere"`
		TaxRegime              *string  `json:"tax_regime"`
		MonthlyRevenueEstimate *float64 `json:"monthly_revenue_estimate"`
		CJMLevel               *int     `json:"cjm_level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "bad json")
		return
	}
	_, errMsg := s.store.AdminUpdateProfile(id, memory.AdminProfilePatch{
		DisplayName:            body.DisplayName,
		FullName:               body.FullName,
		Phone:                  body.Phone,
		Email:                  body.Email,
		City:                   body.City,
		BusinessSphere:         body.BusinessSphere,
		TaxRegime:              body.TaxRegime,
		MonthlyRevenueEstimate: body.MonthlyRevenueEstimate,
		CJMLevel:               body.CJMLevel,
	})
	if errMsg == "not found" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	if errMsg != "" {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", adminErr(errMsg))
		return
	}
	s.writeAdminUser(w, id)
}

func (s *Server) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !s.store.AdminDelete(id) {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (s *Server) adminPatchPiggy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Enabled *bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Enabled == nil {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "нужно поле enabled")
		return
	}
	if _, ok := s.store.AdminSetPiggy(id, *body.Enabled); !ok {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	s.writeAdminUser(w, id)
}

func (s *Server) adminCreateDraft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Amount    float64 `json:"amount"`
		Purpose   string  `json:"purpose"`
		PayeeName string  `json:"payee_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "bad json")
		return
	}
	if body.Amount <= 0 {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "сумма должна быть больше нуля")
		return
	}
	purpose := strings.TrimSpace(body.Purpose)
	if purpose == "" {
		purpose = "Платёж со стенда"
	}
	payee := strings.TrimSpace(body.PayeeName)
	if payee == "" {
		payee = "Получатель"
	}
	if _, ok := s.store.AdminCreateDraft(id, domain.PaymentDraft{
		Amount: body.Amount, Purpose: purpose, PayeeName: payee, RiskLevel: "green",
	}); !ok {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	s.writeAdminUser(w, id)
}

func (s *Server) adminConfirmDraft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	draftID := chi.URLParam(r, "draftID")
	draft, errMsg := s.store.AdminConfirmDraft(id, draftID)
	if errMsg == "not found" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	if errMsg == "draft not found" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "черновик не найден")
		return
	}
	if errMsg == "cancelled" {
		writeErr(w, http.StatusConflict, "CONFLICT", "черновик уже отменён")
		return
	}
	_ = draft
	s.writeAdminUser(w, id)
}

func (s *Server) adminCancelDraft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	draftID := chi.URLParam(r, "draftID")
	_, errMsg := s.store.AdminCancelDraft(id, draftID)
	if errMsg == "not found" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	if errMsg == "draft not found" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "черновик не найден")
		return
	}
	if errMsg == "confirmed" {
		writeErr(w, http.StatusConflict, "CONFLICT", "черновик уже подтверждён")
		return
	}
	s.writeAdminUser(w, id)
}

func (s *Server) writeAdminUser(w http.ResponseWriter, id string) {
	detail, ok := s.store.AdminUser(id)
	if !ok {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "кабинет не найден")
		return
	}
	writeJSON(w, http.StatusOK, clipDetail(detail))
}

func clipDetail(detail memory.AdminUserDetail) memory.AdminUserDetail {
	for i := range detail.Chat {
		detail.Chat[i].Content = clip(detail.Chat[i].Content, 240)
	}
	return detail
}

func adminErr(msg string) string {
	switch msg {
	case "name and phone required":
		return "имя и телефон обязательны"
	case "name required":
		return "имя обязательно"
	case "phone required":
		return "телефон обязателен"
	case "phone already registered":
		return "этот телефон уже зарегистрирован"
	case "bad tax regime":
		return "режим: NPD, IP или UNDECIDED"
	case "bad cjm level":
		return "уровень CJM от 1 до 5"
	default:
		return msg
	}
}

func clip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
