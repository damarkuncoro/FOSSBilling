import { adminSupportRepository as repo } from '../repositories/admin_support.repository';
import type { SupportTicket } from '@/types/api';

export class AdminSupportService {
  listTickets = () => repo.list();
  getTicketDetail = async (id: number): Promise<SupportTicket> => {
    const { ticket, messages } = await repo.get(id);
    return { ...ticket, messages };
  };
  replyTicket = (id: number, msg: string) => repo.reply(id, msg);
  closeTicket = (id: number) => repo.close(id);

  filterByStatus(ts: SupportTicket[], st?: string) {
    if (!st || st === 'all') return ts;
    return ts.filter(t => t.status.toLowerCase() === st.toLowerCase());
  }
}

export const adminSupportService = new AdminSupportService();
