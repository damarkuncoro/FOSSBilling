import { adminInvoiceRepository, AdminInvoiceRepository } from '../repositories/admin_invoice.repository';
import type { Invoice } from '@/types/api';

export class AdminInvoiceService {
  constructor(private repo: AdminInvoiceRepository = adminInvoiceRepository) {}

  listInvoices = () => this.repo.list();
  getInvoiceDetail = (id: number) => this.repo.get(id);
  createInvoice = (clientId: number, items: any[]) => {
    if (!items || items.length === 0) throw new Error('Invoice must contain at least one line item');
    return this.repo.create({ client_id: clientId, items });
  };
  refundInvoice = (id: number) => this.repo.refund(id);
  deleteInvoice = (id: number) => this.repo.delete(id);

  filterByStatus(invs: Invoice[], st?: string) {
    if (!st || st === 'all') return invs;
    return invs.filter(i => (i.status || '').toLowerCase() === st.toLowerCase());
  }
}

export const adminInvoiceService = new AdminInvoiceService();
