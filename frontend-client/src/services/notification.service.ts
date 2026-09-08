import { api } from '@/lib/api';
import { Notification } from '@/types/api';

export const notificationService = {
  async listNotifications(): Promise<Notification[]> {
    return api.getNotifications();
  },

  async markAsRead(id: number): Promise<void> {
    await api.markNotificationRead(id);
  },

  async markAllAsRead(): Promise<void> {
    await api.markAllNotificationsRead();
  },
};
