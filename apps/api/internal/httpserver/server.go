package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/adapters/mockbank"
	"github.com/alphabank-case-champ/copilot-api/internal/agents"
	"github.com/alphabank-case-champ/copilot-api/internal/calc"
	"github.com/alphabank-case-champ/copilot-api/internal/config"
	"github.com/alphabank-case-champ/copilot-api/internal/domain"
	"github.com/alphabank-case-champ/copilot-api/internal/guard"
	"github.com/alphabank-case-champ/copilot-api/internal/llm"
	"github.com/alphabank-case-champ/copilot-api/internal/pdftext"
	"github.com/alphabank-case-champ/copilot-api/internal/rag"
	"github.com/alphabank-case-champ/copilot-api/internal/tools"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type Server struct {
	cfg    config.Config
	store  *memory.Store
	orch   *agents.Orchestrator
	log    *slog.Logger
}

func New(cfg config.Config, store *memory.Store, log *slog.Logger) http.Handler {
	var chatLLM *llm.Client
	if !cfg.DemoOffline && cfg.LLMAPIKey != "" {
		chatLLM = llm.New(cfg.LLMBaseURL, cfg.LLMAPIKey, cfg.LLMModel)
	}
	s := &Server{
		cfg: cfg, store: store, log: log,
		orch: &agents.Orchestrator{Store: store, LLM: chatLLM, Offline: cfg.DemoOffline},
	}
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/api/v1/health", s.health)
	r.Get("/api/v1/ready", s.ready)
	r.Post("/api/v1/register", s.register)
	r.Post("/api/v1/login", s.login)

	r.Group(func(r chi.Router) {
		r.Use(s.auth)
		r.Get("/api/v1/me", s.me)
		r.Patch("/api/v1/me", s.patchMe)
		r.Patch("/api/v1/me/onboarding", s.patchOnboarding)
		r.Get("/api/v1/transactions", s.transactions)
		r.Get("/api/v1/piggy", s.piggy)
		r.Patch("/api/v1/piggy", s.patchPiggy)
		r.Get("/api/v1/alerts", s.alerts)
		r.Post("/api/v1/chat", s.chat)
		r.Get("/api/v1/payments", s.listPayments)
		r.Post("/api/v1/payments/draft", s.createDraft)
		r.Post("/api/v1/payments/draft/{id}/confirm", s.confirmDraft)
		r.Post("/api/v1/payments/draft/{id}/cancel", s.cancelDraft)
		r.Post("/api/v1/documents/legal", s.legalUpload)
		r.Get("/api/v1/tasks", s.tasks)
		r.Get("/api/v1/demo/script", s.tasks) // compat alias
		r.Get("/api/v1/calendar/tax", s.taxCalendar)
		r.Get("/api/v1/report/month", s.monthReport)
		r.Get("/api/v1/piggy/contributions", s.piggyContributions)
		r.Get("/api/v1/home", s.home)
	})
	return r
}

func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		token := strings.TrimPrefix(h, "Bearer ")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		uid, ok := s.store.ResolveToken(token)
		if !ok {
			writeErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			return
		}
		ctx := memory.WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) scope(r *http.Request) *memory.UserScope {
	return s.store.ForContext(r.Context())
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
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
		writeErr(w, 400, "VALIDATION_ERROR", "bad json")
		return
	}
	token, profile, errMsg := s.store.Register(memory.RegisterInput{
		DisplayName:            body.DisplayName,
		Phone:                  body.Phone,
		Email:                  body.Email,
		BusinessSphere:         body.BusinessSphere,
		City:                   body.City,
		MonthlyRevenueEstimate: body.MonthlyRevenueEstimate,
		TaxRegime:              body.TaxRegime,
	})
	if errMsg != "" {
		code := 400
		if strings.Contains(errMsg, "already") {
			code = 409
		}
		writeErr(w, code, "VALIDATION_ERROR", errMsg)
		return
	}
	writeJSON(w, 201, map[string]any{"token": token, "profile": profile})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Phone) == "" {
		writeErr(w, 400, "VALIDATION_ERROR", "phone required")
		return
	}
	token, profile, errMsg := s.store.Login(body.Phone)
	if errMsg != "" {
		writeErr(w, 404, "NOT_FOUND", "Аккаунт не найден. Зарегистрируйтесь — это займёт пару минут.")
		return
	}
	writeJSON(w, 200, map[string]any{"token": token, "profile": profile})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "copilot-api"})
}

func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ready", "store": "memory",
		"offline": s.cfg.DemoOffline, "llm_configured": s.cfg.LLMAPIKey != "",
		"kb_chunks": rag.Default().Len(),
	})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.scope(r).Profile())
}

func (s *Server) patchOnboarding(w http.ResponseWriter, r *http.Request) {
	s.patchMe(w, r)
}

func (s *Server) patchMe(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DisplayName            string  `json:"display_name"`
		BusinessSphere         string  `json:"business_sphere"`
		City                   string  `json:"city"`
		MonthlyRevenueEstimate float64 `json:"monthly_revenue_estimate"`
		TaxRegime              string  `json:"tax_regime"`
		Email                  string  `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, "VALIDATION_ERROR", "bad json")
		return
	}
	sess := s.scope(r)
	p := sess.Profile()
	if body.DisplayName != "" {
		p.DisplayName = body.DisplayName
	}
	if body.BusinessSphere != "" {
		p.BusinessSphere = body.BusinessSphere
	}
	if body.City != "" {
		p.City = body.City
	}
	if body.MonthlyRevenueEstimate > 0 {
		p.MonthlyRevenueEstimate = body.MonthlyRevenueEstimate
	}
	if body.TaxRegime != "" {
		p.TaxRegime = body.TaxRegime
	}
	if body.Email != "" {
		p.Email = body.Email
	}
	// Soft progress: filled profile bumps toward level 2, not jump to 4
	if p.CJMLevel < 2 && p.BusinessSphere != "" && p.City != "" && p.MonthlyRevenueEstimate > 0 {
		p.CJMLevel = 2
	}
	sess.UpdateProfile(p)
	writeJSON(w, 200, p)
}

func (s *Server) transactions(w http.ResponseWriter, r *http.Request) {
	items := s.scope(r).Transactions()
	sort.Slice(items, func(i, j int) bool {
		return items[i].BookedAt.After(items[j].BookedAt)
	})
	var in, out float64
	for _, t := range items {
		if t.Direction == "in" {
			in += t.Amount
		} else {
			out += t.Amount
		}
	}
	writeJSON(w, 200, map[string]any{
		"items": items,
		"summary": map[string]any{
			"income": in, "expense": out, "net": round2(in - out), "currency": "RUB",
		},
	})
}

func (s *Server) piggy(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.scope(r).Piggy())
}

func (s *Server) patchPiggy(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Enabled     *bool    `json:"enabled"`
		RatePercent *float64 `json:"rate_percent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, "VALIDATION_ERROR", "bad json")
		return
	}
	sess := s.scope(r)
	p := sess.Piggy()
	if body.Enabled != nil {
		p.Enabled = *body.Enabled
	}
	if body.RatePercent != nil {
		rate := *body.RatePercent
		if rate < 1 {
			rate = 1
		}
		if rate > 15 {
			rate = 15
		}
		p.RatePercent = rate
	}
	prof := sess.Profile()
	if p.Enabled {
		p.LastContribution = round2(prof.MonthlyRevenueEstimate * p.RatePercent / 100 / 30)
	}
	ceiling := calc.CalculateTax(calc.TaxInput{
		Regime: calc.RegimeNPD, Income: prof.MonthlyRevenueEstimate, ConservativeCeiling: true,
	})
	p.TargetAmount = ceiling.TaxAmount
	sess.SetPiggy(p)
	if p.Enabled && prof.CJMLevel < 5 {
		prof.CJMLevel = 5
		sess.UpdateProfile(prof)
	}
	writeJSON(w, 200, p)
}

func (s *Server) alerts(w http.ResponseWriter, r *http.Request) {
	sess := s.scope(r)
	p := sess.Profile()
	piggy := sess.Piggy()
	ceiling := calc.CalculateTax(calc.TaxInput{
		Regime: calc.RegimeNPD, Income: p.MonthlyRevenueEstimate, ConservativeCeiling: true,
	})
	gap := round2(ceiling.TaxAmount - piggy.Balance)
	if gap < 0 {
		gap = 0
	}
	items := []map[string]any{}
	if sess.IsEmptyCabinet() {
		items = append(items, map[string]any{
			"id": "empty-start", "severity": "medium", "title": "Кабинет только открыт",
			"text": "Укажите оборот в профиле и посчитайте налог — так появятся цифры на главной и в копилке.",
			"cta":  "Профиль", "href": "/app/profile",
		})
	} else {
		items = append(items, map[string]any{
			"id": "enp-gap", "severity": map[bool]string{true: "high", false: "medium"}[gap > 0],
			"title": "Разрыв копилки ЕНП",
			"text":  fmt.Sprintf("До суммы «с запасом» %s ₽ не хватает %s ₽.", formatInt(int(ceiling.TaxAmount)), formatInt(int(gap))),
			"cta":   "Открыть копилку", "href": "/app/piggy",
		})
	}
	items = append(items,
		map[string]any{
			"id": "counterparty", "severity": "medium", "title": "Проверка перед арендой",
			"text": "Крупный исходящий платёж без проверки ИНН повышает риск 115-ФЗ.",
			"cta":  "Проверить контрагента",
			"href": "/app/chat?q=" + urlQuery("Проверь ИНН 7707083893"),
		},
		map[string]any{
			"id": "legalization", "severity": "low", "title": "Мой путь",
			"text": fmt.Sprintf("Уровень пути %d/5 — закройте шаги на трекборде.", p.CJMLevel),
			"cta":  "Открыть путь", "href": "/app/journey",
		},
	)
	writeJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	sess := s.scope(r)
	p := sess.Profile()
	income := p.MonthlyRevenueEstimate
	if income <= 0 {
		income = sess.IncomeThisMonth()
	}
	period := time.Now().Format("2006-01")
	mixed := calc.CalculateTax(calc.TaxInput{Regime: calc.RegimeNPD, Income: income, Period: period})
	ceiling := calc.CalculateTax(calc.TaxInput{Regime: calc.RegimeNPD, Income: income, Period: period, ConservativeCeiling: true})
	piggy := sess.Piggy()
	gap := round2(ceiling.TaxAmount - piggy.Balance)
	if gap < 0 {
		gap = 0
	}
	taxHref := "/app/chat?q=" + urlQuery("Сколько мне отложить на налог за этот месяц?")
	innHref := "/app/chat?q=" + urlQuery("Проверь ИНН 7707083893")
	due := time.Date(time.Now().Year(), time.Now().Month(), 28, 0, 0, 0, 0, time.Now().Location())
	if time.Now().Day() > 28 {
		due = due.AddDate(0, 1, 0)
	}
	daysLeft := int(due.Sub(time.Now()).Hours() / 24)
	if daysLeft < 0 {
		daysLeft = 0
	}
	empty := sess.IsEmptyCabinet()
	insights := []map[string]any{}
	if empty {
		insights = append(insights, map[string]any{
			"id":   "empty",
			"text": "Кабинет пустой: операций ещё нет. Налог ниже считается от оборота в профиле — уточните цифры и посчитайте налог в чате.",
			"why":  "Так дашборды заполнятся вашими сценариями, а не чужими.",
			"cta":  "Рассчитать налог",
			"href": taxHref,
		})
	} else {
		insights = append(insights,
			map[string]any{
				"id": "1",
				"text": fmt.Sprintf(
					"В копилке %s ₽, а отложить с запасом нужно %s ₽ — не хватает %s ₽ до условного срока ЕНП.",
					formatInt(int(piggy.Balance)), formatInt(int(ceiling.TaxAmount)), formatInt(int(gap)),
				),
				"why":  "Так вы не снимаете налог из оборота в последний день.",
				"cta":  "Рассчитать налог",
				"href": taxHref,
			},
			map[string]any{
				"id":   "2",
				"text": "Перед крупной оплатой аренды стоит проверить контрагента по ИНН и сразу собрать платёжку.",
				"why":  "Светофор 115-ФЗ снижает риск блокировки платежа.",
				"cta":  "Проверить ИНН",
				"href": innHref,
			},
			map[string]any{
				"id":   "3",
				"text": "Чтобы быстрее выходить на безубыточность, принимайте оплату по СБП или эквайрингу на РКО.",
				"why":  "Меньше кассовых разрывов между услугой и деньгами на счёте.",
				"cta":  "Спросить про СБП",
				"href": "/app/chat?q=" + urlQuery("Как принимать оплату по СБП и эквайрингу?"),
			},
		)
	}
	writeJSON(w, 200, map[string]any{
		"profile":          p,
		"tax_due":          ceiling.TaxAmount,
		"tax_estimate":     mixed.TaxAmount,
		"tax_rate_label":   mixed.RateLabel,
		"period":           period,
		"piggy":            piggy,
		"piggy_gap":        gap,
		"tax_due_date":     due.Format("2006-01-02"),
		"tax_days_left":    daysLeft,
		"assumptions":      mixed.Assumptions,
		"empty_cabinet":    empty,
		"insights":         insights,
	})
}

func urlQuery(s string) string {
	return url.QueryEscape(s)
}

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message string               `json:"message"`
		History []domain.ChatMessage `json:"history"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
		writeErr(w, 400, "VALIDATION_ERROR", "message required")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, 500, "STREAM_UNSUPPORTED", "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	events := s.orch.HandleCtx(r.Context(), body.Message, body.History)
	for _, ev := range events {
		payload, _ := json.Marshal(ev.Data)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Event, payload)
		flusher.Flush()
		time.Sleep(12 * time.Millisecond)
	}
}

func (s *Server) createDraft(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Amount    float64 `json:"amount"`
		Purpose   string  `json:"purpose"`
		PayeeName string  `json:"payee_name"`
		PayeeINN  string  `json:"payee_inn"`
		RiskLevel string  `json:"risk_level"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	sess := s.scope(r)
	if body.Amount <= 0 {
		body.Amount = round2(sess.Profile().MonthlyRevenueEstimate * 0.06)
	}
	if safe, ok := guard.SanitizeDraftAmount(body.Amount); !ok {
		writeErr(w, 400, "VALIDATION_ERROR", "amount out of allowed demo range")
		return
	} else {
		body.Amount = safe
	}
	if body.Purpose == "" {
		body.Purpose = "Налог НПД"
	}
	// Block purpose strings that look like policy abuse.
	if blocked := guard.CheckIngress(body.Purpose); blocked != nil {
		writeErr(w, 400, blocked.Code, blocked.Message)
		return
	}
	payee := strings.TrimSpace(body.PayeeName)
	if payee == "" {
		payee = "УФК по Республике Татарстан"
	}
	risk := strings.TrimSpace(body.RiskLevel)
	if risk == "" {
		risk = "green"
	}
	if body.PayeeINN != "" {
		rsk := mockbank.CheckCounterpartyRisk(body.PayeeINN)
		risk = rsk.Level
		sess.SetLastRisk(rsk)
		if body.PayeeName == "" {
			payee = "Контрагент ИНН " + body.PayeeINN
		}
	}
	d := sess.CreateDraft(domain.PaymentDraft{
		Amount: body.Amount, Purpose: body.Purpose,
		PayeeName: payee, PayeeINN: body.PayeeINN, RiskLevel: risk,
	})
	writeJSON(w, 200, d)
}

func (s *Server) confirmDraft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sess := s.scope(r)
	d, ok := sess.ConfirmDraft(id)
	if !ok {
		writeErr(w, 404, "NOT_FOUND", "draft not found")
		return
	}
	p := sess.Profile()
	if p.CJMLevel < 5 {
		p.CJMLevel = 5
		sess.UpdateProfile(p)
	}
	writeJSON(w, 200, d)
}

func (s *Server) cancelDraft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	d, ok := s.scope(r).CancelDraft(id)
	if !ok {
		writeErr(w, 404, "NOT_FOUND", "draft not found or already confirmed")
		return
	}
	writeJSON(w, 200, d)
}

func (s *Server) legalUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeErr(w, 400, "VALIDATION_ERROR", "multipart required, max 10MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeErr(w, 400, "VALIDATION_ERROR", "file required")
		return
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		writeErr(w, 400, "VALIDATION_ERROR", "only PDF allowed")
		return
	}
	raw, err := io.ReadAll(io.LimitReader(file, 10<<20+1))
	if err != nil {
		writeErr(w, 500, "UPLOAD_FAILED", err.Error())
		return
	}
	if len(raw) > 10<<20 {
		writeErr(w, 400, "VALIDATION_ERROR", "file exceeds 10MB")
		return
	}
	if len(raw) < 4 || string(raw[:4]) != "%PDF" {
		writeErr(w, 400, "VALIDATION_ERROR", "file is not a valid PDF")
		return
	}
	name := uuid.NewString() + ".pdf"
	dst := filepath.Join(s.cfg.UploadDir, name)
	if err := os.MkdirAll(s.cfg.UploadDir, 0o777); err != nil {
		// still try write — volume may already exist
		_ = err
	}
	if err := os.WriteFile(dst, raw, 0o644); err != nil {
		// Fallback for distroless/nonroot when volume perms are wrong.
		altDir := "/tmp/uploads"
		_ = os.MkdirAll(altDir, 0o777)
		dst = filepath.Join(altDir, name)
		if err2 := os.WriteFile(dst, raw, 0o644); err2 != nil {
			writeErr(w, 500, "UPLOAD_FAILED", err.Error())
			return
		}
	}
	extracted := pdftext.ExtractDetailed(raw)
	scanText := extracted.Text
	if scanText == "" {
		scanText = extracted.Markdown
	}
	result := (&tools.Registry{Store: s.scope(r)}).ScanLegalDocument(header.Filename, scanText)
	flags, _ := result["flags"].([]domain.LegalFlag)
	excerpt, _ := result["excerpt"].(string)
	mode := fmt.Sprint(result["mode"])
	if extracted.Method != "" {
		mode = mode + "+" + extracted.Method
	}
	sess := s.scope(r)
	docMD := extracted.Markdown
	if docMD == "" {
		docMD = pdftext.FormatForLLM(scanText, extracted.Pages, extracted.Method)
	}
	sess.SetLastDocument(domain.LegalDocument{
		ID: name, Name: header.Filename, Excerpt: excerpt, FullText: scanText,
		Markdown: docMD, Flags: flags, Mode: mode, Pages: extracted.Pages, Quality: extracted.Quality,
	})
	sess.AppendChat("user", "Загружен договор: "+header.Filename)
	narrative := "Разобрала загруженный договор. Ниже — на что обратить внимание до подписания. Это предварительный разбор, не юридическое заключение."
	if extracted.Quality == "poor" {
		narrative = "Текст из PDF извлечён слабо (возможно, скан без текстового слоя). Разбор по шаблону ниже — загрузите текстовый PDF или DOC/PDF с копируемым текстом для точного анализа. " + narrative
	}
	llmBody := docMD
	if len([]rune(llmBody)) > 16000 {
		llmBody = string([]rune(llmBody)[:16000]) + "\n…[обрезано]"
	}
	if !s.cfg.DemoOffline && s.orch.LLM != nil && strings.TrimSpace(scanText) != "" {
		system := `Ты юридический ассистент продукта «Альфа-Бизнес: Старт».
Документ уже преобразован в markdown для тебя — опирайся на него как на источник истины.
Структура ответа:
1) Главный вывод для предпринимателя
2) Риски с цитатами из текста
3) Что попросить изменить в договоре
4) Дисклеймер (не замена юристу)
Не выдумывай пункты, которых нет в тексте.`
		user := fmt.Sprintf("Файл: %s\nКачество извлечения: %s (%s), страниц≈%d\nФлаги эвристики: %v\n\nDOCUMENT_MARKDOWN:\n%s",
			header.Filename, extracted.Quality, extracted.Method, extracted.Pages, flags, llmBody)
		if text, err := s.orch.LLM.CompleteTemp(r.Context(), system, user, 0.45); err == nil {
			text = strings.TrimSpace(text)
			if text != "" && guard.CheckEgress(text) {
				narrative = text
			}
		}
	}
	sess.AppendChat("assistant", narrative)
	writeJSON(w, 200, map[string]any{
		"document_id":     name,
		"document_name":   header.Filename,
		"component":       "LegalFlagsList",
		"mode":            mode,
		"flags":           flags,
		"excerpt":         excerpt,
		"narrative":       narrative,
		"text_chars":      extracted.Chars,
		"pages":           extracted.Pages,
		"extract_method":  extracted.Method,
		"extract_quality": extracted.Quality,
		"extract_warnings": extracted.Warnings,
		"disclaimer":      result["disclaimer"],
	})
}

func (s *Server) tasks(w http.ResponseWriter, r *http.Request) {
	p := s.scope(r).Profile()
	onboardQ := fmt.Sprintf(
		"Привет! Меня зовут %s, я занимаюсь «%s» в городе %s, оборот около %.0f в месяц.",
		p.DisplayName, p.BusinessSphere, p.City, p.MonthlyRevenueEstimate,
	)
	items := []map[string]string{
		{"id": "1", "label": "Уточнить профиль", "text": onboardQ, "href": "/app/chat?q=" + urlQuery(onboardQ)},
		{"id": "2", "label": "Рассчитать налог", "text": "Сколько мне отложить на налог за этот месяц?", "href": "/app/chat?q=" + urlQuery("Сколько мне отложить на налог за этот месяц?")},
		{"id": "3", "label": "Посчитать безубыточность", "text": "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.", "href": "/app/chat?q=" + urlQuery("Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.")},
		{"id": "4", "label": "Сформировать платёжку", "text": "Сформируй платёжку.", "href": "/app/chat?q=" + urlQuery("Сформируй платёжку.")},
		{"id": "5", "label": "Настроить копилку", "text": "Включи копилку 6%.", "href": "/app/piggy"},
		{"id": "6", "label": "Проверить контрагента", "text": "Проверь ИНН 1650987654", "href": "/app/chat?q=" + urlQuery("Проверь ИНН 1650987654")},
	}
	if p.CJMLevel <= 2 {
		items = []map[string]string{
			{"id": "1", "label": "Заполнить профиль", "text": "Откройте профиль и проверьте сферу, город и оборот", "href": "/app/profile"},
			{"id": "2", "label": "Рассчитать налог", "text": "Сколько мне отложить на налог за этот месяц?", "href": "/app/chat?q=" + urlQuery("Сколько мне отложить на налог за этот месяц?")},
			{"id": "3", "label": "Пройти путь · уровень 1", "text": "Онбординг", "href": "/app/journey"},
			{"id": "4", "label": "Включить копилку", "text": "Включи копилку 6%.", "href": "/app/chat?q=" + urlQuery("Включи копилку 6%.")},
		}
	}
	writeJSON(w, 200, map[string]any{
		"title": "Рекомендуемые действия",
		"items": items,
	})
}

func (s *Server) listPayments(w http.ResponseWriter, r *http.Request) {
	items := s.scope(r).ListDrafts()
	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	writeJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) taxCalendar(w http.ResponseWriter, r *http.Request) {
	sess := s.scope(r)
	p := sess.Profile()
	piggy := sess.Piggy()
	now := time.Now()
	// ENP due ~ 28th next month for previous period — simplified: 28th of current month
	due := time.Date(now.Year(), now.Month(), 28, 0, 0, 0, 0, now.Location())
	if now.Day() > 28 {
		due = due.AddDate(0, 1, 0)
	}
	days := int(due.Sub(now).Hours() / 24)
	if days < 0 {
		days = 0
	}
	ceiling := calc.CalculateTax(calc.TaxInput{
		Regime: calc.RegimeNPD, Income: p.MonthlyRevenueEstimate, ConservativeCeiling: true,
	})
	gap := round2(ceiling.TaxAmount - piggy.Balance)
	if gap < 0 {
		gap = 0
	}
	status := "on_track"
	if gap > 0 && days <= 10 {
		status = "attention"
	}
	if gap > ceiling.TaxAmount*0.5 && days <= 5 {
		status = "urgent"
	}
	writeJSON(w, 200, domain.TaxCalendar{
		PeriodLabel: monthRU(now), DueDate: due.Format("2006-01-02"), DaysLeft: days,
		AmountDue: ceiling.TaxAmount, AmountSaved: piggy.Balance, AmountGap: gap, Status: status,
	})
}

func (s *Server) monthReport(w http.ResponseWriter, r *http.Request) {
	sess := s.scope(r)
	items := sess.Transactions()
	now := time.Now()
	var in, out float64
	byCat := map[string]float64{}
	countIn := 0
	for _, t := range items {
		if t.BookedAt.Month() != now.Month() || t.BookedAt.Year() != now.Year() {
			continue
		}
		if t.Direction == "in" {
			in += t.Amount
			countIn++
		} else {
			out += t.Amount
			byCat[t.Category] += t.Amount
		}
	}
	tax := calc.CalculateTax(calc.TaxInput{Regime: calc.RegimeNPD, Income: sess.Profile().MonthlyRevenueEstimate})
	writeJSON(w, 200, map[string]any{
		"period":          now.Format("2006-01"),
		"income":          round2(in),
		"expense":         round2(out),
		"net":             round2(in - out),
		"service_count":   countIn,
		"expense_by_category": byCat,
		"tax_estimate":    tax.TaxAmount,
		"currency":        "RUB",
		"empty":           len(items) == 0,
	})
}

func (s *Server) piggyContributions(w http.ResponseWriter, r *http.Request) {
	items := s.scope(r).Contributions()
	sort.Slice(items, func(i, j int) bool {
		return items[i].BookedAt.After(items[j].BookedAt)
	})
	writeJSON(w, 200, map[string]any{"items": items})
}

func monthRU(t time.Time) string {
	months := []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	return fmt.Sprintf("%s %d", months[int(t.Month())-1], t.Year())
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, errCode, msg string) {
	writeJSON(w, code, map[string]any{"error": map[string]any{"code": errCode, "message": msg, "retryable": false}})
}

func round2(v float64) float64 { return float64(int(v*100+0.5)) / 100 }

func formatInt(n int) string {
	s := fmt.Sprintf("%d", n)
	return s
}
