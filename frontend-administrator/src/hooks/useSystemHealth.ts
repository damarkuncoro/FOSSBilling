import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminSystemService } from '@/services/admin_system.service';
import type { SystemStatus } from '@/types/api';

export function useSystemHealth() {
  const queryClient = useQueryClient();
  const [actionMessage, setActionMessage] = useState<{ text: string; success: boolean } | null>(null);

  const { data: status = {
    engine_version: 'v0.7.0-NextGen',
    go_version: 'go1.22.0',
    database_type: 'PostgreSQL 16',
    database_size: '124MB',
    active_sessions: 42,
    cron_last_run: '2 minutes ago',
    cron_status: 'healthy',
    system_load: 'CPU: 12% / RAM: 45%',
    uptime: '14 days, 6 hours'
  } as SystemStatus, isLoading: loading, refetch: fetchStatus } = useQuery({
    queryKey: ['admin', 'system', 'status'],
    queryFn: () => adminSystemService.getSystemStatus(),
    refetchInterval: 5000,
  });

  const cronMutation = useMutation({
    mutationFn: () => adminSystemService.triggerCron(),
    onSuccess: () => {
      setActionMessage({ text: 'Cron tasks executed successfully.', success: true });
      queryClient.invalidateQueries({ queryKey: ['admin', 'system', 'status'] });
    },
    onError: (err: any) => setActionMessage({ text: `Cron failed: ${err.message}`, success: false }),
  });

  const cacheMutation = useMutation({
    mutationFn: () => adminSystemService.clearCache(),
    onSuccess: () => setActionMessage({ text: 'System cache cleared.', success: true }),
    onError: (err: any) => setActionMessage({ text: `Failed to clear cache: ${err.message}`, success: false }),
  });

  const handleExportBackup = async () => {
    try {
      const data = await adminSystemService.exportBackup();
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `fossbilling-backup-${new Date().toISOString()}.json`;
      a.click();
    } catch (err) {
      console.error(err);
    }
  };

  const cronTasks = [
    { name: 'Invoice Generator', schedule: 'Every 5 minutes', status: 'Active' },
    { name: 'Service Provisioner', schedule: 'Every 1 minute', status: 'Active' },
    { name: 'Support Ticket Auto-Close', schedule: 'Daily at 00:00', status: 'Active' },
    { name: 'System Backup', schedule: 'Daily at 02:00', status: 'Active' },
    { name: 'Currency Rate Sync', schedule: 'Daily at 04:00', status: 'Active' },
  ];

  return {
    status,
    loading,
    runningCron: cronMutation.isPending,
    clearingCache: cacheMutation.isPending,
    actionMessage,
    setActionMessage,
    fetchStatus,
    handleRunCron: () => cronMutation.mutate(),
    handleClearCache: () => cacheMutation.mutate(),
    handleExportBackup,
    cronTasks,
  };
}
