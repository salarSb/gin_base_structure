# Base‑Structure — Gin • Gorm • Redis • PostgreSQL 🚀

A production‑ready Golang starter that ships with:

* **Gin** HTTP framework
* **Gorm** ORM (PostgreSQL driver, auto‑migrations)
* JWT authentication with Redis blacklist / revocation
* **Zap** structured logging (JSON + Lumberjack rotation)
* Swagger / OpenAPI 3 docs (Swaggo)
* Docker services (Postgres, Redis, PgAdmin)
* Singleton config loaded via **`--config` flag** or **`CONFIG_FILE`** env‑var

---

## ✨ Highlights

| Layer | Library / Tool | Purpose |
|-------|----------------|---------|
| HTTP server | **Gin‑Gonic** | Fast router & middleware |
| ORM | **Gorm** | PostgreSQL driver, migrations |
| Cache / rate‑limit | **Redis 7** | OTP + token revocation |
| Auth | **golang‑jwt/jwt** | Access & refresh tokens |
| Validation | **validator/v10** | Custom tags: `ir_mobile`, `password` |
| Config | **Viper** + dotenv | Singleton, CLI/env override |
| Logging | **Zap** + Lumberjack | JSON logs, rotation |
| Docs | **Swaggo** + Swagger‑UI | Live at `/swagger/` |

---

## 🚀 Quick start (local dev)

```bash
git clone <repo_url> my-api && cd my-api

# 1) secrets + config templates (never commit real values)
cp .env.example .env
cp src/config/config-development-example.yml src/config/config-development.yml
cp docker/redis/redis_example.conf docker/redis/redis.conf

# 2) start Postgres, Redis, PgAdmin
docker compose -f docker/docker-compose.yml up -d

# 3) Go deps + Swagger CLI
go mod download
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./src/cmd/main.go -o ./docs        # generate docs

# 4) run the API (port 5005 by default)
go run ./src/cmd --config src/config/config-development.yml
# → http://localhost:5005/swagger/
```

---

## 🔧 Configuration

### Selecting the config file

Priority order:

1. **CLI flag** — `--config /path/to/app.yml`
2. **Environment variable** — `CONFIG_FILE=/app/config/app.yml`

If neither is set the binary exits with: *"no config specified"*.

`.env` is auto‑located by walking up from the current directory (works from any CWD, tests, or containers).

---

## 🐳 Docker‑Compose stack

| Service | Host → Container | Volume |
|---------|------------------|--------|
| postgres | 5432 → 5432 | `postgres` |
| pgAdmin  | 8090 → 80  | `pgadmin` |
| redis    | 6379 → 6379 | `redis` |

```bash
docker compose -f docker/docker-compose.yml up -d   # start
docker compose -f docker/docker-compose.yml down    # stop
```

---

## Swagger workflow

```bash
# 1) install once
go install github.com/swaggo/swag/cmd/swag@latest

# 2) regenerate after editing handler comments
swag init -g ./src/cmd/main.go -o ./docs
```

* UI     → **`/swagger/index.html`**
* JSON   → **`/swagger/doc.json`**

CI snippet:

```bash
swag init -g ./src/cmd/main.go -o ./docs
git diff --exit-code ./docs
```

---

## 🗂 Project layout

```
base_structure/
├─ docker/                    # compose + redis.conf
│   └─ redis/
├─ src/
│   ├─ cmd/                   # main.go entry
│   ├─ api/
│   │   ├─ handlers/ routers/ middlewares/
│   │   └─ dto/ helper/
│   ├─ config/                # YAMLs + singleton loader
│   ├─ data/                  # db + cache
│   ├─ services/              # business logic
│   └─ pkg/                   # logging, utils
├─ docs/                      # swagger‑generated
└─ go.mod / go.sum
```

---

## 🧪 Testing

* Unit tests live beside the code (`*_test.go`).
* Custom validators registered in `TestMain`.
* Singleton config reload stub available via build‑tag `testtools` if needed.

```bash
go test ./...                         # run all
go test ./src/api/handlers            # single package
```

---

## Makefile helpers (optional)

```bash
make run     # go run ./src/cmd --config ...
make swag    # swag init
make test    # go test ./...
```

---

## FAQ

| Problem | Remedy |
|---------|--------|
| `no config specified` | Pass `--config` or set `CONFIG_FILE` |
| `swag: command not found` | `$GOPATH/bin` missing from `$PATH`; reinstall CLI |
| Redis `WRONGPASS` | Ensure `.env REDIS_PASSWORD` matches `docker/redis/redis.conf` |
| Postgres connection refused | Wait until container is healthy; verify creds |
| Swagger UI shows old routes | Run `swag init` & hard‑refresh browser |

Happy building 🚀
