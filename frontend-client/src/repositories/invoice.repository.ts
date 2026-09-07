import { request } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export interface IInvoiceRepository {
  listInvoices(limit?: number, offset?: number): Promise<Invoice[]>;
  getInvoice(id: number): Promise<Invoice>;
  payWithBalance(id: number): Promise<any>;
  payWithGateway(id: number, gateway: string): Promise<{ redirect_url: string }>;
  depositFunds(amount: number, gateway?: string): Promise<{ invoice_id: number; redirect_url?: string }>;
}

export class InvoiceRepository implements IInvoiceRepository {
  async listInvoices(limit = 100, offset = 0): Promise<Invoice[]> {
    return request<Invoice[]>(`/client/invoices?limit=${limit}&offset=${offset}`);
  }

  async getInvoice(id: number): Promise<Invoice> {
    return request<Invoice>(`/client/invoices/${id}`);
  }

  async payWithBalance(id: number): Promise<any> {
    return request(`/client/invoices/${id}/pay-balance`, {
      method: 'POST',
    });
  }

  async payWithGateway(id: number, gateway: string): Promise<{ redirect_url: string }> {
    return request<{ redirect_url: string }>(`/client/invoices/${id}/pay-gateway`, {
      method: 'POST',
      body: JSON.stringify({ gateway }),
    });
  }

  async depositFunds(amount: number, gateway = 'midtrans'): Promise<{ invoice_id: number; redirect_url?: string }> {
    return request<{ invoice_id: number; redirect_url?: string }>('/client/funds/deposit', {
      method: 'POST',
      body: JSON.stringify({ amount, gateway }),
    });
  }
}

export const invoiceRepository = new InvoiceRepository();
