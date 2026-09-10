import { request } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export class InvoiceRepository {
  list = (l = 100, o = 0) => request<Invoice[]>(`/client/invoices?limit=${l}&offset=${o}`);
  get = (id: number) => request<Invoice>(`/client/invoices/${id}`);
  payBalance = (id: number) => request(`/client/invoices/${id}/pay-balance`, { method: 'POST' });
  payGateway = (id: number, gateway: string) => request<{ redirect_url: string }>(`/client/invoices/${id}/pay-gateway`, { method: 'POST', body: JSON.stringify({ gateway }) });
  deposit = (amount: number, gateway = 'midtrans') => request<{ invoice_id: number; redirect_url?: string }>('/client/funds/deposit', { method: 'POST', body: JSON.stringify({ amount, gateway }) });
}

export const invoiceRepository = new InvoiceRepository();
