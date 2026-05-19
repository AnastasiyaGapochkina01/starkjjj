package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    _ "github.com/lib/pq"
)

type Product struct {
    ID          int64     `json:"id"`
    SKU         string    `json:"sku"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    Category    string    `json:"category"`
    ImageURL    string    `json:"image_url"`
    UpdatedAt   time.Time `json:"updated_at"`
}

var page = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Catalog</title>
<style>
:root { color-scheme: dark; --bg:#0b1020; --panel:#121931; --muted:#8f9bb3; --text:#edf2ff; --line:#27314f; --accent:#4fd1c5; --danger:#f59e0b; }
*{box-sizing:border-box} body{margin:0;font-family:Inter,system-ui,sans-serif;background:linear-gradient(180deg,#0b1020,#111827);color:var(--text)}
.wrapper{max-width:1180px;margin:0 auto;padding:24px}
.header{display:flex;justify-content:space-between;align-items:end;gap:16px;margin-bottom:24px;flex-wrap:wrap}
.h1{font-size:32px;font-weight:800;margin:0}.sub{color:var(--muted);margin-top:6px}
.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(280px,1fr));gap:18px}
.card{background:rgba(18,25,49,.88);border:1px solid var(--line);border-radius:18px;overflow:hidden;box-shadow:0 10px 30px rgba(0,0,0,.24)}
.img{width:100%;height:190px;object-fit:cover;background:#1f2937}.body{padding:18px}.meta{display:flex;justify-content:space-between;gap:10px;color:var(--muted);font-size:13px}
.name{font-size:22px;font-weight:700;margin:10px 0 8px}.desc{color:#c7d2e6;min-height:48px}.footer{display:flex;justify-content:space-between;align-items:center;margin-top:16px;gap:12px}
.price{font-size:24px;font-weight:800}.stock{padding:8px 10px;border-radius:999px;border:1px solid var(--line);font-size:13px}
.ok{color:#86efac}.warn{color:#fcd34d}.empty{padding:40px;border:1px dashed var(--line);border-radius:18px;color:var(--muted)}
.badge{padding:7px 10px;background:rgba(79,209,197,.12);border:1px solid rgba(79,209,197,.35);border-radius:999px;color:var(--accent);font-size:13px}
.code{font-family:ui-monospace,SFMono-Regular,monospace;background:#0a0f1f;border:1px solid var(--line);padding:7px 10px;border-radius:10px;color:#c4b5fd}
</style>
</head>
<body>
<div class="wrapper">
  <div class="header">
    <div>
      <h1 class="h1">Каталог товаров</h1>
      <div class="sub">Go + PostgreSQL. Остатки обновляет отдельный сервис по расписанию.</div>
    </div>
    <div style="display:flex;gap:10px;flex-wrap:wrap;align-items:center">
      <span class="badge">/api/products</span>
      <span class="code">APP_PORT={{.Port}}</span>
    </div>
  </div>
  {{if .Products}}
  <div class="grid">
    {{range .Products}}
    <article class="card">
      <img class="img" src="{{.ImageURL}}" alt="{{.Name}}" loading="lazy">
      <div class="body">
        <div class="meta"><span>{{.Category}}</span><span>{{.SKU}}</span></div>
        <div class="name">{{.Name}}</div>
        <div class="desc">{{.Description}}</div>
        <div class="footer">
          <div class="price">{{printf "%.2f" .Price}} ₽</div>
          <div class="stock {{if lt .Stock 5}}warn{{else}}ok{{end}}">Остаток: {{.Stock}}</div>
        </div>
      </div>
    </article>
    {{end}}
  </div>
  {{else}}
    <div class="empty">В каталоге пока нет товаров.</div>
  {{end}}
</div>
</body>
</html>`))

func env(key, fallback string) string {
    if v := os.Getenv(key); v != "" { return v }
    return fallback
}

func openDB() (*sql.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        env("DB_HOST", "localhost"), env("DB_PORT", "5432"), env("DB_USER", "catalog"), env("DB_PASSWORD", "catalog"), env("DB_NAME", "catalog"))
    db, err := sql.Open("postgres", dsn)
    if err != nil { return nil, err }
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(30 * time.Minute)
    return db, db.Ping()
}

func loadProducts(db *sql.DB) ([]Product, error) {
    rows, err := db.Query(`select id, sku, name, description, price, stock, category, coalesce(image_url, ''), updated_at from products order by id`)
    if err != nil { return nil, err }
    defer rows.Close()
    products := make([]Product, 0)
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.ImageURL, &p.UpdatedAt); err != nil { return nil, err }
        products = append(products, p)
    }
    return products, rows.Err()
}

func main() {
    db, err := openDB()
    if err != nil { log.Fatalf("db: %v", err) }
    defer db.Close()

    port := env("APP_PORT", "8080")

    http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
        products, err := loadProducts(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        _ = json.NewEncoder(w).Encode(products)
    })

    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        products, err := loadProducts(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        _ = page.Execute(w, map[string]any{"Products": products, "Port": strconv.Itoa(mustAtoi(port))})
    })

    log.Printf("catalog app listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustAtoi(v string) int { n, _ := strconv.Atoi(v); return n }
