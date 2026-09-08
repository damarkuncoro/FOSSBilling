# Security Policy

We take the security of our billing and automation platform seriously. This document outlines how to report vulnerabilities and our commitment to resolving them.

---

## 🛡️ Reporting a Vulnerability

If you discover a security vulnerability within FOSSBilling, please **do not open a public issue**. Instead, follow these steps:

1. **Email:** Send a detailed report to `security@fossbilling.org`.
2. **Details:** Include a proof-of-concept (PoC), steps to reproduce, and the potential impact.
3. **Response:** Our core team will acknowledge your report within 48 hours.

---

## 🔒 Security Lifecycle

- **Coordinated Disclosure:** We ask that you give us reasonable time to investigate and patch the issue before making any public disclosure.
- **Patches:** We aim to release a security patch within 14 days for critical vulnerabilities.
- **Advisories:** Once a patch is released, we will publish a GitHub Security Advisory (GHSA) to notify the community.

---

## 🛑 Out of Scope

The following are generally considered out of scope for security reports:
- Brute-force attacks (ensure you have rate limiting enabled).
- Attacks requiring physical access to the server.
- Spam reports (unless it bypasses the Anti-Spam module filters).
- Issues in third-party extensions not maintained by the core team.

---

## 🛠️ Security Hardening

For production environments, always follow these best practices:
1. **Change the `JWT_SECRET`:** Never use the default key provided in `.env.example`.
2. **Use HTTPS:** Ensure all traffic between your portals and the API is encrypted.
3. **DB Access:** Restrict PostgreSQL access to the API container IP only.
4. **Regular Updates:** Keep your Docker images updated to the latest minor version.
