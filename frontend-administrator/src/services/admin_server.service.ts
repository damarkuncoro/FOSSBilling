import { AdminServerRepository, adminServerRepository, IAdminServerRepository } from '../repositories/admin_server.repository';
import type { ServerItem } from '@/types/api';

export class AdminServerService {
  constructor(private repo: IAdminServerRepository = adminServerRepository) {}

  async listServers(): Promise<ServerItem[]> {
    return this.repo.getServers();
  }

  async getServer(id: number): Promise<ServerItem> {
    if (!id || id <= 0) {
      throw new Error('Valid server ID is required');
    }
    return this.repo.getServer(id);
  }

  async createServer(dto: Partial<ServerItem>): Promise<ServerItem> {
    if (!dto.name || !dto.name.trim()) {
      throw new Error('Server name is required');
    }
    if (!dto.hostname && !dto.ip) {
      throw new Error('Server hostname or IP address is required');
    }
    return this.repo.createServer({
      ...dto,
      manager: dto.manager || 'cpanel',
      status: dto.status || 'online',
      active_accounts: dto.active_accounts || 0,
      max_accounts: Number(dto.max_accounts) || 100,
    });
  }

  async updateServer(id: number, dto: Partial<ServerItem>): Promise<ServerItem> {
    if (!id || id <= 0) {
      throw new Error('Valid server ID is required');
    }
    return this.repo.updateServer(id, dto);
  }

  async deleteServer(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid server ID is required');
    }
    return this.repo.deleteServer(id);
  }

  async testConnection(id: number): Promise<{ success: boolean; message: string }> {
    if (!id || id <= 0) {
      throw new Error('Valid server ID is required');
    }
    return this.repo.testServerConnection(id);
  }
}

export const adminServerService = new AdminServerService();
