import React from 'react';
import { Download, FileCode, ArrowDownToLine } from 'lucide-react';
import { useDownloads } from '../hooks/useDownloads';

export const Downloads: React.FC = () => {
  const { downloads, downloadingId, triggerDownload } = useDownloads();

  return (
    <div className="max-w-4xl mx-auto px-4 py-8 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
          <Download className="w-7 h-7 text-indigo-600" /> Downloads & Digital Assets
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Download software releases, configuration files, and extensions included in your active subscriptions.
        </p>
      </div>

      <div className="space-y-4">
        {downloads.map((item) => (
          <div
            key={item.id}
            className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
          >
            <div className="flex items-start gap-3.5">
              <div className="p-3 bg-indigo-50 text-indigo-600 rounded-xl">
                <FileCode className="w-6 h-6" />
              </div>
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-bold text-base text-gray-900">{item.title}</span>
                  <span className="px-2 py-0.5 rounded font-mono text-[11px] bg-gray-100 text-gray-700">
                    v{item.version}
                  </span>
                </div>
                <p className="text-xs text-gray-500">{item.description}</p>
                <div className="flex items-center gap-3 text-xs text-gray-400 pt-1">
                  <span>Size: {item.file_size}</span>
                  <span>•</span>
                  <span>Category: {item.category}</span>
                </div>
              </div>
            </div>

            <button
              onClick={() => triggerDownload(item)}
              disabled={downloadingId === item.id}
              className="inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-xs rounded-xl shadow-sm transition-all shrink-0"
            >
              <ArrowDownToLine className={`w-4 h-4 ${downloadingId === item.id ? 'animate-bounce' : ''}`} />
              {downloadingId === item.id ? 'Starting Download...' : 'Download File'}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
