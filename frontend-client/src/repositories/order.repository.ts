import { request } from '../lib/api/client';
import type { Order } from '@/types/api';

export class OrderRepository {
  list = (l = 100, o = 0) => request<Order[]>(`/client/orders?limit=${l}&offset=${o}`);
  get = (id: number) => request<Order>(`/client/orders/${id}`);
  sync = (id: number) => request(`/client/orders/${id}/sync`, { method: 'POST' });
  changePassword = (id: number, p: string) => request(`/client/orders/${id}/change-password`, { method: 'POST', body: JSON.stringify({ password: p }) });
}

export const orderRepository = new OrderRepository();
