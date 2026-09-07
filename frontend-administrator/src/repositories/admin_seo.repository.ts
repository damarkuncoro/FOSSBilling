import { apiClient } from '@/lib/api';

export interface AdminSeoRepositoryInterface {
  getSettings(): Promise<any>;
  updateSettings(settings: any): Promise<any>;
}

export class AdminSeoRepository implements AdminSeoRepositoryInterface {
  async getSettings(): Promise<any> {
    const res = await apiClient.get<any>('/admin/seo');
    return res.data?.data || res.data;
  }

  async updateSettings(settings: any): Promise<any> {
    const res = await apiClient.put<any>('/admin/seo', settings);
    return res.data?.data || res.data;
  }
}
