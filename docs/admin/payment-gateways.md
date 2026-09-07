# Payment Gateways & Processing

FOSSBilling includes an extensible payment processing subsystem supporting both synchronous card checkouts, digital wallets, QRIS, and offline bank transfers.

---

## 💳 Supported Payment Gateways

### 1. Stripe (`stripe`)
- **Payment Methods:** Credit/Debit Cards, Apple Pay, Google Pay, SEPA Direct Debit.
- **Configuration:** `Secret Key` (`sk_live_...`), `Publishable Key` (`pk_live_...`), `Webhook Secret` (`whsec_...`).
- **Webhook Endpoint:** `POST /api/v1/guest/payment/webhook/stripe`

### 2. PayPal Express Checkout (`paypal`)
- **Payment Methods:** PayPal balance, connected bank accounts, Pay in 4.
- **Configuration:** `Client ID`, `Client Secret`, `Sandbox Mode` (true/false).
- **Webhook Endpoint:** `POST /api/v1/guest/payment/webhook/paypal`

### 3. Midtrans (`midtrans`)
- **Payment Methods:** QRIS (GoPay, OVO, Dana, LinkAja), Virtual Accounts (BCA, Mandiri, BNI, BRI), Credit Cards.
- **Configuration:** `Server Key`, `Client Key`, `Production Mode` (true/false).
- **Notification Endpoint:** `POST /api/v1/guest/payment/webhook/midtrans`

### 4. Bank Transfer / Custom Manual Gateway (`bank_transfer`)
- **Payment Methods:** Manual bank transfer / wire transfer.
- **Configuration:** Bank Name, Account Number, Account Holder, Instructions.
- **Workflow:** Generates pending transaction; staff approve with 1-click upon receipt confirmation.

---

## 💱 Multi-Currency & Exchange Rates

1. Navigate to **Currencies** (`/currencies`) in Admin Portal.
2. Define a **Default Currency** (e.g. `USD` or `IDR`).
3. Add secondary currencies with automated or custom exchange rates.
4. Prices in invoices and storefronts are dynamically converted and formatted according to client currency.
