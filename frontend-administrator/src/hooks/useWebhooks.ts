import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminWebhookService } from '@/services/admin_webhook.service';

export const AVAILABLE_EVENTS = [
  { id: 'client_registered', label: 'Client Registered' },
  { id: 'invoice_paid', label: 'Invoice Paid' },
  { id: 'order_activated', label: 'Order Activated' },
  { id: 'order_suspended', label: 'Order Suspended' },
  { id: 'ticket_opened', label: 'Ticket Opened' },
  { id: 'system_low_stock', label: 'System Low Stock' },
];

export function useWebhooks() {
  const queryClient = useQueryClient();
  const [isAddModal, setIsAddModal] = useState(false);

  const { data: webhooks = [] } = useQuery({ queryKey: ['admin', 'webhooks'], queryFn: () => adminWebhookService.listWebhooks().catch(() => []) });
  const { data: logs = [] } = useQuery({ queryKey: ['admin', 'webhooks', 'logs'], queryFn: () => adminWebhookService.listLogs().catch(() => []) });

  const createMutation = useMutation({
    mutationFn: (d: any) => adminWebhookService.createWebhook(d.name, d.url, d.events),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'webhooks'] }); setIsAddModal(false); },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminWebhookService.deleteWebhook(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'webhooks'] }),
  });

  const testMutation = useMutation({
    mutationFn: (id: number) => adminWebhookService.triggerTestPing(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'webhooks', 'logs'] }),
  });

  return {
    webhooks, logs, isAddModal, setIsAddModal,
    testingWebhookId: testMutation.isPending ? -1 : null,
    createWebhook: createMutation.mutate,
    deleteWebhook: (id: number) => { if (confirm('Delete?')) deleteMutation.mutate(id); },
    triggerTestPayload: testMutation.mutate,
    toggleWebhook: (id: number) => { /* logic */ }
  };
}
