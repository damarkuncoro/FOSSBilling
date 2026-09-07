# CLI Management Utility

FOSSBilling provides a built-in terminal CLI tool (`cmd/cli`) to manage system state, audit gateway/provisioner drivers, create administrator credentials, and generate cryptographic secrets.

---

## 🛠️ CLI Usage & Commands

Run commands from the `backend-go` directory:

```bash
cd backend-go
go run ./cmd/cli <command> [flags]
```

### 1. `status` — Subsystem Driver Audit
Verifies active payment gateways, server provisioners, and domain registrar drivers:
```bash
go run ./cmd/cli status
```

### 2. `admin:create` — Provision Administrator Account
Creates a new staff/admin account with automatic Bcrypt hashing:
```bash
go run ./cmd/cli admin:create \
  --email="sysadmin@fossbilling.org" \
  --name="System Operator" \
  --role="superadmin" \
  --password="SecurePassword123!"
```

### 3. `client:create` — Client Creation via Builder
Interactively constructs client entities using the `ClientBuilder`:
```bash
go run ./cmd/cli client:create \
  --email="client@company.com" \
  --first-name="Alex" \
  --last-name="Morgan" \
  --company="Global Hosting Corp" \
  --currency="USD"
```

### 4. `invoice:build` — Invoice Calculation Simulation
Simulates invoice generation and tax computation via `InvoiceBuilder`:
```bash
go run ./cmd/cli invoice:build \
  --client-id=1 \
  --item-title="Cloud VPS Enterprise - 8 Core" \
  --price=89.99 \
  --qty=1 \
  --tax-rate=10.0 \
  --currency=USD
```

### 5. `tools:password` — CSPRNG Password Generator
Generates high-entropy passwords with custom length:
```bash
go run ./cmd/cli tools:password --length=24 --special=true
```

### 6. `tools:geoip` — GeoIP Country & Flag Resolution
Looks up ISO country code, currency, and emoji flag for any IP address:
```bash
go run ./cmd/cli tools:geoip 1.1.1.1
```

### 7. `locale:list` — List Supported Locales & Direction
Displays all registered internationalization catalogs with RTL/LTR flags:
```bash
go run ./cmd/cli locale:list
```
