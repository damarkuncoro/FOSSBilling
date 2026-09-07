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
