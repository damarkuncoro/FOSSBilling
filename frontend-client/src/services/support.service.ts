import { supportRepository as repo } from '../repositories/support.repository';

export class SupportService {
  listTickets = () => repo.list();
  getTicketDetail = async (id: number) => { const { ticket, messages } = await repo.get(id); return { ...ticket, messages }; };

  openTicket = (sub: string, msg: string, pri = 'medium') => {
    if (!sub.trim()) throw new Error('Subject is required');
    if (!msg.trim()) throw new Error('Message content is required');
    return repo.open({ subject: sub, message: msg, priority: pri });
  };

  replyTicket = (id: number, msg: string) => repo.reply(id, msg);
  closeTicket = (id: number) => repo.close(id);
}

export const supportService = new SupportService();
