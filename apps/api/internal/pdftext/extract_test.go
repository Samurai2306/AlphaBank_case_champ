package pdftext

import (
	"strings"
	"testing"
)

func TestExtractSampleContract(t *testing.T) {
	text := Extract(SampleContractBytes())
	if !strings.Contains(strings.ToLower(text), "аренд") {
		t.Fatalf("expected rental text, got %q", text)
	}
	if !HasText(SampleContractBytes()) {
		t.Fatal("HasText")
	}
}

func TestExtractDetailedMarkdown(t *testing.T) {
	res := ExtractDetailed(SampleContractBytes())
	if res.Chars < 40 {
		t.Fatalf("chars=%d method=%s warnings=%v", res.Chars, res.Method, res.Warnings)
	}
	if !strings.Contains(res.Markdown, "Извлечённый текст") {
		t.Fatalf("markdown envelope missing: %q", res.Markdown)
	}
	if res.Quality == "poor" {
		t.Fatalf("expected usable quality, got %s", res.Quality)
	}
}

func TestFormatForLLM(t *testing.T) {
	md := FormatForLLM("Договор аренды кабинета", 1, "heuristic")
	if !strings.Contains(md, "heuristic") || !strings.Contains(md, "Договор") {
		t.Fatalf("bad markdown: %q", md)
	}
}
