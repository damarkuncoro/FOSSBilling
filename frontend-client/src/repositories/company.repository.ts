import { request } from '../lib/api/client';
import type { PublicCompanyInfo } from '@/types/api';

export interface ICompanyRepository {
  getCompanyInfo(): Promise<PublicCompanyInfo>;
  listCurrencies(): Promise<any[]>;
}

export class CompanyRepository implements ICompanyRepository {
  async getCompanyInfo(): Promise<PublicCompanyInfo> {
    return request<PublicCompanyInfo>('/guest/company');
  }

  async listCurrencies(): Promise<any[]> {
    return request<any[]>('/guest/currencies');
  }
}

export const companyRepository = new CompanyRepository();
