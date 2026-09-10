import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notificationService } from '@/services/notification.service';
import { useClientAuth } from '@/lib/auth';

export function useNotifications() {
  const queryClient = useQueryClient();
  const { isAuthenticated } = useClientAuth();

  const { data: notifications = [], isLoading: loading } = useQuery({
    queryKey: ['client', 'notifications'],
    queryFn: () => notificationService.listNotifications(),
    enabled: isAuthenticated,
    refetchInterval: 60000,
  });

  const markMutation = useMutation({
    mutationFn: (id: number) => notificationService.markAsRead(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'notifications'] }),
  });

  const markAllMutation = useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'notifications'] }),
  });

  return {
    notifications, loading,
    unreadCount: notifications.filter((n: any) => !n.is_read).length,
    markAsRead: (id: number) => markMutation.mutate(id),
    markAllAsRead: () => markAllMutation.mutate(),
  };
}
