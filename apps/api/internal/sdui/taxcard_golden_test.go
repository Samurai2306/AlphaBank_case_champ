package sdui_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/sdui"
)

func TestTaxCardGolden(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "packages", "shared", "fixtures", "tax_card.json")
	raw, err := os.ReadFile(root)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	var env sdui.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Component != "TaxCard" {
		t.Fatalf("component=%s", env.Component)
	}
	amount, ok := env.Props["tax_amount"].(float64)
	if !ok || amount != 7740 {
		t.Fatalf("tax_amount=%v", env.Props["tax_amount"])
	}
	ceiling, ok := env.Props["ceiling_amount"].(float64)
	if !ok || ceiling != 10800 {
		t.Fatalf("ceiling_amount=%v", env.Props["ceiling_amount"])
	}
}
