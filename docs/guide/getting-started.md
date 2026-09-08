# Getting Started with FOSSBilling

Welcome to the **FOSSBilling Next-Gen Ecosystem**!

FOSSBilling is a modern, high-performance, open-source billing, subscription, and client management platform designed for web hosts, domain registrars, SaaS products, and digital agencies.

---

## 🌟 What is FOSSBilling?

FOSSBilling streamlines your entire recurring billing and cloud service provisioning pipeline:
- **Order Management & Cart:** Automated client checkout, promo codes, custom order forms, and fraud scoring.
- **Automated Invoicing & Subscriptions:** Recurring invoice generation (14-day advance), tax rules, multiple currencies, and automatic overdue suspensions.
- **Server & Hosting Automation:** Automatic provisioning and lifecycle management for cPanel, Plesk, DirectAdmin, HestiaCP, CentOS Web Panel (CWP), and custom webhooks.
- **Domain Registration:** Instant domain availability lookups, transfers, and DNS management via Namecheap, ResellerClub, Internet.bs, and custom providers.
- **Customer Support Desk:** Multi-department ticket management, threaded messages, file attachments, and automated ticket auto-close.

---

## 🏛️ Ecosystem Architecture

The platform is split into three decoupled, cloud-native services:

```text
┌────────────────────────────────────────────────────────┐
│                   React Frontends                      │
│   • Administrator Portal (:3000) [Vite + Tailwind]     │
│   • Customer Client Portal (:3001) [Vite + Tailwind]   │
└──────────────────────────┬─────────────────────────────┘
                           │ HTTP REST JSON / OpenAPI 3.0
┌──────────────────────────▼─────────────────────────────┐
│                 Go Backend Engine (:8080)              │
│   • REST API Server (Clean Architecture)               │
│   • Background Worker Scheduler (Daemon)               │
│   • Terminal Administration Management CLI             │
└──────────────────────────┬─────────────────────────────┘
                           │
             ┌─────────────┴─────────────┐
             ▼                           ▼
    PostgreSQL 16 (Primary)       Redis 7 (Cache/Queue)
```

---

## 🚀 Quick Launch

You can spin up the complete ecosystem in under 2 minutes using Docker:

```bash
# 1. Clone repository
git clone https://github.com/damarkuncoro/FOSSBilling.git
cd FOSSBilling

# 2. Launch production stack
docker compose -f deploy/docker-compose.prod.yml up -d --build
```

Access your portals:
- **Customer Portal:** [`http://localhost:3001`](http://localhost:3001) (`client@fossbilling.org` / `Password123!`)
- **Administrator Portal:** [`http://localhost:3000`](http://localhost:3000) (`admin@fossbilling.org` / `admin123`)
- **Go REST API & Interactive Docs:** [`http://localhost:8080/docs`](http://localhost:8080/docs)

---

## 📖 Operational Guides & Next Steps

Ready to set up your store and start serving clients? Follow our dedicated guides:

- 📘 **[Complete Step-by-Step Usage Guide](/guide/step-by-step-guide):** Comprehensive, phase-by-phase walkthrough covering Admin setup, payment gateways, server provisioning, product creation, customer ordering, automated provisioning, and ticket helpdesk.
- ⚙️ **[Administrator Guide](/admin/overview):** Deep dive into administrative capabilities, currencies, products, and fraud prevention.
- 🛍️ **[Customer Portal Guide](/client/overview):** Understand the client storefront, domain search, cart checkout, invoice management, and service control.
- 💻 **[CLI Management](/guide/cli):** Learn how to run migrations, seeds, and background cron jobs from your terminal.

