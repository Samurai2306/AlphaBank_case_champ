package rag

import (
	"strings"
	"testing"
)

func TestLoadEmbeddedCorpus(t *testing.T) {
	idx, err := LoadEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	if idx.Len() < 80 {
		t.Fatalf("want at least 80 chunks, got %d", idx.Len())
	}
}

func TestRetrieveGKAnd422(t *testing.T) {
	hits := Default().Retrieve("Статья 23 ГК РФ ИП и закон 422-ФЗ самозанятые", 5)
	if len(hits) == 0 {
		t.Fatal("expected law hits")
	}
	joined := ""
	for _, h := range hits {
		joined += h.ID + " " + h.Text
	}
	if !strings.Contains(joined, "23") && !strings.Contains(strings.ToLower(joined), "422") {
		t.Fatalf("expected GK/422 chunk in hits: %+v", hits)
	}
}

func TestRetrieveAlfaRKO(t *testing.T) {
	hits := Default().Retrieve("тарифы РКО Альфа-Банк ноль за обслуживание", 3)
	if len(hits) == 0 {
		t.Fatal("expected Alfa RKO hits")
	}
}

func TestRetrieveNPD(t *testing.T) {
	idx := Default()
	hits := idx.Retrieve("Что такое НПД и какая ставка?", 3)
	if LowConfidence(hits) {
		t.Fatalf("expected confident hits, got %#v", hits)
	}
	found := false
	for _, h := range hits {
		if h.ID == "kb:npd-rates-demo#2026" || strings.Contains(h.Text, "4%") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected NPD rates in hits: %+v", hits)
	}
}

func TestRetrieveCashback(t *testing.T) {
	hits := Default().Retrieve("Как работает кэшбек категорий?", 3)
	if len(hits) == 0 {
		t.Fatal("expected cashback hits")
	}
}
