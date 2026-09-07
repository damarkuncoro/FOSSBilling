# Migration Roadmap: PHP to Go

This document tracks the progress of migrating FOSSBilling from the legacy PHP backend to the modern Go Clean Architecture backend.

## Status Overview

- **Core Engine:** 🏗️ In Progress
- **Admin Portal:** ✅ Migrated to React
- **Client Portal:** ✅ Migrated to React
- **Database:** 🔄 Migrating from MySQL/MariaDB to PostgreSQL (Primary)

## Migrated Components (Go)

| Component | Status | Source (PHP) | Target (Go) |
| :--- | :--- | :--- | :--- |
| Auth System | ✅ Done | `src/library/Auth` | `backend-go/pkg/auth` |
| Client Management | ✅ Done | `src/modules/Client` | `backend-go/core/usecase/client` |
| Invoice Generation | ✅ Done | `src/modules/Invoice` | `backend-go/core/builder/invoice` |
| Order Processing | 🏗️ In Progress | `src/modules/Order` | `backend-go/core/usecase/order` |
| Support Tickets | ⏳ Pending | `src/modules/Support` | `backend-go/core/usecase/support` |

## How to Contribute to Migration

1. **Identify a Module:** Choose a module in `backend-php/src/modules` that hasn't been migrated.
2. **Define Domain Entities:** Create the corresponding Go structs in `backend-go/core/domain`.
3. **Implement Usecases:** Port the business logic from PHP services to Go usecases.
4. **Write Tests:** Ensure 100% test coverage in `tests-backend-go`.
5. **Update API:** Add the new endpoints to `backend-go/cmd/api` and update `openapi.json`.

## Legacy Compatibility

For now, the legacy system resides in `backend-php/`. We aim for eventual 100% parity before deprecating the PHP components.
