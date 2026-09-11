# 📦 FOSSBilling Next-Gen Installation Guide

This guide provides instructions for deploying the modernized FOSSBilling stack (Go + React).

## 🏗️ Architecture Overview

The system consists of four main components:
1.  **Backend API (Go)**: Core business logic and RESTful services.
2.  **Worker (Go)**: Background tasks for billing, provisioning, and maintenance.
3.  **Administrator Portal (React)**: Management dashboard for staff.
4.  **Client Portal (React)**: Self-service area for customers.

---

## 🐳 Option 1: Fast Deployment (Docker Compose)

The easiest way to get started is using Docker Compose.

### 1. Prerequisites
- Docker & Docker Compose (v2.0+)
- Domain name (optional for local testing)

### 2. Configuration
Copy the example environment file and adjust the values:
```bash
cp backend-go/.env.example .env
```
Key variables to set:
- `DATABASE_URL`: Postgres connection string.
- `REDIS_URL`: Redis for caching and distributed locking.
- `JWT_SECRET`: A long, random string.
- `ALLOWED_ADMIN_COUNTRIES`: ISO codes (e.g., `ID,US`) for Geofencing.

### 3. Launch
```bash
docker compose -f deploy/docker-compose.yml up --build -d
```
The following services will be available:
- **Client Portal**: [http://localhost:3001](http://localhost:3001)
- **Admin Portal**: [http://localhost:3000](http://localhost:3000)
- **API Server**: [http://localhost:8080](http://localhost:8080)
- **API Docs**: [http://localhost:8080/docs](http://localhost:8080/docs)
- **Metrics**: [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## 🛠️ Option 2: Manual Installation (Development)

### 1. Backend (Go)
1.  Navigate to `backend-go/`.
2.  Run `go mod download`.
3.  Execute migrations: `go run ./cmd/cli db:migrate`.
4.  Start API: `go run ./cmd/api`.
5.  Start Worker: `go run ./cmd/worker`.

### 2. Frontend (React)
1.  Navigate to `frontend-administrator/` or `frontend-client/`.
2.  Run `npm install`.
3.  Copy `.env.example` to `.env` and set `VITE_API_URL`.
4.  Start development server: `npm run dev`.

---

## 🔄 Migrating from Legacy (PHP)

If you are using the old PHP version of FOSSBilling, you can migrate your data using the built-in importer:

```bash
# Inside backend-go folder
go run ./cmd/cli db:import:legacy --source "user:pass@tcp(host:3306)/legacy_db_name"
```
This will import Clients and Invoices into the new PostgreSQL database.

---

## 🛡️ Production Hardening

1.  **SSL**: Always use Nginx as a reverse proxy with SSL (Template provided in `deploy/nginx/ssl.conf.template`).
2.  **Redis**: Enable Redis for production to ensure **Distributed Locking** works for your workers.
3.  **Monitoring**: Integrate the `/metrics` endpoint with Prometheus/Grafana.
4.  **Backups**: Use the built-in automated backup service (runs daily at 02:00 UTC).

---

## 🆘 Support
Refer to the `docs/` folder for API documentation or join our community forums.
