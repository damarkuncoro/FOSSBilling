import { orderRepository as repo } from '../repositories/order.repository';
import type { Order } from '@/types/api';

export class OrderService {
  listClientOrders = (l = 100, o = 0) => repo.list(l, o);
  getOrderDetail = (id: number) => repo.get(id);
  syncServiceStatus = (id: number) => repo.sync(id);
  changeServicePassword = (id: number, p: string) => repo.changePassword(id, p);
  filterByStatus = (os: Order[], st?: string) => (!st || st === 'all') ? os : os.filter(o => o.status.toLowerCase() === st.toLowerCase());
}

export const orderService = new OrderService();
