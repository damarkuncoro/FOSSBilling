import { useState } from 'react';
import type { SeoSettings } from '../types/seo';

const initialSettings: SeoSettings = {
  site_title: 'FOSSBilling - NextGen Cloud Hosting & Domains',
  meta_description: 'Fast, secure NVMe cloud hosting, domain registrations, and dedicated servers with 99.9% uptime SLA.',
  meta_keywords: 'cloud hosting, domains, vps, dedicated servers, billing management',
  og_image_url: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?q=80&w=1200&auto=format&fit=crop',
  twitter_handle: '@fossbilling',
  google_analytics_id: 'G-XYZ998877',
  robots_txt: `User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /api/\n\nSitemap: https://myhosting.com/sitemap.xml`,
  sitemap_auto_generate: true,
  canonical_url: 'https://myhosting.com',
};

export function useSeoSettings() {
  const [settings, setSettings] = useState<SeoSettings>(initialSettings);
  const [activeTab, setActiveTab] = useState<'meta' | 'robots' | 'sitemap'>('meta');
  const [isSaved, setIsSaved] = useState(false);

  const updateSettings = (key: keyof SeoSettings, value: any) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = () => {
    setIsSaved(true);
    setTimeout(() => setIsSaved(false), 2500);
  };

  return {
    settings,
    activeTab,
    setActiveTab,
    isSaved,
    updateSettings,
    handleSave,
  };
}
