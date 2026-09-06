import { request } from '../lib/api/client';
import type { SupportTicket } from '@/types/api';

export interface IAdminSupportRepository {
  listTickets(): Promise<SupportTicket[]>;
  getTicket(id: number): Promise<SupportTicket>;
  replyTicket(id: number, content: string): Promise<any>;
  closeTicket(id: number): Promise<any>;
}

export class AdminSupportRepository implements IAdminSupportRepository {
  async listTickets(): Promise<SupportTicket[]> {
    return request<SupportTicket[]>('/admin/support/tickets');
  }

  async getTicket(id: number): Promise<SupportTicket> {
    return request<SupportTicket>(`/admin/support/tickets/${id}`);
  }

  async replyTicket(id: number, content: string): Promise<any> {
    return request(`/admin/support/tickets/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  }

  async closeTicket(id: number): Promise<any> {
    return request(`/admin/support/tickets/${id}/close`, {
      method: 'POST',
    });
  }
}

export const adminSupportRepository = new AdminSupportRepository();
