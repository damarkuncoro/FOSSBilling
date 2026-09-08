import { api } from '@/lib/api';
import { HostingPlan } from '@/types/api';

export const catalogService = {
  async listProducts(): Promise<HostingPlan[]> {
    // Standardize backend product to frontend HostingPlan
    const res = await api.get<{ data: any[] }>('/admin/products'); // Temporary using admin endpoint for now, or create guest one
    return res.data.data.map(p => ({
      id: p.id,
      title: p.name || p.title,
      description: p.description,
      price: p.price_monthly || 0,
      period: '1M',
      type: p.type,
      form_id: p.form_id,
      features: [],
    }));
  }
};
