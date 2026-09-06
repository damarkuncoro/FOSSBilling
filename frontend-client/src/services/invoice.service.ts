import { InvoiceRepository, invoiceRepository, IInvoiceRepository } from '../repositories/invoice.repository';
import { API_BASE, getStoredClientToken } from '../lib/api/client';
import type { Invoice } from '@/types/api';

export class InvoiceService {
  constructor(private repo: IInvoiceRepository = invoiceRepository) {}

  async listClientInvoices(limit = 100, offset = 0): Promise<Invoice[]> {
    return this.repo.listInvoices(limit, offset);
  }

  async getInvoiceDetail(id: number): Promise<Invoice> {
    if (!id || id <= 0) {
      throw new Error('Valid invoice ID is required');
    }
    return this.repo.getInvoice(id);
  }

  async payWithBalance(id: number): Promise<any> {
    return this.repo.payWithBalance(id);
  }

  async depositFunds(amount: number, gateway = 'midtrans'): Promise<{ invoice_id: number; redirect_url?: string }> {
    if (!amount || amount <= 0) {
      throw new Error('Deposit amount must be greater than zero');
    }
    return this.repo.depositFunds(amount, gateway);
  }

  getPdfDownloadUrl(id: number): string {
    const token = getStoredClientToken();
    return `${API_BASE}/client/invoices/${id}/pdf${token ? `?token=${encodeURIComponent(token)}` : ''}`;
  }

  filterByStatus(invoices: Invoice[], status?: string): Invoice[] {
    if (!status || status === 'all') return invoices;
    return invoices.filter((inv) => inv.status.toLowerCase() === status.toLowerCase());
  }
}

export const invoiceService = new InvoiceService();
