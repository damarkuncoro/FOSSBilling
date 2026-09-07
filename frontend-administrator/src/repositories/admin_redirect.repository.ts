import { apiClient } from '@/lib/api';

export interface AdminRedirectRepositoryInterface {
  listRedirects(): Promise<any[]>;
  createRedirect(data: any): Promise<any>;
  deleteRedirect(id: number): Promise<any>;
}

export class AdminRedirectRepository implements AdminRedirectRepositoryInterface {
  async listRedirects(): Promise<any[]> {
    const res = await apiClient.get<any>('/admin/redirects');
    return res.data?.data || res.data || [];
  }

  async createRedirect(data: any): Promise<any> {
    const res = await apiClient.post<any>('/admin/redirects', data);
    return res.data?.data || res.data;
  }

  async deleteRedirect(id: number): Promise<any> {
    const res = await apiClient.delete<any>(`/admin/redirects/${id}`);
    return res.data?.data || res.data;
  }
}
