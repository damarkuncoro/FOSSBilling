# Configuration Guide

FOSSBilling uses environment variables for configuration. When running via Docker, these are typically defined in your `docker-compose.yml` or a `.env` file.

---

## ⚙️ Core Configuration Variables

These variables control the connection to essential services and core security settings.

| Variable | Default | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Environment mode (`production`, `development`, `test`) |
| `PORT` | `8080` | The port the Go API server listens on |
| `DATABASE_URL` | - | PostgreSQL DSN: `postgres://user:pass@host:5432/db?sslmode=disable` |
| `REDIS_URL` | - | Redis DSN: `redis://host:6379` |
| `JWT_SECRET` | - | A 32+ character string used to sign authentication tokens |
| `APP_URL` | `http://localhost:8080` | The base URL where the API is publicly accessible |
| `DEFAULT_CURRENCY` | `USD` | Default 3-letter ISO code for pricing |

---

## 🔒 Security Configuration

It is critical to change the `JWT_SECRET` in production. If this key is compromised, attackers can forge authentication tokens for any user or administrator.

```bash
# Generate a secure secret
openssl rand -base64 32
```

---

## 📧 Mailer Configuration

FOSSBilling supports multiple mail transports (SMTP, Sendmail). These are currently configured within the system settings database, but can be overridden via environment variables in future updates.

- **SMTP Host:** Server address for outgoing mail.
- **SMTP Port:** 587 (TLS) or 465 (SSL).
- **Authentication:** Username and password for your mail provider.

---

## 🛠️ Infrastructure Tuning

For high-traffic deployments, you can tune the following:

- **DB Max Connections:** Managed via `pgxpool` configuration in code.
- **Redis Eviction Policy:** Recommended `allkeys-lru` for caching efficiency.
- **Worker Concurrency:** Controls how many background tasks the worker processes simultaneously.
