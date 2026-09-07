# Developer Guide: Go Clean Architecture

FOSSBilling is structured following **Domain-Driven Clean Architecture** to maximize maintainability, testability, and decoupling.

---

## 🏛️ Directory Structure & Separation of Concerns

```text
backend-go/
├── core/
│   ├── domain/        # 1. Pure Go business entities, interfaces & value objects
│   ├── usecase/       # 2. Application business logic (no HTTP or SQL dependencies)
│   ├── repository/    # 3. Data access (PostgreSQL with SQLX & In-memory Mocks)
│   ├── service/       # 4. Multi-driver gateways, provisioners & mailer
│   ├── listener/      # 5. Event Bus subscribers for domain decoupling
│   └── handler/       # 6. HTTP REST handlers & JSON responses
├── pkg/               # Shared packages (auth, decimal, pdf, events, i18n, tools)
└── cmd/               # Binary entrypoints: api, worker, cli, demo
```

---

## 🧪 Dependency Injection & Testing

All use cases accept domain repository and service interfaces:

```go
type BillingService struct {
    invoiceRepo domain.InvoiceRepository
    clientRepo  domain.ClientRepository
    eventBus    *events.EventBus
}

func NewBillingService(
    invoiceRepo domain.InvoiceRepository,
    clientRepo domain.ClientRepository,
    eventBus *events.EventBus,
) *BillingService {
    return &BillingService{
        invoiceRepo: invoiceRepo,
        clientRepo:  clientRepo,
        eventBus:    eventBus,
    }
}
```

This design allows 100% unit testing using in-memory mock repositories without requiring a live PostgreSQL instance.
