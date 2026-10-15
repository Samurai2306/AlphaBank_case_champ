package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alphabank-case-champ/copilot-api/internal/domain"
)

// Persister keeps cabinets across Vercel processes. Each API instance only
// has its own memory, so a registration on one instance is invisible to the
// admin list on another until it is written here.
type Persister interface {
	Load() ([]byte, error)
	Save([]byte) error
}

type cabinetSnap struct {
	Users   []snapUser    `json:"users"`
	Deleted []deletedUser `json:"deleted,omitempty"`
}

type deletedUser struct {
	UserID string `json:"user_id"`
	Rev    int64  `json:"rev"`
}

type snapUser struct {
	Rev           int64                      `json:"rev"`
	Profile       domain.Profile             `json:"profile"`
	Tokens        []string                   `json:"tokens"`
	Txns          []domain.Transaction       `json:"txns,omitempty"`
	Drafts        []domain.PaymentDraft      `json:"drafts,omitempty"`
	Piggy         domain.EnpPiggy            `json:"piggy"`
	Contributions []domain.PiggyContribution `json:"contributions,omitempty"`
	Chat          []domain.ChatMessage       `json:"chat,omitempty"`
}

func (s *Store) AttachPersisterFromEnv() error {
	token := strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN"))
	if token != "" {
		storeID := storeIDFromToken(token)
		if storeID == "" {
			return fmt.Errorf("BLOB_READ_WRITE_TOKEN has no store id")
		}
		s.persist = &blobPersister{
			token:   token,
			storeID: storeID,
			client:  &http.Client{Timeout: 12 * time.Second},
		}
		s.reload()
		return nil
	}
	if path := strings.TrimSpace(os.Getenv("CABINET_FILE")); path != "" {
		s.persist = filePersister{path: path}
		s.reload()
	}
	return nil
}

func (s *Store) Persistent() bool {
	return s.persist != nil
}

func (s *Store) touchLocked(d *userData) {
	if d == nil {
		return
	}
	d.rev = time.Now().UnixNano()
	s.dirty = true
}

func (s *Store) unlockAndFlush() {
	dirty := s.dirty
	s.dirty = false
	s.mu.Unlock()
	if dirty && s.persist != nil {
		s.flush()
	}
}

func (s *Store) reload() {
	if s.persist == nil {
		return
	}
	raw, err := s.persist.Load()
	if err != nil || len(raw) == 0 {
		if err != nil {
			log.Printf("cabinet snapshot load: %v", err)
		}
		return
	}
	s.mu.Lock()
	s.mergeLocked(raw)
	s.mu.Unlock()
}

func (s *Store) flush() {
	if s.persist == nil {
		return
	}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		remote, err := s.persist.Load()
		if err != nil {
			lastErr = err
			continue
		}
		s.mu.Lock()
		if len(remote) > 0 {
			s.mergeLocked(remote)
		}
		payload, mErr := json.Marshal(s.snapshotLocked())
		ids := make(map[string]struct{}, len(s.users))
		for id := range s.users {
			ids[id] = struct{}{}
		}
		s.mu.Unlock()
		if mErr != nil {
			log.Printf("cabinet snapshot encode: %v", mErr)
			return
		}
		if err := s.persist.Save(payload); err != nil {
			lastErr = err
			continue
		}
		again, err := s.persist.Load()
		if err != nil {
			lastErr = err
			continue
		}
		if snapshotCovers(again, ids) {
			return
		}
	}
	if lastErr != nil {
		log.Printf("cabinet snapshot save: %v", lastErr)
	}
}

func (s *Store) snapshotLocked() cabinetSnap {
	out := cabinetSnap{Users: make([]snapUser, 0, len(s.users))}
	for id, d := range s.users {
		if d == nil {
			continue
		}
		drafts := make([]domain.PaymentDraft, 0, len(d.drafts))
		for _, draft := range d.drafts {
			drafts = append(drafts, draft)
		}
		out.Users = append(out.Users, snapUser{
			Rev:           d.rev,
			Profile:       d.profile,
			Tokens:        s.tokensForLocked(id),
			Txns:          append([]domain.Transaction(nil), d.txns...),
			Drafts:        drafts,
			Piggy:         d.piggy,
			Contributions: append([]domain.PiggyContribution(nil), d.contributions...),
			Chat:          append([]domain.ChatMessage(nil), d.chatTurns...),
		})
	}
	sort.Slice(out.Users, func(i, j int) bool {
		return out.Users[i].Profile.UserID < out.Users[j].Profile.UserID
	})
	for id, rev := range s.tombstones {
		if id == "" || rev == 0 {
			continue
		}
		out.Deleted = append(out.Deleted, deletedUser{UserID: id, Rev: rev})
	}
	sort.Slice(out.Deleted, func(i, j int) bool {
		return out.Deleted[i].UserID < out.Deleted[j].UserID
	})
	return out
}

func (s *Store) tokensForLocked(id string) []string {
	out := make([]string, 0, 2)
	for tok, uid := range s.tokens {
		if uid == id && tok != "" {
			out = append(out, tok)
		}
	}
	sort.Strings(out)
	return out
}

func (s *Store) mergeLocked(raw []byte) {
	var snap cabinetSnap
	if err := json.Unmarshal(raw, &snap); err != nil {
		log.Printf("cabinet snapshot decode: %v", err)
		return
	}
	if s.tombstones == nil {
		s.tombstones = map[string]int64{}
	}
	for _, tomb := range snap.Deleted {
		if tomb.UserID == "" {
			continue
		}
		if prev, ok := s.tombstones[tomb.UserID]; !ok || tomb.Rev > prev {
			s.tombstones[tomb.UserID] = tomb.Rev
		}
	}
	for _, su := range snap.Users {
		id := su.Profile.UserID
		if id == "" {
			continue
		}
		if rev, dead := s.tombstones[id]; dead && rev >= su.Rev {
			continue
		}
		delete(s.tombstones, id)
		remote := su.toUser()
		local := s.users[id]
		if local == nil || su.Rev >= local.rev {
			s.users[id] = remote
		}
		for _, tok := range su.Tokens {
			if tok != "" {
				s.tokens[tok] = id
			}
		}
		if phone := NormalizePhone(s.users[id].profile.Phone); phone != "" {
			s.phones[phone] = id
		}
	}
	for id, rev := range s.tombstones {
		if d := s.users[id]; d != nil && rev >= d.rev {
			s.dropUserLocked(id)
		}
	}
	if s.mashaTok != "" {
		if _, ok := s.users[MashaUserID]; ok {
			s.tokens[s.mashaTok] = MashaUserID
			s.tokens["demo"] = MashaUserID
			s.phones[NormalizePhone(MashaPhone)] = MashaUserID
		}
	}
}

func (su snapUser) toUser() *userData {
	drafts := map[string]domain.PaymentDraft{}
	for _, d := range su.Drafts {
		if d.ID == "" {
			continue
		}
		drafts[d.ID] = d
	}
	return &userData{
		profile:       su.Profile,
		txns:          append([]domain.Transaction(nil), su.Txns...),
		drafts:        drafts,
		piggy:         su.Piggy,
		contributions: append([]domain.PiggyContribution(nil), su.Contributions...),
		chatTurns:     append([]domain.ChatMessage(nil), su.Chat...),
		rev:           su.Rev,
	}
}

func snapshotCovers(raw []byte, ids map[string]struct{}) bool {
	if len(raw) == 0 {
		return false
	}
	var snap cabinetSnap
	if err := json.Unmarshal(raw, &snap); err != nil {
		return false
	}
	got := map[string]struct{}{}
	for _, u := range snap.Users {
		got[u.Profile.UserID] = struct{}{}
	}
	for id := range ids {
		if _, ok := got[id]; !ok {
			return false
		}
	}
	return true
}

type filePersister struct {
	path string
}

func (f filePersister) Load() ([]byte, error) {
	b, err := os.ReadFile(f.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return b, err
}

func (f filePersister) Save(b []byte) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return err
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, f.path)
}

type blobPersister struct {
	token   string
	storeID string
	client  *http.Client
}

func storeIDFromToken(token string) string {
	parts := strings.Split(token, "_")
	if len(parts) < 4 {
		return ""
	}
	return strings.TrimPrefix(parts[3], "store_")
}

func (b *blobPersister) Load() ([]byte, error) {
	url := fmt.Sprintf("https://%s.private.blob.vercel-storage.com/cabinets.json?cache=0", b.storeID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("blob get status %d", res.StatusCode)
	}
	return body, nil
}

func (b *blobPersister) Save(payload []byte) error {
	url := "https://vercel.com/api/blob/?pathname=cabinets.json"
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("x-api-version", "12")
	req.Header.Set("x-vercel-blob-store-id", b.storeID)
	req.Header.Set("x-vercel-blob-access", "private")
	req.Header.Set("x-content-type", "application/json")
	req.Header.Set("x-add-random-suffix", "0")
	req.Header.Set("x-allow-overwrite", "1")
	req.Header.Set("x-cache-control-max-age", "60")
	res, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resp, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("blob put status %d", res.StatusCode)
	}
	_ = resp
	return nil
}
