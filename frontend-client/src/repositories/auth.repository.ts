import { request } from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export class AuthRepository {
  login = (e: string, p: string) => request<{ token: string; client: ClientProfile }>('/guest/auth/login', { method: 'POST', body: JSON.stringify({ email: e, password: p }) });
  register = (d: any) => request<{ token: string; client: ClientProfile }>('/guest/auth/register', { method: 'POST', body: JSON.stringify(d) });
  getProfile = () => request<ClientProfile>('/client/profile');
  updateProfile = (d: any) => request<ClientProfile>('/client/profile', { method: 'PUT', body: JSON.stringify(d) });
  changePassword = (o: string, n: string) => request<{ message: string }>('/client/profile/change-password', { method: 'POST', body: JSON.stringify({ current_password: o, new_password: n }) });
}

export const authRepository = new AuthRepository();
