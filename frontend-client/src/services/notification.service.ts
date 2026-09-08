import { api } from '@/lib/api';
import { Notification } from '@/types/api';

export const notificationService = {
  async listNotifications(): Promise<Notification[]> {
    const res = await api.get<{ data: Notification[] }>('/client/notifications');
    return res.data.data;
  },

  async markAsRead(id: number): Promise<void> {
    await api.put(`/client/notifications/${id}/read`, {});
  },

  async markAllAsRead(): Promise<void> {
    await api.post('/client/notifications/mark-all-read', {});
  },
};
