import { request } from '../lib/api/client';
import type { DashboardStats } from '@/types/api';

export interface IAdminStatsRepository {
  getDashboardStats(): Promise<DashboardStats>;
  getIncomeSummary(): Promise<any>;
}

export class AdminStatsRepository implements IAdminStatsRepository {
  async getDashboardStats(): Promise<DashboardStats> {
    return request<DashboardStats>('/admin/stats/dashboard');
  }

  async getIncomeSummary(): Promise<any> {
    return request<any>('/admin/stats/income');
  }
}

export const adminStatsRepository = new AdminStatsRepository();
