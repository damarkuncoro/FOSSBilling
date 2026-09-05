import React, { useState } from 'react';
import { X, Key, RefreshCw } from 'lucide-react';
import type { SoftwareLicense } from '../../types/licenses';

interface AddLicenseDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Partial<SoftwareLicense>) => void;
  onGenerateKey: () => string;
}

export const AddLicenseDialog: React.FC<AddLicenseDialogProps> = ({
  isOpen,
  onClose,
  onSubmit,
  onGenerateKey,
}) => {
  const [clientName, setClientName] = useState('');
  const [productTitle, setProductTitle] = useState('SaaS Enterprise License');
  const [licenseKey, setLicenseKey] = useState(onGenerateKey());
  const [domain, setDomain] = useState('');
  const [ip, setIp] = useState('');
  const [maxInstances, setMaxInstances] = useState(1);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      client_name: clientName || 'New Client',
      product_title: productTitle,
      license_key: licenseKey,
      licensed_domain: domain,
      licensed_ip: ip,
      max_instances: maxInstances,
    });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-indigo-600 font-semibold">
            <Key className="w-5 h-5" />
            <h3 className="text-gray-900 font-bold">Issue New License</h3>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase tracking-wider mb-1">
              Client Name
            </label>
            <input
              type="text"
              required
              placeholder="e.g. John Doe / PT Example"
              value={clientName}
              onChange={(e) => setClientName(e.target.value)}
              className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase tracking-wider mb-1">
              Product / Software
            </label>
            <input
              type="text"
              required
              value={productTitle}
              onChange={(e) => setProductTitle(e.target.value)}
              className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="block text-xs font-semibold text-gray-700 uppercase tracking-wider">
                License Key
              </label>
              <button
                type="button"
                onClick={() => setLicenseKey(onGenerateKey())}
                className="text-xs text-indigo-600 hover:underline flex items-center gap-1"
              >
                <RefreshCw className="w-3 h-3" /> Regenerate
              </button>
            </div>
            <input
              type="text"
              required
              value={licenseKey}
              onChange={(e) => setLicenseKey(e.target.value)}
              className="w-full px-3.5 py-2 font-mono border border-gray-200 rounded-lg text-sm bg-gray-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-semibold text-gray-700 uppercase tracking-wider mb-1">
                Lock Domain (Optional)
              </label>
              <input
                type="text"
                placeholder="app.client.com"
                value={domain}
                onChange={(e) => setDomain(e.target.value)}
                className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-gray-700 uppercase tracking-wider mb-1">
                Lock IP (Optional)
              </label>
              <input
                type="text"
                placeholder="192.0.2.1"
                value={ip}
                onChange={(e) => setIp(e.target.value)}
                className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
              />
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm"
            >
              Issue License
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
