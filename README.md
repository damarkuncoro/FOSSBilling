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

<a href="https://fossbilling.org/downloads/"><img src="https://raw.githubusercontent.com/FOSSBilling/fossbilling.org/main/public/img/gh/download-button.png" alt="Download button" width="400"/></a>

[![Go CI & Tests](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/backend-go.yml/badge.svg)](https://github.com/damarkuncoro/FOSSBilling/actions/workflows/backend-go.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8.svg?logo=go)](https://golang.org)
[![PHP Version](https://img.shields.io/badge/PHP-8.3%2B-777BB4.svg?logo=php)](https://php.net)
[![Discord](https://img.shields.io/discord/747432407757488179?color=%237289FA&logo=discord&logoColor=%23FFF)](https://fossbilling.org/discord)

</div>

---

**FOSSBilling** is a free, open-source billing, subscription, and client management solution designed for hosting businesses, digital goods sellers, and online service providers.

This repository features the **Next-Generation Cloud-Native Golang Backend** alongside the legacy PHP engine, offering lightning-fast execution, Clean Architecture, complete REST API endpoints, and a dedicated testing suite.

---

## 🏗️ Repository Architecture

```text
.
├── backend-go/          # High-performance Next-Gen Golang REST API & Daemon
│   ├── cmd/
│   │   ├── api/         # HTTP REST API server
│   │   ├── demo/        # End-to-End live business simulation
│   │   └── worker/      # Scheduled background cron daemon
│   ├── core/            # Domain Entities, UseCases, Repositories & Handlers
│   ├── docs/            # OpenAPI 3.0 specification
│   ├── pkg/             # Shared packages (Auth, Decimal, PDF, Events, Mailer)
│   ├── Dockerfile       # Production multi-stage Docker build
│   └── docker-compose.yml
│
├── tests-backend-go/    # Centralized Golang Test Suites (100% Passing)
│   ├── Unit/            # Unit tests for Handlers, Services, UseCases, Pkg
│   └── E2E/             # Live REST API integration tests
│
├── backend-php/         # Modular PHP Backend
├── tests-backend-php/   # PHP Pest / PHPUnit test suites
└── frontend/            # Shared browser assets & themes
```

---

## ⚡ Quick Start: Golang Backend

### 1. Run the API Server
```bash
cd backend-go
go run ./cmd/api
```
Server runs at `http://localhost:8080`.
* 📖 **Interactive Scalar API Docs:** [http://localhost:8080/docs](http://localhost:8080/docs)
* ⚙️ **OpenAPI 3.0 Spec:** `http://localhost:8080/openapi.json`
* 🩺 **Health Check:** `http://localhost:8080/health`

### 2. Run the End-to-End Business Simulation
Simulates the entire customer journey (Registration, Multi-Currency, Invoicing, Tax, Payment Webhook, cPanel/DirectAdmin/Plesk Provisioning, Signed Digital Downloads, API Keys, Mass Mailer, and Analytics Dashboard):
```bash
cd backend-go
go run ./cmd/demo/main.go
```

### 3. Run Full-Stack with Docker Compose
Starts API, Background Worker, PostgreSQL 16 (with auto-migrations), and Redis 7:
```bash
cd backend-go
docker compose up -d
```

---

## 🧪 Testing

All Golang tests are organized inside `tests-backend-go/`:
```bash
# Run all Go unit & E2E tests
cd tests-backend-go
go test -v ./...

# Or using Makefile from backend-go
cd backend-go
make test
```

---

## 🚀 Key Modules & Endpoints

| Category | Endpoints | Description |
| :--- | :--- | :--- |
| **Auth & Profile** | `/api/v1/guest/auth/register`<br>`/api/v1/guest/auth/login`<br>`/api/v1/client/profile` | JWT-based client authentication and profile management |
| **Multi-Currency** | `/api/v1/guest/currencies`<br>`/api/v1/admin/currencies` | Dynamic exchange rates, custom price formatting, base currency rules |
| **News & Announcements**| `/api/v1/guest/news`<br>`/api/v1/guest/news/{slug}`<br>`/api/v1/admin/news` | SEO friendly slug publishing, drafts, and staff authoring |
| **Digital Downloads** | `/api/v1/client/downloads/{id}/link`<br>`/api/v1/client/downloads/{id}/file` | Secure HMAC-SHA256 signed digital delivery with expiration TTL |
| **API Keys** | `/api/v1/client/api-keys` | Token generator (`fb_*`) with secret hashing for third-party bots |
| **Billing & Cart** | `/api/v1/guest/cart/checkout`<br>`/api/v1/client/invoices` | Automated tax calculator, coupons, prorata, and PDF generation |
| **Provisioning** | cPanel, DirectAdmin, Plesk, License | Multi-driver server provisioning & hosting account automation |
| **Mass Mailer** | `/api/v1/admin/mass-mail` | Broadcast newsletter and batch notification system |

---

## 📜 License

FOSSBilling is released under the [Apache 2.0 License](LICENSE).
