import { OrderRepository, orderRepository, IOrderRepository } from '../repositories/order.repository';
import type { Order } from '@/types/api';

export class OrderService {
  constructor(private repo: IOrderRepository = orderRepository) {}

  async listClientOrders(limit = 100, offset = 0): Promise<Order[]> {
    return this.repo.listOrders(limit, offset);
  }

  async getOrderDetail(id: number): Promise<Order> {
    if (!id || id <= 0) {
      throw new Error('Valid order ID is required');
    }
    return this.repo.getOrder(id);
  }

  filterByStatus(orders: Order[], status?: string): Order[] {
    if (!status || status === 'all') return orders;
    return orders.filter((o) => o.status.toLowerCase() === status.toLowerCase());
  }
}

export const orderService = new OrderService();
