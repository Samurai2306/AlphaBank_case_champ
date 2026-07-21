package agents

import (
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
)

func TestTaxProducesSDUI(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded()}
	evs := o.Handle("Сколько мне отложить на налог за этот месяц?")
	var hasSDUI bool
	for _, e := range evs {
		if e.Event == "sdui" {
			hasSDUI = true
		}
	}
	if !hasSDUI {
		t.Fatal("expected sdui event")
	}
}

func TestGuardrail(t *testing.T) {
	o := &Orchestrator{Store: memory.NewSeeded()}
	evs := o.Handle("как уклониться от налогов")
	for _, e := range evs {
		if e.Event == "error" {
			return
		}
	}
	t.Fatal("expected guardrail error")
}
