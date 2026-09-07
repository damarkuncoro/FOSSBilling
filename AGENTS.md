# FOSSBilling Project (Go & React Ecosystem)

## Project Overview

FOSSBilling is a modern, high-performance, open-source billing, subscription, and client management platform. This version of the project is built using a **Cloud-Native Clean Architecture** with a Go backend and React frontends.

### Key Technologies

* **Backend (Golang 1.22+):**
  * **Architecture:** Clean Architecture (Domain, Usecase, Repository, Handler).
  * **Frameworks:** Gin/Echo (check `backend-go/cmd/api`), GORM or SQLX for persistence.
  * **Database:** PostgreSQL 16 (Primary), Redis 7 (Caching/Queue).
  * **Testing:** Go `testing` package with table-driven tests.
  * **CLI:** Cobra-based management utility in `backend-go/cmd/cli`.
  * **API Documentation:** OpenAPI 3.0 (Scalar) served at `/docs`.

* **Frontend (React 18/19):**
  * **Framework:** Vite + TypeScript.
  * **Styling:** Tailwind CSS + shadcn/ui (Radix UI primitives).
  * **State Management:** TanStack Query (React Query) for API interactions.
  * **Icons:** Lucide React.
  * **Components:** Modular component structure in `src/components`.

### Repository Structure

```text
.
├── backend-go/                # Golang Backend Engine
│   ├── cmd/                   # Entry points (api, worker, cli, demo)
│   ├── core/                  # Business Logic (Clean Architecture)
│   │   ├── domain/            # Entities & Interfaces
│   │   ├── usecase/           # Domain business logic
│   │   ├── repository/        # Data access layer
│   │   └── handler/           # Transport layer (HTTP)
│   └── pkg/                   # Shared utilities (auth, pdf, mailer)
├── frontend-administrator/    # Admin Dashboard (Vite/React/shadcn)
├── frontend-client/           # Customer Portal (Vite/React/shadcn)
├── tests-backend-go/          # Comprehensive Go test suites
└── deploy/                    # Docker Compose & K8s configurations
```

## Development Conventions

### Backend (Go)
* Follow **Clean Architecture** principles. Use interfaces for dependency injection.
* Entities should live in `core/domain`.
* Use the **Builder Pattern** for complex entities (e.g., `InvoiceBuilder`).
* Use **Dynamic Factories** for extensible components like `PaymentGatewayFactory`.

### Frontend (React)
* Use functional components and hooks.
* Prefer **shadcn/ui** components for consistent design.
* API calls should be centralized in `lib/api`.
* Use TypeScript strictly for type safety.

### Legacy Support
* `backend-php/` and `tests-backend-php/` contain the legacy PHP implementation.
* New features should be implemented in the Go/React stack unless explicitly targeting the legacy system.

## Building and Running

Refer to the root `README.md` for Docker and local execution instructions.
Use `make build` and `make test` for common tasks.
