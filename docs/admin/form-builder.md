# Form Builder & Custom Order Fields

The Form Builder (`/form-builder`) allows administrators to design dynamic custom input forms attached to specific products or checkout flows.

---

## 📝 Creating Custom Forms

1. Go to **Form Builder** in the Administrator Portal.
2. Click **Create Form** (e.g., `Dedicated Server Custom Options` or `Domain Extra Info`).
3. Add Fields:
   - **Text / Textarea:** For server hostnames, root passwords, or custom instructions.
   - **Select / Dropdown:** For Operating System choices (Ubuntu 24.04, Debian 12, AlmaLinux 9).
   - **Checkbox / Radio:** For add-on backups or control panel licenses.
4. Set validation rules (Required, regex patterns, minimum/maximum lengths).
5. Attach the form to a product in **Products & Catalog**. During checkout, the client portal will dynamically render these fields.
