import { AdminWebhookRepository, adminWebhookRepository, IAdminWebhookRepository } from '../repositories/admin_webhook.repository';
import type { WebhookEndpoint, WebhookDeliveryLog } from '../types/webhooks';

export class AdminWebhookService {
  constructor(private repo: IAdminWebhookRepository = adminWebhookRepository) {}

  async listWebhooks(): Promise<WebhookEndpoint[]> {
    return this.repo.getWebhooks();
  }

  async createWebhook(name: string, url: string, events: string[]): Promise<WebhookEndpoint> {
    if (!name.trim()) {
      throw new Error('Webhook name is required');
    }
    if (!url.trim() || !url.startsWith('http')) {
      throw new Error('Valid webhook URL (HTTP/HTTPS) is required');
    }
    if (!events || events.length === 0) {
      throw new Error('At least one event trigger must be selected');
    }
    const secret = `whsec_${Math.random().toString(36).substring(2)}${Math.random().toString(36).substring(2)}`;
    return this.repo.createWebhook({
      name: name.trim(),
      url: url.trim(),
      events,
      secret,
      is_active: true,
      total_deliveries: 0,
    });
  }

  async deleteWebhook(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid webhook ID is required');
    }
    return this.repo.deleteWebhook(id);
  }

  async triggerTestPing(id: number): Promise<{ success: boolean; log: WebhookDeliveryLog }> {
    if (!id || id <= 0) {
      throw new Error('Valid webhook ID is required');
    }
    return this.repo.triggerTestPing(id);
  }

  async listLogs(): Promise<WebhookDeliveryLog[]> {
    return this.repo.getLogs();
  }
}

export const adminWebhookService = new AdminWebhookService();
