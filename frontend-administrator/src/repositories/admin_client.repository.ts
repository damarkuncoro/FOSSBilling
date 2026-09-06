import { request } from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export interface IAdminClientRepository {
  listClients(): Promise<ClientProfile[]>;
  getClient(id: number): Promise<ClientProfile>;
  createClient(dto: Partial<ClientProfile>): Promise<ClientProfile>;
  updateClient(id: number, dto: Partial<ClientProfile>): Promise<ClientProfile>;
  deleteClient(id: number): Promise<any>;
}

export class AdminClientRepository implements IAdminClientRepository {
  async listClients(): Promise<ClientProfile[]> {
    return request<ClientProfile[]>('/admin/clients');
  }

  async getClient(id: number): Promise<ClientProfile> {
    return request<ClientProfile>(`/admin/clients/${id}`);
  }

  async createClient(dto: Partial<ClientProfile>): Promise<ClientProfile> {
    return request<ClientProfile>('/admin/clients', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateClient(id: number, dto: Partial<ClientProfile>): Promise<ClientProfile> {
    return request<ClientProfile>(`/admin/clients/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteClient(id: number): Promise<any> {
    return request<any>(`/admin/clients/${id}`, {
      method: 'DELETE',
    });
  }
}

export const adminClientRepository = new AdminClientRepository();
