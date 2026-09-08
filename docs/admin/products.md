# Product Types & Catalog

FOSSBilling supports a wide range of product types, each with its own specialized automation and delivery logic.

---

## 📦 Core Product Types

When creating a product, you must select one of the following types to enable specific automation features.

### 1. 🖥️ Hosting
Designed for shared hosting or VPS instances.
- **Automation:** Connects to cPanel, DirectAdmin, or HestiaCP.
- **Workflow:** Generates a hosting account on the remote server once the invoice is paid.
- **Client View:** Displays server IP, login credentials, and control panel links.

### 2. 🌐 Domains
Handles domain registration and renewals.
- **Automation:** Connects to registrars like Namecheap or ResellerClub.
- **Workflow:** Performs real-time availability checks and registers the domain upon payment.
- **Client View:** Allows DNS management, WHOIS privacy toggles, and EPP code retrieval.

### 3. 💾 Downloadable Products
Digital assets like software, themes, or ebooks.
- **Delivery:** Generates a secure, expiring HMAC-signed download link.
- **Security:** Links are tied to the client's session and expire to prevent link sharing.
- **Client View:** A persistent "Downloads" tab where they can access their files.

### 4. 🔑 Licenses
Software license keys or serial numbers.
- **Delivery:** Generates a unique key based on a secure salt (e.g., `FOSS-XXXX-XXXX-XXXX`).
- **Validation:** Provides an API for your software to "ping" FOSSBilling for license status check.
- **Client View:** Displays the key and activation status.

---

## ⚙️ Custom Order Forms

You can attach a **Custom Form** to any product to gather technical details from the customer during the checkout process (e.g., Hostname, Root Password, or Server OS).

---

## 💰 Pricing Tiers

FOSSBilling supports flexible billing cycles:
- **Free:** No charge for the service.
- **One-time:** A single payment for lifetime access (ideal for downloads).
- **Recurring:** Monthly, Quarterly, or Annual billing with automated renewal invoices.

---

## 📈 Stock Management

Optionally set a stock limit for physical or limited-resource products. The system will automatically mark the product as "Sold Out" once the limit is reached.
