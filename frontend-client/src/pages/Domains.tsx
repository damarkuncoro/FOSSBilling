import React from 'react';
import { Globe, Search, CheckCircle2, XCircle, ShoppingCart } from 'lucide-react';
import { useClientDomains } from '../hooks/useClientDomains';
import { useCart } from '../lib/cart';
import { ManageDnsDialog } from '../components/domains/ManageDnsDialog';

export const Domains: React.FC = () => {
  const {
    domains,
    search,
    setSearch,
    checkQuery,
    setCheckQuery,
    checkResult,
    isSearching,
    editingDomain,
    setEditingDomain,
    checkAvailability,
    updateNameservers,
    toggleAutoRenew,
  } = useClientDomains();

  const { addItem } = useCart();

  const handleAddToCart = () => {
    if (!checkResult) return;
    addItem({
      id: `cart_domain_${Date.now()}`,
      product_id: 99,
      title: `Domain Registration: ${checkResult.domain}`,
      type: 'domain',
      domain_name: checkResult.domain,
      period: '1Y',
      price: checkResult.price,
    });
  };

  return (
    <div className="max-w-5xl mx-auto px-4 py-8 space-y-8">
      {/* Domain Lookup Card */}
      <div className="bg-gradient-to-br from-indigo-900 to-slate-900 rounded-3xl p-8 text-white shadow-xl space-y-6">
        <div className="max-w-xl space-y-2">
          <span className="text-xs font-semibold uppercase tracking-wider text-indigo-300">WHOIS Checker</span>
          <h2 className="text-2xl font-black tracking-tight">Register Your Next Domain</h2>
          <p className="text-xs text-indigo-200">Instant registration with free DNS management and WHOIS privacy protection.</p>
        </div>

        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <input
              type="text"
              placeholder="Find your new domain (e.g. mycompany.com)..."
              value={checkQuery}
              onChange={(e) => setCheckQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && checkAvailability()}
              className="w-full px-4 py-3.5 bg-white/10 border border-white/20 rounded-xl text-white placeholder:text-gray-400 text-sm backdrop-blur-md outline-none focus:ring-2 focus:ring-indigo-400"
            />
          </div>
          <button
            onClick={checkAvailability}
            disabled={isSearching}
            className="px-6 py-3.5 bg-indigo-500 hover:bg-indigo-600 text-white font-bold text-sm rounded-xl transition-all shadow-lg flex items-center justify-center gap-2"
          >
            <Search className="w-4 h-4" /> {isSearching ? 'Checking...' : 'Check Availability'}
          </button>
        </div>

        {checkResult && (
          <div className="p-4 bg-white/10 rounded-2xl border border-white/20 backdrop-blur-md flex items-center justify-between animate-in fade-in">
            <div className="flex items-center gap-3">
              {checkResult.available ? <CheckCircle2 className="w-6 h-6 text-emerald-400" /> : <XCircle className="w-6 h-6 text-rose-400" />}
              <div>
                <span className="font-mono font-bold text-base">{checkResult.domain}</span>
                <p className="text-xs text-indigo-200">
                  {checkResult.available ? `Available for $${checkResult.price}/year` : 'Already registered by another owner.'}
                </p>
              </div>
            </div>
            {checkResult.available && (
              <button
                onClick={handleAddToCart}
                className="inline-flex items-center gap-1.5 px-4 py-2 bg-emerald-500 hover:bg-emerald-600 text-white font-semibold text-xs rounded-xl shadow-sm transition-colors"
              >
                <ShoppingCart className="w-4 h-4" /> Add to Cart
              </button>
            )}
          </div>
        )}
      </div>

      {/* Active Domains Table */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-bold text-gray-900 flex items-center gap-2">
            <Globe className="w-5 h-5 text-indigo-600" /> My Active Domains
          </h3>
          <input
            type="text"
            placeholder="Filter domains..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3.5 py-1.5 bg-white border border-gray-200 rounded-xl text-xs outline-none focus:border-indigo-500"
          />
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <table className="w-full text-left text-sm text-gray-600">
            <thead className="bg-gray-50 text-xs uppercase font-semibold text-gray-500 border-b border-gray-100">
              <tr>
                <th className="px-6 py-4">Domain Name</th>
                <th className="px-6 py-4">Nameservers</th>
                <th className="px-6 py-4">Expires</th>
                <th className="px-6 py-4">Auto-Renew</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {domains.map((d) => (
                <tr key={d.id} className="hover:bg-gray-50/70">
                  <td className="px-6 py-4 font-mono font-bold text-gray-900">{d.domain_name}</td>
                  <td className="px-6 py-4 text-xs font-mono text-gray-500">
                    {d.nameservers.slice(0, 2).join(', ')}
                  </td>
                  <td className="px-6 py-4 text-xs text-gray-500">
                    {new Date(d.expires_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4">
                    <button
                      onClick={() => toggleAutoRenew(d.id)}
                      className={`px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                        d.auto_renew ? 'bg-emerald-50 text-emerald-700' : 'bg-gray-100 text-gray-500'
                      }`}
                    >
                      {d.auto_renew ? 'Enabled' : 'Disabled'}
                    </button>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => setEditingDomain(d)}
                      className="px-3 py-1 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 rounded-lg text-xs font-medium"
                    >
                      Manage DNS
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <ManageDnsDialog
        domain={editingDomain}
        onClose={() => setEditingDomain(null)}
        onSave={updateNameservers}
      />
    </div>
  );
};
