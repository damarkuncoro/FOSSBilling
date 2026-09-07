# Invoicing & Subscription Engine

FOSSBilling features an automated billing engine that handles one-time charges, recurring subscriptions, tax calculations, dynamic promos, and overdue suspensions.

---

## 📅 Subscription Lifecycle & Automated Worker

1. **Advance Invoicing (14 Days Prior):**
   - The background scheduler daemon (`cmd/worker`) scans active subscriptions and generates upcoming renewal invoices 14 days before the service due date.
2. **Invoice Reminder Notifications:**
   - Automated email templates notify clients of pending due dates and provide direct 1-click payment links.
3. **Automatic Settlement on Balance:**
   - If a client has sufficient credit in their balance, the invoice is settled automatically upon generation.
4. **Grace Period & Suspension (7 Days Post Due Date):**
   - If an invoice remains unpaid 7 days past the due date, the worker marks the service as `suspended` and executes remote suspension via the assigned server provisioner (e.g., WHM/Plesk API).
5. **Instant Reactivation upon Payment:**
   - As soon as the overdue invoice is paid, the event listener triggers automatic unsuspension.

---

## 🧾 PDF Invoices & Receipts

- Invoices feature downloadable PDF generation powered by `pkg/pdf`.
- Customizable company headers, logos, tax identification numbers, and terms of service.
