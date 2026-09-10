import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { orderService } from '@/services/order.service';
import { downloadService } from '@/services/download.service';

export function useClientServices() {
  const queryClient = useQueryClient();
  const [downloadLink, setDownloadLink] = useState<string | null>(null);
  const [downloadModal, setDownloadModal] = useState(false);

  const { data: orders = [], isLoading: loading, refetch } = useQuery({
    queryKey: ['client', 'orders'],
    queryFn: () => orderService.listClientOrders(),
  });

  const syncMutation = useMutation({
    mutationFn: (id: number) => orderService.syncServiceStatus(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'orders'] }); alert('Synchronized.'); },
  });

  const passwordMutation = useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) => orderService.changeServicePassword(id, password),
    onSuccess: () => alert('Updated.'),
  });

  return {
    orders, loading, downloadLink, downloadModal, setDownloadModal,
    fetchServices: () => refetch(),
    handleGetDownload: async (id: number) => { try { const url = await downloadService.getSecureDownloadUrl(id); setDownloadLink(url); setDownloadModal(true); } catch (err: any) { alert(err.message); } },
    handleSyncStatus: (id: number) => syncMutation.mutateAsync(id),
    handleChangePassword: (id: number, password: string) => passwordMutation.mutateAsync({ id, password }),
  };
}
