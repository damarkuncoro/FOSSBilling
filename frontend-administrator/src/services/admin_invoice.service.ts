import { AdminInvoiceRepository, adminInvoiceRepository, IAdminInvoiceRepository } from '../repositories/admin_invoice.repository';
import type { Invoice } from '@/types/api';

export class AdminInvoiceService {
  constructor(private repo: IAdminInvoiceRepository = adminInvoiceRepository) {}

  async listInvoices(): Promise<Invoice[]> {
    return this.repo.listInvoices();
  }

  async getInvoiceDetail(id: number): Promise<Invoice> {
    if (!id || id <= 0) {
      throw new Error('Valid invoice ID is required');
    }
    return this.repo.getInvoice(id);
  }

  async createInvoice(clientId: number, items: Array<{ title: string; price: number; quantity?: number }>): Promise<Invoice> {
    if (!clientId || clientId <= 0) {
      throw new Error('Client ID is required');
    }
    if (!items || items.length === 0) {
      throw new Error('Invoice must contain at least one line item');
    }
    return this.repo.createInvoice({
      client_id: clientId,
      items,
    });
  }

  filterByStatus(invoices: Invoice[], status?: string): Invoice[] {
    if (!status || status === 'all') return invoices;
    return invoices.filter((inv) => inv.status.toLowerCase() === status.toLowerCase());
  }
}

export const adminInvoiceService = new AdminInvoiceService();
