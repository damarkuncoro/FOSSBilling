# Frontend Portals Architecture (React + Vite)

The frontend applications (`frontend-administrator` and `frontend-client`) are modern single-page applications (SPAs) built with React, TypeScript, Tailwind CSS, and shadcn/ui.

---

## 🛠️ Technology Stack

- **Framework:** Vite + React 18/19 + TypeScript
- **Styling:** Tailwind CSS + Lucide Icons + shadcn/ui (Radix UI primitives)
- **Data Fetching:** TanStack React Query + Axios / Fetch Client
- **Routing:** React Router v6
- **Testing:** Vitest + React Testing Library

---

## 📂 Modular Structure

```text
frontend-administrator/src/
├── components/
│   ├── ui/             # Reusable shadcn/ui primitives (Button, Card, Dialog, Table...)
│   ├── layout/         # Layout shells, SidebarNav, ProtectedRoute
│   └── common/         # LanguageSwitcher, CookieConsentBanner, ThemeToggle
├── lib/
│   ├── api/            # Typed API client facade
│   ├── auth.tsx        # AuthContext provider & token storage
│   └── i18n.tsx        # Multi-locale provider & translation hook
├── pages/              # Route pages (Dashboard, Clients, Orders, Invoices...)
└── services/           # Data services with React Query integration
```
