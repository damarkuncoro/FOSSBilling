# Security & Anti-Spam

FOSSBilling includes enterprise-grade security features to protect your platform from automated abuse, fraudulent orders, and unauthorized access.

---

## 🛡️ The Anti-Spam Module

The Anti-Spam module is designed to block malicious bots while providing a smooth experience for real users.

### Core Protections
- **Disposable Email Blocker:** Automatically blocks registration from known temporary email providers (e.g., Mailinator, 10MinuteMail).
- **IP Reputation (StopForumSpam):** Cross-references new signups against the global StopForumSpam database.
- **Honeypot Fields:** Invisible form fields that, when filled out by bots, trigger an immediate rejection of the submission.
- **Rate Limiting:** Protects every API endpoint from brute-force attacks and DDoS by limiting requests per IP address.

### Captcha Integration
FOSSBilling supports modern, privacy-friendly captchas:
- **Cloudflare Turnstile:** Our recommended choice for invisible challenge-response.
- **Google reCAPTCHA v3:** Standard protection against automated interactions.

---

## 🔑 API Keys

For developers and external integrations, FOSSBilling provides a secure API Key system.

### Personal API Keys
Clients and Staff can generate API keys from their respective profile settings. These keys allow programmatic access to the platform without sharing passwords.
- **Auto-Rotation:** Keys can be set to expire after a certain period.
- **IP Restriction:** (Optional) Lock a key to a specific IP address or CIDR range.
- **Scoped Permissions:** Keys inherit the permissions of the user who created them.

### Managing Keys as an Admin
Administrators can view and revoke any active API keys from the **System > Security > API Keys** section if suspicious activity is detected.

---

## 📝 Audit Logs

The system maintains a comprehensive audit trail of all staff activities.
- **Who:** The staff member who performed the action.
- **What:** The specific change made (e.g., "Updated invoice #123").
- **When:** Accurate UTC timestamp.
- **Where:** The IP address and device fingerprint of the staff member.

Audit logs are immutable and provide critical data for compliance and security investigations.
