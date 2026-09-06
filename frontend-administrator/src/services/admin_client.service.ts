import { AdminClientRepository, adminClientRepository, IAdminClientRepository } from '../repositories/admin_client.repository';
import type { ClientProfile } from '@/types/api';

export class AdminClientService {
  constructor(private repo: IAdminClientRepository = adminClientRepository) {}

  async listClients(): Promise<ClientProfile[]> {
    return this.repo.listClients();
  }

  async getClientDetail(id: number): Promise<ClientProfile> {
    if (!id || id <= 0) {
      throw new Error('Valid client ID is required');
    }
    return this.repo.getClient(id);
  }

  async createClient(dto: Partial<ClientProfile>): Promise<ClientProfile> {
    if (!dto.email || !dto.email.includes('@')) {
      throw new Error('A valid email address is required');
    }
    return this.repo.createClient(dto);
  }

  async updateClient(id: number, dto: Partial<ClientProfile>): Promise<ClientProfile> {
    if (!id || id <= 0) {
      throw new Error('Valid client ID is required');
    }
    return this.repo.updateClient(id, dto);
  }

  async deleteClient(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid client ID is required');
    }
    return this.repo.deleteClient(id);
  }

  filterClients(clients: ClientProfile[], query: string): ClientProfile[] {
    if (!query.trim()) return clients;
    const q = query.toLowerCase();
    return clients.filter(
      (c) =>
        c.email?.toLowerCase().includes(q) ||
        c.first_name?.toLowerCase().includes(q) ||
        c.last_name?.toLowerCase().includes(q) ||
        c.company?.toLowerCase().includes(q)
    );
  }
}

export const adminClientService = new AdminClientService();
