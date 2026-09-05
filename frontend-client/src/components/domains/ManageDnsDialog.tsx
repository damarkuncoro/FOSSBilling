import React, { useState, useEffect } from 'react';
import { X, Server, Save } from 'lucide-react';
import type { DomainRecord } from '../../types/clientModules';

interface ManageDnsDialogProps {
  domain: DomainRecord | null;
  onClose: () => void;
  onSave: (id: number, nameservers: string[]) => void;
}

export const ManageDnsDialog: React.FC<ManageDnsDialogProps> = ({
  domain,
  onClose,
  onSave,
}) => {
  const [ns1, setNs1] = useState('');
  const [ns2, setNs2] = useState('');
  const [ns3, setNs3] = useState('');
  const [ns4, setNs4] = useState('');

  useEffect(() => {
    if (domain) {
      setNs1(domain.nameservers[0] || 'ns1.fossbilling.org');
      setNs2(domain.nameservers[1] || 'ns2.fossbilling.org');
      setNs3(domain.nameservers[2] || '');
      setNs4(domain.nameservers[3] || '');
    }
  }, [domain]);

  if (!domain) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const nsList = [ns1.trim(), ns2.trim(), ns3.trim(), ns4.trim()].filter(Boolean);
    onSave(domain.id, nsList);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-indigo-600 font-semibold">
            <Server className="w-5 h-5" />
            <h3 className="text-gray-900 font-bold">Manage Nameservers (DNS)</h3>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="my-3 p-3 bg-gray-50 rounded-xl border border-gray-100 text-xs text-gray-600">
          Domain: <strong className="font-mono text-gray-900">{domain.domain_name}</strong>
        </div>

        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Nameserver 1</label>
            <input
              type="text"
              required
              value={ns1}
              onChange={(e) => setNs1(e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Nameserver 2</label>
            <input
              type="text"
              required
              value={ns2}
              onChange={(e) => setNs2(e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Nameserver 3 (Optional)</label>
            <input
              type="text"
              value={ns3}
              onChange={(e) => setNs3(e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Nameserver 4 (Optional)</label>
            <input
              type="text"
              value={ns4}
              onChange={(e) => setNs4(e.target.value)}
              className="w-full px-3.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none focus:border-indigo-500"
            />
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
            <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg">
              Cancel
            </button>
            <button type="submit" className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">
              <Save className="w-4 h-4" /> Save Changes
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
