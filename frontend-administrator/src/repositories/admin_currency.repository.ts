import { request } from '../lib/api/client';

export interface CurrencyItem {
  code: string;
  title: string;
  conversion_rate: number;
  format: string;
  is_default?: boolean;
}

export interface IAdminCurrencyRepository {
  getCurrencies(): Promise<CurrencyItem[]>;
  createCurrency(dto: Partial<CurrencyItem>): Promise<CurrencyItem>;
  setDefaultCurrency(code: string): Promise<any>;
  deleteCurrency(code: string): Promise<any>;
}

export class AdminCurrencyRepository implements IAdminCurrencyRepository {
  async getCurrencies(): Promise<CurrencyItem[]> {
    return request<CurrencyItem[]>('/admin/currencies');
  }

  async createCurrency(dto: Partial<CurrencyItem>): Promise<CurrencyItem> {
    return request<CurrencyItem>('/admin/currencies', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async setDefaultCurrency(code: string): Promise<any> {
    return request<any>(`/admin/currencies/${code}/default`, {
      method: 'PUT',
    });
  }

  async deleteCurrency(code: string): Promise<any> {
    return request<any>(`/admin/currencies/${code}`, {
      method: 'DELETE',
    });
  }
}

export const adminCurrencyRepository = new AdminCurrencyRepository();
