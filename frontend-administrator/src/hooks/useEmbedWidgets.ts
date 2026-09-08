import { useState, useEffect } from 'react';
import type { WidgetConfig } from '../types/embedWidgets';
import { catalogApi } from '../lib/api/catalog';

const initialConfig: WidgetConfig = {
  product_id: 1,
  product_title: 'Select a Product',
  button_text: 'Order Now',
  button_color: '#4f46e5',
  text_color: '#ffffff',
  border_radius: 8,
  layout: 'button',
  action_type: 'popup',
  show_price: true,
  price_display: '',
};

export function useEmbedWidgets() {
  const [config, setConfig] = useState<WidgetConfig>(initialConfig);
  const [products, setProducts] = useState<any[]>([]);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    catalogApi.getProducts().then(res => {
      setProducts(res || []);
      if (res && res.length > 0) {
        const first = res[0];
        setConfig(prev => ({
          ...prev,
          product_id: first.id,
          product_title: first.title,
          price_display: first.price_monthly ? `$${first.price_monthly}/mo` : 'Free',
        }));
      }
    });
  }, []);

  const updateConfig = (key: keyof WidgetConfig, val: any) => {
    if (key === 'product_id') {
      const p = products.find(x => x.id === Number(val));
      if (p) {
        setConfig(prev => ({
          ...prev,
          product_id: p.id,
          product_title: p.title,
          price_display: p.price_monthly ? `$${p.price_monthly}/mo` : 'Free',
        }));
        return;
      }
    }
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
    products,
    copied,
    updateConfig,
    generateEmbedCode,
    copyToClipboard,
  };
}
