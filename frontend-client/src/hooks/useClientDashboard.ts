import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';

export function useClientDashboard() {
  const { user, balance } = useClientAuth();

  const { data: orders = [], isLoading: ordersLoading, refetch: fetchOrders } = useQuery({
    queryKey: ['client', 'orders'],
    queryFn: () => api.getOrders(),
  });

  const { data: invoices = [], isLoading: invoicesLoading, refetch: fetchInvoices } = useQuery({
    queryKey: ['client', 'invoices'],
    queryFn: () => api.getInvoices(),
  });

  const { data: tickets = [], isLoading: ticketsLoading, refetch: fetchTickets } = useQuery({
    queryKey: ['client', 'support', 'tickets'],
    queryFn: () => api.getTickets(),
  });

  const unpaidInvoices = invoices.filter((inv) => inv.status === 'unpaid');
  const activeOrders = orders.filter((o) => o.status === 'active');

  return {
    user,
    balance,
    orders,
    invoices,
    tickets,
    loading: ordersLoading || invoicesLoading || ticketsLoading,
    unpaidInvoices,
    activeOrders,
    fetchDashboardData: () => {
      fetchOrders();
      fetchInvoices();
      fetchTickets();
    },
  };
}
