import React from 'react';
import { Sliders, ShieldCheck } from 'lucide-react';
import type { CookieConsentSettings } from '../../types/cookieConsent';

interface BannerConfigCardProps {
  settings: CookieConsentSettings;
  onChange: (key: keyof CookieConsentSettings, val: any) => void;
}

export const BannerConfigCard: React.FC<BannerConfigCardProps> = ({ settings, onChange }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
      <div className="flex items-center justify-between pb-3 border-b border-gray-100">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Sliders className="w-4 h-4 text-indigo-600" /> Cookie Banner Settings
        </h3>
        <label className="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            checked={settings.is_enabled}
            onChange={(e) => onChange('is_enabled', e.target.checked)}
            className="sr-only peer"
          />
          <div className="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-indigo-600"></div>
          <span className="ml-2 text-xs font-semibold text-gray-700">
            {settings.is_enabled ? 'Active' : 'Disabled'}
          </span>
        </label>
      </div>

      <div>
        <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Consent Notice Message</label>
        <textarea
          rows={2}
          value={settings.message}
          onChange={(e) => onChange('message', e.target.value)}
          className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Banner Position</label>
          <select
            value={settings.banner_position}
            onChange={(e) => onChange('banner_position', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs bg-white outline-none focus:border-indigo-500"
          >
            <option value="bottom">Bottom Full Width</option>
            <option value="top">Top Full Width</option>
            <option value="floating_left">Floating Box (Bottom Left)</option>
            <option value="floating_right">Floating Box (Bottom Right)</option>
          </select>
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Color Theme</label>
          <select
            value={settings.theme}
            onChange={(e) => onChange('theme', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs bg-white outline-none focus:border-indigo-500"
          >
            <option value="dark">Dark Slate</option>
            <option value="light">Clean Light</option>
            <option value="indigo">Indigo Accent</option>
          </select>
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Accept Button Text</label>
          <input
            type="text"
            value={settings.accept_button_text}
            onChange={(e) => onChange('accept_button_text', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
          />
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Decline Button Text</label>
          <input
            type="text"
            value={settings.decline_button_text}
            onChange={(e) => onChange('decline_button_text', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
          />
        </div>
      </div>
    </div>
  );
};
