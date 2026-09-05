import React from 'react';
import { Code2, Copy, CheckCircle2, Eye } from 'lucide-react';
import type { WidgetConfig } from '../../types/embedWidgets';

interface WidgetCodePreviewProps {
  config: WidgetConfig;
  embedCode: string;
  copied: boolean;
  onCopy: () => void;
}

export const WidgetCodePreview: React.FC<WidgetCodePreviewProps> = ({
  config,
  embedCode,
  copied,
  onCopy,
}) => {
  return (
    <div className="space-y-6">
      {/* Live Visual Preview */}
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
        <div className="flex items-center gap-2 text-xs font-semibold text-gray-500 uppercase tracking-wider pb-3 border-b border-gray-100">
          <Eye className="w-4 h-4 text-indigo-600" /> Interactive Button Preview
        </div>

        <div className="p-8 bg-gray-50/80 border border-dashed border-gray-200 rounded-xl flex items-center justify-center">
          <button
            style={{
              backgroundColor: config.button_color,
              color: config.text_color,
              borderRadius: `${config.border_radius}px`,
              padding: '12px 24px',
              fontWeight: 600,
              fontSize: '14px',
              border: 'none',
              cursor: 'pointer',
              boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
            }}
          >
            {config.button_text}
          </button>
        </div>
      </div>

      {/* Embed Code Output */}
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-xs font-semibold text-gray-700 uppercase tracking-wider flex items-center gap-1.5">
            <Code2 className="w-4 h-4 text-indigo-600" /> Embed HTML / JS Snippet
          </h4>
          <button
            onClick={onCopy}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 rounded-lg text-xs font-medium transition-colors"
          >
            {copied ? <CheckCircle2 className="w-3.5 h-3.5 text-emerald-600" /> : <Copy className="w-3.5 h-3.5" />}
            {copied ? 'Copied to Clipboard!' : 'Copy Snippet'}
          </button>
        </div>

        <pre className="p-4 bg-gray-900 text-emerald-400 font-mono text-xs rounded-xl overflow-x-auto leading-relaxed border border-gray-800">
          {embedCode}
        </pre>
      </div>
    </div>
  );
};
