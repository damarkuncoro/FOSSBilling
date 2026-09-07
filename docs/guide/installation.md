# Installation & Deployment Guide

FOSSBilling can be deployed via Docker Compose, Kubernetes, or standalone bare-metal environments.

---

## 🐳 Docker Compose Deployment (Recommended)

### Prerequisites
- Docker Engine 24.0+
- Docker Compose v2.20+

### Production Launch
```bash
# 1. Clone the project
git clone https://github.com/damarkuncoro/FOSSBilling.git
cd FOSSBilling

# 2. Build and start containers
docker compose -f deploy/docker-compose.prod.yml up -d --build
```

### Container Stack Topology

| Service | Port | Description |
| :--- | :--- | :--- |
| `fossbilling-api` | `:8080` | High-concurrency Go REST API server |
| `fossbilling-worker` | - | Background task scheduler daemon |
| `fossbilling-admin-portal` | `:3000` | Administrator React Portal (Nginx Alpine) |
| `fossbilling-client-portal` | `:3001` | Customer React Portal (Nginx Alpine) |
| `fossbilling-postgres` | `:5433` | PostgreSQL 16 primary database |
| `fossbilling-redis` | `:6379` | Redis 7 cache and job queue |

---

## 💻 Local Development Setup

### 1. Prerequisites
- **Go 1.22+**
- **Node.js 20+** & **npm**
- **PostgreSQL 16**

### 2. Run Locally
```bash
# Run unit & integration test suites
make test-all

# Start Go Backend API server
cd backend-go
go run ./cmd/api

# Start Admin Portal
cd frontend-administrator
npm install
npm run dev

# Start Client Portal
cd frontend-client
npm install
npm run dev
```
