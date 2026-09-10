import { useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminInvoiceService } from '@/services/admin_invoice.service';
import { adminClientService } from '@/services/admin_client.service';
import { adminStatsService } from '@/services/admin_stats.service';
import type { Invoice } from '@/types/api';

export function useInvoices() {
  const queryClient = useQueryClient();

  const { data: invoices = [], isLoading: invoicesLoading, refetch: fetchInvoices } = useQuery({
    queryKey: ['admin', 'invoices'],
    queryFn: () => adminInvoiceService.listInvoices(),
  });

  const { data: clients = [], isLoading: clientsLoading } = useQuery({
    queryKey: ['admin', 'clients'],
    queryFn: () => adminClientService.listClients(),
  });

  const { data: stats = null } = useQuery({
    queryKey: ['admin', 'stats', 'dashboard'],
    queryFn: () => adminStatsService.getDashboardMetrics(),
  });

  const createMutation = useMutation({
    mutationFn: ({ clientId, items }: { clientId: number; items: any[] }) =>
      adminInvoiceService.createInvoice(clientId, items),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'invoices'] }),
  });

  const refundMutation = useMutation({
    mutationFn: (id: number) => adminInvoiceService.refundInvoice(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'invoices'] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminInvoiceService.deleteInvoice(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'invoices'] }),
  });

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
    const link = document.createElement('a');
    link.setAttribute('href', encodeURI(csvContent));
    link.setAttribute('download', `FOSSBilling-Invoices-${new Date().toISOString().slice(0, 10)}.csv`);
    link.click();
  }, [invoices]);

  return {
    stats,
    invoices,
    clients,
    loading: invoicesLoading || clientsLoading,
    fetchInvoices,
    exportInvoicesToCSV,
    createInvoice: (clientId: number, items: any[]) => createMutation.mutateAsync({ clientId, items }),
    refundInvoice: refundMutation.mutateAsync,
    deleteInvoice: async (id: number) => {
      if (confirm('Are you sure?')) await deleteMutation.mutateAsync(id);
    },
  };
}
