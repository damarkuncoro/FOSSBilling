import { request } from '../lib/api/client';
import type { Order } from '@/types/api';

export class AdminOrderRepository {
  list = () => request<Order[]>('/admin/orders');
  get = (id: number) => request<Order>(`/admin/orders/${id}`);
  activate = (id: number) => request(`/admin/orders/${id}/activate`, { method: 'POST' });
  suspend = (id: number, reason: string) => request(`/admin/orders/${id}/suspend`, { method: 'POST', body: JSON.stringify({ reason }) });
  unsuspend = (id: number) => request(`/admin/orders/${id}/unsuspend`, { method: 'POST' });
  cancel = (id: number, reason = 'Canceled') => request(`/admin/orders/${id}/cancel`, { method: 'POST', body: JSON.stringify({ reason }) });
}

export const adminOrderRepository = new AdminOrderRepository();
