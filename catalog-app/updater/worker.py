import os
import random
import time
from datetime import datetime, timezone
from croniter import croniter
import psycopg


def env(name, default):
    return os.getenv(name, default)

DB_DSN = (
    f"host={env('DB_HOST', 'localhost')} "
    f"port={env('DB_PORT', '5432')} "
    f"dbname={env('DB_NAME', 'catalog')} "
    f"user={env('DB_USER', 'catalog')} "
    f"password={env('DB_PASSWORD', 'catalog')}"
)
SCHEDULE = env('SCHEDULE', '*/1 * * * *')
HANG_PATTERN = [x.strip() in {'1', 'true', 'yes', 'y'} for x in env('HANG_PATTERN', '1,1,1,0,0').split(',') if x.strip()]
HANG_SECONDS = int(env('HANG_SECONDS', '180'))
CPU_BURN_SECONDS = int(env('CPU_BURN_SECONDS', '45'))
STOCK_MIN_DELTA = int(env('STOCK_MIN_DELTA', '-3'))
STOCK_MAX_DELTA = int(env('STOCK_MAX_DELTA', '7'))
BATCH_SIZE = int(env('BATCH_SIZE', '8'))
RUN_NO = 0


def wait_until(ts: datetime):
    while True:
        now = datetime.now(timezone.utc)
        diff = (ts - now).total_seconds()
        if diff <= 0:
            return
        time.sleep(min(diff, 1))


def hog_host(seconds: int):
    print(f"[updater] hanging for {seconds}s and burning cpu", flush=True)
    deadline = time.time() + seconds
    burn_until = time.time() + min(seconds, CPU_BURN_SECONDS)
    x = 0
    while time.time() < burn_until:
        x = (x * 3 + random.randint(1, 9)) % 10000019
    while time.time() < deadline:
        time.sleep(1)
    return x


def refresh_stock(conn):
    with conn.cursor() as cur:
        cur.execute("select id, stock from products order by random() limit %s", (BATCH_SIZE,))
        rows = cur.fetchall()
        for product_id, stock in rows:
            delta = random.randint(STOCK_MIN_DELTA, STOCK_MAX_DELTA)
            new_stock = max(0, stock + delta)
            cur.execute(
                "update products set stock = %s, updated_at = now() where id = %s",
                (new_stock, product_id),
            )
    conn.commit()


def should_hang(run_no: int) -> bool:
    if not HANG_PATTERN:
        return False
    return HANG_PATTERN[(run_no - 1) % len(HANG_PATTERN)]


def main():
    global RUN_NO
    itr = croniter(SCHEDULE, datetime.now(timezone.utc))
    while True:
        next_run = itr.get_next(datetime)
        print(f"[updater] next run at {next_run.isoformat()}", flush=True)
        wait_until(next_run)
        RUN_NO += 1
        print(f"[updater] run #{RUN_NO} started", flush=True)
        if should_hang(RUN_NO):
            hog_host(HANG_SECONDS)
        with psycopg.connect(DB_DSN, autocommit=False) as conn:
            refresh_stock(conn)
        print(f"[updater] run #{RUN_NO} finished", flush=True)


if __name__ == '__main__':
    main()
