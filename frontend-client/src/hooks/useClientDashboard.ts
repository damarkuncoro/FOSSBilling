import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { Order, Invoice, SupportTicket } from '@/types/api';

export function useClientDashboard() {
  const { user, balance } = useClientAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchDashboardData = useCallback(() => {
    setLoading(true);
    Promise.all([
      api.getOrders().catch(() => []),
      api.getInvoices().catch(() => []),
      api.getTickets().catch(() => []),
    ]).then(([ordersData, invoicesData, ticketsData]) => {
      setOrders(ordersData || []);
      setInvoices(invoicesData || []);
      setTickets(ticketsData || []);
      setLoading(false);
    });
  }, []);

  useEffect(() => {
    fetchDashboardData();
  }, [fetchDashboardData]);

  const unpaidInvoices = invoices.filter((inv) => inv.status === 'unpaid');
  const activeOrders = orders.filter((o) => o.status === 'active');

  return {
    user,
    balance,
    orders,
    invoices,
    tickets,
    loading,
    unpaidInvoices,
    activeOrders,
    fetchDashboardData,
  };
}
