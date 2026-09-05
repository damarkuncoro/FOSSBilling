import React from 'react';
import {
  LayoutDashboard,
  Users,
  Package,
  FileText,
  LifeBuoy,
  Coins,
  Newspaper,
  Building2,
  ShieldAlert,
  Boxes,
  Server,
  CreditCard,
  Tag,
  Send,
  Mail,
  Activity,
  FileCheck2,
  Globe,
  BookOpen,
  Layers,
  ShieldCheck,
  BarChart3,
  Key,
  SlidersHorizontal,
  Radio,
  Search,
  Code2,
  ArrowRightLeft,
} from 'lucide-react';

export interface NavGroup {
  groupName?: string;
  items: {
    name: string;
    path: string;
    icon: React.ComponentType<{ className?: string }>;
  }[];
}

export const navGroups: NavGroup[] = [
  {
    groupName: 'Core Operations',
    items: [
      { name: 'Dashboard', path: '/', icon: LayoutDashboard },
      { name: 'Clients', path: '/clients', icon: Users },
      { name: 'Orders', path: '/orders', icon: Package },
      { name: 'Invoices & Billing', path: '/invoices', icon: FileText },
      { name: 'Financial Reports', path: '/reports', icon: BarChart3 },
      { name: 'Support Tickets', path: '/support', icon: LifeBuoy },
    ],
  },
  {
    groupName: 'Store & Provisioning',
    items: [
      { name: 'Products & Services', path: '/products', icon: Boxes },
      { name: 'Software Licenses', path: '/licenses', icon: Key },
      { name: 'Form Builder', path: '/form-builder', icon: SlidersHorizontal },
      { name: 'Domain & TLDs', path: '/domains', icon: Globe },
      { name: 'Hosting Servers', path: '/servers', icon: Server },
      { name: 'Payment & Tax', path: '/gateways', icon: CreditCard },
      { name: 'Coupons & Promo', path: '/coupons', icon: Tag },
      { name: 'Currencies', path: '/currencies', icon: Coins },
    ],
  },
  {
    groupName: 'Content & Growth',
    items: [
      { name: 'Pages & KB (CMS)', path: '/pages', icon: BookOpen },
      { name: 'News & Articles', path: '/news', icon: Newspaper },
      { name: 'Mass Mailer', path: '/mass-mail', icon: Send },
      { name: 'Email & SMTP', path: '/email-templates', icon: Mail },
      { name: 'SEO & Webmaster', path: '/seo', icon: Search },
      { name: 'Embed Widgets', path: '/embed-widgets', icon: Code2 },
    ],
  },
  {
    groupName: 'System & Security',
    items: [
      { name: 'Company Settings', path: '/company', icon: Building2 },
      { name: 'Extensions Hub', path: '/extensions', icon: Layers },
      { name: 'Staff & Security', path: '/staff-security', icon: ShieldAlert },
      { name: 'Anti-Spam Shield', path: '/antispam', icon: ShieldCheck },
      { name: 'Event Webhooks', path: '/webhooks', icon: Radio },
      { name: 'Cookie Consent GDPR', path: '/cookie-consent', icon: ShieldCheck },
      { name: 'URL Redirects', path: '/redirects', icon: ArrowRightLeft },
      { name: 'System & Cron', path: '/system', icon: Activity },
      { name: 'Audit Logs', path: '/audit-logs', icon: FileCheck2 },
    ],
  },
];
