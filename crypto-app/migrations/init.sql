CREATE TABLE IF NOT EXISTS crypto_prices (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    vs_currency TEXT NOT NULL,
    price NUMERIC(20,8) NOT NULL,
    market_cap NUMERIC(24,2),
    volume_24h NUMERIC(24,2),
    change_24h NUMERIC(10,4),
    source TEXT NOT NULL DEFAULT 'freecryptoapi',
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_crypto_prices_symbol_fetched_at
    ON crypto_prices (symbol, fetched_at DESC);