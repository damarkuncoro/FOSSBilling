import React from 'react';
import { X, Key, ShieldCheck, Clock } from 'lucide-react';
import type { SoftwareLicense, LicenseValidationLog } from '../../types/licenses';

interface LicenseDetailsModalProps {
  license: SoftwareLicense | null;
  logs: LicenseValidationLog[];
  onClose: () => void;
}

export const LicenseDetailsModal: React.FC<LicenseDetailsModalProps> = ({
  license,
  logs,
  onClose,
}) => {
  if (!license) return null;

  const licLogs = logs.filter((log) => log.license_key === license.license_key);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-2xl w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-indigo-600 font-semibold">
            <Key className="w-5 h-5" />
            <h3 className="text-gray-900 font-bold">License Details & Logs</h3>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="mt-4 grid grid-cols-2 gap-4 bg-gray-50 p-4 rounded-xl text-sm">
          <div>
            <span className="text-xs text-gray-500 block">License Key</span>
            <span className="font-mono font-bold text-indigo-600 text-sm">{license.license_key}</span>
          </div>
          <div>
            <span className="text-xs text-gray-500 block">Status</span>
            <span className="font-semibold capitalize text-gray-800">{license.status}</span>
          </div>
          <div>
            <span className="text-xs text-gray-500 block">Client</span>
            <span className="font-medium text-gray-800">{license.client_name}</span>
          </div>
          <div>
            <span className="text-xs text-gray-500 block">Product</span>
            <span className="font-medium text-gray-800">{license.product_title}</span>
          </div>
          <div>
            <span className="text-xs text-gray-500 block">Locked Domain</span>
            <span className="font-mono text-xs text-gray-700">{license.licensed_domain || 'Unrestricted'}</span>
          </div>
          <div>
            <span className="text-xs text-gray-500 block">Locked IP</span>
            <span className="font-mono text-xs text-gray-700">{license.licensed_ip || 'Unrestricted'}</span>
          </div>
        </div>

        <div className="mt-6">
          <h4 className="text-xs font-semibold text-gray-700 uppercase tracking-wider mb-3 flex items-center gap-1.5">
            <ShieldCheck className="w-4 h-4 text-emerald-600" /> Recent Validation API Calls
          </h4>
          {licLogs.length === 0 ? (
            <p className="text-xs text-gray-400 italic py-4 text-center bg-gray-50 rounded-lg">No validation requests logged yet.</p>
          ) : (
            <div className="divide-y divide-gray-100 border border-gray-100 rounded-lg overflow-hidden text-xs">
              {licLogs.map((log) => (
                <div key={log.id} className="p-3 flex items-center justify-between hover:bg-gray-50">
                  <div className="flex items-center gap-2">
                    <span className="px-2 py-0.5 rounded font-mono bg-emerald-50 text-emerald-700 font-semibold">
                      {log.result}
                    </span>
                    <span className="text-gray-600 font-mono">{log.domain} ({log.ip_address})</span>
                  </div>
                  <span className="text-gray-400 flex items-center gap-1">
                    <Clock className="w-3 h-3" /> {new Date(log.created_at).toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="mt-6 flex justify-end">
          <button onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg">
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
