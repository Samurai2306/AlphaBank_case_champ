package agents

import (
	"testing"

	"github.com/alphabank-case-champ/copilot-api/internal/rag"
)

func TestLooksLikeWidgetGarbage(t *testing.T) {
	if !looksLikeWidgetGarbage("Смотрите UnitEconomicsChart ниже") {
		t.Fatal("expected detect widget name")
	}
	if !looksLikeWidgetGarbage("график: (вставка графика)") {
		t.Fatal("expected detect placeholder")
	}
	if looksLikeWidgetGarbage("Точка безубыточности 2 клиента в день при аренде 40 000 ₽") {
		t.Fatal("clean answer flagged")
	}
}

func TestGroundedFallbackHasNextSteps(t *testing.T) {
	hits := []rag.Hit{
		{ID: "kb:npd", Title: "НПД", Text: "Ставка 4% и 6%.", Source: "demo"},
		{ID: "kb:enp", Title: "ЕНП", Text: "Уплата до 28-го числа.", Source: "demo"},
	}
	text := groundedFallbackAnswer(hits)
	if !containsAny(normalizeRU(text), "посчитать налог", "копил") {
		t.Fatalf("expected next steps: %s", text)
	}
}
