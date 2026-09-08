import { useState, useEffect, useCallback } from 'react';
import { request } from '../lib/api/client';

export function useAdminNotifications() {
  const [notifications, setNotifications] = useState<any[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);

  const fetchNotifications = useCallback(async () => {
    try {
      const res = await request<any[]>('/admin/notifications');
      setNotifications(res || []);
      setUnreadCount((res || []).filter((n: any) => !n.is_read).length);
    } catch (err) {
      console.error('Failed to fetch admin notifications');
    }
  }, []);

  useEffect(() => {
    fetchNotifications();
    // Poll every 30 seconds
    const interval = setInterval(fetchNotifications, 30000);
    return () => clearInterval(interval);
  }, [fetchNotifications]);

  const markAsRead = async (id: number) => {
    try {
      await request(`/admin/notifications/${id}/read`, { method: 'PUT' });
      setNotifications(prev => prev.map(n => n.id === id ? { ...n, is_read: true } : n));
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (err) {
      console.error('Failed to mark notification as read');
    }
  };

  const markAllRead = async () => {
    try {
      await request('/admin/notifications/mark-all-read', { method: 'POST' });
      setNotifications(prev => prev.map(n => ({ ...n, is_read: true })));
      setUnreadCount(0);
    } catch (err) {
      console.error('Failed to mark all notifications as read');
    }
  };

  return {
    notifications,
    unreadCount,
    loading,
    markAsRead,
    markAllRead,
    refresh: fetchNotifications,
  };
}
