# crypto-poller

Небольшое приложение на Go, которое:
- ходит во внешнее API FreeCryptoAPI (https://freecryptoapi.com/);
- забирает котировки по списку монет;
- сохраняет результат в PostgreSQL;
- повторяет опрос по таймеру.

## Что используется
- Go 1.22
- PostgreSQL 16
- pgx/v5
## Переменные окружения
- `DATABASE_URL` — строка подключения к Postgres.
- `FREECRYPTOAPI_KEY` — API key FreeCryptoAPI.
- `POLL_INTERVAL` — интервал опроса, например `30s`.
- `SYMBOLS` — список монет через запятую, например `BTC,ETH,SOL`.
- `VS_CURRENCY` — валюта котировки, например `USD`.
- `APP_PORT` — порт healthcheck-сервера.

## Проверка
- healthcheck: `GET http://localhost:8080/healthz`
- посмотреть данные:
  ```sql
  SELECT symbol, vs_currency, price, fetched_at
  FROM crypto_prices
  ORDER BY fetched_at DESC
  LIMIT 20;
  ```

## Как это работает
Приложение по расписанию вызывает FreeCryptoAPI, парсит ответ, а затем вставляет и нормализованные поля, и исходный JSON в таблицу `crypto_prices`.