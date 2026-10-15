package memory

import (
	"path/filepath"
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/domain"
)

func TestAdminDirectoryIncludesMasha(t *testing.T) {
	s := NewSeeded()
	items := s.AdminDirectory()
	if len(items) != 1 {
		t.Fatalf("users=%d", len(items))
	}
	if items[0].Profile.UserID != MashaUserID {
		t.Fatalf("user=%s", items[0].Profile.UserID)
	}
	if items[0].TxnCount == 0 || items[0].PiggyBalance <= 0 {
		t.Fatalf("card=%+v", items[0])
	}
	detail, ok := s.AdminUser(MashaUserID)
	if !ok || detail.DraftCount == 0 {
		t.Fatal("expected masha drafts")
	}
	if _, ok := s.AdminUser("missing"); ok {
		t.Fatal("missing user should be absent")
	}
}

func TestRegisterSurvivesNewStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cabinets.json")
	t.Setenv("CABINET_FILE", path)
	t.Setenv("BLOB_READ_WRITE_TOKEN", "")

	first := NewSeeded()
	if err := first.AttachPersisterFromEnv(); err != nil {
		t.Fatal(err)
	}
	token, profile, errMsg := first.Register(RegisterInput{
		DisplayName: "Лена", Phone: "+7 900 000-11-22", City: "Москва", BusinessSphere: "Бьюти",
	})
	if errMsg != "" || token == "" {
		t.Fatalf("register: %s", errMsg)
	}

	second := NewSeeded()
	if err := second.AttachPersisterFromEnv(); err != nil {
		t.Fatal(err)
	}
	items := second.AdminDirectory()
	if len(items) != 2 {
		t.Fatalf("users=%d", len(items))
	}
	found := false
	for _, item := range items {
		if item.Profile.UserID == profile.UserID && item.Profile.DisplayName == "Лена" {
			found = true
		}
	}
	if !found {
		t.Fatal("lena missing after reload")
	}
	uid, ok := second.ResolveToken(token)
	if !ok || uid != profile.UserID {
		t.Fatalf("token resolve uid=%s ok=%v", uid, ok)
	}
	if _, _, errMsg := second.Register(RegisterInput{DisplayName: "Лена", Phone: "89000001122"}); errMsg == "" {
		t.Fatal("same phone should be rejected")
	}
}

func TestAdminMutations(t *testing.T) {
	s := NewSeeded()
	_, profile, errMsg := s.Register(RegisterInput{
		DisplayName: "Лена", Phone: "+7 900 111-22-33", City: "Москва", BusinessSphere: "Брови",
	})
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	name := "Елена"
	phone := "+7 900 111-22-44"
	city := "Казань"
	regime := "IP"
	level := 3
	updated, errMsg := s.AdminUpdateProfile(profile.UserID, AdminProfilePatch{
		DisplayName: &name, Phone: &phone, City: &city, TaxRegime: &regime, CJMLevel: &level,
	})
	if errMsg != "" || updated.DisplayName != "Елена" || updated.TaxRegime != "IP" || updated.CJMLevel != 3 {
		t.Fatalf("update profile=%+v err=%s", updated, errMsg)
	}
	if _, _, errMsg := s.Register(RegisterInput{DisplayName: "Другая", Phone: phone}); errMsg == "" {
		t.Fatal("updated phone should be taken")
	}

	draft, ok := s.AdminCreateDraft(profile.UserID, domain.PaymentDraft{Amount: 1500, Purpose: "Аренда кресла"})
	if !ok || draft.Status != "draft" {
		t.Fatalf("draft=%+v ok=%v", draft, ok)
	}
	confirmed, errMsg := s.AdminConfirmDraft(profile.UserID, draft.ID)
	if errMsg != "" || confirmed.Status != "confirmed" {
		t.Fatalf("confirm status=%s err=%s", confirmed.Status, errMsg)
	}
	if _, errMsg := s.AdminCancelDraft(profile.UserID, draft.ID); errMsg != "confirmed" {
		t.Fatalf("cancel confirmed err=%s", errMsg)
	}
	piggy, ok := s.AdminSetPiggy(profile.UserID, true)
	if !ok || !piggy.Enabled {
		t.Fatal("piggy should turn on")
	}
	if !s.AdminDelete(profile.UserID) {
		t.Fatal("delete")
	}
	if _, ok := s.AdminUser(profile.UserID); ok {
		t.Fatal("deleted user still visible")
	}
	if _, _, errMsg := s.Register(RegisterInput{DisplayName: "Снова", Phone: phone}); errMsg != "" {
		t.Fatalf("phone should free after delete: %s", errMsg)
	}
}

func TestAdminDeleteSurvivesReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cabinets.json")
	t.Setenv("CABINET_FILE", path)
	t.Setenv("BLOB_READ_WRITE_TOKEN", "")

	first := NewSeeded()
	if err := first.AttachPersisterFromEnv(); err != nil {
		t.Fatal(err)
	}
	_, profile, errMsg := first.Register(RegisterInput{
		DisplayName: "Оля", Phone: "+7 900 222-33-44",
	})
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if !first.AdminDelete(profile.UserID) || !first.AdminDelete(MashaUserID) {
		t.Fatal("delete")
	}

	second := NewSeeded()
	if err := second.AttachPersisterFromEnv(); err != nil {
		t.Fatal(err)
	}
	if _, ok := second.AdminUser(profile.UserID); ok {
		t.Fatal("deleted cabinet restored")
	}
	if _, ok := second.AdminUser(MashaUserID); ok {
		t.Fatal("deleted seed restored")
	}
	if _, ok := second.ResolveToken(second.MashaToken()); ok {
		t.Fatal("masha token still resolves")
	}
}
