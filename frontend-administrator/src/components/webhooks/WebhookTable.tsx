import React from 'react';
import { Radio, CheckCircle2, XCircle, Send, Trash2, Shield, Clock } from 'lucide-react';
import type { WebhookEndpoint } from '../../types/webhooks';

interface WebhookTableProps {
  webhooks: WebhookEndpoint[];
  testingId: number | null;
  onToggle: (id: number) => void;
  onDelete: (id: number) => void;
  onTest: (id: number) => void;
}

export const WebhookTable: React.FC<WebhookTableProps> = ({
  webhooks,
  testingId,
  onToggle,
  onDelete,
  onTest,
}) => {
  if (webhooks.length === 0) {
    return (
      <div className="text-center py-12 bg-white rounded-xl border border-gray-100 shadow-sm">
        <Radio className="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <h3 className="text-base font-semibold text-gray-800">No Webhooks Configured</h3>
        <p className="text-sm text-gray-500 mt-1">Set up automated event webhooks for Discord, Slack, or custom endpoints.</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-gray-600">
          <thead className="bg-gray-50 text-xs uppercase font-semibold text-gray-500 border-b border-gray-100">
            <tr>
              <th className="px-6 py-4">Endpoint Name & URL</th>
              <th className="px-6 py-4">Subscribed Events</th>
              <th className="px-6 py-4">Deliveries</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {webhooks.map((wh) => (
              <tr key={wh.id} className="hover:bg-gray-50/70 transition-colors">
                <td className="px-6 py-4">
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-indigo-50 text-indigo-600">
                      <Radio className="w-4 h-4" />
                    </div>
                    <div>
                      <div className="font-semibold text-gray-900">{wh.name}</div>
                      <div className="font-mono text-xs text-gray-500 max-w-sm truncate mt-0.5">{wh.url}</div>
                      <div className="flex items-center gap-1.5 text-[11px] font-mono text-gray-400 mt-1">
                        <Shield className="w-3 h-3 text-emerald-500" /> {wh.secret.substring(0, 16)}...
                      </div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-wrap gap-1 max-w-xs">
                    {wh.events.map((ev) => (
                      <span key={ev} className="px-2 py-0.5 rounded text-[11px] font-mono bg-gray-100 text-gray-700">
                        {ev}
                      </span>
                    ))}
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="text-xs font-semibold text-gray-800">{wh.total_deliveries} sent</div>
                  {wh.last_delivery_at && (
                    <div className="text-[11px] text-gray-400 flex items-center gap-1 mt-0.5">
                      <Clock className="w-3 h-3" /> {new Date(wh.last_delivery_at).toLocaleTimeString()}
                    </div>
                  )}
                </td>
                <td className="px-6 py-4">
                  <button
                    onClick={() => onToggle(wh.id)}
                    className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold transition-colors ${
                      wh.is_active
                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-100 hover:bg-emerald-100'
                        : 'bg-gray-100 text-gray-500 border border-gray-200 hover:bg-gray-200'
                    }`}
                  >
                    {wh.is_active ? <CheckCircle2 className="w-3.5 h-3.5" /> : <XCircle className="w-3.5 h-3.5" />}
                    {wh.is_active ? 'ACTIVE' : 'PAUSED'}
                  </button>
                </td>
                <td className="px-6 py-4 text-right">
                  <div className="flex items-center justify-end gap-2">
                    <button
                      onClick={() => onTest(wh.id)}
                      disabled={testingId === wh.id}
                      className="inline-flex items-center gap-1 px-3 py-1.5 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 rounded-lg text-xs font-medium transition-colors"
                      title="Send Test Ping Payload"
                    >
                      <Send className={`w-3.5 h-3.5 ${testingId === wh.id ? 'animate-spin' : ''}`} />
                      {testingId === wh.id ? 'Sending...' : 'Test Ping'}
                    </button>
                    <button
                      onClick={() => onDelete(wh.id)}
                      className="p-1.5 text-gray-400 hover:text-rose-600 rounded-lg transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
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
