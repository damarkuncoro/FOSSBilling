# REST API & OpenAPI 3.0 Documentation

FOSSBilling provides a comprehensive, type-safe REST API documenting all guest, customer, and administrative endpoints.

---

## 📖 Live OpenAPI Console (Scalar)

The Go backend engine directly serves the interactive **Scalar UI** at:
- **Interactive Console:** `http://localhost:8080/docs`
- **OpenAPI 3.0 Specification:** `http://localhost:8080/openapi.json`

---

## 🔐 Authentication & Roles

FOSSBilling uses standard Bearer Token authorization:

```http
Authorization: Bearer <jwt_or_session_token>
```

### Route Namespaces
- `/api/v1/guest/*` — Public endpoints (Storefront, Cart, Domain lookup, Locales, Webhooks, Login, Register).
- `/api/v1/client/*` — Authenticated client endpoints (Invoices, Services, Tickets, Deposit, Profile).
- `/api/v1/admin/*` — Authenticated staff/admin endpoints (Clients, Orders, Gateways, Servers, Extensions, System Health).

---

## ⚡ Live Notification Stream (WebSocket)

Admin portal uses a persistent WebSocket connection for real-time alerts.

- **Endpoint:** `GET /api/v1/admin/system/ws?token=<token>`
- **Events:**
  - `invoice_paid`: Triggered on new revenue.
  - `ticket_opened`: Alert for new support requests.

---

## 🩺 System Health & Monitoring

Standard health check endpoint for DevOps orchestrators (Kubernetes, Docker).

- **Endpoint:** `GET /api/v1/guest/system/health`
- **JSON Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-09-10T...",
  "services": {
    "api": "running",
    "cache": "active",
    "database": "healthy"
  }
}
```

---

## 📦 Standard API Response Format

```json
{
  "success": true,
  "data": {
    "id": 1,
    "status": "active"
  },
  "error": null
}
```
