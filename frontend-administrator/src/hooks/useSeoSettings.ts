import { useState, useEffect } from 'react';
import type { SeoSettings } from '../types/seo';
import { systemApi } from '../lib/api/system';

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
  const [seoInfo, setSeoInfo] = useState<any>(null);
  const [pingLoading, setPingLoading] = useState(false);

  useEffect(() => {
    fetchSeoInfo();
  }, []);

  const fetchSeoInfo = async () => {
    try {
      const info = await systemApi.getSeoInfo();
      setSeoInfo(info);
    } catch (err) {
      console.error('Failed to fetch SEO info');
    }
  };

  const updateSettings = (key: keyof SeoSettings, value: any) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = () => {
    setIsSaved(true);
    setTimeout(() => setIsSaved(false), 2500);
  };

  const handlePing = async () => {
    setPingLoading(true);
    try {
      const res = await systemApi.pingSearchEngines();
      if (res.success) {
        alert('Sitemap update notified to Google and Bing!');
        fetchSeoInfo();
      }
    } catch (err) {
      alert('Failed to ping search engines');
    } finally {
      setPingLoading(false);
    }
  };

  return {
    settings,
    seoInfo,
    activeTab,
    setActiveTab,
    isSaved,
    pingLoading,
    updateSettings,
    handleSave,
    handlePing,
  };
}
