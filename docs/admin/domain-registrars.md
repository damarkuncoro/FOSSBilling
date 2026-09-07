# Domain Registrars & TLD Management

FOSSBilling allows automated domain registration, transfer, renewal, WHOIS contact updates, and nameserver synchronization.

---

## 🌐 Supported Registrar Adapters

| Registrar | Adapter ID | Supported Operations |
| :--- | :--- | :--- |
| **Namecheap** | `namecheap` | Register, Renew, Transfer, Nameservers, Lock, EPP Code |
| **ResellerClub** | `resellerclub` | Register, Renew, Transfer, Contact Management, DNS |
| **Internet.bs** | `internetbs` | Real-time registration, renewals, private WHOIS |
| **Email Registrar** | `email` | Generates automated ticket/email to registrar desk |
| **Custom Webhook** | `custom` | Dispatches JSON webhook payload to custom API endpoint |

---

## 🏷️ Setting Up Top-Level Domains (TLDs)

1. Navigate to **Domains** (`/domains`) in the Administrator Portal.
2. Click **Add New TLD** (e.g. `.com`, `.id`, `.io`, `.net`).
3. Set Pricing:
   - Registration Price (1–10 years)
   - Renewal Price
   - Transfer Price
4. Assign the default **Registrar Driver** to automate registrations upon invoice settlement.
