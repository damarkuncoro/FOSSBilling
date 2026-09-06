import { request } from '../lib/api/client';
import type { CompanySettings } from '@/types/api';

export interface IAdminCompanyRepository {
  getCompany(): Promise<CompanySettings>;
  updateCompany(dto: Partial<CompanySettings>): Promise<CompanySettings>;
}

export class AdminCompanyRepository implements IAdminCompanyRepository {
  async getCompany(): Promise<CompanySettings> {
    return request<CompanySettings>('/admin/company');
  }

  async updateCompany(dto: Partial<CompanySettings>): Promise<CompanySettings> {
    return request<CompanySettings>('/admin/company', {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }
}

export const adminCompanyRepository = new AdminCompanyRepository();
