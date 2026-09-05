import { useState, useEffect } from 'react';
import { api } from '@/lib/api';

export function useOrders() {
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<number | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const data = await api.getOrders();
      setOrders(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleActivate = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.activateOrder(id);
      setMessage(`Order #${id} successfully activated and provisioned!`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error activating order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  const handleSuspend = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.suspendOrder(id, 'Admin manual suspension');
      setMessage(`Order #${id} has been suspended.`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error suspending order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  const handleUnsuspend = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.unsuspendOrder(id);
      setMessage(`Order #${id} has been reactivated.`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error unsuspending order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  return {
    orders,
    loading,
    actionLoading,
    message,
    fetchOrders,
    handleActivate,
    handleSuspend,
    handleUnsuspend,
  };
}
