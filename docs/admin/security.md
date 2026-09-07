# Security, Anti-Spam & Fraud Prevention

FOSSBilling is hardened against abuse, fraudulent orders, brute-force attacks, and cross-site scripting (XSS).

---

## 🛡️ Active Security Modules

### 1. Disposable Email Address Detection
- Validates client registration against a continuously updated blacklist of 1,000+ temporary/throwaway email domains (`pkg/security/disposable_email.go`).
- Blocks throwaway accounts while allowing legitimate personal and corporate emails.

### 2. IP Blocklist & Honeypot Protection
- **Honeypot Form Fields:** Invisible hidden fields catch automated spam bots attempting registration or contact submissions.
- **StopForumSpam Integration:** Cross-references client IP and email against global abuse databases.
- **CIDR Blocklist:** Manual and automatic IP subnet blocking.

### 3. SHA-256 Client Fingerprinting
- Generates a deterministic device/session fingerprint based on User-Agent, Accept headers, encoding preferences, and IP address (`pkg/security/fingerprint.go`).
- Detects credential stuffing and session hijacking attempts.

### 4. Input Sanitization & Anti-XSS
- All incoming HTML and text inputs are sanitized with anti-XSS stripping rules (`pkg/security/sanitizer.go`).

### 5. Audit Logging & Role-Based Access Control (RBAC)
- All staff administrative actions (client updates, server reboots, invoice edits, currency changes) are logged immutably in the `audit_logs` table.
