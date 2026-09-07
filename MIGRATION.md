# Migration Roadmap: PHP to Go

This document tracks the progress of migrating FOSSBilling from the legacy PHP backend to the modern Go Clean Architecture backend.

## Status Overview

- **Core Engine:** ✅ Core Migrated (Stabilizing)
- **Admin Portal:** ✅ Migrated to React
- **Client Portal:** ✅ Migrated to React
- **Database:** 🔄 Migrating from MySQL/MariaDB to PostgreSQL (Primary)

## Migrated Components (Go)

| Component | Status | Source (PHP) | Target (Go) |
| :--- | :--- | :--- | :--- |
| Auth System | ✅ Done | `src/library/Auth` | `backend-go/pkg/auth` |
| Client Management | ✅ Done | `src/modules/Client` | `backend-go/core/usecase/client` |
| Invoice Generation | ✅ Done | `src/modules/Invoice` | `backend-go/core/builder/invoice` |
| Order Processing | ✅ Done | `src/modules/Order` | `backend-go/core/usecase/order` |
| Support Tickets | ✅ Done | `src/modules/Support` | `backend-go/core/usecase/support` |
| News & Knowledgebase| ✅ Done | `src/modules/News` | `backend-go/core/usecase/news` |
| Staff Management | ✅ Done | `src/modules/Staff` | `backend-go/core/usecase/staff` |
| Currency System | ✅ Done | `src/modules/Currency`| `backend-go/core/usecase/currency` |
| Promo & Coupons | ✅ Done | `src/modules/Promo` | `backend-go/core/usecase/cart` |
| API Key Management | ✅ Done | `src/modules/ApiKey`| `backend-go/core/usecase/apikey` |
| Mass Mailer | ✅ Done | `src/modules/Email` | `backend-go/core/usecase/massmail` |
| Product Catalog | ✅ Done | `src/modules/Product` | `backend-go/core/usecase/catalog` |
| System Settings | ✅ Done | `src/modules/System` | `backend-go/core/usecase/system` |
| Activity Logs | ✅ Done | `src/modules/Activity`| `backend-go/core/usecase/activity` |
| In-app Notif | ✅ Done | `src/modules/Notification`| `backend-go/core/usecase/notification` |
| Payment Gateways| ✅ Done | `library/Payment` | `backend-go/core/service/payment` |

## How to Contribute to Migration

1. **Identify a Module:** Choose a module in `backend-php/src/modules` that hasn't been migrated.
2. **Define Domain Entities:** Create the corresponding Go structs in `backend-go/core/domain`.
3. **Implement Usecases:** Port the business logic from PHP services to Go usecases.
4. **Write Tests:** Ensure 100% test coverage in `tests-backend-go`.
5. **Update API:** Add the new endpoints to `backend-go/cmd/api` and update `openapi.json`.

## Legacy Compatibility

For now, the legacy system resides in `backend-php/`. We aim for eventual 100% parity before deprecating the PHP components.
