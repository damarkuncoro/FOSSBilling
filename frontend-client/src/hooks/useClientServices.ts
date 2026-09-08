import { useState, useEffect, useCallback } from 'react';
import { orderService } from '@/services/order.service';
import { downloadService } from '@/services/download.service';
import { Order } from '@/types/api';

export function useClientServices() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [downloadLink, setDownloadLink] = useState<string | null>(null);
  const [downloadModal, setDownloadModal] = useState(false);

  const fetchServices = useCallback(async () => {
    setLoading(true);
    try {
      const data = await orderService.listClientOrders();
      setOrders(data || []);
    } catch (err) {
      console.error('Failed to fetch services:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchServices();
  }, [fetchServices]);

  const handleGetDownload = async (id: number) => {
    try {
      const url = await downloadService.getSecureDownloadUrl(id);
      setDownloadLink(url);
      setDownloadModal(true);
    } catch (err: any) {
      alert(`Download generation failed: ${err.message}`);
    }
  };

  const handleSyncStatus = async (id: number) => {
    try {
      await orderService.syncServiceStatus(id);
      alert('Service status synchronized with remote provider.');
      fetchServices();
    } catch (err: any) {
      alert(`Sync failed: ${err.message}`);
    }
  };

  const handleChangePassword = async (id: number, password: string) => {
    try {
      await orderService.changeServicePassword(id, password);
      alert('Service password updated successfully.');
      fetchServices();
    } catch (err: any) {
      alert(`Password change failed: ${err.message}`);
    }
  };

  return {
    orders,
    loading,
    downloadLink,
    downloadModal,
    setDownloadModal,
    fetchServices,
    handleGetDownload,
    handleSyncStatus,
    handleChangePassword,
  };
}
