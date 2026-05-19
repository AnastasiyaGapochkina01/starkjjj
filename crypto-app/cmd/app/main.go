package main

import (
    "context"
    "log"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "crypto-poller/internal/config"
    "crypto-poller/internal/db"
    "crypto-poller/internal/freecrypto"
    "crypto-poller/internal/service"
)

func main() {
    cfg := config.Load()
    if cfg.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is required")
    }
    if cfg.APIKey == "" {
        slog.Warn("FREECRYPTOAPI_KEY is empty; requests may fail depending on plan")
    }

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    pool, err := db.New(ctx, cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("db init error: %v", err)
    }
    defer pool.Close()

    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            _, _ = w.Write([]byte("ok"))
        })
        addr := ":" + cfg.AppPort
        slog.Info("health server started", "addr", addr)
        if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
            slog.Error("http server error", "error", err)
            os.Exit(1)
        }
    }()

    poller := &service.Poller{
        DB:         pool,
        Client:     freecrypto.New(cfg.APIKey),
        Symbols:    cfg.Symbols,
        VsCurrency: cfg.VsCurrency,
        Interval:   cfg.PollInterval,
    }

    if err := poller.Run(ctx); err != nil && err != context.Canceled {
        log.Fatalf("poller stopped: %v", err)
    }
}