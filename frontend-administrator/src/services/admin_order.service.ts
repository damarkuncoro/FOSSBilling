import { AdminOrderRepository, adminOrderRepository, IAdminOrderRepository } from '../repositories/admin_order.repository';
import type { Order } from '@/types/api';

export class AdminOrderService {
  constructor(private repo: IAdminOrderRepository = adminOrderRepository) {}

  async listOrders(): Promise<Order[]> {
    return this.repo.listOrders();
  }

  async getOrderDetail(id: number): Promise<Order> {
    if (!id || id <= 0) {
      throw new Error('Valid order ID is required');
    }
    return this.repo.getOrder(id);
  }

  async activateOrder(id: number): Promise<any> {
    return this.repo.activateOrder(id);
  }

  async suspendOrder(id: number, reason: string): Promise<any> {
    if (!reason.trim()) {
      throw new Error('Suspension reason is required');
    }
    return this.repo.suspendOrder(id, reason.trim());
  }

  async unsuspendOrder(id: number): Promise<any> {
    return this.repo.unsuspendOrder(id);
  }

  async cancelOrder(id: number, reason?: string): Promise<any> {
    return this.repo.cancelOrder(id, reason);
  }

  filterByStatus(orders: Order[], status?: string): Order[] {
    if (!status || status === 'all') return orders;
    return orders.filter((o) => o.status.toLowerCase() === status.toLowerCase());
  }
}

export const adminOrderService = new AdminOrderService();
