package service

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "strings"
    "time"

    "crypto-poller/internal/freecrypto"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Poller struct {
    DB         *pgxpool.Pool
    Client     *freecrypto.Client
    Symbols    []string
    VsCurrency string
    Interval   time.Duration
}

func (p *Poller) Run(ctx context.Context) error {
    if err := p.fetchAndStore(ctx); err != nil {
        slog.Error("initial fetch failed", "error", err)
    }

    ticker := time.NewTicker(p.Interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            if err := p.fetchAndStore(ctx); err != nil {
                slog.Error("poll failed", "error", err)
            }
        }
    }
}

func (p *Poller) fetchAndStore(ctx context.Context) error {
    assets, err := p.Client.GetData(ctx, p.Symbols, p.VsCurrency)
    if err != nil {
        return fmt.Errorf("get data: %w", err)
    }

    for _, asset := range assets {
        raw, err := json.Marshal(asset.Raw)
        if err != nil {
            return fmt.Errorf("marshal raw for %s: %w", asset.Symbol, err)
        }

        _, err = p.DB.Exec(ctx, `
            INSERT INTO crypto_prices (
                symbol, vs_currency, price, market_cap, volume_24h, change_24h, raw
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, strings.ToUpper(asset.Symbol), p.VsCurrency, asset.Price, asset.MarketCap, asset.Volume24h, asset.Change24h, raw)
        if err != nil {
            return fmt.Errorf("insert %s: %w", asset.Symbol, err)
        }
    }

    slog.Info("prices saved", "count", len(assets), "symbols", strings.Join(p.Symbols, ","))
    return nil
}