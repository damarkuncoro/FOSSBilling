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
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8.svg?logo=go)](https://golang.org)
[![Vite](https://img.shields.io/badge/Vite-8.x-646CFF.svg?logo=vite)](https://vitejs.dev)
[![React](https://img.shields.io/badge/React-18%2F19-61DAFB.svg?logo=react)](https://react.dev)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.4%2Fshadcn-38B2AC.svg?logo=tailwind-css)](https://ui.shadcn.com)

</div>

---

**FOSSBilling** is a modern, high-performance, open-source billing, subscription, and client management platform designed for hosting providers, domain registrars, SaaS businesses, and digital goods merchants.

This repository features the **100% PHP-Free Modern Cloud-Native Ecosystem**:
- ⚡ **Backend Engine:** Golang 1.22+ Clean Architecture REST API & Background Worker (`backend-go`)
- 🖥️ **Administrator Portal:** Vite + React + TypeScript + shadcn/ui Dashboard (`frontend-administrator`)
- 🛍️ **Customer Client Portal:** Vite + React + TypeScript + shadcn/ui Storefront & Client Hub (`frontend-client`)
- 🧪 **Automated Testing:** 25/25 Comprehensive Unit & E2E Test Suites (`tests-backend-go`)
- 🐳 **One-Click Deployment:** Multi-container Docker Compose with PostgreSQL 16 & Redis 7 (`deploy/`)

---

## 🏗️ Repository Architecture

```text
.
├── backend-go/                # High-Performance Golang Backend (:8080)
│   ├── cmd/
│   │   ├── api/               # HTTP REST API server
│   │   ├── demo/              # End-to-End live business simulation
│   │   └── worker/            # Scheduled cron & background job worker
│   ├── core/                  # Domain Entities, UseCases, Repositories & Handlers
│   ├── docs/                  # OpenAPI 3.0 specification & Scalar reference
│   ├── pkg/                   # Shared packages (Auth, Decimal, PDF, Events, Mailer)
│   ├── migrations/            # PostgreSQL schema migrations
│   └── Dockerfile             # Multi-stage production build
│
├── frontend-administrator/    # Modern Staff/Admin Portal (:3000)
│   ├── src/                   # React + TypeScript + Tailwind + shadcn/ui
│   │   ├── components/ui/     # Radix UI primitives (Button, Card, Dialog, Table, etc.)
│   │   ├── pages/             # Dashboard, Clients, Orders, Invoices, Support, etc.
│   │   └── lib/               # API client & Staff auth context
│   └── Dockerfile             # Alpine Nginx production image
│
├── frontend-client/           # Customer Storefront & Client Portal (:3001)
│   ├── src/                   # React + TypeScript + Tailwind + shadcn/ui
│   │   ├── pages/             # Storefront, Domain Checker, Cart, Services, Invoices
│   │   └── lib/               # API client, Cart state & Customer auth context
│   └── Dockerfile             # Alpine Nginx production image
│
├── tests-backend-go/          # Centralized Go Test Suites (25/25 Passing)
├── deploy/                    # Unified Docker Compose orchestration
└── Makefile                   # Quick-start workflow automation commands
```

---

## ⚡ Quick Start: Running the Ecosystem

### 1. Run via Docker Compose (Recommended)
Launch all 6 services with one command:
```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

### 2. Run Locally in Development Mode
Start each service in dedicated terminal tabs:

```bash
# 1. Start Golang REST API Server (Port :8080)
cd backend-go
go run ./cmd/api

# 2. Start Admin Portal (Port :3000)
cd frontend-administrator
npm run dev

# 3. Start Customer Portal (Port :3001)
cd frontend-client
npm run dev
```

---

## 🌐 Live URLs & Access

| Portal | URL | Default Credentials | Description |
| :--- | :--- | :--- | :--- |
| 🛍️ **Customer Portal** | `http://localhost:3001` | `client@fossbilling.org` / `client123` | Storefront, Domain Lookup, Cart, Hosting Dashboard, Invoices |
| 🖥️ **Administrator Portal** | `http://localhost:3000` | `admin@fossbilling.org` / `admin123` | Executive MRR/ARR Dashboard, Service Provisioning, Currencies |
| ⚡ **Golang REST API** | `http://localhost:8080` | - | High-performance JSON REST API |
| 📖 **Scalar API Documentation** | `http://localhost:8080/docs` | - | Interactive OpenAPI 3.0 Reference & Test Console |

---

## 🧪 Testing & Verification

Run the entire 25-suite unit and integration test suite:
```bash
cd tests-backend-go
go test -v ./...
```

Run frontend production builds:
```bash
# Build Admin Portal
cd frontend-administrator && npm run build

# Build Customer Portal
cd frontend-client && npm run build
```

---

## 🚀 Key Modules & Capabilities

- **Automated Multi-Driver Provisioning:** Instant creation and management for cPanel/WHM, DirectAdmin, Plesk, and Enterprise Serial Licenses.
- **Secure Digital Downloads:** HMAC-SHA256 encrypted temporary download links with expiration TTL.
- **Tax Engine & Multi-Currency:** Real-time PPN 11% tax calculation, multi-currency conversion (USD, IDR, EUR, SGD), and promo discount vouchers (`MERDEKA20`).
- **Support Helpdesk:** Threaded two-way conversation between customers and staff with priority escalation.
- **Mass Mailer & Campaigns:** Broadcast transactional and promotional campaigns to verified customer segments.
- **Audit Logging:** Immutable security audit trail recording staff actions, entity modifications, and IP addresses.

---

## 📜 License

FOSSBilling is released under the [Apache 2.0 License](LICENSE).
