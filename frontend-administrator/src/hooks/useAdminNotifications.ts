import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { request } from '../lib/api/client';

export function useAdminNotifications() {
  const queryClient = useQueryClient();

  const { data: notifications = [], isLoading: loading } = useQuery({
    queryKey: ['admin', 'notifications'],
    queryFn: () => request<any[]>('/admin/notifications'),
    refetchInterval: 30000,
  });

  const markMutation = useMutation({
    mutationFn: (id: number) => request(`/admin/notifications/${id}/read`, { method: 'PUT' }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'notifications'] }),
  });

  const markAllMutation = useMutation({
    mutationFn: () => request('/admin/notifications/mark-all-read', { method: 'POST' }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'notifications'] }),
  });

  return {
    notifications,
    unreadCount: notifications.filter((n: any) => !n.is_read).length,
    loading,
    markAsRead: markMutation.mutate,
    markAllRead: markAllMutation.mutate,
  };
}
