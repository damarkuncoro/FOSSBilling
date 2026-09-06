import { request } from '../lib/api/client';
import type { StaffProfile } from '@/types/api';

export interface IAdminAuthRepository {
  login(email: string, password: string): Promise<{ token: string; staff: StaffProfile }>;
  getProfile(): Promise<StaffProfile>;
  getAuditLogs(): Promise<any[]>;
}

export class AdminAuthRepository implements IAdminAuthRepository {
  async login(email: string, password: string): Promise<{ token: string; staff: StaffProfile }> {
    return request<{ token: string; staff: StaffProfile }>('/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async getProfile(): Promise<StaffProfile> {
    return request<StaffProfile>('/admin/auth/profile');
  }

  async getAuditLogs(): Promise<any[]> {
    return request<any[]>('/admin/audit-logs');
  }
}

export const adminAuthRepository = new AdminAuthRepository();
