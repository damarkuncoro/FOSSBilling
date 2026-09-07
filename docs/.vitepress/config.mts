import { defineConfig } from 'vitepress';

export default defineConfig({
  title: 'FOSSBilling Docs',
  description: 'Documentation for FOSSBilling Cloud-Native Billing, Hosting, & Automation Platform',
  lang: 'en-US',
  cleanUrls: true,
  lastUpdated: false,
  ignoreDeadLinks: true,
  themeConfig: {
    siteTitle: 'FOSSBilling Docs',
    logo: {
      light: 'https://raw.githubusercontent.com/FOSSBilling/branding/refs/heads/main/logo-svg/fossb_logo-black_text.svg',
      dark: 'https://raw.githubusercontent.com/FOSSBilling/branding/refs/heads/main/logo-svg/fossb_logo-white_text.svg',
      alt: 'FOSSBilling Logo'
    },
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Administrator', link: '/admin/overview' },
      { text: 'Client Portal', link: '/client/overview' },
      { text: 'Developer & API', link: '/developer/architecture' },
      { text: 'API Specs (Scalar)', link: 'http://localhost:8080/docs' },
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/guide/getting-started' },
            { text: 'Architecture Overview', link: '/guide/architecture' },
            { text: 'Installation & Docker', link: '/guide/installation' },
            { text: 'CLI Management Utility', link: '/guide/cli' },
          ],
        },
      ],
      '/admin/': [
        {
          text: 'Administrator Guide',
          items: [
            { text: 'Overview & Dashboard', link: '/admin/overview' },
            { text: 'Server Provisioning', link: '/admin/servers' },
            { text: 'Payment Gateways', link: '/admin/payment-gateways' },
            { text: 'Domain Registrars', link: '/admin/domain-registrars' },
            { text: 'Invoicing & Subscriptions', link: '/admin/invoicing' },
            { text: 'Security & Anti-Spam', link: '/admin/security' },
            { text: 'Extensions & Marketplace', link: '/admin/extensions' },
            { text: 'Form Builder & Customization', link: '/admin/form-builder' },
            { text: 'Themes & Widgets', link: '/admin/themes-widgets' },
          ],
        },
      ],
      '/client/': [
        {
          text: 'Customer Portal Guide',
          items: [
            { text: 'Client Portal Overview', link: '/client/overview' },
            { text: 'Storefront & Ordering', link: '/client/storefront' },
            { text: 'Managing Invoices & Deposits', link: '/client/invoices' },
            { text: 'Licenses & Downloads', link: '/client/licenses-downloads' },
            { text: 'Helpdesk & Support Tickets', link: '/client/support' },
          ],
        },
      ],
      '/developer/': [
        {
          text: 'Developer & Engineering Guide',
          items: [
            { text: 'Clean Architecture (Go)', link: '/developer/architecture' },
            { text: 'REST API & OpenAPI', link: '/developer/api' },
            { text: 'Background Worker & Scheduler', link: '/developer/worker' },
            { text: 'Custom Provisioners & Drivers', link: '/developer/provisioners' },
            { text: 'Localization & i18n', link: '/developer/i18n' },
            { text: 'Frontend Architecture (React)', link: '/developer/frontend' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/damarkuncoro/FOSSBilling' },
      { icon: 'discord', link: 'https://fossbilling.org/discord' },
    ],
    search: {
      provider: 'local',
    },
    footer: {
      message: 'Released under the Apache 2.0 License.',
      copyright: 'Copyright © 2026 FOSSBilling Team & Contributors',
    },
  },
});
