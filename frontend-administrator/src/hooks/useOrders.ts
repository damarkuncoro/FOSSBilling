import { useState, useEffect } from 'react';
import { adminOrderService } from '@/services/admin_order.service';
import type { Order } from '@/types/api';

export function useOrders() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<number | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const data = await adminOrderService.listOrders();
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
      await adminOrderService.activateOrder(id);
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
      await adminOrderService.suspendOrder(id, 'Admin manual suspension');
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
      await adminOrderService.unsuspendOrder(id);
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
