import { AdminAuthRepository, adminAuthRepository, IAdminAuthRepository } from '../repositories/admin_auth.repository';
import { getStoredToken, setStoredToken, removeStoredToken } from '../lib/api/client';
import type { StaffProfile } from '@/types/api';

export class AdminAuthService {
  constructor(private repo: IAdminAuthRepository = adminAuthRepository) {}

  isAuthenticated(): boolean {
    return !!getStoredToken();
  }

  getToken(): string | null {
    return getStoredToken();
  }

  async login(email: string, password: string): Promise<{ token: string; staff: StaffProfile }> {
    const res = await this.repo.login(email.trim(), password);
    if (res?.token) {
      setStoredToken(res.token);
      localStorage.setItem('fossbilling_admin_user', JSON.stringify(res.staff));
    }
    return res;
  }

  logout(): void {
    removeStoredToken();
  }

  async getProfile(): Promise<StaffProfile> {
    return this.repo.getProfile();
  }

  async getAuditLogs(): Promise<any[]> {
    return this.repo.getAuditLogs();
  }
}

export const adminAuthService = new AdminAuthService();
