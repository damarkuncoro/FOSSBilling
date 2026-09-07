import { apiClient } from '@/lib/api';

export interface AdminExtensionRepositoryInterface {
  listExtensions(): Promise<any[]>;
  listMarketplace(): Promise<any[]>;
  activateExtension(id: string): Promise<any>;
  deactivateExtension(id: string): Promise<any>;
  installExtension(id: string): Promise<any>;
  uninstallExtension(id: string): Promise<any>;
  getConfig(id: string): Promise<any>;
  updateConfig(id: string, config: any): Promise<any>;
}

export class AdminExtensionRepository implements AdminExtensionRepositoryInterface {
  async listExtensions(): Promise<any[]> {
    const res = await apiClient.get<any>('/admin/extensions');
    return res.data?.data || res.data || [];
  }

  async listMarketplace(): Promise<any[]> {
    const res = await apiClient.get<any>('/admin/extensions/marketplace');
    return res.data?.data || res.data || [];
  }

  async activateExtension(id: string): Promise<any> {
    const res = await apiClient.post<any>(`/admin/extensions/${id}/activate`);
    return res.data?.data || res.data;
  }

  async deactivateExtension(id: string): Promise<any> {
    const res = await apiClient.post<any>(`/admin/extensions/${id}/deactivate`);
    return res.data?.data || res.data;
  }

  async installExtension(id: string): Promise<any> {
    const res = await apiClient.post<any>(`/admin/extensions/${id}/install`);
    return res.data?.data || res.data;
  }

  async uninstallExtension(id: string): Promise<any> {
    const res = await apiClient.post<any>(`/admin/extensions/${id}/uninstall`);
    return res.data?.data || res.data;
  }

  async getConfig(id: string): Promise<any> {
    const res = await apiClient.get<any>(`/admin/extensions/${id}/config`);
    return res.data?.data || res.data;
  }

  async updateConfig(id: string, config: any): Promise<any> {
    const res = await apiClient.put<any>(`/admin/extensions/${id}/config`, config);
    return res.data?.data || res.data;
  }
}
