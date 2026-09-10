import { request } from '../lib/api/client';
import type { SystemInfo, SystemStatusInfo } from '@/types/api';

export class AdminSystemRepository {
  getSystemInfo = () => request<SystemInfo>('/admin/system/info');
  getSystemStatus = () => request<SystemStatusInfo>('/admin/system/status');
  triggerCron = () => request<{ success: boolean; message: string }>('/admin/system/cron', { method: 'POST' });
  clearCache = () => request<{ success: boolean; message: string }>('/admin/system/clear-cache', { method: 'POST' });
  listNews = () => request<any[]>('/admin/news');
  createNews = (d: any) => request('/admin/news', { method: 'POST', body: JSON.stringify(d) });
  deleteNews = (id: number) => request(`/admin/news/${id}`, { method: 'DELETE' });
  exportBackup = () => request('/admin/system/backup', { method: 'POST' });
  getAuditLogs = (l = 10, o = 0) => request<any[]>(`/admin/system/audit-logs?limit=${l}&offset=${o}`);
  getActivityTrend = (days: number) => request<Record<string, number>>(`/admin/activity/trend?days=${days}`);
}

export const adminSystemRepository = new AdminSystemRepository();
