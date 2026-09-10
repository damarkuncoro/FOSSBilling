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
| `COMPANY_NAME` | `FOSSBilling` | Name of your hosting company for branding |
| `ADMIN_PORTAL_URL` | `http://localhost:3000` | URL to access the staff dashboard |
| `CLIENT_PORTAL_URL` | `http://localhost:3001` | URL for the customer portal |

---

## 🔐 Security & 2FA

FOSSBilling Next-Gen supports TOTP-based Two-Factor Authentication.

- **Issuer Name:** The name that appears in authenticator apps (e.g. Google Authenticator) is controlled by `COMPANY_NAME`.
- **JWT Lifespan:** Configured in `AuthUsecase` (default 24 hours).

---

## 📧 Notification Channels

### Email (SMTP)
| Variable | Description |
| :--- | :--- |
| `MAIL_DRIVER` | `smtp` or `mock` |
| `MAIL_HOST` | SMTP server address |
| `MAIL_PORT` | SMTP server port |
| `MAIL_USER` | SMTP username |
| `MAIL_PASS` | SMTP password |
| `MAIL_FROM_ADDRESS` | Sender email (e.g. billing@yourcompany.com) |
| `MAIL_FROM_NAME` | Sender name shown in inbox |

### Telegram Alerts
| Variable | Description |
| :--- | :--- |
| `TELEGRAM_BOT_TOKEN` | API token from @BotFather |
| `TELEGRAM_CHAT_ID` | Your numeric Chat ID (use @userinfobot to find it) |

---

## 🌐 External API Integrations

| Variable | Provider | Purpose |
| :--- | :--- | :--- |
| `CLOUDFLARE_TOKEN` | Cloudflare | DNS Zone & Record management |
| `NAMECHEAP_API_USER` | Namecheap | Automated domain registration |
| `NAMECHEAP_API_KEY` | Namecheap | Domain API authentication |
| `STRIPE_SECRET_KEY` | Stripe | Global credit card processing |
| `MIDTRANS_SERVER_KEY`| Midtrans | Local Indonesian payments (QRIS/VA) |

---

## 🛠️ Infrastructure Tuning

For high-traffic deployments, you can tune the following:

- **DB Max Connections:** Managed via `pgxpool` configuration in code.
- **Redis Eviction Policy:** Recommended `allkeys-lru` for caching efficiency.
- **Worker Concurrency:** Controls how many background tasks the worker processes simultaneously.
