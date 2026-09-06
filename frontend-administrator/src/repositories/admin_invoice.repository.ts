import { request } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export interface IAdminInvoiceRepository {
  listInvoices(): Promise<Invoice[]>;
  getInvoice(id: number): Promise<Invoice>;
  createInvoice(dto: {
    client_id: number;
    items: Array<{ title: string; price: number; quantity?: number }>;
  }): Promise<Invoice>;
}

export class AdminInvoiceRepository implements IAdminInvoiceRepository {
  async listInvoices(): Promise<Invoice[]> {
    return request<Invoice[]>('/admin/invoices');
  }

  async getInvoice(id: number): Promise<Invoice> {
    return request<Invoice>(`/admin/invoices/${id}`);
  }

  async createInvoice(dto: {
    client_id: number;
    items: Array<{ title: string; price: number; quantity?: number }>;
  }): Promise<Invoice> {
    return request<Invoice>('/admin/invoices', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }
}

export const adminInvoiceRepository = new AdminInvoiceRepository();
