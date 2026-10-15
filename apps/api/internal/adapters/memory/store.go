package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alphabank-case-champ/copilot-api/internal/domain"
	"github.com/google/uuid"
)

const (
	MashaUserID = "usr_masha_kazan"
	MashaPhone  = "+79175551234"
	MashaToken  = "demo-masha-token"
)

type ctxKey int

const userIDKey ctxKey = 1

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(userIDKey).(string); ok {
		return v
	}
	return ""
}

type userData struct {
	profile       domain.Profile
	txns          []domain.Transaction
	drafts        map[string]domain.PaymentDraft
	piggy         domain.EnpPiggy
	contributions []domain.PiggyContribution
	lastDoc       *domain.LegalDocument
	lastRisk      *domain.RiskResult
	chatTurns     []domain.ChatMessage
	rev           int64
}

// Store is a multi-user in-memory directory (token → user).
type Store struct {
	mu         sync.RWMutex
	users      map[string]*userData
	tokens     map[string]string // token → userID
	phones     map[string]string // normalized phone → userID
	mashaTok   string
	persist    Persister
	dirty      bool
	tombstones map[string]int64 // userID → rev; newer than a snapshot user means the cabinet stays deleted
}

// UserScope is a per-user view used by handlers, tools, and orchestrator.
type UserScope struct {
	root *Store
	id   string
}

func NewSeeded() *Store {
	return NewSeededWithToken(MashaToken)
}

func NewSeededWithToken(mashaToken string) *Store {
	if mashaToken == "" {
		mashaToken = MashaToken
	}
	now := time.Now()
	s := &Store{
		users:      map[string]*userData{},
		tokens:     map[string]string{},
		phones:     map[string]string{},
		mashaTok:   mashaToken,
		tombstones: map[string]int64{},
	}
	masha := &userData{
		drafts: map[string]domain.PaymentDraft{},
		profile: domain.Profile{
			UserID:                 MashaUserID,
			DisplayName:            "Маша",
			FullName:               "Соколова Мария Андреевна",
			Phone:                  "+7 917 555-12-34",
			Email:                  "masha.nails@example.ru",
			INN:                    "165012345678",
			BusinessSphere:         "Маникюр",
			City:                   "Казань",
			Address:                "г. Казань, ул. Баумана, 58",
			MonthlyRevenueEstimate: 180000,
			TaxRegime:              "NPD",
			OKVED:                  []string{"96.02"},
			CJMLevel:               4,
			PersonaKey:             "masha_nails",
			AccountMasked:          "···4582",
			BankName:               "Альфа-Банк",
			RegisteredAt:           now.AddDate(0, -4, -12).Format("2006-01-02"),
			NPDStatus:              "active",
		},
		piggy: domain.EnpPiggy{
			Enabled:          true,
			RatePercent:      6,
			Balance:          5400,
			Currency:         "RUB",
			LastContribution: 600,
			TargetAmount:     10800,
		},
		txns:          seedTxns(now),
		contributions: seedContributions(now),
	}
	prev := domain.PaymentDraft{
		ID: uuid.NewString(), Amount: 9600, Currency: "RUB",
		Purpose: "Налог НПД за предыдущий период", PayeeName: "УФК по Республике Татарстан",
		Status: "confirmed", RiskLevel: "green",
		CreatedAt: now.AddDate(0, -1, -5), UpdatedAt: now.AddDate(0, -1, -5),
	}
	masha.drafts[prev.ID] = prev
	masha.rev = 1
	s.users[MashaUserID] = masha
	s.tokens[mashaToken] = MashaUserID
	s.tokens["demo"] = MashaUserID
	s.phones[NormalizePhone(masha.profile.Phone)] = MashaUserID
	return s
}

func NormalizePhone(p string) string {
	var b strings.Builder
	for _, r := range p {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	s := b.String()
	if strings.HasPrefix(s, "8") && len(s) == 11 {
		s = "7" + s[1:]
	}
	if len(s) == 10 {
		s = "7" + s
	}
	return s
}

func (s *Store) For(userID string) *UserScope {
	return &UserScope{root: s, id: userID}
}

func (s *Store) ForContext(ctx context.Context) *UserScope {
	uid := UserIDFromContext(ctx)
	if uid == "" {
		uid = MashaUserID
	}
	return s.For(uid)
}

func (s *Store) ResolveToken(token string) (userID string, ok bool) {
	s.mu.RLock()
	uid, ok := s.tokens[token]
	s.mu.RUnlock()
	if ok || s.persist == nil {
		return uid, ok
	}
	s.reload()
	s.mu.RLock()
	defer s.mu.RUnlock()
	uid, ok = s.tokens[token]
	return uid, ok
}

func (s *Store) MashaToken() string {
	return s.mashaTok
}

type RegisterInput struct {
	DisplayName            string
	Phone                  string
	Email                  string
	BusinessSphere         string
	City                   string
	MonthlyRevenueEstimate float64
	TaxRegime              string
}

func (s *Store) Register(in RegisterInput) (token string, profile domain.Profile, errMsg string) {
	phone := NormalizePhone(in.Phone)
	if in.DisplayName == "" || phone == "" {
		return "", domain.Profile{}, "name and phone required"
	}
	if in.BusinessSphere == "" {
		in.BusinessSphere = "Услуги"
	}
	if in.City == "" {
		in.City = "Город"
	}
	if in.MonthlyRevenueEstimate <= 0 {
		in.MonthlyRevenueEstimate = 50000
	}
	regime := strings.ToUpper(strings.TrimSpace(in.TaxRegime))
	switch regime {
	case "NPD", "IP", "UNDECIDED":
	default:
		regime = "NPD"
	}

	s.reload()
	s.mu.Lock()
	defer s.unlockAndFlush()
	if _, exists := s.phones[phone]; exists {
		return "", domain.Profile{}, "phone already registered"
	}
	uid := "usr_" + uuid.NewString()[:8]
	token = "tok_" + uuid.NewString()
	now := time.Now()
	p := domain.Profile{
		UserID:                 uid,
		DisplayName:            strings.TrimSpace(in.DisplayName),
		FullName:               strings.TrimSpace(in.DisplayName),
		Phone:                  formatPhoneDisplay(phone),
		Email:                  strings.TrimSpace(in.Email),
		BusinessSphere:         strings.TrimSpace(in.BusinessSphere),
		City:                   strings.TrimSpace(in.City),
		MonthlyRevenueEstimate: in.MonthlyRevenueEstimate,
		TaxRegime:              regime,
		OKVED:                  []string{},
		CJMLevel:               1,
		PersonaKey:             "custom",
		AccountMasked:          "···0000",
		BankName:               "Альфа-Банк",
		RegisteredAt:           now.Format("2006-01-02"),
		NPDStatus:              map[bool]string{true: "active", false: "none"}[regime == "NPD"],
	}
	u := &userData{
		profile:       p,
		txns:          nil,
		drafts:        map[string]domain.PaymentDraft{},
		contributions: nil,
		piggy: domain.EnpPiggy{
			Enabled: false, RatePercent: 6, Balance: 0, Currency: "RUB",
			LastContribution: 0, TargetAmount: 0,
		},
	}
	s.touchLocked(u)
	s.users[uid] = u
	s.tokens[token] = uid
	s.phones[phone] = uid
	return token, p, ""
}

func (s *Store) Login(phone string) (token string, profile domain.Profile, errMsg string) {
	norm := NormalizePhone(phone)
	s.reload()
	s.mu.Lock()
	defer s.unlockAndFlush()
	uid, ok := s.phones[norm]
	if !ok {
		return "", domain.Profile{}, "user not found"
	}
	u := s.users[uid]
	// Issue a fresh token each login (except keep masha demo token stable for examples)
	token = "tok_" + uuid.NewString()
	if uid == MashaUserID {
		token = s.mashaTok
	} else {
		s.tokens[token] = uid
		s.touchLocked(u)
	}
	return token, u.profile, ""
}

func formatPhoneDisplay(norm string) string {
	if len(norm) == 11 && norm[0] == '7' {
		return "+7 " + norm[1:4] + " " + norm[4:7] + "-" + norm[7:9] + "-" + norm[9:11]
	}
	return "+" + norm
}

func (u *UserScope) data() *userData {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	return u.root.users[u.id]
}

func (u *UserScope) Profile() domain.Profile {
	d := u.data()
	if d == nil {
		return domain.Profile{}
	}
	return d.profile
}

func (u *UserScope) UpdateProfile(p domain.Profile) {
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	d := u.root.users[u.id]
	if d == nil {
		return
	}
	d.profile = p
	u.root.touchLocked(d)
}

func (u *UserScope) Transactions() []domain.Transaction {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	d := u.root.users[u.id]
	if d == nil {
		return nil
	}
	cp := make([]domain.Transaction, len(d.txns))
	copy(cp, d.txns)
	return cp
}

func (u *UserScope) IncomeThisMonth() float64 {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	d := u.root.users[u.id]
	if d == nil {
		return 0
	}
	var sum float64
	now := time.Now()
	for _, t := range d.txns {
		if t.Direction == "in" && t.BookedAt.Month() == now.Month() && t.BookedAt.Year() == now.Year() {
			sum += t.Amount
		}
	}
	if sum == 0 {
		return d.profile.MonthlyRevenueEstimate
	}
	return sum
}

func (u *UserScope) CreateDraft(d domain.PaymentDraft) domain.PaymentDraft {
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	ud := u.root.users[u.id]
	if ud == nil {
		return d
	}
	now := time.Now()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	d.Status = "draft"
	if d.Currency == "" {
		d.Currency = "RUB"
	}
	d.CreatedAt = now
	d.UpdatedAt = now
	ud.drafts[d.ID] = d
	u.root.touchLocked(ud)
	return d
}

func (u *UserScope) ConfirmDraft(id string) (domain.PaymentDraft, bool) {
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	ud := u.root.users[u.id]
	if ud == nil {
		return domain.PaymentDraft{}, false
	}
	d, ok := ud.drafts[id]
	if !ok {
		return domain.PaymentDraft{}, false
	}
	if d.Status == "cancelled" || d.Status == "confirmed" {
		return d, d.Status == "confirmed"
	}
	d.Status = "confirmed"
	d.UpdatedAt = time.Now()
	ud.drafts[id] = d
	u.root.touchLocked(ud)
	return d, true
}

func (u *UserScope) CancelDraft(id string) (domain.PaymentDraft, bool) {
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	ud := u.root.users[u.id]
	if ud == nil {
		return domain.PaymentDraft{}, false
	}
	d, ok := ud.drafts[id]
	if !ok {
		return domain.PaymentDraft{}, false
	}
	if d.Status == "confirmed" {
		return d, false
	}
	d.Status = "cancelled"
	d.UpdatedAt = time.Now()
	ud.drafts[id] = d
	u.root.touchLocked(ud)
	return d, true
}

func (u *UserScope) GetDraft(id string) (domain.PaymentDraft, bool) {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return domain.PaymentDraft{}, false
	}
	d, ok := ud.drafts[id]
	return d, ok
}

func (u *UserScope) ListDrafts() []domain.PaymentDraft {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return nil
	}
	out := make([]domain.PaymentDraft, 0, len(ud.drafts))
	for _, d := range ud.drafts {
		out = append(out, d)
	}
	return out
}

func (u *UserScope) Piggy() domain.EnpPiggy {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return domain.EnpPiggy{}
	}
	return ud.piggy
}

func (u *UserScope) SetPiggy(p domain.EnpPiggy) {
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	ud := u.root.users[u.id]
	if ud == nil {
		return
	}
	ud.piggy = p
	u.root.touchLocked(ud)
}

func (u *UserScope) Contributions() []domain.PiggyContribution {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return nil
	}
	cp := make([]domain.PiggyContribution, len(ud.contributions))
	copy(cp, ud.contributions)
	return cp
}

func (u *UserScope) IsEmptyCabinet() bool {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return true
	}
	return len(ud.txns) == 0 && !ud.piggy.Enabled && ud.piggy.Balance == 0
}

func (u *UserScope) SetLastRisk(r domain.RiskResult) {
	u.root.mu.Lock()
	defer u.root.mu.Unlock()
	ud := u.root.users[u.id]
	if ud == nil {
		return
	}
	cp := r
	ud.lastRisk = &cp
}

func (u *UserScope) LastRisk() (domain.RiskResult, bool) {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	ud := u.root.users[u.id]
	if ud == nil || ud.lastRisk == nil {
		return domain.RiskResult{}, false
	}
	return *ud.lastRisk, true
}

func (u *UserScope) SetLastDocument(doc domain.LegalDocument) {
	u.root.mu.Lock()
	defer u.root.mu.Unlock()
	d := u.root.users[u.id]
	if d == nil {
		return
	}
	cp := doc
	d.lastDoc = &cp
}

func (u *UserScope) LastDocument() (domain.LegalDocument, bool) {
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	d := u.root.users[u.id]
	if d == nil || d.lastDoc == nil {
		return domain.LegalDocument{}, false
	}
	return *d.lastDoc, true
}

func (u *UserScope) AppendChat(role, content string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}
	u.root.mu.Lock()
	defer u.root.unlockAndFlush()
	d := u.root.users[u.id]
	if d == nil {
		return
	}
	d.chatTurns = append(d.chatTurns, domain.ChatMessage{Role: role, Content: content})
	u.root.touchLocked(d)
	if len(d.chatTurns) > 16 {
		d.chatTurns = d.chatTurns[len(d.chatTurns)-16:]
	}
}

func (u *UserScope) RecentChat(n int) []domain.ChatMessage {
	if n <= 0 {
		n = 6
	}
	u.root.mu.RLock()
	defer u.root.mu.RUnlock()
	d := u.root.users[u.id]
	if d == nil || len(d.chatTurns) == 0 {
		return nil
	}
	start := 0
	if len(d.chatTurns) > n {
		start = len(d.chatTurns) - n
	}
	cp := make([]domain.ChatMessage, len(d.chatTurns)-start)
	copy(cp, d.chatTurns[start:])
	return cp
}

type AdminUserCard struct {
	Profile      domain.Profile `json:"profile"`
	TxnCount     int            `json:"txn_count"`
	DraftCount   int            `json:"draft_count"`
	ChatCount    int            `json:"chat_count"`
	PiggyBalance float64        `json:"piggy_balance"`
	PiggyOn      bool           `json:"piggy_on"`
}

type AdminUserDetail struct {
	AdminUserCard
	Drafts []domain.PaymentDraft `json:"drafts"`
	Chat   []domain.ChatMessage  `json:"chat"`
}

func (s *Store) AdminDirectory() []AdminUserCard {
	s.reload()
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]AdminUserCard, 0, len(s.users))
	for _, d := range s.users {
		if d == nil {
			continue
		}
		out = append(out, adminCard(d))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Profile.RegisteredAt > out[j].Profile.RegisteredAt
	})
	return out
}

func (s *Store) AdminUser(id string) (AdminUserDetail, bool) {
	s.reload()
	s.mu.RLock()
	defer s.mu.RUnlock()
	d := s.users[id]
	if d == nil {
		return AdminUserDetail{}, false
	}
	drafts := make([]domain.PaymentDraft, 0, len(d.drafts))
	for _, draft := range d.drafts {
		drafts = append(drafts, draft)
	}
	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].CreatedAt.After(drafts[j].CreatedAt)
	})
	chat := append([]domain.ChatMessage(nil), d.chatTurns...)
	return AdminUserDetail{
		AdminUserCard: adminCard(d),
		Drafts:        drafts,
		Chat:          chat,
	}, true
}

func adminCard(d *userData) AdminUserCard {
	return AdminUserCard{
		Profile:      d.profile,
		TxnCount:     len(d.txns),
		DraftCount:   len(d.drafts),
		ChatCount:    len(d.chatTurns),
		PiggyBalance: d.piggy.Balance,
		PiggyOn:      d.piggy.Enabled,
	}
}

type AdminProfilePatch struct {
	DisplayName            *string
	FullName               *string
	Phone                  *string
	Email                  *string
	City                   *string
	BusinessSphere         *string
	TaxRegime              *string
	MonthlyRevenueEstimate *float64
	CJMLevel               *int
}

func (s *Store) AdminUpdateProfile(id string, patch AdminProfilePatch) (domain.Profile, string) {
	s.reload()
	s.mu.Lock()
	d := s.users[id]
	if d == nil {
		s.mu.Unlock()
		return domain.Profile{}, "not found"
	}
	p := d.profile
	if patch.DisplayName != nil {
		name := strings.TrimSpace(*patch.DisplayName)
		if name == "" {
			s.mu.Unlock()
			return domain.Profile{}, "name required"
		}
		p.DisplayName = name
	}
	if patch.FullName != nil {
		p.FullName = strings.TrimSpace(*patch.FullName)
	}
	if patch.Email != nil {
		p.Email = strings.TrimSpace(*patch.Email)
	}
	if patch.City != nil {
		city := strings.TrimSpace(*patch.City)
		if city != "" {
			p.City = city
		}
	}
	if patch.BusinessSphere != nil {
		sphere := strings.TrimSpace(*patch.BusinessSphere)
		if sphere != "" {
			p.BusinessSphere = sphere
		}
	}
	if patch.MonthlyRevenueEstimate != nil && *patch.MonthlyRevenueEstimate > 0 {
		p.MonthlyRevenueEstimate = *patch.MonthlyRevenueEstimate
	}
	if patch.TaxRegime != nil && strings.TrimSpace(*patch.TaxRegime) != "" {
		regime := strings.ToUpper(strings.TrimSpace(*patch.TaxRegime))
		switch regime {
		case "NPD", "IP", "UNDECIDED":
			p.TaxRegime = regime
		default:
			s.mu.Unlock()
			return domain.Profile{}, "bad tax regime"
		}
		if regime == "NPD" && p.NPDStatus == "" {
			p.NPDStatus = "active"
		}
	}
	if patch.CJMLevel != nil {
		level := *patch.CJMLevel
		if level < 1 || level > 5 {
			s.mu.Unlock()
			return domain.Profile{}, "bad cjm level"
		}
		p.CJMLevel = level
	}
	if patch.Phone != nil {
		phone := NormalizePhone(*patch.Phone)
		if phone == "" {
			s.mu.Unlock()
			return domain.Profile{}, "phone required"
		}
		if owner, ok := s.phones[phone]; ok && owner != id {
			s.mu.Unlock()
			return domain.Profile{}, "phone already registered"
		}
		old := NormalizePhone(d.profile.Phone)
		if old != "" && old != phone && s.phones[old] == id {
			delete(s.phones, old)
		}
		s.phones[phone] = id
		p.Phone = formatPhoneDisplay(phone)
	}
	d.profile = p
	s.touchLocked(d)
	s.unlockAndFlush()
	return p, ""
}

func (s *Store) AdminDelete(id string) bool {
	s.reload()
	s.mu.Lock()
	if s.users[id] == nil {
		s.mu.Unlock()
		return false
	}
	s.dropUserLocked(id)
	if s.tombstones == nil {
		s.tombstones = map[string]int64{}
	}
	s.tombstones[id] = time.Now().UnixNano()
	s.dirty = true
	s.unlockAndFlush()
	return true
}

func (s *Store) dropUserLocked(id string) {
	d := s.users[id]
	if d != nil {
		if phone := NormalizePhone(d.profile.Phone); phone != "" && s.phones[phone] == id {
			delete(s.phones, phone)
		}
	}
	delete(s.users, id)
	for tok, uid := range s.tokens {
		if uid == id {
			delete(s.tokens, tok)
		}
	}
}

func (s *Store) AdminSetPiggy(id string, enabled bool) (domain.EnpPiggy, bool) {
	if _, ok := s.AdminUser(id); !ok {
		return domain.EnpPiggy{}, false
	}
	scope := s.For(id)
	piggy := scope.Piggy()
	piggy.Enabled = enabled
	if piggy.Currency == "" {
		piggy.Currency = "RUB"
	}
	if piggy.RatePercent == 0 {
		piggy.RatePercent = 6
	}
	scope.SetPiggy(piggy)
	return scope.Piggy(), true
}

func (s *Store) AdminCreateDraft(id string, d domain.PaymentDraft) (domain.PaymentDraft, bool) {
	if _, ok := s.AdminUser(id); !ok {
		return domain.PaymentDraft{}, false
	}
	if d.Amount <= 0 {
		return domain.PaymentDraft{}, false
	}
	return s.For(id).CreateDraft(d), true
}

func (s *Store) AdminConfirmDraft(id, draftID string) (domain.PaymentDraft, string) {
	if _, ok := s.AdminUser(id); !ok {
		return domain.PaymentDraft{}, "not found"
	}
	current, ok := s.For(id).GetDraft(draftID)
	if !ok {
		return domain.PaymentDraft{}, "draft not found"
	}
	if current.Status == "cancelled" {
		return current, "cancelled"
	}
	if current.Status == "confirmed" {
		return current, ""
	}
	updated, ok := s.For(id).ConfirmDraft(draftID)
	if !ok {
		return domain.PaymentDraft{}, "draft not found"
	}
	return updated, ""
}

func (s *Store) AdminCancelDraft(id, draftID string) (domain.PaymentDraft, string) {
	if _, ok := s.AdminUser(id); !ok {
		return domain.PaymentDraft{}, "not found"
	}
	updated, ok := s.For(id).CancelDraft(draftID)
	if !ok {
		current, exists := s.For(id).GetDraft(draftID)
		if !exists {
			return domain.PaymentDraft{}, "draft not found"
		}
		if current.Status == "confirmed" {
			return current, "confirmed"
		}
		return current, "draft not found"
	}
	return updated, ""
}

// Backward-compat helpers for tests that still call Store methods (default: Маша).
func (s *Store) Profile() domain.Profile            { return s.For(MashaUserID).Profile() }
func (s *Store) UpdateProfile(p domain.Profile)     { s.For(MashaUserID).UpdateProfile(p) }
func (s *Store) Transactions() []domain.Transaction { return s.For(MashaUserID).Transactions() }
func (s *Store) IncomeThisMonth() float64           { return s.For(MashaUserID).IncomeThisMonth() }
func (s *Store) CreateDraft(d domain.PaymentDraft) domain.PaymentDraft {
	return s.For(MashaUserID).CreateDraft(d)
}
func (s *Store) ConfirmDraft(id string) (domain.PaymentDraft, bool) {
	return s.For(MashaUserID).ConfirmDraft(id)
}
func (s *Store) CancelDraft(id string) (domain.PaymentDraft, bool) {
	return s.For(MashaUserID).CancelDraft(id)
}
func (s *Store) GetDraft(id string) (domain.PaymentDraft, bool) {
	return s.For(MashaUserID).GetDraft(id)
}
func (s *Store) ListDrafts() []domain.PaymentDraft         { return s.For(MashaUserID).ListDrafts() }
func (s *Store) Piggy() domain.EnpPiggy                    { return s.For(MashaUserID).Piggy() }
func (s *Store) SetPiggy(p domain.EnpPiggy)                { s.For(MashaUserID).SetPiggy(p) }
func (s *Store) Contributions() []domain.PiggyContribution { return s.For(MashaUserID).Contributions() }

func seedContributions(now time.Time) []domain.PiggyContribution {
	out := make([]domain.PiggyContribution, 0, 9)
	for i := 0; i < 9; i++ {
		out = append(out, domain.PiggyContribution{
			ID: uuid.NewString(), Amount: 600,
			BookedAt: now.AddDate(0, 0, -i*3), SourceRef: "auto:incoming",
		})
	}
	return out
}

func seedTxns(now time.Time) []domain.Transaction {
	type row struct {
		name, cat, dir string
		amt            float64
		daysAgo        int
		inn, comment   string
	}
	rows := []row{
		{"Анна · маникюр", "services", "in", 2500, 0, "", "Классика + покрытие"},
		{"Дарья · педикюр", "services", "in", 3200, 1, "", ""},
		{"Ирина · наращивание", "services", "in", 4500, 1, "", ""},
		{"Салон «Нежность»", "services", "in", 12000, 2, "1651234567", "Аутсорс 4 мастера"},
		{"Елена · комбо", "services", "in", 3800, 2, "", ""},
		{"Аренда кабинета SoftSpace", "rent", "out", 40000, 3, "1650987654", "Июль"},
		{"Ольга · дизайн", "services", "in", 2900, 3, "", ""},
		{"Мария · коррекция", "services", "in", 2100, 4, "", ""},
		{"Расходники NailPro", "supplies", "out", 4200, 5, "7701234567", "Гель-лак, типсы"},
		{"Клиент через Авито", "services", "in", 1800, 5, "", ""},
		{"Студия BeautyHub", "services", "in", 9000, 6, "1651112233", "Сменный день"},
		{"Подарочный сертификат", "services", "in", 5000, 7, "", ""},
		{"Реклама VK Ads", "ads", "out", 890, 8, "", "Продвижение кабинета"},
		{"Анна · повтор", "services", "in", 2500, 9, "", ""},
		{"Дарья · маникюр", "services", "in", 2300, 10, "", ""},
		{"Комиссия эквайринга", "fees", "out", 420, 10, "", "Альфа-Касса"},
		{"Ирина · снятие", "services", "in", 3500, 11, "", ""},
		{"Елена · педикюр", "services", "in", 3000, 12, "", ""},
		{"Расходники Ozone", "supplies", "out", 1680, 13, "", "Салфетки, антисептик"},
		{"Мария · дизайн", "services", "in", 3100, 14, "", ""},
		{"Ольга · маникюр", "services", "in", 2400, 15, "", ""},
		{"Салон «Нежность»", "services", "in", 8000, 16, "1651234567", ""},
		{"Связь МТС Бизнес", "ops", "out", 650, 17, "", ""},
		{"Анна · педикюр", "services", "in", 3200, 18, "", ""},
		{"Дарья · комбо", "services", "in", 4000, 19, "", ""},
		{"Ирина · маникюр", "services", "in", 2600, 20, "", ""},
		{"Аренда доп. час", "rent", "out", 2500, 21, "1650987654", ""},
		{"Елена · коррекция", "services", "in", 2000, 22, "", ""},
		{"Поступление от клиента", "services", "in", 1500, 24, "", "Перевод СБП"},
		{"Реклама Telegram", "ads", "out", 500, 26, "", ""},
	}
	maxAgo := now.Day() - 1
	if maxAgo < 0 {
		maxAgo = 0
	}
	out := make([]domain.Transaction, 0, len(rows))
	for _, r := range rows {
		ago := r.daysAgo
		if ago > maxAgo {
			ago = maxAgo
		}
		out = append(out, domain.Transaction{
			ID: uuid.NewString(), Amount: r.amt, Direction: r.dir,
			CounterpartyName: r.name, CounterpartyINN: r.inn,
			BookedAt: now.AddDate(0, 0, -ago), Category: r.cat, Comment: r.comment,
		})
	}
	return out
}
