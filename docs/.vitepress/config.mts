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
      { text: 'Step-by-Step Guide', link: '/guide/step-by-step-guide' },
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
            { text: 'Step-by-Step Usage Guide', link: '/guide/step-by-step-guide' },
            { text: 'Architecture Overview', link: '/guide/architecture' },
            { text: 'Installation & Building', link: '/guide/installation' },
            { text: 'Configuration Guide', link: '/guide/configuration' },
            { text: 'Localization & i18n', link: '/guide/localization' },
            { text: 'CLI Management Utility', link: '/guide/cli' },
          ],
        },
        {
          text: 'Operations',
          items: [
            { text: 'Maintenance & Troubleshooting', link: '/guide/maintenance' },
            { text: 'Security Policy', link: '/guide/security-policy' },
            { text: 'Changelog', link: '/guide/changelog' },
          ],
        },
      ],
      '/admin/': [
        {
          text: 'Administrator Guide',
          items: [
            { text: 'Overview & Dashboard', link: '/admin/overview' },
            { text: 'Step-by-Step Usage Guide ➔', link: '/guide/step-by-step-guide' },
            { text: 'Company Information', link: '/admin/company' },
            { text: 'Product Types & Catalog', link: '/admin/products' },
            { text: 'Server Provisioning', link: '/admin/servers' },
            { text: 'Payment Gateways', link: '/admin/payment-gateways' },
            { text: 'Domain Registrars', link: '/admin/domain-registrars' },
            { text: 'Invoicing & Email', link: '/admin/invoicing' },
            { text: 'Email Templates', link: '/admin/email-templates' },
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
            { text: 'Step-by-Step Usage Guide ➔', link: '/guide/step-by-step-guide' },
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
            { text: 'Custom Provisioners', link: '/developer/provisioners' },
            { text: 'Custom Payment Gateways', link: '/developer/gateways' },
            { text: 'Customizing Invoice PDFs', link: '/developer/pdf' },
            { text: 'Localization & i18n', link: '/developer/i18n' },
            { text: 'Frontend Architecture', link: '/developer/frontend' },
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
