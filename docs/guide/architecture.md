# Architecture & Design Overview

FOSSBilling is built upon the **Domain-Driven Clean Architecture** and **SOLID** principles, ensuring that business rules remain completely decoupled from transport layers, databases, and third-party APIs.

---

## 🏗️ Layered Clean Architecture

```text
├── core/
│   ├── domain/        # Core business entities & repository interfaces (Pure Go)
│   ├── usecase/       # Application business use cases (Invoice, Cart, Auth, Staff...)
│   ├── repository/    # Persistence implementation (PostgreSQL & Memory Mock)
│   ├── service/       # Provisioners (cPanel, Plesk...), Payment Gateways, Notifications
│   ├── listener/      # Event-driven subscribers (Event Bus)
│   └── handler/       # Presentation layer (HTTP REST handlers & JSON responses)
├── pkg/               # Reusable utility libraries (Auth, Decimal Money, PDF, Mailer, i18n)
└── cmd/               # Executable entrypoints (api, worker, cli, demo)
```

---

## 🎯 Design Patterns Used

### 1. Fluent Entity Builders
Complex business entities use type-safe builders:
- **`InvoiceBuilder`**: Builds immutable invoices with auto-calculated subtotals, compound/flat tax rates, line items, and dynamic currency conversions.
- **`ClientBuilder`**: Validates credentials, sets default currency, creates default client balance, and hashes passwords securely.
- **`OrderBuilder`**: Manages billing cycle periods, renewals, and custom form configuration values.

### 2. Dynamic Gateway & Provisioner Factories
- **`PaymentGatewayFactory`**: Instantiates and executes gateways dynamically based on payment method (`midtrans`, `stripe`, `paypal`, `bank_transfer`, `custom`).
- **`ProvisionerFactory`**: Directs service creation, suspension, unsuspension, and password updates across hosting control panels (`cpanel`, `plesk`, `directadmin`, `hestia`, `cwp`, `custom`).

### 3. Decoupled Event Bus
The asynchronous `EventBus` publishes domain events without coupling producers to consumers:
- `client.created` ➔ Dispatches welcome email template & initializes client balance.
- `invoice.paid` ➔ Triggers automatic order provisioning & receipt generation.
- `order.activated` ➔ Dispatches login details to the customer.

### 4. Zero-Loss Decimal Money Arithmetic
All monetary calculations throughout the system use the `decimal.Money` fixed-point arithmetic package (`pkg/decimal`), preventing floating-point precision rounding errors common in billing platforms.
