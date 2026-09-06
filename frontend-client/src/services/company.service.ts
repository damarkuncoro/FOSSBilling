import { CompanyRepository, companyRepository, ICompanyRepository } from '../repositories/company.repository';
import type { PublicCompanyInfo } from '@/types/api';

export class CompanyService {
  constructor(private repo: ICompanyRepository = companyRepository) {}

  async getCompanyInfo(): Promise<PublicCompanyInfo> {
    return this.repo.getCompanyInfo();
  }

  async listActiveCurrencies(): Promise<any[]> {
    return this.repo.listCurrencies();
  }
}

export const companyService = new CompanyService();
