import { request } from '../lib/api/client';

export interface CurrencyItem { code: string; title: string; conversion_rate: number; format: string; is_default: boolean }

export class AdminCurrencyRepository {
  list = () => request<CurrencyItem[]>('/admin/currencies');
  create = (d: any) => request<CurrencyItem>('/admin/currencies', { method: 'POST', body: JSON.stringify(d) });
  update = (c: string, d: any) => request<CurrencyItem>(`/admin/currencies/${c}`, { method: 'PUT', body: JSON.stringify(d) });
  delete = (c: string) => request(`/admin/currencies/${c}`, { method: 'DELETE' });
  setDefault = (c: string) => request(`/admin/currencies/${c}/default`, { method: 'POST' });
  sync = () => request('/admin/currencies/sync', { method: 'POST' });
}

export const adminCurrencyRepository = new AdminCurrencyRepository();
