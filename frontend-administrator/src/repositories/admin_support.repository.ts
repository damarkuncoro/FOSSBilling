import { request } from '../lib/api/client';
import type { SupportTicket } from '@/types/api';

export interface IAdminSupportRepository {
  listTickets(): Promise<SupportTicket[]>;
  getTicket(id: number): Promise<{ ticket: SupportTicket; messages: any[] }>;
  replyTicket(id: number, content: string): Promise<any>;
  closeTicket(id: number): Promise<any>;
}

export class AdminSupportRepository implements IAdminSupportRepository {
  async listTickets(): Promise<SupportTicket[]> {
    return request<SupportTicket[]>('/admin/support/tickets');
  }

  async getTicket(id: number): Promise<{ ticket: SupportTicket; messages: any[] }> {
    return request<{ ticket: SupportTicket; messages: any[] }>(`/admin/support/tickets/${id}`);
  }

  async replyTicket(id: number, message: string): Promise<any> {
    return request(`/admin/support/tickets/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ message }),
    });
  }

  async closeTicket(id: number): Promise<any> {
    return request(`/admin/support/tickets/${id}/close`, {
      method: 'POST',
    });
  }
}

export const adminSupportRepository = new AdminSupportRepository();
