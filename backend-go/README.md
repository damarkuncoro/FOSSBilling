# FOSSBilling Next-Gen Golang Backend

High-performance, Cloud-Native backend for FOSSBilling built with **Clean Architecture**, **Domain-Driven Design (DDD)**, and **100% test coverage**.

---

## 🏛️ Architecture Overview

```text
backend-go/
├── cmd/
│   ├── api/          # HTTP REST API Server with Go 1.22 enhanced ServeMux
│   ├── demo/         # E2E Business Lifecycle Simulation
│   └── worker/       # Background Cron & Batch Task Daemon
│
├── core/
│   ├── domain/       # Core Business Entities & Repository Interfaces
│   ├── repository/   # Data Access Layer (PostgreSQL Pool & In-Memory Fallback)
│   ├── usecase/      # Application Business Rules (Auth, Billing, Orders, etc.)
│   ├── service/      # Infrastructure Services (Mailer, Provisioners, Scheduler)
│   └── handler/      # HTTP Handlers (Guest, Client, Admin) & Middleware
│
├── docs/             # OpenAPI 3.0 Specification (openapi.json)
├── pkg/              # Shared Utilities (Auth JWT, Decimal Money, PDF, Events)
└── migrations/       # SQL Schema DDL (000001_init_schema.up.sql)
```

---

## ⚙️ Configuration & Environment Variables

Create a `.env` file or provide environment variables:

| Variable | Default | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | HTTP listening port |
| `APP_ENV` | `development` | Environment (`development` / `production`) |
| `DATABASE_URL` | `postgres://...` | PostgreSQL connection string (auto-fallback to In-Memory if offline) |
| `JWT_SECRET` | `secret-key-32-chars` | Secret key for signing JWT tokens and HMAC URLs |
| `TURNSTILE_SECRET`| `""` | Optional Cloudflare Turnstile captcha secret |

---

## 🚀 Running the Services

### 1. API Server
```bash
go run ./cmd/api
```
* **Interactive Docs (Scalar):** [http://localhost:8080/docs](http://localhost:8080/docs)
* **OpenAPI Spec:** `http://localhost:8080/openapi.json`
* **Health Check:** `http://localhost:8080/health`

### 2. Background Worker & Scheduler
```bash
go run ./cmd/worker
```
Handles automatic renewal invoice generation and overdue service auto-suspension.

### 3. Docker Compose Stack
```bash
docker compose up --build -d
```

---

## 🧪 Testing & Simulation

### 1. Unit & Integration Tests
All test suites are centralized in `tests-backend-go/`:
```bash
# Run all unit and integration tests
make test
```

### 2. Ultra-Deep Business Simulation
Perform a complete end-to-end business lifecycle simulation including Fraud Prevention, Domain Lifecycle, Advanced Billing, Provisioning, and Support:
```bash
go run ./cmd/demo/deep_simulation.go
```
**Tested Scenarios:**
* **Identity:** Registration with Anti-Spam (IP Blacklist) and TOTP 2FA.
* **Domains:** Availability checks and DNS record management (A record).
* **Billing:** Promo codes, PPN 11% Tax, Webhook payment settlement, and Client Balance credits.
* **Provisioning:** Automated cPanel account creation and License Key generation.
* **CRM:** News publishing, Mass Mail campaigns, and Support Ticket lifecycle.
* **Lifecycle:** Prorata upgrade calculations, HMAC-secured digital asset links, and Automated Termination.

---

## 📑 Default Demo Credentials (Standalone Mode)

When running without PostgreSQL, the server automatically starts in **In-Memory Fallback Mode** with seeded credentials:

* **Superadmin:** `admin@fossbilling.org` / `admin123`
* **Client:** `client@fossbilling.org` / `client123`
