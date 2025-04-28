# Base Structure – Go (Gin + Gorm) Web API

> Production‑ready boilerplate with PostgreSQL, Redis, JWT auth, Zap logging and Swagger‑UI.

---

## ✨ Features

| Layer | Tech | Notes |
|-------|------|-------|
| **HTTP server** | [Gin‑Gonic](https://github.com/gin-gonic/gin) | ultra‑fast router & middleware support |
| **ORM** | [Gorm](https://gorm.io) + Postgres driver | auto‑migrations in `/src/data/db/migrations` |
| **Caching / Rate‑limit / Black‑list** | Redis v7 | singleton initialized via `src/data/cache` |
| **Validation** | go‑playground/validator v10 | custom tags: `ir_mobile`, `password` |
| **Config** | Viper + env files | copy `.env.example` → `.env` |
| **Logs** | Uber Zap + Lumberjack | JSON logs + daily rotation |
| **Auth** | HMAC‑JWT (golang‑jwt) | access & refresh tokens, Redis blacklist, middleware in `src/api/middlewares` |
| **Docs** | Swaggo (OpenAPI 3) + Swagger‑UI | live at `/swagger/index.html` |

---

## 🖥️ Local development

```bash
# 1 Clone repo and enter it
$ git clone <your‑fork‑url> project && cd project

# 2 Copy samples ➜ edit as needed
$ cp docker/redis/redis_example.conf docker/redis/redis.conf
$ cp .env.example .env

# 3 Start backing services (Postgres + Redis + pgAdmin)
$ docker compose -f docker/docker-compose.yml up -d

# 4 Resolve Go modules & install Swag CLI (one‑time)
$ go mod download
$ go install github.com/swaggo/swag/cmd/swag@latest

# 5 Generate Swagger docs
$ swag init -g ./src/cmd/main.go -o ./docs

# 6 Run the API
$ go run ./src/cmd
# -> http://localhost:8080/swagger/index.html
```

### Make targets (optional)

If you use the provided `Makefile` skeleton:

```bash
make swag        # regen docs
make run         # go run ./src/cmd
make test        # run unit tests
```

---

## ⚙️ Dockerized deployment

```bash
$ docker compose -f docker/docker-compose.yml up -d --build
```

* API container definition is left to you (multi‑stage build in `Dockerfile`).
* `postgres`, `redis`, and `pgadmin` services are defined in *docker/docker-compose.yml*.
  * Exposed ports:
    * **Postgres** → `5432`
    * **Redis**    → `6379`
    * **pgAdmin**  → `8090`

### Environment variables

The stack reads sensitive values from **`.env`** (git‑ignored). Populate at least:

```env
APP_ENV=development
SERVER_PORT=8080

# Postgres
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=webapi
DATABASE_URL=postgres://app:secret@postgres:5432/webapi?sslmode=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=password          # remember to mirror in redis.conf

# JWT
JWT_SECRET=supersecret
JWT_REFRESH_SECRET=anothersecret
ACCESS_TOKEN_EXPIRE_MIN=15
REFRESH_TOKEN_EXPIRE_MIN=10080   # =7 days
```

> **Tip:** pgAdmin credentials are `PGADMIN_DEFAULT_EMAIL` and `PGADMIN_DEFAULT_PASSWORD` in `.env`.

---

## 🗂️ Project layout (high‑level)

```
├── src
│   ├── cmd/            # main.go (entry‑point)
│   ├── api/            # presentation layer
│   │   ├── routers/
│   │   ├── handlers/
│   │   ├── dto/
│   │   └── middlewares/
│   ├── services/       # business logic
│   ├── data
│   │   ├── db/         # gorm models + migrations
│   │   └── cache/      # redis singleton helpers
│   └── pkg/            # reusable helpers (logging, service_errors, …)
└── docker/
    ├── docker-compose.yml
    └── redis/redis.conf
```

---

## 📚 Key Go dependencies

```text
github.com/gin-gonic/gin           # HTTP server
github.com/swaggo/swag             # Swagger generator
github.com/swaggo/gin-swagger      # Swagger‑UI middleware
gorm.io/gorm                       # ORM core
gorm.io/driver/postgres            # Postgres driver
github.com/golang-jwt/jwt          # JWT auth
github.com/go-redis/redis/v7       # Redis client
github.com/go-playground/validator # struct validation
go.uber.org/zap                    # structured logs
github.com/spf13/viper             # config loader
```

The full list is in `go.mod`; indirect dependencies are pulled automatically via `go mod tidy`.

---

## 🚀 Swagger / OpenAPI workflow

| Step | Command |
|------|---------|
| **Install CLI** | `go install github.com/swaggo/swag/cmd/swag@latest` |
| **Generate / update docs** | `swag init -g ./src/cmd/main.go -o ./docs` |
| **Serve UI** | visit `http://localhost:<PORT>/swagger/index.html` |

> ⚠️  Re‑run `swag init` whenever you change handler annotations so the JSON/YAML stays in sync.

---

## 🧪 Tests & linting (optional)

```bash
go test ./...                 # unit tests
# go vet ./...                # static analysis (built‑in)
# golangci-lint run            # if you use golangci‑lint
```

---

## 🆘 Troubleshooting

| Symptom | Fix |
|---------|-----|
| `swag: command not found` | `$GOPATH/bin` not in `$PATH`; run install step again |
| Swagger UI shows old routes | `swag init` then browser hard refresh |
| Redis `WRONGPASS` error | Make sure `docker/redis/redis.conf` and `.env` have matching `REDIS_PASSWORD` |
| Postgres refuses connection | Wait 2‑3 s after `docker compose up` or add `depends_on` in your API service |

Enjoy building on the **Base Structure**! ✌🏻

