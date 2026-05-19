package freecrypto

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

type Client struct {
    apiKey     string
    httpClient *http.Client
    baseURL    string
}

type Asset struct {
    Symbol    string
    Price     float64
    MarketCap float64
    Volume24h float64
    Change24h float64
    Raw       map[string]any
}

func New(apiKey string) *Client {
    return &Client{
        apiKey: apiKey,
        baseURL: "https://api.freecryptoapi.com/v1",
        httpClient: &http.Client{Timeout: 15 * time.Second},
    }
}

func (c *Client) GetData(ctx context.Context, symbols []string, vsCurrency string) ([]Asset, error) {
    q := url.Values{}
    q.Set("symbol", strings.Join(symbols, ","))
    q.Set("convert", strings.ToUpper(vsCurrency))

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getData?"+q.Encode(), nil)
    if err != nil {
        return nil, err
    }
    if c.apiKey != "" {
        req.Header.Set("x-api-key", c.apiKey)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode >= 300 {
        return nil, fmt.Errorf("freecryptoapi status %d: %s", resp.StatusCode, string(body))
    }

    var payload map[string]any
    if err := json.Unmarshal(body, &payload); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    return normalize(symbols, payload), nil
}

func normalize(requested []string, payload map[string]any) []Asset {
    out := make([]Asset, 0, len(requested))
    for _, symbol := range requested {
        upper := strings.ToUpper(symbol)
        raw, ok := payload[upper]
        if !ok {
            continue
        }
        item, ok := raw.(map[string]any)
        if !ok {
            continue
        }
        out = append(out, Asset{
            Symbol:    upper,
            Price:     pickFloat(item, "price", "price_usd", "value"),
            MarketCap: pickFloat(item, "market_cap", "marketCap"),
            Volume24h: pickFloat(item, "volume_24h", "volume24h", "volume"),
            Change24h: pickFloat(item, "change_24h", "change24h", "percent_change_24h"),
            Raw:       item,
        })
    }
    return out
}

func pickFloat(m map[string]any, keys ...string) float64 {
    for _, key := range keys {
        v, ok := m[key]
        if !ok || v == nil {
            continue
        }
        switch t := v.(type) {
        case float64:
            return t
        case int:
            return float64(t)
        }
    }
    return 0
}