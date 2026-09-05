import React from 'react';
import { History, CheckCircle2, XCircle, Clock } from 'lucide-react';
import type { ConsentLog } from '../../types/cookieConsent';

interface ConsentLogsTableProps {
  logs: ConsentLog[];
}

export const ConsentLogsTable: React.FC<ConsentLogsTableProps> = ({ logs }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6">
      <div className="flex items-center gap-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4 pb-3 border-b border-gray-100">
        <History className="w-4 h-4 text-indigo-600" /> Recent User Consent Records (GDPR Audit Trail)
      </div>

      <div className="divide-y divide-gray-100">
        {logs.map((log) => (
          <div key={log.id} className="py-2.5 flex items-center justify-between text-xs">
            <div className="flex items-center gap-3">
              {log.decision === 'accepted' ? (
                <CheckCircle2 className="w-4 h-4 text-emerald-500" />
              ) : (
                <XCircle className="w-4 h-4 text-gray-400" />
              )}
              <span className="font-mono text-gray-700">{log.ip_address}</span>
              <span className="px-1.5 py-0.5 rounded bg-gray-100 text-[10px] font-bold text-gray-600">
                {log.country}
              </span>
              <span className="capitalize font-semibold text-gray-800">
                {log.decision}
              </span>
            </div>
            <div className="text-gray-400 flex items-center gap-1 text-[11px]">
              <Clock className="w-3 h-3" />
              {new Date(log.created_at).toLocaleTimeString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
