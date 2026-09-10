import { adminOrderRepository, AdminOrderRepository } from '../repositories/admin_order.repository';
import type { Order } from '@/types/api';

export class AdminOrderService {
  constructor(private repo: AdminOrderRepository = adminOrderRepository) {}

  listOrders = () => this.repo.list();
  getOrderDetail = (id: number) => this.repo.get(id);
  activateOrder = (id: number) => this.repo.activate(id);
  suspendOrder = (id: number, reason: string) => {
    if (!reason || !reason.trim()) throw new Error('Suspension reason is required');
    return this.repo.suspend(id, reason);
  };
  unsuspendOrder = (id: number) => this.repo.unsuspend(id);
  cancelOrder = (id: number, reason?: string) => this.repo.cancel(id, reason);

  filterByStatus(os: Order[], st?: string) {
    if (!st || st === 'all') return os;
    return os.filter(o => o.status.toLowerCase() === st.toLowerCase());
  }
}

export const adminOrderService = new AdminOrderService();
