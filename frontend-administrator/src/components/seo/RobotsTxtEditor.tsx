import React from 'react';
import { Bot, Info } from 'lucide-react';

interface RobotsTxtEditorProps {
  robotsTxt: string;
  onChange: (val: string) => void;
}

export const RobotsTxtEditor: React.FC<RobotsTxtEditorProps> = ({ robotsTxt, onChange }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Bot className="w-4 h-4 text-indigo-600" /> Live robots.txt Editor
        </h3>
        <span className="text-xs font-mono text-gray-400">/robots.txt</span>
      </div>

      <div className="p-3 bg-blue-50 text-blue-800 rounded-lg text-xs flex items-start gap-2">
        <Info className="w-4 h-4 shrink-0 mt-0.5" />
        <span>Controls which search engine bots are allowed to crawl your client area, storefront, and documentation.</span>
      </div>

      <textarea
        rows={10}
        value={robotsTxt}
        onChange={(e) => onChange(e.target.value)}
        className="w-full p-4 font-mono text-xs bg-gray-900 text-emerald-400 rounded-xl border border-gray-800 outline-none focus:ring-2 focus:ring-indigo-500"
      />
    </div>
  );
};
