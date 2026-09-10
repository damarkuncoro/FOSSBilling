import { request } from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export class AdminClientRepository {
  list = () => request<ClientProfile[]>('/admin/clients');
  get = (id: number) => request<ClientProfile>(`/admin/clients/${id}`);
  create = (d: any) => request<ClientProfile>('/admin/clients', { method: 'POST', body: JSON.stringify(d) });
  update = (id: number, d: any) => request<ClientProfile>(`/admin/clients/${id}`, { method: 'PUT', body: JSON.stringify(d) });
  delete = (id: number) => request(`/admin/clients/${id}`, { method: 'DELETE' });
  impersonate = (id: number) => request<{ token: string }>(`/admin/clients/${id}/impersonate`, { method: 'POST' });
}

export const adminClientRepository = new AdminClientRepository();
