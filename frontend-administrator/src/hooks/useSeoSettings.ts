import { useState, useEffect } from 'react';
import type { SeoSettings } from '../types/seo';
import { systemApi } from '../lib/api/system';

const initialSettings: SeoSettings = {
  site_title: 'FOSSBilling - Cloud Hosting & Domains',
  meta_description: 'Fast, secure cloud hosting and domain registrations.',
  meta_keywords: 'cloud hosting, domains, billing management',
  og_image_url: '',
  twitter_handle: '',
  google_analytics_id: '',
  robots_txt: `User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /api/\n`,
  sitemap_auto_generate: true,
  canonical_url: '',
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
