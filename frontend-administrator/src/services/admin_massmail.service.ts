import { AdminMassMailRepository, adminMassMailRepository, IAdminMassMailRepository } from '../repositories/admin_massmail.repository';

export class AdminMassMailService {
  constructor(private repo: IAdminMassMailRepository = adminMassMailRepository) {}

  async sendMassMail(subject: string, body: string, clientGroupId?: number): Promise<{ queued: number }> {
    if (!subject.trim()) {
      throw new Error('Email subject is required');
    }
    if (!body.trim()) {
      throw new Error('Email body content is required');
    }
    return this.repo.sendBroadcast({
      subject: subject.trim(),
      body: body.trim(),
      client_group_id: clientGroupId,
    });
  }
}

export const adminMassMailService = new AdminMassMailService();
