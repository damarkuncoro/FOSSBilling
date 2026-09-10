import { request } from '../lib/api/client';
import type { SupportTicket } from '@/types/api';

export class SupportRepository {
  list = () => request<SupportTicket[]>('/client/support/tickets');
  get = (id: number) => request<{ ticket: SupportTicket; messages: any[] }>(`/client/support/tickets/${id}`);
  open = (d: any) => request<SupportTicket>('/client/support/tickets', { method: 'POST', body: JSON.stringify({ ...d, helpdesk_id: d.helpdesk_id || 1 }) });
  reply = (id: number, message: string) => request(`/client/support/tickets/${id}/reply`, { method: 'POST', body: JSON.stringify({ message }) });
  close = (id: number) => request(`/client/support/tickets/${id}/close`, { method: 'POST' });
}

export const supportRepository = new SupportRepository();
