package config

import "testing"

func TestHTTPAddrPrefersHTTPAddr(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("PORT", "3000")
	if got := httpAddr(); got != ":9090" {
		t.Fatalf("httpAddr() = %q", got)
	}
}

func TestHTTPAddrUsesPORT(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("PORT", "3333")
	if got := httpAddr(); got != ":3333" {
		t.Fatalf("httpAddr() = %q", got)
	}
}

func TestUploadDirOnVercel(t *testing.T) {
	t.Setenv("UPLOAD_DIR", "")
	t.Setenv("VERCEL", "1")
	if got := uploadDir(); got != "/tmp/uploads" {
		t.Fatalf("uploadDir() = %q", got)
	}
}

func TestVercelEnablesZenLLM(t *testing.T) {
	t.Setenv("VERCEL", "1")
	t.Setenv("DEMO_OFFLINE", "")
	t.Setenv("LLM_API_KEY", "")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_PROVIDER", "")
	t.Setenv("LLM_MODEL", "")
	t.Setenv("LLM_MODELS", "")
	t.Setenv("DATABASE_URL", "")
	cfg := Load()
	if cfg.DemoOffline {
		t.Fatal("expected live LLM on Vercel")
	}
	if cfg.LLMAPIKey != "public" || cfg.LLMProvider != "opencode-zen" {
		t.Fatalf("key=%q provider=%q", cfg.LLMAPIKey, cfg.LLMProvider)
	}
	eps := cfg.LLMEndpoints()
	if len(eps) < 3 {
		t.Fatalf("want zen model failover, got %d", len(eps))
	}
}

func TestExplicitOfflineWinsOnVercel(t *testing.T) {
	t.Setenv("VERCEL", "1")
	t.Setenv("DEMO_OFFLINE", "1")
	t.Setenv("LLM_API_KEY", "public")
	cfg := Load()
	if !cfg.DemoOffline {
		t.Fatal("DEMO_OFFLINE=1 must stay offline")
	}
}
