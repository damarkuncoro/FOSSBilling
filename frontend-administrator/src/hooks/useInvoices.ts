import { useState, useEffect, useCallback } from 'react';
import { adminInvoiceService } from '@/services/admin_invoice.service';
import { adminClientService } from '@/services/admin_client.service';
import { adminStatsService } from '@/services/admin_stats.service';
import type { Invoice, ClientProfile, DashboardStats } from '@/types/api';

export function useInvoices() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [clients, setClients] = useState<ClientProfile[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchInvoices = useCallback(async () => {
    setLoading(true);
    try {
      const [invData, clientData, statsData] = await Promise.allSettled([
        adminInvoiceService.listInvoices(),
        adminClientService.listClients(),
        adminStatsService.getDashboardMetrics(),
      ]);

      if (invData.status === 'fulfilled' && Array.isArray(invData.value)) {
        setInvoices(invData.value);
      }
      if (clientData.status === 'fulfilled' && Array.isArray(clientData.value)) {
        setClients(clientData.value);
      }
      if (statsData.status === 'fulfilled') {
        setStats(statsData.value);
      }
    } catch (err) {
      console.error('Failed to fetch invoices/clients:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const exportInvoicesToCSV = useCallback(() => {
    if (!invoices.length) return;
    const headers = ['Invoice ID', 'Serie Nr', 'Client ID', 'Subtotal', 'Tax', 'Total', 'Currency', 'Status', 'Created At', 'Due At'];
    const rows = invoices.map((inv: any) => [
      inv.id,
      inv.nr || inv.serie_nr || `INV${inv.id}`,
      inv.client_id,
      inv.subtotal,
      inv.tax,
      inv.total,
      inv.currency,
      inv.status,
      inv.created_at,
      inv.due_at,
    ]);

    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n');
    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `FOSSBilling-Invoices-${new Date().toISOString().slice(0, 10)}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }, [invoices]);

  const createInvoice = useCallback(async (arg1: any, arg2?: any) => {
    if (typeof arg1 === 'object' && arg1 !== null) {
      return adminInvoiceService.createInvoice(arg1.client_id, arg1.items);
    }
    return adminInvoiceService.createInvoice(arg1, arg2);
  }, []);

  useEffect(() => {
    fetchInvoices();
  }, [fetchInvoices]);

  return {
    stats,
    invoices,
    clients,
    loading,
    fetchInvoices,
    exportInvoicesToCSV,
    createInvoice,
  };
}
