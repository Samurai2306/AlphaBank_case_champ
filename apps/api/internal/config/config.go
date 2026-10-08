package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/llm"
)

type Config struct {
	HTTPAddr    string
	DemoToken   string
	CORSOrigins []string
	LLMBaseURL  string
	LLMAPIKey   string
	LLMModel    string
	LLMProvider string
	UploadDir   string
	DemoOffline bool
}

func Load() Config {
	key := env("LLM_API_KEY", "")
	base := env("LLM_BASE_URL", "")
	provider := env("LLM_PROVIDER", "")
	model := env("LLM_MODEL", "")
	vercel := os.Getenv("VERCEL") != ""

	// Hosted demo (Vercel): turn on free OpenCode Zen unless DEMO_OFFLINE=1.
	if vercel && !envIsTrue("DEMO_OFFLINE") {
		if key == "" {
			key = "public"
		}
		if base == "" || base == "https://api.openai.com/v1" {
			base = "https://opencode.ai/zen/v1"
		}
		if provider == "" {
			provider = "opencode-zen"
		}
		if model == "" || model == "gpt-4o-mini" || model == "deepseek-v4-flash-free" {
			model = "big-pickle"
		}
	}
	if provider == "" && strings.Contains(strings.ToLower(base), "opencode.ai/zen") {
		provider = "opencode-zen"
	}
	if key == "" && provider == "opencode-zen" {
		key = "public"
	}
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4o-mini"
	}

	offlineDefault := true
	if vercel && (key != "" || gatewayLikely()) {
		offlineDefault = false
	}
	offline := envBool("DEMO_OFFLINE", offlineDefault)
	if key == "" && !gatewayLikely() {
		offline = true
	}

	return Config{
		HTTPAddr:    httpAddr(),
		DemoToken:   env("DEMO_TOKEN", "demo-masha-token"),
		CORSOrigins: splitCSV(env("CORS_ORIGINS", "http://localhost:3000")),
		LLMBaseURL:  base,
		LLMAPIKey:   key,
		LLMModel:    model,
		LLMProvider: provider,
		UploadDir:   uploadDir(),
		DemoOffline: offline,
	}
}

// httpAddr prefers HTTP_ADDR, then PORT (Vercel/Railway/Render), then :8080.
func httpAddr() string {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		return v
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func uploadDir() string {
	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		return v
	}
	if os.Getenv("VERCEL") != "" {
		return "/tmp/uploads"
	}
	return "./data/uploads"
}

// LLMEndpoints builds the failover chain: primary model(s) + optional second provider.
func (c Config) LLMEndpoints() []llm.Endpoint {
	var out []llm.Endpoint
	// No-key OpenAI-compatible route. Zen is locked to the official client,
	// AI Gateway needs a card, and Pollinations answers 402 from Vercel.
	if !envIsTrue("LLM_DISABLE_VIREONIX") {
		out = append(out, llm.Endpoint{
			Name: "vireonix", BaseURL: "https://vireonix.ai/v1",
			APIKey: "none", Model: "auto", Provider: "vireonix",
		})
	}
	if !envIsTrue("LLM_DISABLE_POLLINATIONS") {
		polModels := splitCSV(env("POLLINATIONS_MODELS", ""))
		if len(polModels) == 0 {
			// openai-fast answers inside a serverless budget. The larger
			// "openai" model often returns 402 from Vercel before any text.
			polModels = []string{"openai-fast"}
		}
		for i, m := range polModels {
			name := "pollinations"
			if i > 0 {
				name = "pollinations-" + sanitizeName(m)
			}
			out = append(out, llm.Endpoint{
				Name: name, BaseURL: "https://text.pollinations.ai/v1",
				APIKey: "anonymous", Model: m, Provider: "pollinations",
			})
		}
	}
	if gatewayLikely() {
		gwModels := splitCSV(env("AI_GATEWAY_MODELS", ""))
		if len(gwModels) == 0 {
			gwModels = []string{"google/gemini-2.5-flash-lite", "openai/gpt-5-nano", "google/gemini-2.5-flash"}
		}
		for i, m := range gwModels {
			name := "gateway"
			if i > 0 {
				name = "gateway-" + sanitizeName(m)
			}
			out = append(out, llm.Endpoint{
				Name: name, BaseURL: "https://ai-gateway.vercel.sh/v1",
				APIKey: gatewayKeyPlaceholder(), Model: m, Provider: "vercel-ai-gateway",
			})
		}
	}

	key := c.LLMAPIKey
	base := c.LLMBaseURL
	provider := c.LLMProvider
	if provider == "" {
		provider = env("LLM_PROVIDER", "")
	}
	if provider == "" && strings.Contains(strings.ToLower(base), "opencode.ai/zen") {
		provider = "opencode-zen"
	}
	if key == "" && provider == "opencode-zen" {
		key = "public"
	}
	if key != "" && base != "" {
		models := splitCSV(env("LLM_MODELS", ""))
		if len(models) == 0 {
			if provider == "opencode-zen" {
				models = []string{"big-pickle", "mimo-v2.5-free", "nemotron-3.5-lightning-free", "ling-3.0-flash-fin-free"}
			} else {
				models = []string{c.LLMModel}
			}
		}
		for i, m := range models {
			name := "primary"
			if i > 0 {
				name = "primary-" + sanitizeName(m)
			}
			out = append(out, llm.Endpoint{
				Name: name, BaseURL: base, APIKey: key, Model: m, Provider: provider,
			})
		}
	}
	// Optional second provider (Pollinations / DeepSeek / etc.).
	fbURL := env("LLM_FALLBACK_BASE_URL", "")
	fbKey := env("LLM_FALLBACK_API_KEY", "")
	fbModel := env("LLM_FALLBACK_MODEL", "")
	fbProvider := env("LLM_FALLBACK_PROVIDER", "")
	fbModels := splitCSV(env("LLM_FALLBACK_MODELS", ""))
	if fbURL != "" && fbKey != "" {
		if len(fbModels) == 0 {
			if fbModel == "" {
				fbModel = "openai-fast"
			}
			fbModels = []string{fbModel}
		}
		for i, m := range fbModels {
			name := "fallback"
			if i > 0 {
				name = "fallback-" + sanitizeName(m)
			}
			out = append(out, llm.Endpoint{
				Name: name, BaseURL: fbURL, APIKey: fbKey, Model: m, Provider: fbProvider,
			})
		}
	}
	return out
}

func sanitizeName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, " ", "-")
	if len(s) > 32 {
		return s[:32]
	}
	return s
}

func gatewayLikely() bool {
	return os.Getenv("VERCEL") != "" ||
		os.Getenv("AI_GATEWAY_API_KEY") != "" ||
		os.Getenv("VERCEL_OIDC_TOKEN") != ""
}

func gatewayKeyPlaceholder() string {
	if v := os.Getenv("AI_GATEWAY_API_KEY"); v != "" {
		return v
	}
	if v := os.Getenv("VERCEL_OIDC_TOKEN"); v != "" {
		return v
	}
	return "oidc"
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envBool(k string, def bool) bool {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envIsTrue(k string) bool {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
