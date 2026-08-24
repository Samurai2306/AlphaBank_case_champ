package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/alphabank-case-champ/copilot-api/internal/llm"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisURL    string
	DemoToken   string
	CORSOrigins []string
	LLMBaseURL  string
	LLMAPIKey   string
	LLMModel    string
	UploadDir   string
	DemoOffline bool
	LogLevel    string
	UseMemory   bool
}

func Load() Config {
	offline := envBool("DEMO_OFFLINE", true)
	db := env("DATABASE_URL", "")
	key := env("LLM_API_KEY", "")
	return Config{
		HTTPAddr:    env("HTTP_ADDR", ":8080"),
		DatabaseURL: db,
		RedisURL:    env("REDIS_URL", ""),
		DemoToken:   env("DEMO_TOKEN", "demo-masha-token"),
		CORSOrigins: splitCSV(env("CORS_ORIGINS", "http://localhost:3000")),
		LLMBaseURL:  env("LLM_BASE_URL", "https://api.openai.com/v1"),
		LLMAPIKey:   key,
		LLMModel:    env("LLM_MODEL", "gpt-4o-mini"),
		UploadDir:   env("UPLOAD_DIR", "./data/uploads"),
		DemoOffline: offline || key == "",
		LogLevel:    env("LOG_LEVEL", "info"),
		UseMemory:   envBool("USE_MEMORY_STORE", true) || db == "",
	}
}

// LLMEndpoints builds the failover chain: primary model(s) + optional second provider.
func (c Config) LLMEndpoints() []llm.Endpoint {
	if c.LLMAPIKey == "" || c.LLMBaseURL == "" {
		return nil
	}
	models := splitCSV(env("LLM_MODELS", ""))
	if len(models) == 0 {
		models = []string{c.LLMModel}
	}
	provider := env("LLM_PROVIDER", "")
	var out []llm.Endpoint
	for i, m := range models {
		name := "primary"
		if i > 0 {
			name = "primary-" + sanitizeName(m)
		}
		out = append(out, llm.Endpoint{
			Name: name, BaseURL: c.LLMBaseURL, APIKey: c.LLMAPIKey, Model: m, Provider: provider,
		})
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
