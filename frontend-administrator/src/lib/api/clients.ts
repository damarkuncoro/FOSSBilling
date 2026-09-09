import { request } from './client';

export const clientsApi = {
  getClients: () => request<any[]>('/admin/clients'),
  getClient: (id: number) => request<any>(`/admin/clients/${id}`),
  createClient: (data: any) => request<any>('/admin/clients', { method: 'POST', body: JSON.stringify(data) }),
  updateClient: (id: number, data: any) => request<any>(`/admin/clients/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteClient: (id: number) => request<any>(`/admin/clients/${id}`, { method: 'DELETE' }),
  impersonateClient: (id: number) => request<{ token: string }>(`/admin/clients/${id}/impersonate`, { method: 'POST' }),
  getOrders: () => request<any[]>('/admin/orders'),
  activateOrder: (id: number) => request<any>(`/admin/orders/${id}/activate`, { method: 'POST' }),
  suspendOrder: (id: number, reason = 'Administrative suspension') =>
    request<any>(`/admin/orders/${id}/suspend`, { method: 'POST', body: JSON.stringify({ reason }) }),
  unsuspendOrder: (id: number) => request<any>(`/admin/orders/${id}/unsuspend`, { method: 'POST' }),
};
