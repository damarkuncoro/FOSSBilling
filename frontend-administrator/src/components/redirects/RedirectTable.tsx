import React from 'react';
import { ArrowRightLeft, ArrowRight, CheckCircle2, XCircle, Trash2, Eye } from 'lucide-react';
import type { UrlRedirect } from '../../types/redirects';

interface RedirectTableProps {
  redirects: UrlRedirect[];
  onToggle: (id: number) => void;
  onDelete: (id: number) => void;
}

export const RedirectTable: React.FC<RedirectTableProps> = ({
  redirects,
  onToggle,
  onDelete,
}) => {
  if (redirects.length === 0) {
    return (
      <div className="text-center py-12 bg-white rounded-xl border border-gray-100 shadow-sm">
        <ArrowRightLeft className="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <h3 className="text-base font-semibold text-gray-800">No URL Redirects Found</h3>
        <p className="text-sm text-gray-500 mt-1">Create your first 301 or 302 redirect rule.</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-gray-600">
          <thead className="bg-gray-50 text-xs uppercase font-semibold text-gray-500 border-b border-gray-100">
            <tr>
              <th className="px-6 py-4">Source Path</th>
              <th className="px-6 py-4">Target Destination</th>
              <th className="px-6 py-4">Type</th>
              <th className="px-6 py-4">Hits</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {redirects.map((r) => (
              <tr key={r.id} className="hover:bg-gray-50/70 transition-colors">
                <td className="px-6 py-4 font-mono font-medium text-gray-900">
                  {r.source_path}
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-1.5 font-mono text-xs text-indigo-600 max-w-sm truncate">
                    <ArrowRight className="w-3.5 h-3.5 text-gray-400 shrink-0" />
                    {r.target_url}
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-0.5 rounded text-xs font-mono font-bold ${
                    r.status_code === 301 ? 'bg-purple-50 text-purple-700' : 'bg-blue-50 text-blue-700'
                  }`}>
                    {r.status_code} {r.status_code === 301 ? 'Permanent' : 'Temporary'}
                  </span>
                </td>
                <td className="px-6 py-4 font-medium text-xs text-gray-800 flex items-center gap-1 mt-3">
                  <Eye className="w-3.5 h-3.5 text-gray-400" /> {r.hits_count.toLocaleString()}
                </td>
                <td className="px-6 py-4">
                  <button
                    onClick={() => onToggle(r.id)}
                    className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${
                      r.is_active
                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-100'
                        : 'bg-gray-100 text-gray-500 border border-gray-200'
                    }`}
                  >
                    {r.is_active ? <CheckCircle2 className="w-3.5 h-3.5" /> : <XCircle className="w-3.5 h-3.5" />}
                    {r.is_active ? 'ACTIVE' : 'DISABLED'}
                  </button>
                </td>
                <td className="px-6 py-4 text-right">
                  <button
                    onClick={() => onDelete(r.id)}
                    className="p-1.5 text-gray-400 hover:text-rose-600 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
