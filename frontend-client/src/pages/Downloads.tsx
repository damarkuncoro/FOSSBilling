import React from 'react';
import { Download, FileCode, ArrowDownToLine, FolderArchive, ShoppingCart } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useDownloads } from '../hooks/useDownloads';

export const Downloads: React.FC = () => {
  const { downloads, loading, downloadingId, triggerDownload } = useDownloads();

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

      {loading ? (
        <div className="bg-white rounded-2xl border border-gray-100 p-12 text-center text-gray-400">
          <div className="animate-spin w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full mx-auto mb-3" />
          <p className="text-sm font-medium">Loading digital assets...</p>
        </div>
      ) : downloads.length === 0 ? (
        <div className="bg-white rounded-2xl border border-dashed border-gray-200 p-12 text-center space-y-4">
          <div className="w-14 h-14 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center mx-auto shadow-sm">
            <FolderArchive className="w-7 h-7" />
          </div>
          <div className="space-y-1">
            <h3 className="text-base font-bold text-gray-900">No Downloadable Assets</h3>
            <p className="text-xs text-gray-500 max-w-md mx-auto">
              You do not have any active downloadable files or assets linked to your active services.
            </p>
          </div>
          <div className="pt-2">
            <Link
              to="/order"
              className="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold rounded-xl shadow-sm transition-all"
            >
              <ShoppingCart className="w-4 h-4" /> Browse Digital Products
            </Link>
          </div>
        </div>
      ) : (
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
      )}
    </div>
  );
};

