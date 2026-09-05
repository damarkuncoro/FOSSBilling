import React from 'react';
import { Sliders, Code2 } from 'lucide-react';
import type { WidgetConfig } from '../../types/embedWidgets';

interface WidgetGeneratorCardProps {
  config: WidgetConfig;
  onChange: (key: keyof WidgetConfig, val: any) => void;
}

export const WidgetGeneratorCard: React.FC<WidgetGeneratorCardProps> = ({ config, onChange }) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6 space-y-4">
      <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2 pb-3 border-b border-gray-100">
        <Sliders className="w-4 h-4 text-indigo-600" /> Widget Customizer
      </h3>

      <div>
        <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Target Product</label>
        <select
          value={config.product_id}
          onChange={(e) => {
            const id = Number(e.target.value);
            const title = id === 1 ? 'Cloud NVMe Pro Hosting' : id === 2 ? 'Dedicated Baremetal' : 'cPanel Starter';
            onChange('product_id', id);
            onChange('product_title', title);
          }}
          className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:border-indigo-500"
        >
          <option value={1}>Cloud NVMe Pro Hosting ($9.99/mo)</option>
          <option value={2}>Dedicated Baremetal Server ($129.00/mo)</option>
          <option value={3}>cPanel Starter Shared ($3.99/mo)</option>
        </select>
      </div>

      <div>
        <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Button Call-to-Action Text</label>
        <input
          type="text"
          value={config.button_text}
          onChange={(e) => onChange('button_text', e.target.value)}
          className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm outline-none focus:border-indigo-500"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Button Color</label>
          <div className="flex items-center gap-2">
            <input
              type="color"
              value={config.button_color}
              onChange={(e) => onChange('button_color', e.target.value)}
              className="w-9 h-9 rounded-lg border border-gray-200 p-0.5 cursor-pointer"
            />
            <input
              type="text"
              value={config.button_color}
              onChange={(e) => onChange('button_color', e.target.value)}
              className="flex-1 px-2.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none"
            />
          </div>
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Text Color</label>
          <div className="flex items-center gap-2">
            <input
              type="color"
              value={config.text_color}
              onChange={(e) => onChange('text_color', e.target.value)}
              className="w-9 h-9 rounded-lg border border-gray-200 p-0.5 cursor-pointer"
            />
            <input
              type="text"
              value={config.text_color}
              onChange={(e) => onChange('text_color', e.target.value)}
              className="flex-1 px-2.5 py-2 font-mono text-xs border border-gray-200 rounded-lg outline-none"
            />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Layout Type</label>
          <select
            value={config.layout}
            onChange={(e) => onChange('layout', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs bg-white outline-none focus:border-indigo-500"
          >
            <option value="button">Single Action Button</option>
            <option value="iframe_checkout">Embedded Checkout iFrame</option>
          </select>
        </div>

        <div>
          <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Action Type</label>
          <select
            value={config.action_type}
            onChange={(e) => onChange('action_type', e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs bg-white outline-none focus:border-indigo-500"
          >
            <option value="popup">Modal Checkout Popup</option>
            <option value="redirect">Direct Redirect to Cart</option>
          </select>
        </div>
      </div>
    </div>
  );
};
