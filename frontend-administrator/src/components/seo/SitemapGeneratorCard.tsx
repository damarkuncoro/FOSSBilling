import React from 'react';
import { Map, RefreshCw, ExternalLink, CheckCircle2 } from 'lucide-react';

interface SitemapGeneratorCardProps {
  autoGenerate: boolean;
  onToggleAuto: (val: boolean) => void;
}

export const SitemapGeneratorCard: React.FC<SitemapGeneratorCardProps> = ({
  autoGenerate,
  onToggleAuto,
}) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Map className="w-4 h-4 text-indigo-600" /> XML Sitemap Generator
        </h3>
        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-100">
          <CheckCircle2 className="w-3.5 h-3.5" /> Indexed & Ready
        </span>
      </div>

      <div className="p-4 bg-gray-50 rounded-xl border border-gray-100 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <div className="text-xs font-semibold text-gray-800">Public Sitemap Location</div>
          <div className="font-mono text-xs text-indigo-600 mt-0.5">https://myhosting.com/sitemap.xml</div>
          <div className="text-[11px] text-gray-400 mt-1">Includes 18 products, 4 custom pages, and 12 announcements.</div>
        </div>
        <div className="flex items-center gap-2">
          <button className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white border border-gray-200 text-gray-700 hover:bg-gray-50 rounded-lg text-xs font-medium shadow-sm">
            <ExternalLink className="w-3.5 h-3.5" /> View XML
          </button>
          <button className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-600 text-white hover:bg-indigo-700 rounded-lg text-xs font-medium shadow-sm">
            <RefreshCw className="w-3.5 h-3.5" /> Rebuild Now
          </button>
        </div>
      </div>

      <div className="flex items-center justify-between pt-2">
        <div>
          <div className="text-sm font-medium text-gray-900">Automated Daily Sitemap Sync</div>
          <p className="text-xs text-gray-500">Automatically update sitemap when new products, pages, or news articles are published.</p>
        </div>
        <input
          type="checkbox"
          checked={autoGenerate}
          onChange={(e) => onToggleAuto(e.target.checked)}
          className="w-5 h-5 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500"
        />
      </div>
    </div>
  );
};
