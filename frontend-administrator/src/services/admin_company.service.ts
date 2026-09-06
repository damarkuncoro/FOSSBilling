import { AdminCompanyRepository, adminCompanyRepository, IAdminCompanyRepository } from '../repositories/admin_company.repository';
import type { CompanySettings } from '@/types/api';

export class AdminCompanyService {
  constructor(private repo: IAdminCompanyRepository = adminCompanyRepository) {}

  async getCompany(): Promise<CompanySettings> {
    return this.repo.getCompany();
  }

  async updateCompany(dto: Partial<CompanySettings>): Promise<CompanySettings> {
    if (dto.email && !dto.email.includes('@')) {
      throw new Error('Valid company email address is required');
    }
    return this.repo.updateCompany(dto);
  }
}

export const adminCompanyService = new AdminCompanyService();
