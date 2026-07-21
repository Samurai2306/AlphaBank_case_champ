package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
	"github.com/alphabank-case-champ/copilot-api/internal/config"
	"github.com/alphabank-case-champ/copilot-api/internal/httpserver"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	store := memory.NewSeededWithToken(cfg.DemoToken)
	_ = os.MkdirAll(cfg.UploadDir, 0o755)

	h := httpserver.New(cfg, store, log)
	log.Info("copilot-api starting", "addr", cfg.HTTPAddr, "offline", cfg.DemoOffline)
	if err := http.ListenAndServe(cfg.HTTPAddr, h); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
