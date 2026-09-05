import React from 'react';
import { Globe, Save, CheckCircle2, Bot, Map } from 'lucide-react';
import { useSeoSettings } from '../hooks/useSeoSettings';
import { MetaTagsTab } from '../components/seo/MetaTagsTab';
import { RobotsTxtEditor } from '../components/seo/RobotsTxtEditor';
import { SitemapGeneratorCard } from '../components/seo/SitemapGeneratorCard';

export const SeoSettings: React.FC = () => {
  const {
    settings,
    activeTab,
    setActiveTab,
    isSaved,
    updateSettings,
    handleSave,
  } = useSeoSettings();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <Globe className="w-7 h-7 text-indigo-600" /> SEO & Webmaster Tools
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Manage search engine indexing, OpenGraph metadata, XML sitemaps, and robots.txt.
          </p>
        </div>
        <button
          onClick={handleSave}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          {isSaved ? <CheckCircle2 className="w-4 h-4 text-emerald-300" /> : <Save className="w-4 h-4" />}
          {isSaved ? 'Settings Saved!' : 'Save SEO Settings'}
        </button>
      </div>

      <div className="flex border-b border-gray-200">
        <button
          onClick={() => setActiveTab('meta')}
          className={`pb-3 px-4 text-sm font-semibold flex items-center gap-2 border-b-2 transition-all ${
            activeTab === 'meta'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          <Globe className="w-4 h-4" /> Meta Tags & Social Share
        </button>
        <button
          onClick={() => setActiveTab('robots')}
          className={`pb-3 px-4 text-sm font-semibold flex items-center gap-2 border-b-2 transition-all ${
            activeTab === 'robots'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          <Bot className="w-4 h-4" /> robots.txt Editor
        </button>
        <button
          onClick={() => setActiveTab('sitemap')}
          className={`pb-3 px-4 text-sm font-semibold flex items-center gap-2 border-b-2 transition-all ${
            activeTab === 'sitemap'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          <Map className="w-4 h-4" /> XML Sitemap
        </button>
      </div>

      {activeTab === 'meta' && (
        <MetaTagsTab settings={settings} onChange={updateSettings} />
      )}

      {activeTab === 'robots' && (
        <RobotsTxtEditor
          robotsTxt={settings.robots_txt}
          onChange={(val) => updateSettings('robots_txt', val)}
        />
      )}

      {activeTab === 'sitemap' && (
        <SitemapGeneratorCard
          autoGenerate={settings.sitemap_auto_generate}
          onToggleAuto={(val) => updateSettings('sitemap_auto_generate', val)}
        />
      )}
    </div>
  );
};
