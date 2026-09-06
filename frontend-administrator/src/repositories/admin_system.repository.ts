import { request } from '../lib/api/client';
import type { SystemInfo, SystemStatusInfo } from '@/types/api';

export interface IAdminSystemRepository {
  getSystemInfo(): Promise<SystemInfo>;
  getSystemStatus(): Promise<SystemStatusInfo>;
  triggerCron(): Promise<{ success: boolean; message: string }>;
  clearCache(): Promise<{ success: boolean; message: string }>;
  listCurrencies(): Promise<any[]>;
  updateCurrency(code: string, rate: number): Promise<any>;
  listNews(): Promise<any[]>;
  createNews(dto: { title: string; content: string; slug?: string }): Promise<any>;
  deleteNews(id: number): Promise<any>;
}

export class AdminSystemRepository implements IAdminSystemRepository {
  async getSystemInfo(): Promise<SystemInfo> {
    return request<SystemInfo>('/admin/system/info');
  }

  async getSystemStatus(): Promise<SystemStatusInfo> {
    return request<SystemStatusInfo>('/admin/system/status');
  }

  async triggerCron(): Promise<{ success: boolean; message: string }> {
    return request<{ success: boolean; message: string }>('/admin/system/cron', {
      method: 'POST',
    });
  }

  async clearCache(): Promise<{ success: boolean; message: string }> {
    return request<{ success: boolean; message: string }>('/admin/system/clear-cache', {
      method: 'POST',
    });
  }

  async listCurrencies(): Promise<any[]> {
    return request<any[]>('/admin/currencies');
  }

  async updateCurrency(code: string, rate: number): Promise<any> {
    return request(`/admin/currencies/${code}`, {
      method: 'PUT',
      body: JSON.stringify({ rate }),
    });
  }

  async listNews(): Promise<any[]> {
    return request<any[]>('/admin/news');
  }

  async createNews(dto: { title: string; content: string; slug?: string }): Promise<any> {
    return request('/admin/news', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async deleteNews(id: number): Promise<any> {
    return request(`/admin/news/${id}`, {
      method: 'DELETE',
    });
  }
}

export const adminSystemRepository = new AdminSystemRepository();
