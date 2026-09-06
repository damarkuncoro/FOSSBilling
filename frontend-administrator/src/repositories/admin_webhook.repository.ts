import { request } from '../lib/api/client';
import type { WebhookEndpoint, WebhookDeliveryLog } from '../types/webhooks';

export interface IAdminWebhookRepository {
  getWebhooks(): Promise<WebhookEndpoint[]>;
  createWebhook(dto: Partial<WebhookEndpoint>): Promise<WebhookEndpoint>;
  updateWebhook(id: number, dto: Partial<WebhookEndpoint>): Promise<WebhookEndpoint>;
  deleteWebhook(id: number): Promise<any>;
  getLogs(): Promise<WebhookDeliveryLog[]>;
  triggerTestPing(id: number): Promise<{ success: boolean; log: WebhookDeliveryLog }>;
}

export class AdminWebhookRepository implements IAdminWebhookRepository {
  async getWebhooks(): Promise<WebhookEndpoint[]> {
    return request<WebhookEndpoint[]>('/admin/webhooks');
  }

  async createWebhook(dto: Partial<WebhookEndpoint>): Promise<WebhookEndpoint> {
    return request<WebhookEndpoint>('/admin/webhooks', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateWebhook(id: number, dto: Partial<WebhookEndpoint>): Promise<WebhookEndpoint> {
    return request<WebhookEndpoint>(`/admin/webhooks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteWebhook(id: number): Promise<any> {
    return request<any>(`/admin/webhooks/${id}`, {
      method: 'DELETE',
    });
  }

  async getLogs(): Promise<WebhookDeliveryLog[]> {
    return request<WebhookDeliveryLog[]>('/admin/webhooks/logs');
  }

  async triggerTestPing(id: number): Promise<{ success: boolean; log: WebhookDeliveryLog }> {
    return request<{ success: boolean; log: WebhookDeliveryLog }>(`/admin/webhooks/${id}/test`, {
      method: 'POST',
    });
  }
}

export const adminWebhookRepository = new AdminWebhookRepository();
