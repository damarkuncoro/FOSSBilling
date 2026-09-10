import { request } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export class AdminInvoiceRepository {
  list = () => request<Invoice[]>('/admin/invoices');
  get = (id: number) => request<Invoice>(`/admin/invoices/${id}`);
  create = (d: any) => request<Invoice>('/admin/invoices', { method: 'POST', body: JSON.stringify(d) });
  refund = (id: number) => request(`/admin/invoices/${id}/refund`, { method: 'POST' });
  delete = (id: number) => request(`/admin/invoices/${id}`, { method: 'DELETE' });
}

export const adminInvoiceRepository = new AdminInvoiceRepository();
