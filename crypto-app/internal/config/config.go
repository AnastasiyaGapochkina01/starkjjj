package config

import (
    "os"
    "strings"
    "time"
)

type Config struct {
    AppPort       string
    DatabaseURL   string
    APIKey        string
    PollInterval  time.Duration
    Symbols       []string
    VsCurrency    string
}

func Load() Config {
    interval := 30 * time.Second
    if v := os.Getenv("POLL_INTERVAL"); v != "" {
        if parsed, err := time.ParseDuration(v); err == nil {
            interval = parsed
        }
    }

    symbols := []string{"BTC", "ETH", "SOL"}
    if v := os.Getenv("SYMBOLS"); v != "" {
        parts := strings.Split(v, ",")
        symbols = make([]string, 0, len(parts))
        for _, p := range parts {
            p = strings.TrimSpace(strings.ToUpper(p))
            if p != "" {
                symbols = append(symbols, p)
            }
        }
    }

    vs := os.Getenv("VS_CURRENCY")
    if vs == "" {
        vs = "USD"
    }

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    return Config{
        AppPort:      port,
        DatabaseURL:  os.Getenv("DATABASE_URL"),
        APIKey:       os.Getenv("FREECRYPTOAPI_KEY"),
        PollInterval: interval,
        Symbols:      symbols,
        VsCurrency:   strings.ToUpper(vs),
    }
}