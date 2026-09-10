import { invoiceRepository as repo } from '../repositories/invoice.repository';
import { API_BASE, getStoredClientToken } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export class InvoiceService {
  listClientInvoices = (l = 100, o = 0) => repo.list(l, o);
  getInvoiceDetail = (id: number) => repo.get(id);
  payWithBalance = (id: number) => repo.payBalance(id);
  payWithGateway = (id: number, g = 'midtrans') => repo.payGateway(id, g);
  depositFunds = (a: number, g = 'midtrans') => repo.deposit(a, g);
  getPdfDownloadUrl = (id: number) => `${API_BASE}/client/invoices/${id}/pdf?token=${getStoredClientToken()}`;
  filterByStatus = (is: Invoice[], st?: string) => (!st || st === 'all') ? is : is.filter(i => i.status.toLowerCase() === st.toLowerCase());
}

export const invoiceService = new InvoiceService();
