import { authRepository as repo } from '../repositories/auth.repository';
import { getStoredClientToken, setStoredClientToken, removeStoredClientToken } from '../lib/api/client';
import type { ClientProfile } from '@/types/api';

export class AuthService {
  isAuthenticated = () => !!getStoredClientToken();
  getToken = () => getStoredClientToken();

  login = async (e: string, p: string) => {
    const r = await repo.login(e.trim(), p);
    if (r?.token) { setStoredClientToken(r.token); localStorage.setItem('fossbilling_client_user', JSON.stringify(r.client)); }
    return r.client;
  };

  register = async (d: any) => {
    const r = await repo.register({ ...d, email: d.email.trim().toLowerCase() });
    if (r?.token) { setStoredClientToken(r.token); localStorage.setItem('fossbilling_client_user', JSON.stringify(r.client)); }
    return r.client;
  };

  logout = () => removeStoredClientToken();
  getProfile = async () => { const p = await repo.getProfile(); localStorage.setItem('fossbilling_client_user', JSON.stringify(p)); return p; };
  updateProfile = async (d: any) => { const p = await repo.updateProfile(d); localStorage.setItem('fossbilling_client_user', JSON.stringify(p)); return p; };
  changePassword = async (o: string, n: string) => { if (n.length < 8) throw new Error('Short'); const r = await repo.changePassword(o, n); return r.message; };
}

export const authService = new AuthService();
