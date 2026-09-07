# Background Worker & Task Scheduler

The FOSSBilling Background Worker (`cmd/worker`) is a standalone background daemon responsible for recurring billing, service suspensions, ticket lifecycle, and system maintenance.

---

## ⏱️ Scheduled Pipeline Tasks

1. **`RunInvoiceRenewalsTask`**:
   - Runs periodically (every 6 hours).
   - Scans subscriptions with due dates within 14 days and generates renewal invoices.
2. **`RunOverdueSuspensionsTask`**:
   - Scans unpaid invoices older than 7 days grace period.
   - Suspends the remote hosting account via the server driver and notifies the client.
3. **`RunTicketAutoCloseTask`**:
   - Checks resolved support tickets inactive for 72 hours and marks them `closed`.
4. **`RunSystemMaintenanceTask`**:
   - Prunes expired session tokens, temporary cache entries, and stale IP blocklist records.

---

## 🚀 Running the Worker Standalone

```bash
cd backend-go
go run ./cmd/worker
```
