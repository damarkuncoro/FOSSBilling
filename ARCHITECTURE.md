# FOSSBilling Next-Gen Architecture & Design Patterns

Dokumen ini menjelaskan arsitektur teknis modern, pemisahan dependensi (*Separation of Concerns*), serta pola desain (*Design Patterns*) yang diimplementasikan pada FOSSBilling Next-Gen (Golang Backend + React TypeScript Frontends).

---

## 🏛️ 1. Prinsip Arsitektur Utama (Clean Architecture)

Sistem dibangun mengikuti prinsip **Clean Architecture** dan **SOLID Principles**:

```
+-------------------------------------------------------------------------+
|                              Transport Layer                            |
|    - HTTP Handlers (Guest, Client, Admin)                               |
|    - Background Task Workers / Cron Jobs                                |
+-------------------------------------------------------------------------+
                                    │
                                    ▼
+-------------------------------------------------------------------------+
|                              Usecase Layer                              |
|    - Auth, Billing, Cart, Order, Support, Stats, MassMail, APIKey, etc. |
+-------------------------------------------------------------------------+
                                    │
                                    ▼
+-------------------------------------------------------------------------+
|                               Domain Layer                              |
|    - Core Entities (Client, Order, Invoice, Staff, Promo, Currency)     |
|    - Domain Repository Interfaces & Value Objects (decimal.Money)        |
+-------------------------------------------------------------------------+
                                    │
                                    ▼
+-------------------------------------------------------------------------+
|                            Infrastructure Layer                         |
|    - PostgreSQL Repositories / In-Memory Mock Repositories              |
|    - Payment Gateways, Provisioning Drivers, Email Notification Drivers |
+-------------------------------------------------------------------------+
```

---

## 🏗️ 2. Builder Patterns

Pola **Builder** digunakan untuk merakit entitas domain yang memiliki banyak parameter opsional, kalkulasi bertingkat, atau lampiran file:

### A. `InvoiceBuilder` (`core/builder/invoice_builder.go`)
Memfasilitasi pembuatan invoice yang kompleks dengan kalkulasi pajak multi-tier, diskon voucher promo, konversi kurs mata uang, due date prorata, dan line items.
```go
invoice, items, err := builder.NewInvoiceBuilder().
    ForClient(clientID).
    WithSerieAndNr("INV", "2026-0001").
    WithCurrency("USD", 1.0).
    WithDueDays(14).
    WithTaxRate(11.0).
    ApplyPromo(promoVoucher).
    AddItem("Cloud VPS Pro (2 vCPU, 4GB RAM)", decimal.FromFloat(99.0), 1, true).
    Build()
```

### B. `OrderBuilder` (`core/builder/order_builder.go`)
Memfasilitasi pembuatan pesanan layanan hosting, VPS, domain, atau lisensi perangkat lunak dengan metadata konfigurasi dinamis.

### C. `ClientBuilder` (`core/builder/client_builder.go`)
Memfasilitasi registrasi klien lengkap dengan hashing password otomatis (`bcrypt`), penentuan status aktif, dan mata uang dasar.

### D. `MessageBuilder` (`pkg/mailer/builder.go`)
Fluent builder untuk merakit pesan email transaksional lengkap dengan header MIME, multi-recipients, dan lampiran file PDF (*invoice attachment*).

---

## 🏭 3. Factory Patterns

Pola **Factory** digunakan untuk menginstansiasi adapter driver pihak ketiga secara dinamis berdasarkan konfigurasi yang tersimpan di basis data tanpa *hardcoded branching*:

### A. `ProvisionerFactory` (`core/service/provisioning/factory.go`)
Menginstansiasi driver server provisioning:
* `cPanel`
* `DirectAdmin`
* `Plesk`
* `License Engine`

### B. `PaymentGatewayFactory` (`core/service/payment/gateways/factory.go`)
Menginstansiasi driver gateway pembayaran:
* `Midtrans` (QRIS, GoPay, Virtual Account)
* `Stripe` (Credit / Debit Card, Apple Pay)
* `PayPal` (Express Checkout)
* `BankTransfer` (Manual Wire Transfer)

---

## 🔌 4. Registry & Driver Adapters

Menggunakan thread-safe Registry untuk registrasi runtime dan lookup driver:

1. **`GatewayRegistry`** (`core/service/payment/registry.go`): Agregasi driver pembayaran.
2. **`ProvisionerRegistry`** (`core/service/provisioning/registry.go`): Agregasi driver server hosting.
3. **`DriverRegistry`** (`core/service/notification/driver.go`): Agregasi driver email (`SMTP`, `Resend`, `SendGrid`, `Mock`).

---

## 📡 5. Event-Driven Architecture (Listeners)

Menggunakan publish-subscribe event bus untuk mengisolasi efek samping (*side-effects*) dari alur usecase utama:

* **`ClientListener`** (`core/listener/client_listener.go`): Mendengarkan event `client.registered` $\to$ Mengirim email sambutan secara otomatis.
* **`InvoiceListener`** (`core/listener/invoice_listener.go`): Mendengarkan event `invoice.paid` $\to$ Mengirimkan kuitansi PDF lunas.
* **`OrderListener`** (`core/listener/order_listener.go`): Mendengarkan event `order.activated` $\to$ Mengirimkan rincian akses layanan aktif kepada pengguna.

---

## 🐳 6. Multi-Environment Deployments

Konfigurasi Docker Compose telah dipisahkan sesuai kebutuhan lingkungan:
* **`deploy/docker-compose.dev.yml`**: Volume live reload untuk development cepat.
* **`deploy/docker-compose.prod.yml`**: Production stack terisolasi (PostgreSQL 16, Redis, Minimal Go Binaries, Nginx SPA).
* **`deploy/docker-compose.test.yml`**: Ephemeral database untuk test pipeline CI/CD.

---

## 🧪 7. Perintah Pengujian & Build

```bash
# Jalankan seluruh unit & integration test suites
make test

# Kompilasi seluruh binary Go dan frontend React bundles
make build

# Eksekusi simulasi bisnis E2E secara langsung
make demo
```
