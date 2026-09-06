<h1 align="center">
  <br>
  <a href="https://fossbilling.org/">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/FOSSBilling/branding/refs/heads/main/logo-svg/fossb_logo-white_text.svg">
      <img alt="FOSSBilling logo" src="https://raw.githubusercontent.com/FOSSBilling/branding/refs/heads/main/logo-svg/fossb_logo-black_text.svg" height="100">
    </picture>
  </a>
  <br>
</h1>

<div align="center">

[![Go CI & Tests](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/backend-go.yml/badge.svg)](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/backend-go.yml)
[![Frontends CI](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/frontends-ci.yml/badge.svg)](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/frontends-ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8.svg?logo=go)](https://golang.org)
[![Vite](https://img.shields.io/badge/Vite-8.x-646CFF.svg?logo=vite)](https://vitejs.dev)
[![React](https://img.shields.io/badge/React-18%2F19-61DAFB.svg?logo=react)](https://react.dev)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.4%2Fshadcn-38B2AC.svg?logo=tailwind-css)](https://ui.shadcn.com)

</div>

---

**FOSSBilling** is a modern, high-performance, open-source billing, subscription, and client management platform designed for hosting providers, domain registrars, SaaS businesses, and digital goods merchants.

This repository features the **100% Modern Cloud-Native Clean Architecture Ecosystem**:
- ⚡ **Backend Engine:** Golang 1.22+ Clean Architecture REST API, Background Worker, & Management CLI (`backend-go`)
- 🖥️ **Administrator Portal:** Vite + React + TypeScript + shadcn/ui Dashboard (`frontend-administrator`)
- 🛍️ **Customer Client Portal:** Vite + React + TypeScript + shadcn/ui Storefront & Client Hub (`frontend-client`)
- 🧪 **Automated Testing:** 30/30 Comprehensive Unit & Integration Test Suites (`tests-backend-go`)
- 🏗️ **Design Patterns:** Fluent Builders (`InvoiceBuilder`, `OrderBuilder`, `ClientBuilder`, `MessageBuilder`) & Dynamic Factories (`ProvisionerFactory`, `PaymentGatewayFactory`)
- 🐳 **Multi-Environment Deployment:** Docker Compose for `dev` (live reload), `prod` (production), and `test` (`deploy/`)

---

## 🏗️ Repository Architecture

```text
.
├── backend-go/                # High-Performance Golang Backend (:8080)
│   ├── cmd/
│   │   ├── api/               # HTTP REST API server
│   │   ├── cli/               # Terminal administration management utility
│   │   ├── demo/              # End-to-End live business simulation
│   │   └── worker/            # Scheduled cron & background job worker
│   ├── core/
│   │   ├── builder/           # Domain entity builders (Invoice, Order, Client)
│   │   ├── domain/            # Core entities, interfaces & value objects
│   │   ├── handler/           # HTTP presentation handlers (guest, client, admin)
│   │   ├── listener/          # Event-driven subscribers (client, invoice, order)
│   │   ├── repository/        # PostgreSQL & memory mock repositories
│   │   ├── service/           # Multi-driver payment, provisioning & notification
│   │   └── usecase/           # 14 isolated domain business usecases
│   ├── pkg/                   # Shared packages (Auth, Decimal, PDF, Events, Mailer)
│   ├── migrations/            # PostgreSQL schema migrations
│   └── Dockerfile             # Multi-stage production build
│
├── frontend-administrator/    # Modern Staff/Admin Portal (:3000)
│   ├── src/                   # React + TypeScript + Tailwind + shadcn/ui
│   │   ├── components/ui/     # Radix UI primitives (Button, Card, Dialog, Table, etc.)
│   │   ├── pages/             # Dashboard, Clients, Orders, Invoices, Support, etc.
│   │   └── lib/api/           # Modular Admin API client
│   └── Dockerfile             # Alpine Nginx production image
│
├── frontend-client/           # Customer Storefront & Client Portal (:3001)
│   ├── src/                   # React + TypeScript + Tailwind + shadcn/ui
│   │   ├── pages/             # Storefront, Domain Checker, Cart, Services, Invoices
│   │   └── lib/api/           # Modular Customer Portal API client
│   └── Dockerfile             # Alpine Nginx production image
│
├── tests-backend-go/          # Centralized Go Test Suites (30/30 Passing)
├── deploy/                    # Docker Compose orchestration (dev, prod, test)
├── ARCHITECTURE.md            # Comprehensive technical design document
└── Makefile                   # Quick-start workflow automation commands
```

---

## ⚡ Quick Start: Running the Ecosystem

### 1. Run via Docker Compose (Recommended)
Launch all services with one command:
```bash
# Production stack with PostgreSQL 16 & Redis 7
docker compose -f deploy/docker-compose.prod.yml up -d --build

# Or Development stack with local volume mounts
docker compose -f deploy/docker-compose.dev.yml up -d --build
```

### 2. Run Locally in Development Mode
```bash
# 1. Run unit tests
make test

# 2. Compile all Go binaries (bin/api, bin/worker, bin/cli) & Frontend bundles
make build

# 3. Execute End-to-End Live Business Simulation
make demo
```

---

## 🌐 Live URLs & Access

| Portal | URL | Default Credentials | Description |
| :--- | :--- | :--- | :--- |
| 🛍️ **Customer Portal** | `http://localhost:3001` | `client@fossbilling.org` / `Password123!` | Storefront, Domain Lookup, Cart, Hosting Dashboard, Invoices |
| 🖥️ **Administrator Portal** | `http://localhost:3000` | `admin@fossbilling.org` / `SuperSecretAdmin123!` | Executive MRR/ARR Dashboard, Service Provisioning, Currencies |
| ⚡ **Golang REST API** | `http://localhost:8080` | - | High-performance JSON REST API |
| 📖 **Scalar API Documentation** | `http://localhost:8080/docs` | - | Interactive OpenAPI 3.0 Reference & Test Console |

---

## 💻 CLI Management Utility

Administrators can perform system operations directly from the terminal:
```bash
# Check runtime version and active system drivers
go run ./backend-go/cmd/cli status

# Create a new client account interactively via ClientBuilder
go run ./backend-go/cmd/cli client:create --email="user@cloud.id" --first-name="Budi" --company="PT Solusi Cloud"
```

---

## 🧪 Testing & Verification

```bash
# Run all 30 Go unit, builder, factory, and integration test suites:
make test

# Run frontend test suites (Vitest):
cd frontend-administrator && npm test -- --run
cd frontend-client && npm test -- --run
```

---

## 📜 Documentation & License

- **Detailed Technical Architecture:** See [`ARCHITECTURE.md`](ARCHITECTURE.md)
- **License:** Released under the [Apache 2.0 License](LICENSE).
