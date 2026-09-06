import { request } from '../lib/api/client';
import type { ServerItem } from '@/types/api';

export interface IAdminServerRepository {
  getServers(): Promise<ServerItem[]>;
  getServer(id: number): Promise<ServerItem>;
  createServer(dto: Partial<ServerItem>): Promise<ServerItem>;
  updateServer(id: number, dto: Partial<ServerItem>): Promise<ServerItem>;
  deleteServer(id: number): Promise<any>;
  testServerConnection(id: number): Promise<{ success: boolean; message: string }>;
}

export class AdminServerRepository implements IAdminServerRepository {
  async getServers(): Promise<ServerItem[]> {
    return request<ServerItem[]>('/admin/servers');
  }

  async getServer(id: number): Promise<ServerItem> {
    return request<ServerItem>(`/admin/servers/${id}`);
  }

  async createServer(dto: Partial<ServerItem>): Promise<ServerItem> {
    return request<ServerItem>('/admin/servers', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateServer(id: number, dto: Partial<ServerItem>): Promise<ServerItem> {
    return request<ServerItem>(`/admin/servers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteServer(id: number): Promise<any> {
    return request<any>(`/admin/servers/${id}`, {
      method: 'DELETE',
    });
  }

  async testServerConnection(id: number): Promise<{ success: boolean; message: string }> {
    return request<{ success: boolean; message: string }>(`/admin/servers/${id}/test`, {
      method: 'POST',
    });
  }
}

export const adminServerRepository = new AdminServerRepository();
