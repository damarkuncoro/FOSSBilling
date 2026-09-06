import { request } from '../lib/api/client';
import type { Order } from '@/types/api';

export interface IOrderRepository {
  listOrders(limit?: number, offset?: number): Promise<Order[]>;
  getOrder(id: number): Promise<Order>;
}

export class OrderRepository implements IOrderRepository {
  async listOrders(limit = 100, offset = 0): Promise<Order[]> {
    return request<Order[]>(`/client/orders?limit=${limit}&offset=${offset}`);
  }

  async getOrder(id: number): Promise<Order> {
    return request<Order>(`/client/orders/${id}`);
  }
}

export const orderRepository = new OrderRepository();
