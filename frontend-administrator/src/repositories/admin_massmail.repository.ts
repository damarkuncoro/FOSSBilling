import { request } from '../lib/api/client';

export interface IAdminMassMailRepository {
  sendBroadcast(dto: { subject: string; body: string; client_group_id?: number }): Promise<{ queued: number }>;
}

export class AdminMassMailRepository implements IAdminMassMailRepository {
  async sendBroadcast(dto: { subject: string; body: string; client_group_id?: number }): Promise<{ queued: number }> {
    return request<{ queued: number }>('/admin/mass-mail', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }
}

export const adminMassMailRepository = new AdminMassMailRepository();
