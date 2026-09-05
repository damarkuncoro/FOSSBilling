import { useState } from 'react';
import type { KbArticle } from '../types/clientModules';

const initialArticles: KbArticle[] = [
  {
    id: 1,
    category: 'Getting Started',
    title: 'How to connect your custom domain to cPanel hosting',
    slug: 'connect-domain-cpanel',
    summary: 'Step-by-step guide to updating DNS A records and nameservers for your domain.',
    content: 'To connect your custom domain, point your nameservers to `ns1.fossbilling.org` and `ns2.fossbilling.org`. DNS propagation typically takes between 1 to 24 hours.',
    views: 1240,
    helpful_count: 88,
    updated_at: '2026-08-10T10:00:00Z',
  },
  {
    id: 2,
    category: 'Security & SSL',
    title: 'Installing Free Let’s Encrypt SSL Certificate',
    slug: 'install-ssl-certificate',
    summary: 'Automate free SSL installation for HTTPS security on all subdomains.',
    content: 'All our hosting plans include automated Let’s Encrypt AutoSSL. It verifies domain DNS and automatically installs within 30 minutes of activation.',
    views: 890,
    helpful_count: 65,
    updated_at: '2026-08-20T14:00:00Z',
  },
  {
    id: 3,
    category: 'Billing & Invoicing',
    title: 'How to make payments using Credit Card, PayPal, or Bank Transfer',
    slug: 'payment-methods-guide',
    summary: 'Understanding invoice payment gateways and automatic recurring subscriptions.',
    content: 'Invoices can be paid instantly via Stripe, PayPal, or Midtrans. Once the transaction completes, your services are provisioned immediately.',
    views: 540,
    helpful_count: 42,
    updated_at: '2026-09-01T08:00:00Z',
  },
];

export function useKnowledgebase() {
  const [articles] = useState<KbArticle[]>(initialArticles);
  const [search, setSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [activeArticle, setActiveArticle] = useState<KbArticle | null>(null);

  const categories = ['all', 'Getting Started', 'Security & SSL', 'Billing & Invoicing', 'E-Mail Setup'];

  const filteredArticles = articles.filter((art) => {
    const matchesCategory = selectedCategory === 'all' || art.category === selectedCategory;
    const matchesSearch =
      art.title.toLowerCase().includes(search.toLowerCase()) ||
      art.summary.toLowerCase().includes(search.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  return {
    articles: filteredArticles,
    categories,
    search,
    setSearch,
    selectedCategory,
    setSelectedCategory,
    activeArticle,
    setActiveArticle,
  };
}
