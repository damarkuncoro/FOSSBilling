import { SupportRepository, supportRepository, ISupportRepository } from '../repositories/support.repository';
import type { SupportTicket } from '@/types/api';

export class SupportService {
  constructor(private repo: ISupportRepository = supportRepository) {}

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

  async openTicket(subject: string, message: string, priority = 'medium'): Promise<SupportTicket> {
    if (!subject.trim()) {
      throw new Error('Subject is required');
    }
    if (!message.trim()) {
      throw new Error('Message content is required');
    }
    return this.repo.openTicket({
      subject: subject.trim(),
      message: message.trim(),
      priority,
    });
  }

  async replyTicket(id: number, content: string): Promise<any> {
    if (!content.trim()) {
      throw new Error('Reply message cannot be empty');
    }
    return this.repo.replyTicket(id, content.trim());
  }

  async closeTicket(id: number): Promise<any> {
    return this.repo.closeTicket(id);
  }
}

export const supportService = new SupportService();
