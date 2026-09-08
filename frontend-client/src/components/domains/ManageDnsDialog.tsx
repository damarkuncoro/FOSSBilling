import React, { useState, useEffect } from 'react';
import { X, Server, Save, LayoutGrid } from 'lucide-react';
import type { DomainRecord } from '../../types/clientModules';
import { ZoneEditor } from './ZoneEditor';

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
  const [activeTab, setActiveTab] = useState<'ns' | 'zone'>('ns');
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
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm">
      <div className="bg-white rounded-3xl max-w-2xl w-full p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100 mb-6">
          <div className="flex items-center gap-2.5 text-indigo-600">
            <div className="p-2 bg-indigo-50 rounded-xl">
              <Server className="w-5 h-5" />
            </div>
            <h3 className="text-gray-900 font-black text-xl">Manage Domain DNS</h3>
          </div>
          <button onClick={onClose} className="p-2 rounded-xl text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-all">
            <X className="w-6 h-6" />
          </button>
        </div>

        <div className="mb-6 p-4 bg-indigo-50/50 rounded-2xl border border-indigo-100 flex items-center justify-between">
           <div>
              <p className="text-[10px] font-bold text-indigo-500 uppercase tracking-widest mb-0.5">Active Domain</p>
              <strong className="font-mono text-base text-indigo-900">{domain.domain_name}</strong>
           </div>
           <Badge variant={domain.status === 'active' ? 'success' : 'warning'} className="h-fit">
              {domain.status.toUpperCase()}
           </Badge>
        </div>

        <div className="flex gap-1 bg-gray-100 p-1 rounded-xl mb-6">
           <button
             onClick={() => setActiveTab('ns')}
             className={`flex-1 py-2 text-xs font-bold rounded-lg transition-all ${activeTab === 'ns' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
           >
              Nameservers
           </button>
           <button
             onClick={() => setActiveTab('zone')}
             className={`flex-1 py-2 text-xs font-bold rounded-lg transition-all ${activeTab === 'zone' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-500 hover:text-gray-700'}`}
           >
              DNS Zone Records
           </button>
        </div>

        {activeTab === 'ns' ? (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase mb-1.5 ml-1">Nameserver 1</label>
                <input
                  type="text"
                  required
                  value={ns1}
                  onChange={(e) => setNs1(e.target.value)}
                  className="w-full px-4 py-2.5 font-mono text-sm border border-gray-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                />
              </div>

              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase mb-1.5 ml-1">Nameserver 2</label>
                <input
                  type="text"
                  required
                  value={ns2}
                  onChange={(e) => setNs2(e.target.value)}
                  className="w-full px-4 py-2.5 font-mono text-sm border border-gray-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                />
              </div>

              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase mb-1.5 ml-1">Nameserver 3 (Optional)</label>
                <input
                  type="text"
                  value={ns3}
                  onChange={(e) => setNs3(e.target.value)}
                  className="w-full px-4 py-2.5 font-mono text-sm border border-gray-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                />
              </div>

              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase mb-1.5 ml-1">Nameserver 4 (Optional)</label>
                <input
                  type="text"
                  value={ns4}
                  onChange={(e) => setNs4(e.target.value)}
                  className="w-full px-4 py-2.5 font-mono text-sm border border-gray-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                />
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-6 border-t border-gray-100">
              <button type="button" onClick={onClose} className="px-6 py-2.5 text-sm font-bold text-gray-500 hover:bg-gray-50 rounded-xl transition-colors">
                Cancel
              </button>
              <button type="submit" className="inline-flex items-center gap-2 px-8 py-2.5 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl shadow-lg shadow-indigo-500/20 transition-all">
                <Save className="w-4 h-4" /> Update Nameservers
              </button>
            </div>
          </form>
        ) : (
          <ZoneEditor domainId={domain.id} />
        )}
      </div>
    </div>
  );
};

const Badge = ({ children, variant, className }: any) => (
  <span className={`px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider border ${
    variant === 'success' ? 'bg-emerald-50 text-emerald-700 border-emerald-100' : 'bg-amber-50 text-amber-700 border-amber-100'
  } ${className}`}>
    {children}
  </span>
);
