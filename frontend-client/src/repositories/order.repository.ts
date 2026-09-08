import { request } from '../lib/api/client';
import type { Order } from '@/types/api';

export interface IOrderRepository {
  listOrders(limit?: number, offset?: number): Promise<Order[]>;
  getOrder(id: number): Promise<Order>;
  syncStatus(id: number): Promise<any>;
  changePassword(id: number, password: string): Promise<any>;
}

export class OrderRepository implements IOrderRepository {
  async listOrders(limit = 100, offset = 0): Promise<Order[]> {
    return request<Order[]>(`/client/orders?limit=${limit}&offset=${offset}`);
  }

  async getOrder(id: number): Promise<Order> {
    return request<Order>(`/client/orders/${id}`);
  }

  async syncStatus(id: number): Promise<any> {
    return request<any>(`/client/orders/${id}/sync`, { method: 'POST' });
  }

  async changePassword(id: number, password: string): Promise<any> {
    return request<any>(`/client/orders/${id}/change-password`, {
      method: 'POST',
      body: JSON.stringify({ password }),
    });
  }
}

export const orderRepository = new OrderRepository();
