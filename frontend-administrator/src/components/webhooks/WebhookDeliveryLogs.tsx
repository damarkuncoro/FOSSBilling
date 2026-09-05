import React from 'react';
import { Activity, Clock, CheckCircle2, AlertCircle } from 'lucide-react';
import type { WebhookDeliveryLog } from '../../types/webhooks';

interface WebhookDeliveryLogsProps {
  logs: WebhookDeliveryLog[];
}

export const WebhookDeliveryLogs: React.FC<WebhookDeliveryLogsProps> = ({ logs }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6">
      <div className="flex items-center gap-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4 pb-3 border-b border-gray-100">
        <Activity className="w-4 h-4 text-indigo-600" /> Recent Webhook Dispatches
      </div>

      {logs.length === 0 ? (
        <div className="text-center py-6 text-gray-400 text-xs italic">
          No dispatch attempts logged yet.
        </div>
      ) : (
        <div className="divide-y divide-gray-100">
          {logs.map((log) => (
            <div key={log.id} className="py-3 flex items-center justify-between text-xs">
              <div className="flex items-center gap-2.5">
                {log.response_code >= 200 && log.response_code < 300 ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-500" />
                ) : (
                  <AlertCircle className="w-4 h-4 text-rose-500" />
                )}
                <div>
                  <div className="flex items-center gap-2 font-mono">
                    <span className="font-bold text-gray-800">{log.event}</span>
                    <span className="px-1.5 py-0.5 rounded bg-gray-100 text-gray-600 text-[10px]">
                      HTTP {log.response_code}
                    </span>
                    <span className="text-gray-400 text-[11px]">{log.execution_time_ms}ms</span>
                  </div>
                  <div className="text-[11px] text-gray-400 truncate max-w-md">{log.url}</div>
                </div>
              </div>

              <div className="text-gray-400 flex items-center gap-1 text-[11px]">
                <Clock className="w-3 h-3" />
                {new Date(log.created_at).toLocaleTimeString()}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
