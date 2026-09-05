import React from 'react';
import { Eye, Shield } from 'lucide-react';
import type { CookieConsentSettings } from '../../types/cookieConsent';

interface BannerLivePreviewProps {
  settings: CookieConsentSettings;
}

export const BannerLivePreview: React.FC<BannerLivePreviewProps> = ({ settings }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
      <div className="flex items-center gap-2 text-xs font-semibold text-gray-500 uppercase tracking-wider pb-3 border-b border-gray-100">
        <Eye className="w-4 h-4 text-indigo-600" /> Live Client Banner Preview
      </div>

      <div className="relative h-44 bg-gray-100/70 border border-dashed border-gray-300 rounded-xl overflow-hidden flex flex-col justify-end p-3">
        <div className="text-center text-xs text-gray-400 my-auto">
          Client Storefront Area Simulated Viewport
        </div>

        {settings.is_enabled && (
          <div
            className={`p-3.5 rounded-xl shadow-lg border text-xs flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 ${
              settings.theme === 'dark'
                ? 'bg-gray-900 text-white border-gray-800'
                : settings.theme === 'indigo'
                ? 'bg-indigo-950 text-white border-indigo-900'
                : 'bg-white text-gray-800 border-gray-200'
            }`}
          >
            <div className="flex items-center gap-2">
              <Shield className="w-4 h-4 text-indigo-400 shrink-0" />
              <p className="line-clamp-2 text-[11px] text-gray-300">{settings.message}</p>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              {settings.show_decline_button && (
                <button className="px-2.5 py-1 text-[11px] font-medium text-gray-300 bg-white/10 hover:bg-white/20 rounded-lg">
                  {settings.decline_button_text}
                </button>
              )}
              <button className="px-3 py-1 text-[11px] font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">
                {settings.accept_button_text}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
