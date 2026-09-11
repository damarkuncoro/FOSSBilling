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

  // Calculate monthly spending from paid invoices
  const spendingTrends = (() => {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const currentMonth = new Date().getMonth();
    const result = [];
    for (let i = 5; i >= 0; i--) {
      const idx = (currentMonth - i + 12) % 12;
      const monthName = months[idx];
      const monthPaid = invoices
        .filter(inv => inv.status === 'paid' && new Date(inv.created_at).getMonth() === idx)
        .reduce((sum, inv) => sum + (inv.total || 0), 0);
      result.push({ month: monthName, amount: monthPaid });
    }
    return result;
  })();

  return {
    user,
    balance,
    orders,
    invoices,
    tickets,
    unpaidInvoices,
    activeOrders,
    spendingTrends,
    loading: ordersLoading || invoicesLoading || ticketsLoading,
    fetchDashboardData: () => {
      fetchOrders();
      fetchInvoices();
      fetchTickets();
    },
  };
}
