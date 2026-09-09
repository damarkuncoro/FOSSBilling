import { AdminSupportRepository, adminSupportRepository, IAdminSupportRepository } from '../repositories/admin_support.repository';
import type { SupportTicket } from '@/types/api';

export class AdminSupportService {
  constructor(private repo: IAdminSupportRepository = adminSupportRepository) {}

  async listTickets(): Promise<SupportTicket[]> {
    return this.repo.listTickets();
  }

  async getTicketDetail(id: number): Promise<SupportTicket> {
    if (!id || id <= 0) {
      throw new Error('Valid ticket ID is required');
    }
    const { ticket, messages } = await this.repo.getTicket(id);
    return {
      ...ticket,
      messages,
    };
  }

  async replyTicket(id: number, message: string): Promise<any> {
    if (!message.trim()) {
      throw new Error('Reply message cannot be empty');
    }
    return this.repo.replyTicket(id, message.trim());
  }

  async closeTicket(id: number): Promise<any> {
    return this.repo.closeTicket(id);
  }

  filterByStatus(tickets: SupportTicket[], status?: string): SupportTicket[] {
    if (!status || status === 'all') return tickets;
    return tickets.filter((t) => t.status.toLowerCase() === status.toLowerCase());
  }
}

export const adminSupportService = new AdminSupportService();
