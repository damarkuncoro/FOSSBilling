import { useState } from 'react';
import type { WidgetConfig } from '../types/embedWidgets';

const initialConfig: WidgetConfig = {
  product_id: 1,
  product_title: 'Cloud NVMe Pro Hosting',
  button_text: 'Order Now - $9.99/mo',
  button_color: '#4f46e5',
  text_color: '#ffffff',
  border_radius: 8,
  layout: 'button',
  action_type: 'popup',
  show_price: true,
  price_display: '$9.99 / month',
};

export function useEmbedWidgets() {
  const [config, setConfig] = useState<WidgetConfig>(initialConfig);
  const [copied, setCopied] = useState(false);

  const updateConfig = (key: keyof WidgetConfig, val: any) => {
    setConfig((prev) => ({ ...prev, [key]: val }));
  };

  const generateEmbedCode = () => {
    if (config.layout === 'iframe_checkout') {
      return `<iframe \n  src="https://billing.myhosting.com/cart/embed?product_id=${config.product_id}" \n  width="100%" \n  height="600" \n  frameborder="0"\n></iframe>`;
    }

    return `<script src="https://billing.myhosting.com/assets/widgets/fossbilling-button.js"></script>\n<button \n  data-fossbilling-btn \n  data-product-id="${config.product_id}" \n  data-action="${config.action_type}" \n  style="background-color: ${config.button_color}; color: ${config.text_color}; border-radius: ${config.border_radius}px; padding: 10px 20px; font-weight: 600; border: none; cursor: pointer;"\n>\n  ${config.button_text}\n</button>`;
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(generateEmbedCode());
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return {
    config,
    copied,
    updateConfig,
    generateEmbedCode,
    copyToClipboard,
  };
}
