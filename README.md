# Base‑Structure — Gin • Gorm • Redis • PostgreSQL

**Production‑ready Golang starter** with JWT auth, Docker‑Compose services, Zap logging and live Swagger UI.

---

## ✨ Highlights

| Layer | Library / Tool | Purpose |
|-------|----------------|---------|
| HTTP server | **Gin‑Gonic** | Fast router, middleware ecosystem |
| ORM | **Gorm** (`gorm.io/gorm`) | PostgreSQL driver, auto‑migrations |
| Caching / blacklist | **Redis 7** (`go-redis/redis/v7`) | OTP rate‑limit & token revocation |
| Auth | **golang‑jwt/jwt** | Access & Refresh tokens (HMAC) |
| Validation | **validator/v10** | Custom tags: `ir_mobile`, `password` |
| Config | **Viper** + dotenv | YAML per environment, env overrides |
| Logging | **Zap** + Lumberjack | JSON logs + file rotation |
| API docs | **Swaggo** + Swagger‑UI | OpenAPI 3 at `/swagger/` |

---

## 🚀 Quick start (development)

```bash
# 1. clone & enter
$ git clone <repo_url> my-api && cd my-api

# 2. create local configs (never commit!)
$ cp .env.example .env                                  # dotenv secrets
$ cp src/config/config-development-example.yml src/config/config-development.yml
$ cp docker/redis/redis_example.conf docker/redis/redis.conf

# 3. spin up Postgres, Redis, pgAdmin
$ docker compose -f docker/docker-compose.yml up -d

# 4. grab Go modules & Swagger CLI
$ go mod download
$ go install github.com/swaggo/swag/cmd/swag@latest

# 5. generate docs
$ swag init -g ./src/cmd/main.go -o ./docs

# 6. run the API
$ go run ./src/cmd                   # default :5005
# → http://localhost:5005/swagger/
```

> Change the port in `src/config/config-development.yml → server.port`.

---

## 🔧 Configuration

### YAML files (<code>src/config/*.yml</code>)

| APP_ENV | File loaded | Tracked in git |
|---------|-------------|----------------|
| `development` | `src/config/config-development.yml` | ❌ (copy from *_example*) |
| `docker` | `/app/config/config-docker.yml` | ✅ |
| `production` | `/config/config-production.yml` | ✅ / secret store |

Each file mirrors the `Config` struct in `src/config/config.go`.

### `.env`

Secrets for Docker services (Postgres user/pwd, pgAdmin login, etc.).  
Example keys: `POSTGRES_PASSWORD`, `PGADMIN_DEFAULT_PASSWORD`, `REDIS_PASSWORD`, `APP_ENV`.

---

## 🐳 Docker‑Compose stack (`docker/docker-compose.yml`)

| Service  | Host → Container | Notes |
|----------|------------------|-------|
| **postgres** | 5432 → 5432 | volume `postgres` |
| **pgadmin4** | 8090 → 80  | volume `pgadmin` |
| **redis** | 6379 → 6379 | uses `docker/redis/redis.conf` |

Start / stop:

```bash
docker compose -f docker/docker-compose.yml up -d   # start background
docker compose -f docker/docker-compose.yml down    # stop & remove
```

---

## Swagger workflow

```bash
# install CLI (once)
go install github.com/swaggo/swag/cmd/swag@latest

# regenerate after editing handler comments
swag init -g ./src/cmd/main.go -o ./docs
```

* UI: **`/swagger/index.html`**
* Raw spec: **`/swagger/doc.json`**

Add to CI:
```bash
swag init -g ./src/cmd/main.go -o ./docs
git diff --exit-code ./docs        # fail PR if docs stale
```

---

## 🗂️ Project layout (top‑level)

```
base_structure/
│
├─ docker/
│   ├─ docker-compose.yml
│   └─ redis/
│       ├─ redis_example.conf
│       └─ redis.conf   # copied, git‑ignored
│
├─ src/
│   ├─ cmd/                # main.go entry
│   ├─ api/
│   │   ├─ handlers/       # Gin handlers + Swagger comments
│   │   ├─ routers/
│   │   ├─ middlewares/
│   │   ├─ dto/
│   │   └─ helper/
│   ├─ config/             # YAMLs + loader code
│   ├─ constants/
│   ├─ data/
│   │   ├─ db/             # Gorm init + migrations
│   │   └─ cache/          # Redis singleton helpers
│   ├─ services/           # business logic (user, token, otp…)
│   └─ pkg/                # logging, util packages
│
├─ docs/                   # swagger auto‑generated
└─ go.mod / go.sum
```

---

## 📚 Key Go dependencies

```text
github.com/gin-gonic/gin           # HTTP router
github.com/swaggo/gin-swagger      # Swagger UI middleware
github.com/swaggo/swag             # OpenAPI generator
gorm.io/gorm & gorm.io/driver/postgres
github.com/golang-jwt/jwt          # JWT auth
github.com/go-redis/redis/v7       # Redis client
go.uber.org/zap                    # logging
github.com/spf13/viper             # config
```

Indirect packages are pulled automatically (`go mod tidy`).

---

## Makefile helpers (optional)

```bash
make swag   # swag init
make run    # go run ./src/cmd
make test   # go test ./...
```

---

## 🆘 FAQ

| Issue | Fix |
|-------|-----|
| `swag: command not found` | `$GOPATH/bin` not on `$PATH`; reinstall CLI |
| UI shows outdated routes | `swag init` then hard‑refresh browser |
| Redis `WRONGPASS` | Ensure `.env REDIS_PASSWORD` == `redis.conf requirepass` |
| Postgres connection refused | Wait 2‑3 s after compose‑up; check creds in YAML & .env |

Happy coding 🚀

