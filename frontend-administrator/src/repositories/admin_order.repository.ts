import { request } from '../lib/api/client';
import type { Order } from '@/types/api';

export interface IAdminOrderRepository {
  listOrders(): Promise<Order[]>;
  getOrder(id: number): Promise<Order>;
  activateOrder(id: number): Promise<any>;
  suspendOrder(id: number, reason: string): Promise<any>;
  unsuspendOrder(id: number): Promise<any>;
  cancelOrder(id: number, reason?: string): Promise<any>;
}

export class AdminOrderRepository implements IAdminOrderRepository {
  async listOrders(): Promise<Order[]> {
    return request<Order[]>('/admin/orders');
  }

  async getOrder(id: number): Promise<Order> {
    return request<Order>(`/admin/orders/${id}`);
  }

  async activateOrder(id: number): Promise<any> {
    return request(`/admin/orders/${id}/activate`, {
      method: 'POST',
    });
  }

  async suspendOrder(id: number, reason: string): Promise<any> {
    return request(`/admin/orders/${id}/suspend`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
  }

  async unsuspendOrder(id: number): Promise<any> {
    return request(`/admin/orders/${id}/unsuspend`, {
      method: 'POST',
    });
  }

  async cancelOrder(id: number, reason = 'Canceled by admin'): Promise<any> {
    return request(`/admin/orders/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
  }
}

export const adminOrderRepository = new AdminOrderRepository();
