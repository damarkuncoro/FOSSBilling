import { AuthRepository, authRepository, IAuthRepository } from '../repositories/auth.repository';
import {
  getStoredClientToken,
  setStoredClientToken,
  removeStoredClientToken,
} from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export class AuthService {
  constructor(private repo: IAuthRepository = authRepository) {}

  isAuthenticated(): boolean {
    return !!getStoredClientToken();
  }

  getToken(): string | null {
    return getStoredClientToken();
  }

  async login(email: string, password: string): Promise<ClientProfile> {
    const res = await this.repo.login(email.trim(), password);
    if (res?.token) {
      setStoredClientToken(res.token);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(res.client));
    }
    return res.client;
  }

  async register(dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }): Promise<ClientProfile> {
    const res = await this.repo.register({
      ...dto,
      email: dto.email.trim().toLowerCase(),
    });
    if (res?.token) {
      setStoredClientToken(res.token);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(res.client));
    }
    return res.client;
  }

  logout(): void {
    removeStoredClientToken();
  }

  async getProfile(): Promise<ClientProfile> {
    const profile = await this.repo.getProfile();
    localStorage.setItem('fossbilling_client_user', JSON.stringify(profile));
    return profile;
  }

  async updateProfile(dto: Partial<ClientProfile>): Promise<ClientProfile> {
    const updated = await this.repo.updateProfile(dto);
    localStorage.setItem('fossbilling_client_user', JSON.stringify(updated));
    return updated;
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<string> {
    if (newPassword.length < 8) {
      throw new Error('New password must be at least 8 characters');
    }
    const res = await this.repo.changePassword(oldPassword, newPassword);
    return res.message;
  }
}

export const authService = new AuthService();
