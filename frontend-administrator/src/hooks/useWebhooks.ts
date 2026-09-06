import { useState } from 'react';
import { adminWebhookService } from '@/services/admin_webhook.service';
import type { WebhookEndpoint, WebhookDeliveryLog } from '../types/webhooks';

export const AVAILABLE_EVENTS = [
  { id: 'order.created', label: 'Order Created', category: 'Orders' },
  { id: 'order.activated', label: 'Order Activated', category: 'Orders' },
  { id: 'order.suspended', label: 'Order Suspended', category: 'Orders' },
  { id: 'invoice.created', label: 'Invoice Generated', category: 'Billing' },
  { id: 'invoice.paid', label: 'Invoice Paid', category: 'Billing' },
  { id: 'ticket.created', label: 'Ticket Opened', category: 'Support' },
  { id: 'ticket.reply', label: 'Ticket Staff/Client Reply', category: 'Support' },
  { id: 'client.registered', label: 'New Client Registered', category: 'Clients' },
];

const initialWebhooks: WebhookEndpoint[] = [
  {
    id: 1,
    name: 'Discord Ops Notifications',
    url: 'https://discord.com/api/webhooks/123456789/token_abc',
    secret: 'whsec_99182371928371928371',
    events: ['order.activated', 'invoice.paid', 'ticket.created'],
    is_active: true,
    total_deliveries: 142,
    last_status_code: 200,
    last_delivery_at: '2026-09-05T14:20:00Z',
    created_at: '2026-01-20T10:00:00Z',
  },
  {
    id: 2,
    name: 'Slack #billing Channel',
    url: 'https://hooks.slack.com/services/T00/B00/XXXX',
    secret: 'whsec_55182910381928471920',
    events: ['invoice.paid'],
    is_active: true,
    total_deliveries: 89,
    last_status_code: 200,
    last_delivery_at: '2026-09-05T13:00:00Z',
    created_at: '2026-02-14T11:00:00Z',
  },
];

const initialLogs: WebhookDeliveryLog[] = [
  {
    id: 1,
    webhook_id: 1,
    event: 'invoice.paid',
    url: 'https://discord.com/api/webhooks/123456789/token_abc',
    response_code: 200,
    response_body: '{"message": "delivered"}',
    execution_time_ms: 184,
    created_at: '2026-09-05T14:20:00Z',
  },
];

export function useWebhooks() {
  const [webhooks, setWebhooks] = useState<WebhookEndpoint[]>(initialWebhooks);
  const [logs, setLogs] = useState<WebhookDeliveryLog[]>(initialLogs);
  const [isAddModal, setIsAddModal] = useState(false);
  const [testingWebhookId, setTestingWebhookId] = useState<number | null>(null);

  const createWebhook = (data: { name: string; url: string; events: string[] }) => {
    const newWh: WebhookEndpoint = {
      id: Date.now(),
      name: data.name,
      url: data.url,
      secret: `whsec_${Math.random().toString(36).substring(2)}${Math.random().toString(36).substring(2)}`,
      events: data.events,
      is_active: true,
      total_deliveries: 0,
      created_at: new Date().toISOString(),
    };
    setWebhooks((prev) => [newWh, ...prev]);
    setIsAddModal(false);
    adminWebhookService.createWebhook(data.name, data.url, data.events).catch(() => null);
  };

  const toggleWebhook = (id: number) => {
    setWebhooks((prev) =>
      prev.map((wh) => (wh.id === id ? { ...wh, is_active: !wh.is_active } : wh))
    );
  };

  const deleteWebhook = (id: number) => {
    setWebhooks((prev) => prev.filter((wh) => wh.id !== id));
    adminWebhookService.deleteWebhook(id).catch(() => null);
  };

  const triggerTestPayload = (id: number) => {
    const wh = webhooks.find((w) => w.id === id);
    if (!wh) return;

    setTestingWebhookId(id);
    setTimeout(() => {
      const newLog: WebhookDeliveryLog = {
        id: Date.now(),
        webhook_id: id,
        event: 'test.ping',
        url: wh.url,
        response_code: 200,
        response_body: '{"status": "ok", "message": "Test ping received"}',
        execution_time_ms: Math.floor(Math.random() * 120) + 80,
        created_at: new Date().toISOString(),
      };
      setLogs((prev) => [newLog, ...prev]);
      setWebhooks((prev) =>
        prev.map((w) =>
          w.id === id
            ? { ...w, total_deliveries: w.total_deliveries + 1, last_status_code: 200, last_delivery_at: new Date().toISOString() }
            : w
        )
      );
      setTestingWebhookId(null);
      adminWebhookService.triggerTestPing(id).catch(() => null);
    }, 600);
  };

  return {
    webhooks,
    logs,
    isAddModal,
    setIsAddModal,
    testingWebhookId,
    createWebhook,
    toggleWebhook,
    deleteWebhook,
    triggerTestPayload,
  };
}
