import React from 'react';
import { Key, Globe, Server, CheckCircle2, AlertTriangle, XCircle, RotateCcw, ShieldAlert, Eye } from 'lucide-react';
import type { SoftwareLicense } from '../../types/licenses';

interface LicenseTableProps {
  licenses: SoftwareLicense[];
  onSelect: (license: SoftwareLicense) => void;
  onToggleStatus: (id: number) => void;
  onReissue: (id: number) => void;
  onResetLock: (id: number) => void;
}

export const LicenseTable: React.FC<LicenseTableProps> = ({
  licenses,
  onSelect,
  onToggleStatus,
  onReissue,
  onResetLock,
}) => {
  if (licenses.length === 0) {
    return (
      <div className="text-center py-12 bg-white rounded-xl border border-gray-100 shadow-sm">
        <Key className="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <h3 className="text-base font-semibold text-gray-800">No Licenses Found</h3>
        <p className="text-sm text-gray-500 mt-1">Generate a new software license key to get started.</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-gray-600">
          <thead className="bg-gray-50 text-xs uppercase font-semibold text-gray-500 border-b border-gray-100">
            <tr>
              <th className="px-6 py-4">License Key & Product</th>
              <th className="px-6 py-4">Client</th>
              <th className="px-6 py-4">Locked Domain / IP</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4">Expires</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {licenses.map((lic) => (
              <tr key={lic.id} className="hover:bg-gray-50/70 transition-colors">
                <td className="px-6 py-4">
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-indigo-50 text-indigo-600">
                      <Key className="w-4 h-4" />
                    </div>
                    <div>
                      <div className="font-mono font-medium text-gray-900">{lic.license_key}</div>
                      <div className="text-xs text-gray-500 mt-0.5">{lic.product_title}</div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 font-medium text-gray-900">{lic.client_name}</td>
                <td className="px-6 py-4">
                  <div className="flex flex-col gap-1 text-xs">
                    <span className="flex items-center gap-1.5 text-gray-700">
                      <Globe className="w-3.5 h-3.5 text-gray-400" />
                      {lic.licensed_domain || <span className="text-gray-400 italic">Any domain</span>}
                    </span>
                    <span className="flex items-center gap-1.5 text-gray-500">
                      <Server className="w-3.5 h-3.5 text-gray-400" />
                      {lic.licensed_ip || <span className="text-gray-400 italic">Any IP</span>}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span
                    className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${
                      lic.status === 'active'
                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-100'
                        : lic.status === 'suspended'
                        ? 'bg-amber-50 text-amber-700 border border-amber-100'
                        : 'bg-rose-50 text-rose-700 border border-rose-100'
                    }`}
                  >
                    {lic.status === 'active' ? (
                      <CheckCircle2 className="w-3.5 h-3.5" />
                    ) : lic.status === 'suspended' ? (
                      <AlertTriangle className="w-3.5 h-3.5" />
                    ) : (
                      <XCircle className="w-3.5 h-3.5" />
                    )}
                    {lic.status.toUpperCase()}
                  </span>
                </td>
                <td className="px-6 py-4 text-xs text-gray-500">
                  {lic.expires_at ? new Date(lic.expires_at).toLocaleDateString() : 'Lifetime'}
                </td>
                <td className="px-6 py-4 text-right">
                  <div className="flex items-center justify-end gap-1.5">
                    <button
                      onClick={() => onSelect(lic)}
                      className="p-1.5 rounded-lg text-gray-500 hover:text-indigo-600 hover:bg-gray-100 transition-colors"
                      title="View Details & Logs"
                    >
                      <Eye className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => onResetLock(lic.id)}
                      className="p-1.5 rounded-lg text-gray-500 hover:text-amber-600 hover:bg-amber-50 transition-colors"
                      title="Reset IP & Domain Lock"
                    >
                      <RotateCcw className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => onReissue(lic.id)}
                      className="p-1.5 rounded-lg text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 transition-colors"
                      title="Re-issue Key"
                    >
                      <Key className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => onToggleStatus(lic.id)}
                      className="p-1.5 rounded-lg text-gray-500 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                      title={lic.status === 'active' ? 'Suspend License' : 'Activate License'}
                    >
                      <ShieldAlert className="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
