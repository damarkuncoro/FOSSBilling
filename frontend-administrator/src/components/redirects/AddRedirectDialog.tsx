import React, { useState } from 'react';
import { X, ArrowRightLeft } from 'lucide-react';

interface AddRedirectDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { source_path: string; target_url: string; status_code: 301 | 302 }) => void;
}

export const AddRedirectDialog: React.FC<AddRedirectDialogProps> = ({
  isOpen,
  onClose,
  onSubmit,
}) => {
  const [sourcePath, setSourcePath] = useState('');
  const [targetUrl, setTargetUrl] = useState('');
  const [statusCode, setStatusCode] = useState<301 | 302>(301);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (sourcePath.trim() && targetUrl.trim()) {
      onSubmit({ source_path: sourcePath, target_url: targetUrl, status_code: statusCode });
      setSourcePath('');
      setTargetUrl('');
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-indigo-600 font-semibold">
            <ArrowRightLeft className="w-5 h-5" />
            <h3 className="text-gray-900 font-bold">Add URL Redirect</h3>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Source Path (From)</label>
            <input
              type="text"
              required
              placeholder="/old-page or /promo"
              value={sourcePath}
              onChange={(e) => setSourcePath(e.target.value)}
              className="w-full px-3.5 py-2 font-mono border border-gray-200 rounded-lg text-xs focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Target Destination (To)</label>
            <input
              type="text"
              required
              placeholder="/new-page or https://external-url.com"
              value={targetUrl}
              onChange={(e) => setTargetUrl(e.target.value)}
              className="w-full px-3.5 py-2 font-mono border border-gray-200 rounded-lg text-xs focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">HTTP Redirect Type</label>
            <select
              value={statusCode}
              onChange={(e) => setStatusCode(Number(e.target.value) as 301 | 302)}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs bg-white outline-none focus:border-indigo-500"
            >
              <option value={301}>301 Moved Permanently (SEO Recommended)</option>
              <option value={302}>302 Found / Temporary Redirect</option>
            </select>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
            <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg">
              Cancel
            </button>
            <button type="submit" className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">
              Save Redirect Rule
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
