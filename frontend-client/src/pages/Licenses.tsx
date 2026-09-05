import React from 'react';
import { Key, Copy, CheckCircle2, Globe, Server, RotateCcw } from 'lucide-react';
import { useClientLicenses } from '../hooks/useClientLicenses';

export const Licenses: React.FC = () => {
  const { licenses, copiedKey, copyKey, resetLock } = useClientLicenses();

  return (
    <div className="max-w-4xl mx-auto px-4 py-8 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
          <Key className="w-7 h-7 text-indigo-600" /> Software Licenses
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          View your purchased software keys, active domain/IP bindings, and reset locks.
        </p>
      </div>

      <div className="space-y-4">
        {licenses.map((lic) => (
          <div
            key={lic.id}
            className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 space-y-4"
          >
            <div className="flex items-center justify-between">
              <div>
                <span className="text-xs font-semibold text-indigo-600 uppercase tracking-wider">Product License</span>
                <h3 className="text-lg font-bold text-gray-900 mt-0.5">{lic.product_title}</h3>
              </div>
              <span className="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-100">
                ACTIVE
              </span>
            </div>

            <div className="p-3.5 bg-gray-50 rounded-xl border border-gray-200 flex items-center justify-between">
              <div>
                <span className="text-[11px] text-gray-500 block uppercase font-semibold">License Key</span>
                <span className="font-mono font-bold text-gray-900 text-sm">{lic.license_key}</span>
              </div>
              <button
                onClick={() => copyKey(lic.license_key)}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white border border-gray-200 hover:bg-gray-50 text-gray-700 rounded-lg text-xs font-medium shadow-sm transition-colors"
              >
                {copiedKey === lic.license_key ? <CheckCircle2 className="w-3.5 h-3.5 text-emerald-600" /> : <Copy className="w-3.5 h-3.5" />}
                {copiedKey === lic.license_key ? 'Copied!' : 'Copy Key'}
              </button>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
              <div className="flex items-center gap-2 text-gray-600">
                <Globe className="w-4 h-4 text-gray-400" />
                <span>Bound Domain: <strong>{lic.licensed_domain || 'Any Domain'}</strong></span>
              </div>
              <div className="flex items-center gap-2 text-gray-600">
                <Server className="w-4 h-4 text-gray-400" />
                <span>Bound IP: <strong>{lic.licensed_ip || 'Any IP'}</strong></span>
              </div>
            </div>

            <div className="pt-3 border-t border-gray-100 flex justify-end">
              <button
                onClick={() => resetLock(lic.id)}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 text-amber-700 bg-amber-50 hover:bg-amber-100 rounded-lg text-xs font-semibold transition-colors"
              >
                <RotateCcw className="w-3.5 h-3.5" /> Re-lock to New Domain / Server
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
