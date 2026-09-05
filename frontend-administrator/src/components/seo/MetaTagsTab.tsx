import React from 'react';
import { Globe, Share2 } from 'lucide-react';
import type { SeoSettings } from '../../types/seo';

interface MetaTagsTabProps {
  settings: SeoSettings;
  onChange: (key: keyof SeoSettings, value: any) => void;
}

export const MetaTagsTab: React.FC<MetaTagsTabProps> = ({ settings, onChange }) => {
  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Globe className="w-4 h-4 text-indigo-600" /> Default Meta Tags & Indexing
        </h3>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Global Meta Title</label>
          <input
            type="text"
            value={settings.site_title}
            onChange={(e) => onChange('site_title', e.target.value)}
            className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
          />
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Meta Description</label>
          <textarea
            rows={2}
            value={settings.meta_description}
            onChange={(e) => onChange('meta_description', e.target.value)}
            className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Canonical Base URL</label>
            <input
              type="text"
              value={settings.canonical_url}
              onChange={(e) => onChange('canonical_url', e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Google Analytics / GTM ID</label>
            <input
              type="text"
              placeholder="G-XXXXXX or UA-XXXXXX"
              value={settings.google_analytics_id}
              onChange={(e) => onChange('google_analytics_id', e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>
        </div>
      </div>

      {/* Social Card Preview */}
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Share2 className="w-4 h-4 text-indigo-600" /> OpenGraph & Social Sharing Preview
        </h3>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">OG Image URL</label>
          <input
            type="text"
            value={settings.og_image_url}
            onChange={(e) => onChange('og_image_url', e.target.value)}
            className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
          />
        </div>

        <div className="p-4 bg-gray-50 rounded-xl max-w-lg border border-gray-200/80">
          <div className="aspect-[1.91/1] w-full rounded-lg overflow-hidden bg-gray-200 mb-3">
            <img src={settings.og_image_url} alt="OG Preview" className="w-full h-full object-cover" />
          </div>
          <div className="text-xs font-semibold text-gray-500 uppercase tracking-wider">{settings.canonical_url.replace('https://', '')}</div>
          <div className="font-bold text-sm text-gray-900 mt-0.5">{settings.site_title}</div>
          <div className="text-xs text-gray-600 mt-1 line-clamp-2">{settings.meta_description}</div>
        </div>
      </div>
    </div>
  );
};
