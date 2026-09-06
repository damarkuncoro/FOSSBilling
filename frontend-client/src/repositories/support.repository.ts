import { request } from '../lib/api/client';
import type { SupportTicket } from '@/types/api';

export interface ISupportRepository {
  listTickets(): Promise<SupportTicket[]>;
  getTicket(id: number): Promise<SupportTicket>;
  openTicket(dto: { subject: string; content: string; priority?: string }): Promise<SupportTicket>;
  replyTicket(id: number, content: string): Promise<any>;
  closeTicket(id: number): Promise<any>;
}

export class SupportRepository implements ISupportRepository {
  async listTickets(): Promise<SupportTicket[]> {
    return request<SupportTicket[]>('/client/support/tickets');
  }

  async getTicket(id: number): Promise<SupportTicket> {
    return request<SupportTicket>(`/client/support/tickets/${id}`);
  }

  async openTicket(dto: { subject: string; content: string; priority?: string }): Promise<SupportTicket> {
    return request<SupportTicket>('/client/support/tickets', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async replyTicket(id: number, content: string): Promise<any> {
    return request(`/client/support/tickets/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  }

  async closeTicket(id: number): Promise<any> {
    return request(`/client/support/tickets/${id}/close`, {
      method: 'POST',
    });
  }
}

export const supportRepository = new SupportRepository();
