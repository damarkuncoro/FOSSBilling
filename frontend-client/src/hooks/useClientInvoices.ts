import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { Invoice } from '@/types/api';

export function useClientInvoices() {
  const { user, balance, refreshProfile } = useClientAuth();
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [payModal, setPayModal] = useState<Invoice | null>(null);
  const [paying, setPaying] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  const fetchInvoices = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getInvoices();
      setInvoices(data || []);
    } catch (err) {
      console.error('Failed to fetch invoices:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchInvoices();
  }, [fetchInvoices]);

  const handlePayBalance = async (id: number) => {
    setPaying(true);
    setMessage(null);
    try {
      await api.payWithBalance(id);
      setMessage(`Invoice #${id} successfully paid with account balance!`);
      setPayModal(null);
      await Promise.all([fetchInvoices(), refreshProfile()]);
    } catch (err: any) {
      alert(`Payment failed: ${err.message}`);
    } finally {
      setPaying(false);
    }
  };

  return {
    user,
    balance,
    invoices,
    loading,
    payModal,
    setPayModal,
    paying,
    message,
    fetchInvoices,
    handlePayBalance,
  };
}
