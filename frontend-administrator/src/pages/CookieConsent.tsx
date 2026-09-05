import React from 'react';
import { ShieldCheck, Save, CheckCircle2 } from 'lucide-react';
import { useCookieConsent } from '../hooks/useCookieConsent';
import { BannerConfigCard } from '../components/cookieconsent/BannerConfigCard';
import { BannerLivePreview } from '../components/cookieconsent/BannerLivePreview';
import { ConsentLogsTable } from '../components/cookieconsent/ConsentLogsTable';

export const CookieConsent: React.FC = () => {
  const {
    settings,
    logs,
    isSaved,
    updateSetting,
    handleSave,
  } = useCookieConsent();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <ShieldCheck className="w-7 h-7 text-indigo-600" /> Cookie Consent & GDPR
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Configure EU GDPR cookie notification banner, tracking options, and view user audit logs.
          </p>
        </div>
        <button
          onClick={handleSave}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          {isSaved ? <CheckCircle2 className="w-4 h-4 text-emerald-300" /> : <Save className="w-4 h-4" />}
          {isSaved ? 'Saved!' : 'Save Banner Settings'}
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <BannerConfigCard settings={settings} onChange={updateSetting} />
        <BannerLivePreview settings={settings} />
      </div>

      <ConsentLogsTable logs={logs} />
    </div>
  );
};
