# Administrator Guide

The FOSSBilling Administrator Portal is your central hub for managing your hosting business. It provides tools to automate service delivery, manage financial records, and provide customer support.

---

## 🚀 Getting Started as an Admin

Once you have installed FOSSBilling, log in to the Admin Portal at `http://localhost:3000`. The default credentials (if seeded) are `admin@fossbilling.org` / `admin123`.

::: tip Comprehensive End-to-End Walkthrough
For the complete step-by-step guide with diagrams and screenshots, see the **[Complete Step-by-Step Usage Guide](/guide/step-by-step-guide)**.
:::

### Essential Setup Checklist
1. **[Company Profile](/admin/company):** Set your company name, logo, support email, phone, physical address, and Tax ID.
2. **Currencies & Exchange Rates:** Define your base currency (e.g. USD, EUR, IDR) and automatic exchange rates.
3. **[Payment Gateways](/admin/payment-gateways):** Enable Stripe, PayPal Express, Midtrans (QRIS/VA), or Offline Bank Wire.
4. **[Server Nodes](/admin/servers):** Connect remote servers (cPanel, HestiaCP, Plesk, DirectAdmin, CWP) with API tokens and test connectivity.
5. **[Domain Registrars](/admin/domain-registrars):** Configure live registrar APIs (Namecheap, ResellerClub, Internet.bs).
6. **[Products & Packages](/admin/products):** Create shared hosting plans, domain pricing tables, license keys, and digital downloads.
7. **[Email Templates](/admin/email-templates):** Customize automated order confirmations, welcome emails, invoice notices, and dunning reminders.

---

## 📊 The Dashboard & Daily Workflow

Your command center for real-time monitoring and day-to-day operations:
- **Revenue Overview:** Track Monthly Recurring Revenue (MRR), total gross revenue, and pending invoice totals.
- **Client Activity:** Real-time stream of new client registrations, logins, and geographic distribution.
- **Operational Queue:** Quick badges for pending orders requiring manual verification, unpaid overdue invoices, and urgent open support tickets.
- **System Diagnostics:** Live health monitoring of the Go REST API engine, Redis caching latency, and background worker scheduler.

---

## 👥 Client Lifecycle Management

Manage customer accounts from initial registration through subscription renewals:
- **Client Profiles:** View and edit client contact info, company details, currency preference, and tax exemption status.
- **Prepay Balance & Credits:** Add manual account balance credits (`+$100.00`), adjust promotional funds, or issue store credit for returns.
- **Service Subscriptions:** View all active, suspended, or canceled hosting accounts, domains, and licenses tied to a client.
- **Security & Audit Trail:** Inspect client login IPs, user-agent fingerprints, failed login attempts, and password reset requests.
- **Internal Staff Notes:** Keep confidential staff notes and account history logs visible only to administrators.

---

## 🧾 Billing, Tax & Invoicing Operations

Automated and manual financial control:
- **Recurring Invoices:** The background daemon scans subscriptions and generates renewal invoices 14 days in advance.
- **Instant Prepay Settlement:** Invoices automatically settle if the customer holds sufficient prepaid credit balance.
- **Tax Rules & VAT Engine:** Configure multi-tier tax rates based on country, state/province, or European Union reverse-charge VAT rules.
- **Refunds & Dispute Handling:** Issue 1-click full or partial refunds directly back to Stripe/PayPal or credit the customer's on-platform balance.
- **Automated Dunning Reminders:** Configurable reminders sent 3 days before due date, on due date, and 3 days post due date.
- **Overdue Suspension:** Automated remote service suspension via server API at 7 days overdue, with instant auto-reactivation upon invoice payment.

---

## 🎧 Support Helpdesk Administration

Deliver enterprise support:
- **Department Routing:** Create specialized departments (e.g. Sales, L1 Tech Support, L2 Cloud Engineering, Billing).
- **Ticket Escalation & Assignment:** Assign tickets to specific operators or teams.
- **Canned Responses:** Accelerate resolution with pre-written templates and variable placeholders.
- **Staff-Only Internal Notes:** Collaborate with colleagues directly inside the ticket thread without exposing internal discussion to the client.

---

## 🛠️ System Tools & CLI Administration

- **Activity & Audit Logs:** Comprehensive logging of all administrator actions (deletions, modifications, logins, impersonations).
- **Cron Worker Management:** Monitor background worker executions or manually trigger billing cycles via `./backend-go/bin/cli cron:run`.
- **Cache Management:** Invalidate Redis cache keys across servers when modifying templates or global configurations.
- **Mass Mail Broadcasting:** Send system maintenance advisories or marketing updates with markdown formatting and variable substitution.

