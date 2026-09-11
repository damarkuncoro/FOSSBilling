import { request } from '../lib/api/client';
import type { DashboardStats } from '@/types/api';

export interface IAdminStatsRepository {
  getDashboardStats(): Promise<DashboardStats>;
  getIncomeSummary(): Promise<any>;
  getRevenueProjections(): Promise<Record<string, number>>;
}

export class AdminStatsRepository implements IAdminStatsRepository {
  async getDashboardStats(): Promise<DashboardStats> {
    return request<DashboardStats>('/admin/stats/dashboard');
  }

  async getIncomeSummary(): Promise<any> {
    return request<any>('/admin/stats/income');
  }

  async getRevenueProjections(): Promise<Record<string, number>> {
    return request<Record<string, number>>('/admin/stats/projections');
  }
}

export const adminStatsRepository = new AdminStatsRepository();
