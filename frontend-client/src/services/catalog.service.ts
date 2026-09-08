import { request } from '@/lib/api/client';
import { HostingPlan } from '@/types/api';

export const catalogService = {
  async listProducts(): Promise<HostingPlan[]> {
    try {
      const res = await request<any[]>('/guest/products');
      return (res || []).map((p: any) => ({
        id: p.id,
        title: p.name || p.title,
        description: p.description,
        price: p.price_monthly || 0,
        period: '1M',
        type: p.type,
        form_id: p.form_id,
        features: p.features || [],
      }));
    } catch {
      return [];
    }
  },
};
