import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminOrderService } from '@/services/admin_order.service';
import type { Order } from '@/types/api';

export function useOrders() {
  const queryClient = useQueryClient();

  const { data: orders = [], isLoading: loading, refetch: fetchOrders } = useQuery({
    queryKey: ['admin', 'orders'],
    queryFn: () => adminOrderService.listOrders(),
  });

  const activateMutation = useMutation({
    mutationFn: (id: number) => adminOrderService.activateOrder(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'orders'] }),
  });

  const suspendMutation = useMutation({
    mutationFn: (id: number) => adminOrderService.suspendOrder(id, 'Admin manual suspension'),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'orders'] }),
  });

  const unsuspendMutation = useMutation({
    mutationFn: (id: number) => adminOrderService.unsuspendOrder(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'orders'] }),
  });

  return {
    orders,
    loading,
    actionLoading: activateMutation.isPending || suspendMutation.isPending || unsuspendMutation.isPending ? -1 : null,
    fetchOrders,
    handleActivate: activateMutation.mutateAsync,
    handleSuspend: suspendMutation.mutateAsync,
    handleUnsuspend: unsuspendMutation.mutateAsync,
  };
}
