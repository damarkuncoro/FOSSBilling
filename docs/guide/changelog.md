# Changelog

All notable changes to the FOSSBilling project will be documented in this file.

---

## [2.0.0-golang] - 2026-09-08

This is a major milestone release transitioning FOSSBilling to a Cloud-Native Go backend and React frontend ecosystem.

### Added
- **Core Engine:** Brand new Golang 1.22+ backend with Clean Architecture.
- **Modern UI:** React-based portals for both Clients and Administrators with Tailwind CSS.
- **Multi-Driver Provisioning:** Native support for cPanel, DirectAdmin, Plesk, HestiaCP, and CWP.
- **Custom Formbuilder:** Ability to create custom technical fields for product checkouts.
- **Notification Center:** Real-time in-app alerts for payment and service status changes.
- **Security:** Integrated honeypot, anti-XSS, and disposable email blocking.

### Changed
- Migrated primary database from MySQL to **PostgreSQL 16**.
- Optimized money arithmetic using `int64` cent-precision to avoid floating point errors.
- Improved API performance with sub-millisecond response times.

---

## [1.x.x] - Legacy
Refer to the `backend-php/` folder for legacy PHP source code and history. New features are exclusively developed on the 2.x Go branch.
