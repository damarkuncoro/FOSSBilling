import { AdminCurrencyRepository, adminCurrencyRepository, IAdminCurrencyRepository, CurrencyItem } from '../repositories/admin_currency.repository';

export class AdminCurrencyService {
  constructor(private repo: IAdminCurrencyRepository = adminCurrencyRepository) {}

  async listCurrencies(): Promise<CurrencyItem[]> {
    return this.repo.getCurrencies();
  }

  async createCurrency(dto: Partial<CurrencyItem>): Promise<CurrencyItem> {
    if (!dto.code || dto.code.trim().length !== 3) {
      throw new Error('A valid 3-letter currency code (e.g. USD, IDR) is required');
    }
    if (!dto.title || !dto.title.trim()) {
      throw new Error('Currency title is required');
    }
    return this.repo.createCurrency({
      code: dto.code.trim().toUpperCase(),
      title: dto.title.trim(),
      conversion_rate: Number(dto.conversion_rate) > 0 ? Number(dto.conversion_rate) : 1.0,
      format: dto.format || '$ {{price}}',
    });
  }

  async setDefault(code: string): Promise<any> {
    if (!code || !code.trim()) {
      throw new Error('Currency code is required');
    }
    return this.repo.setDefaultCurrency(code.trim().toUpperCase());
  }

  async deleteCurrency(code: string): Promise<any> {
    if (!code || !code.trim()) {
      throw new Error('Currency code is required');
    }
    return this.repo.deleteCurrency(code.trim().toUpperCase());
  }
}

export const adminCurrencyService = new AdminCurrencyService();
