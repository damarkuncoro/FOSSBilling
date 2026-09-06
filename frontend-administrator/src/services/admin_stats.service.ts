import { AdminStatsRepository, adminStatsRepository, IAdminStatsRepository } from '../repositories/admin_stats.repository';
import type { DashboardStats } from '@/types/api';

export class AdminStatsService {
  constructor(private repo: IAdminStatsRepository = adminStatsRepository) {}

  async getDashboardStats(): Promise<DashboardStats> {
    return this.repo.getDashboardStats();
  }

  async getDashboardMetrics(): Promise<DashboardStats> {
    return this.getDashboardStats();
  }

  async getIncomeSummary(): Promise<any> {
    return this.repo.getIncomeSummary();
  }
}

export const adminStatsService = new AdminStatsService();
