import { apiClient } from '@/lib/api';

export interface AdminThemeRepositoryInterface {
  listThemes(): Promise<any[]>;
  getActiveTheme(): Promise<any>;
  setActiveTheme(id: string): Promise<any>;
  updateThemeSettings(id: string, settings: any): Promise<any>;
}

export class AdminThemeRepository implements AdminThemeRepositoryInterface {
  async listThemes(): Promise<any[]> {
    const res = await apiClient.get<any>('/admin/themes');
    return res.data?.data || res.data || [];
  }

  async getActiveTheme(): Promise<any> {
    const res = await apiClient.get<any>('/admin/themes/active');
    return res.data?.data || res.data;
  }

  async setActiveTheme(id: string): Promise<any> {
    const res = await apiClient.post<any>(`/admin/themes/${id}/select`);
    return res.data?.data || res.data;
  }

  async updateThemeSettings(id: string, settings: any): Promise<any> {
    const res = await apiClient.put<any>(`/admin/themes/${id}/settings`, settings);
    return res.data?.data || res.data;
  }
}
