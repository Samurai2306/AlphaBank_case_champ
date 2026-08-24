package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFailoverSkipsBadEndpoint(t *testing.T) {
	var hits []string
	mux := http.NewServeMux()
	mux.HandleFunc("/bad/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		hits = append(hits, "bad|"+r.Header.Get("User-Agent"))
		http.Error(w, "nope", http.StatusBadGateway)
	})
	mux.HandleFunc("/good/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		hits = append(hits, "good|"+r.Header.Get("User-Agent"))
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": "ok-from-good"}},
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := NewWithFailover([]Endpoint{
		{Name: "bad", BaseURL: srv.URL + "/bad", APIKey: "k", Model: "m"},
		{Name: "good", BaseURL: srv.URL + "/good", APIKey: "k", Model: "m"},
	})
	c.HTTP = &http.Client{Timeout: 5 * time.Second}

	text, err := c.CompleteTemp(context.Background(), "sys", "user", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok-from-good" {
		t.Fatalf("got %q", text)
	}
	if len(hits) < 2 {
		t.Fatalf("expected failover hits, got %v", hits)
	}
	if !strings.Contains(hits[0], "Chrome/") || !strings.Contains(hits[1], "Chrome/") {
		t.Fatalf("browser user-agent missing: %v", hits)
	}
	st := c.Status()
	if st["last_ok"] != "good" {
		t.Fatalf("status=%v", st)
	}
}

func TestConfiguredRequiresKeyAndURL(t *testing.T) {
	c := NewWithFailover([]Endpoint{
		{Name: "empty", BaseURL: "", APIKey: "k", Model: "m"},
		{Name: "ok", BaseURL: "https://example.com/v1", APIKey: "k", Model: "m"},
	})
	if !c.Configured() {
		t.Fatal("expected configured")
	}
	if len(c.activeEndpoints()) != 1 {
		t.Fatalf("want 1 endpoint")
	}
}

func TestOpenCodeZenHeaders(t *testing.T) {
	var gotAuth, gotUA, gotClient string
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		gotClient = r.Header.Get("x-opencode-client")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": "НПД — налог для самозанятых."}},
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := NewWithFailover([]Endpoint{{
		Name: "zen", BaseURL: srv.URL + "/v1", APIKey: "public",
		Model: "deepseek-v4-flash-free", Provider: "opencode-zen",
	}})
	c.HTTP = &http.Client{Timeout: 5 * time.Second}
	text, err := c.CompleteTemp(context.Background(), "sys", "user", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "НПД") {
		t.Fatalf("text=%q", text)
	}
	if gotAuth != "Bearer public" || gotClient != "cli" || !strings.Contains(gotUA, "opencode/") {
		t.Fatalf("headers auth=%q ua=%q client=%q", gotAuth, gotUA, gotClient)
	}
}
