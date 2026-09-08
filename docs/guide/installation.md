# Installation & Deployment Guide

This guide covers the system requirements and various ways to install and build FOSSBilling.

---

## 🖥️ System Requirements

Before installing FOSSBilling, ensure your environment meets the following minimum requirements:

### Server Hardware
- **CPU:** 1 Core (2+ Cores recommended for high traffic)
- **RAM:** 1 GB (2 GB+ recommended for Redis caching and concurrent worker tasks)
- **Disk:** 500 MB for binaries + database storage space

### Software Dependencies
- **Operating System:** Linux (Ubuntu 22.04+, Debian 11+, CentOS 9), macOS, or Windows (WSL2)
- **Containerization:** Docker 24.0+ and Docker Compose v2.20+ (Recommended)
- **Database:** PostgreSQL 16
- **Caching:** Redis 7

---

## 🐳 Docker Installation (Recommended)

Docker is the fastest way to get FOSSBilling running with all its dependencies correctly configured.

### Quick Start
```bash
# 1. Clone the repository
git clone https://github.com/damarkuncoro/FOSSBilling.git
cd FOSSBilling

# 2. Start the production stack
docker compose -f deploy/docker-compose.yml up -d --build
```

### Accessing the Portals
Once the containers are healthy:
- **Admin Portal:** [http://localhost:3000](http://localhost:3000)
- **Client Portal:** [http://localhost:3001](http://localhost:3001)
- **API Health:** [http://localhost:8080/health](http://localhost:8080/health)

---

## 🔨 Building FOSSBilling

If you prefer to run FOSSBilling on bare metal or custom environments, you can build the binaries and frontend bundles manually using the provided `Makefile`.

### 1. Build Backend Binaries
Requires **Go 1.22+**.
```bash
cd backend-go
# Build API, Worker, and CLI binaries into bin/
make build
```

### 2. Build Frontend Portals
Requires **Node.js 20+** and **npm**.
```bash
# Build Administrator Portal
cd frontend-administrator
npm install
npm run build

# Build Client Portal
cd frontend-client
npm install
npm run build
```

The resulting `dist/` folders can be served via Nginx or any static web server.

---

## 📦 Installing FOSSBilling Standalone

To install FOSSBilling without Docker:

1. **Setup Database:** Create a PostgreSQL 16 database and user.
2. **Setup Redis:** Ensure a Redis 7 instance is reachable.
3. **Configure Environment:** Create a `.env` file in the root directory based on `.env.example`.
4. **Initialize Schema:** Run the SQL migrations found in `backend-go/migrations/` against your Postgres database.
5. **Run Services:**
   - Start the API: `./backend-go/bin/api`
   - Start the Worker: `./backend-go/bin/worker`
   - Serve the `dist/` folders from the frontend builds using a web server.
