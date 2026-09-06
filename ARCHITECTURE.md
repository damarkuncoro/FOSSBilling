# FOSSBilling Next-Gen Architecture & Design Patterns

Dokumen ini menjelaskan arsitektur teknis modern, pemisahan dependensi (*Separation of Concerns*), serta pola desain (*Design Patterns*) yang diimplementasikan pada FOSSBilling Next-Gen (Golang Backend + React TypeScript Frontends).

---

## 🏛️ 1. Prinsip Arsitektur Utama (Clean Architecture & DDD)

Sistem dibangun mengikuti prinsip **Clean Architecture**, **Domain-Driven Design (DDD)**, dan **SOLID Principles** dengan pemisahan tanggung jawab yang ketat antar layer:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           1. Presentation Layer                         │
│   - HTTP Handlers: Guest, Client, Admin (100% Bebas Akses Repo Langsung)│
│   - Background Task Workers / Cron Jobs                                 │
└────────────────────────────────────┬────────────────────────────────────┘
                                     │ (Hanya memanggil Services / Use Cases)
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    2. Service / Use Case Domain Layer                   │
│   - Auth, Billing, Cart, Order, Domain, License, Support, Stats, etc.  │
│   - Mengenkapsulasi seluruh Business Logic, Validasi, dan Orchestration │
└────────────────────────────────────┬────────────────────────────────────┘
                                     │ (Injeksi Repository & Driver Interfaces)
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                            3. Core Domain Layer                         │
│   - Core Entities (Client, Order, Invoice, Staff, Promo, Currency)      │
│   - Repository Interfaces & Value Objects (decimal.Money)               │
└────────────────────────────────────┬────────────────────────────────────┘
                                     │ (Dikonkretkan oleh Infrastructure)
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                       4. Infrastructure & Driver Layer                  │
│   - PostgreSQL Data Repositories & In-Memory Test Repositories          │
│   - Payment Gateways, RDAP/WHOIS Drivers, Provisioners, Mail Drivers   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 🧩 2. Pemisahan Tegas Repository dan Services

Setiap layer memiliki batas yang jelas untuk mencegah kebocoran abstraksi (*leaky abstraction*):

### A. Service & Use Case Layer (`core/usecase/`)
* **`DomainService`** (`core/usecase/domain/domain_service.go`): Mengatur pengecekan ketersediaan domain secara real-time via RDAP/DNS, manajemen nameserver, toggle perpanjangan otomatis (*auto-renew*), dan penerbitan kode transfer EPP.
* **`LicenseService`** (`core/usecase/license/license_service.go`): Mengatur siklus lisensi enterprise, reset IP/Domain Lock, dan regenerasi kunci kriptografis lisensi.
* **`OrderService`** (`core/usecase/order/order_service.go`): Mengatur siklus hidup pesanan (`Activate`, `Suspend`, `Unsuspend`, `Renew`, `Cancel`, `ListByClientID`).
* **`InvoiceService`** (`core/usecase/billing/invoice_service.go`): Mengatur kalkulasi pajak berjenjang, pembuatan invoice, dan pembayaran dengan saldo akun.
* **`CartService`** (`core/usecase/cart/cart_service.go`): Mengatur validasi keranjang belanja, kalkulasi voucher promo, dan alur checkout.

### B. HTTP Presentation Layer (`core/handler/http/`)
* **Nol Akses Repository:** Semua handler HTTP (`guest`, `client`, `admin`) **hanya menerima Service/UseCase** melalui *Constructor Injection*. Tidak ada handler yang mengakses database repository secara langsung.

### C. Frontend Client Architecture (`frontend-client/src/`)
Frontend klien mengadopsi pola Clean Architecture yang selaras:
* **Repository Layer (`src/repositories/`):** Mengabstraksi komunikasi HTTP/REST, serialisasi data, dan penanganan status error HTTP (`domain.repository.ts`, `auth.repository.ts`, `invoice.repository.ts`, `order.repository.ts`, `license.repository.ts`, `support.repository.ts`, `cart.repository.ts`, `download.repository.ts`, `news.repository.ts`, `company.repository.ts`).
* **Service Layer (`src/services/`):** Mengenkapsulasi logika bisnis frontend, sanitasi input, normalisasi domain, validasi password/form, format mata uang, serta orkestrasi multi-repository (`domain.service.ts`, `auth.service.ts`, `invoice.service.ts`, `order.service.ts`, `license.service.ts`, `support.service.ts`, `cart.service.ts`, `download.service.ts`, `news.service.ts`, `company.service.ts`).
* **Presentation Hooks Layer (`src/hooks/`):** React Custom Hooks (`useClientDomains`, `useStorefront`, `useClientInvoices`, dll) hanya mengonsumsi **Services**, bukan raw endpoint API.

### D. Frontend Administrator Architecture (`frontend-administrator/src/`)
Frontend portal admin mengadopsi struktur arsitektural yang konsisten:
* **Repository Layer (`src/repositories/`):** Mengisolasi endpoint administrasi (`admin_auth.repository.ts`, `admin_client.repository.ts`, `admin_order.repository.ts`, `admin_invoice.repository.ts`, `admin_support.repository.ts`, `admin_stats.repository.ts`, `admin_system.repository.ts`, `admin_massmail.repository.ts`).
* **Service Layer (`src/services/`):** Mengenkapsulasi kontrol siklus pesanan admin (aktivasi, penangguhan, pembatalan), pembuatan invoice manual, balasan tiket support, penyesuaian kurs mata uang, dan broadcast email massal (`admin_auth.service.ts`, `admin_client.service.ts`, `admin_order.service.ts`, `admin_invoice.service.ts`, `admin_support.service.ts`, `admin_stats.service.ts`, `admin_system.service.ts`, `admin_massmail.service.ts`).
* **Presentation Hooks Layer (`src/hooks/`):** Custom Hooks admin (`useOrders`, `useClients`, `useInvoices`, `useSupport`, `useCurrencies`, dll) mengonsumsi Services.

---

## 🌐 3. Domain Registrar & RDAP Lookup Driver

Sistem provisioning domain modern mengimplementasikan antarmuka `RegistrarDriver`:

* **`RDAPRegistrarDriver`** (`core/service/provisioning/rdap_registrar.go`):
  * Melakukan lookup ketersediaan domain real-time mengikuti spesifikasi **RFC 7480/7484 (RDAP - Registration Data Access Protocol)**.
  * Mendukung endpoint resmi ICANN, Verisign (`.com`, `.net`), PIR (`.org`), serta PANDI (`.id`, `.co.id`).
  * **Smart DNS Fallback:** Jika registry RDAP tidak merespons, sistem secara otomatis melakukan resolusi DNS (`net.LookupNS` dan `net.LookupIP`) untuk menentukan kepemilikan domain secara akurat.
  * Terintegrasi dengan **`OrderListener`** untuk memicu registrasi domain dan penerbitan kode otorisasi EPP secara instan setelah pembayaran invoice diselesaikan.

---

## 🏗️ 4. Builder Patterns

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

## 🏭 5. Factory Patterns

Pola **Factory** digunakan untuk menginstansiasi adapter driver pihak ketiga secara dinamis berdasarkan konfigurasi yang tersimpan di basis data tanpa *hardcoded branching*:

### A. `ProvisionerFactory` (`core/service/provisioning/factory.go`)
Menginstansiasi driver server provisioning:
* `cPanel`
* `DirectAdmin`
* `Plesk`
* `License Engine`
* `RDAP Registrar`

### B. `PaymentGatewayFactory` (`core/service/payment/gateways/factory.go`)
Menginstansiasi driver gateway pembayaran:
* `Midtrans` (QRIS, GoPay, Virtual Account)
* `Stripe` (Credit / Debit Card, Apple Pay)
* `PayPal` (Express Checkout)
* `BankTransfer` (Manual Wire Transfer)

---

## 🔌 6. Registry & Driver Adapters

Menggunakan thread-safe Registry untuk registrasi runtime dan lookup driver:

1. **`GatewayRegistry`** (`core/service/payment/registry.go`): Agregasi driver pembayaran.
2. **`ProvisionerRegistry`** (`core/service/provisioning/registry.go`): Agregasi driver server hosting dan registrar.
3. **`DriverRegistry`** (`core/service/notification/driver.go`): Agregasi driver email (`SMTP`, `Resend`, `SendGrid`, `Mock`).

---

## 📡 7. Event-Driven Architecture (Listeners)

Menggunakan publish-subscribe event bus untuk mengisolasi efek samping (*side-effects*) dari alur usecase utama:

* **`ClientListener`** (`core/listener/client_listener.go`): Mendengarkan event `client.registered` $\to$ Mengirim email sambutan secara otomatis.
* **`InvoiceListener`** (`core/listener/invoice_listener.go`): Mendengarkan event `invoice.paid` $\to$ Mengirimkan kuitansi PDF lunas.
* **`OrderListener`** (`core/listener/order_listener.go`): Mendengarkan event `order.activated` $\to$ Memproses otomatisasi provisioning server / registrasi domain / penerbitan lisensi dan mengirim notifikasi akses ke pengguna.

---

## 🐳 8. Multi-Environment Deployments & Reverse Proxy

* **Nginx Reverse Proxy:** Dilengkapi dengan `resolver 127.0.0.11 valid=5s` untuk resolusi dynamic upstream IP di jaringan Docker internal agar tidak mengalami *stale IP cache* (502 Bad Gateway) saat API container di-restart.
* **`deploy/docker-compose.dev.yml`**: Volume live reload untuk development cepat.
* **`deploy/docker-compose.prod.yml`**: Production stack terisolasi (PostgreSQL 16, Redis, Minimal Go Binaries, Nginx SPA).
* **`deploy/docker-compose.test.yml`**: Ephemeral database untuk test pipeline CI/CD.

---

## 🧪 9. Perintah Pengujian & Build

```bash
# Jalankan seluruh unit & integration test suites
make test

# Kompilasi seluruh binary Go dan frontend React bundles
make build

# Eksekusi simulasi bisnis E2E secara langsung
make demo
```
