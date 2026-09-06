import { request } from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export interface IAuthRepository {
  login(email: string, password: string): Promise<{ token: string; client: ClientProfile }>;
  register(dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }): Promise<{ token: string; client: ClientProfile }>;
  getProfile(): Promise<ClientProfile>;
  updateProfile(dto: Partial<ClientProfile>): Promise<ClientProfile>;
  changePassword(oldPassword: string, newPassword: string): Promise<{ message: string }>;
}

export class AuthRepository implements IAuthRepository {
  async login(email: string, password: string): Promise<{ token: string; client: ClientProfile }> {
    return request<{ token: string; client: ClientProfile }>('/guest/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async register(dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }): Promise<{ token: string; client: ClientProfile }> {
    return request<{ token: string; client: ClientProfile }>('/guest/auth/register', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async getProfile(): Promise<ClientProfile> {
    return request<ClientProfile>('/client/profile');
  }

  async updateProfile(dto: Partial<ClientProfile>): Promise<ClientProfile> {
    return request<ClientProfile>('/client/profile', {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<{ message: string }> {
    return request<{ message: string }>('/client/profile/change-password', {
      method: 'POST',
      body: JSON.stringify({
        old_password: oldPassword,
        new_password: newPassword,
      }),
    });
  }
}

export const authRepository = new AuthRepository();
