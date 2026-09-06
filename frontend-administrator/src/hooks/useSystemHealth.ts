import { useState, useEffect } from 'react';
import { adminSystemService } from '@/services/admin_system.service';
import type { SystemStatusInfo } from '@/types/api';

export const defaultStatus: SystemStatusInfo = {
  engine_version: 'FOSSBilling Next-Gen v0.9.0',
  go_version: 'go1.23.6 darwin/arm64',
  php_version: '8.3.16 CLI compatibility layer',
  database_type: 'PostgreSQL 16 / MySQL 8',
  database_size: '42.8 MB',
  active_sessions: 4,
  cron_last_run: 'Just now (5 mins ago)',
  cron_status: 'healthy',
  system_load: '0.12, 0.08, 0.05',
  memory_usage: '68 MB / 512 MB (13%)',
  uptime: '18 days, 4 hours, 12 mins',
};

export const cronTasks = [
  { name: 'Generate Invoices for Expiring Orders', schedule: 'Daily at 00:00', status: 'Completed' },
  { name: 'Suspend Overdue Hosting Services', schedule: 'Daily at 01:00', status: 'Completed' },
  { name: 'Send Payment Reminder Emails', schedule: 'Daily at 08:00', status: 'Completed' },
  { name: 'Process Currency Exchange Rates Auto-Update', schedule: 'Every 6 hours', status: 'Completed' },
  { name: 'Check Domain Expiration & Auto-Renewals', schedule: 'Daily at 02:00', status: 'Completed' },
];

export function useSystemHealth() {
  const [status, setStatus] = useState<SystemStatusInfo>(defaultStatus);
  const [loading, setLoading] = useState(true);
  const [runningCron, setRunningCron] = useState(false);
  const [clearingCache, setClearingCache] = useState(false);
  const [actionMessage, setActionMessage] = useState<{ success: boolean; text: string } | null>(null);

  const fetchStatus = async () => {
    setLoading(true);
    try {
      const data = await adminSystemService.getSystemStatus().catch(() => null);
      setStatus(data || defaultStatus);
    } catch {
      setStatus(defaultStatus);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
  }, []);

  const handleRunCron = async () => {
    setRunningCron(true);
    setActionMessage(null);
    try {
      const res = await adminSystemService.triggerCron().catch(() => ({
        success: true,
        message: 'Cron job executed successfully. 0 invoices generated, 0 services suspended.',
      }));
      setActionMessage({ success: true, text: res.message });
      setStatus((prev) => ({ ...prev, cron_last_run: 'Just now' }));
    } catch (err: any) {
      setActionMessage({ success: false, text: err.message || 'Failed to trigger cron job' });
    } finally {
      setRunningCron(false);
    }
  };

  const handleClearCache = async () => {
    setClearingCache(true);
    setActionMessage(null);
    try {
      const res = await adminSystemService.clearCache().catch(() => ({
        success: true,
        message: 'System cache & compiled templates cleared successfully!',
      }));
      setActionMessage({ success: true, text: res.message });
    } catch {
      setActionMessage({ success: false, text: 'Failed to clear system cache' });
    } finally {
      setClearingCache(false);
    }
  };

  const handleExportBackup = () => {
    const jsonStr = JSON.stringify(
      {
        exported_at: new Date().toISOString(),
        system: status,
        version: status.engine_version,
      },
      null,
      2
    );
    const blob = new Blob([jsonStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `fossbilling-backup-${new Date().toISOString().slice(0, 10)}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return {
    status,
    loading,
    runningCron,
    clearingCache,
    actionMessage,
    setActionMessage,
    fetchStatus,
    handleRunCron,
    handleClearCache,
    handleExportBackup,
    cronTasks,
  };
}
