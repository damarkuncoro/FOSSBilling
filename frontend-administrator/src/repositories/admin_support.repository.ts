import { request } from '../lib/api/client';
import type { SupportTicket } from '@/types/api';

export class AdminSupportRepository {
  list = () => request<SupportTicket[]>('/admin/support/tickets');
  get = (id: number) => request<{ ticket: SupportTicket; messages: any[] }>(`/admin/support/tickets/${id}`);
  reply = (id: number, message: string) => request(`/admin/support/tickets/${id}/reply`, { method: 'POST', body: JSON.stringify({ message }) });
  close = (id: number) => request(`/admin/support/tickets/${id}/close`, { method: 'POST' });
}

export const adminSupportRepository = new AdminSupportRepository();
