import {
  LayoutDashboard,
  Store,
  Server,
  Globe,
  Key,
  Download,
  Receipt,
  Headphones,
  BookOpen,
  Newspaper,
  Settings,
  LucideIcon,
} from 'lucide-react';

export interface ClientNavItem {
  name: string;
  path: string;
  icon: LucideIcon;
  requiresAuth?: boolean;
}

export interface ClientNavGroup {
  groupName?: string;
  items: ClientNavItem[];
}

export const clientNavGroups: ClientNavGroup[] = [
  {
    groupName: 'Overview',
    items: [
      { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard, requiresAuth: true },
      { name: 'Order & Store', path: '/', icon: Store },
    ],
  },
  {
    groupName: 'Services & Assets',
    items: [
      { name: 'My Services', path: '/services', icon: Server, requiresAuth: true },
      { name: 'Domains', path: '/domains', icon: Globe },
      { name: 'Licenses', path: '/licenses', icon: Key, requiresAuth: true },
      { name: 'Downloads', path: '/downloads', icon: Download, requiresAuth: true },
    ],
  },
  {
    groupName: 'Billing & Support',
    items: [
      { name: 'Invoices & Funds', path: '/invoices', icon: Receipt, requiresAuth: true },
      { name: 'Support Tickets', path: '/support', icon: Headphones, requiresAuth: true },
    ],
  },
  {
    groupName: 'Resources & Help',
    items: [
      { name: 'Knowledgebase', path: '/kb', icon: BookOpen },
      { name: 'Announcements', path: '/news', icon: Newspaper },
      { name: 'Account Settings', path: '/settings', icon: Settings, requiresAuth: true },
    ],
  },
];
