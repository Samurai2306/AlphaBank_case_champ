package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	RedisURL     string
	DemoToken    string
	CORSOrigins  []string
	LLMBaseURL   string
	LLMAPIKey    string
	LLMModel     string
	UploadDir    string
	DemoOffline  bool
	LogLevel     string
	UseMemory    bool
}

func Load() Config {
	offline := envBool("DEMO_OFFLINE", true)
	db := env("DATABASE_URL", "")
	return Config{
		HTTPAddr:    env("HTTP_ADDR", ":8080"),
		DatabaseURL: db,
		RedisURL:    env("REDIS_URL", ""),
		DemoToken:   env("DEMO_TOKEN", "demo-masha-token"),
		CORSOrigins: splitCSV(env("CORS_ORIGINS", "http://localhost:3000")),
		LLMBaseURL:  env("LLM_BASE_URL", "https://api.openai.com/v1"),
		LLMAPIKey:   env("LLM_API_KEY", ""),
		LLMModel:    env("LLM_MODEL", "gpt-4o-mini"),
		UploadDir:   env("UPLOAD_DIR", "./data/uploads"),
		DemoOffline: offline || env("LLM_API_KEY", "") == "",
		LogLevel:    env("LOG_LEVEL", "info"),
		UseMemory:   envBool("USE_MEMORY_STORE", true) || db == "",
	}
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
