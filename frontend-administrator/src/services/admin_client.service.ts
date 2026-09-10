import { adminClientRepository, AdminClientRepository } from '../repositories/admin_client.repository';
import type { ClientProfile } from '@/types/api';

export class AdminClientService {
  constructor(private repo: AdminClientRepository = adminClientRepository) {}

  listClients = () => this.repo.list();
  getClientDetail = (id: number) => {
    if (!id || id <= 0) throw new Error('Valid client ID is required');
    return this.repo.get(id);
  };
  createClient = (d: any) => this.repo.create(d);
  updateClient = (id: number, d: any) => this.repo.update(id, d);
  deleteClient = (id: number) => this.repo.delete(id);
  impersonateClient = async (id: number) => (await this.repo.impersonate(id)).token;

  filterClients(cls: ClientProfile[], q: string) {
    if (!q.trim()) return cls;
    const lq = q.toLowerCase();
    return cls.filter(c =>
      (c.email || '').toLowerCase().includes(lq) ||
      (c.first_name || '').toLowerCase().includes(lq) ||
      (c.last_name || '').toLowerCase().includes(lq)
    );
  }
}

export const adminClientService = new AdminClientService();
