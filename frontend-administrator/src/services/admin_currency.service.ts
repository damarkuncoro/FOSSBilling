import { adminCurrencyRepository, AdminCurrencyRepository } from '../repositories/admin_currency.repository';

export class AdminCurrencyService {
  constructor(private repo: AdminCurrencyRepository = adminCurrencyRepository) {}

  listCurrencies = () => this.repo.list();
  createCurrency = (d: any) => {
    if (!d.code || d.code.length !== 3) throw new Error('A valid 3-letter currency code (e.g. USD, IDR) is required');
    return this.repo.create({ ...d, code: d.code.toUpperCase() });
  };
  updateCurrency = (c: string, d: any) => this.repo.update(c, d);
  deleteCurrency = (c: string) => this.repo.delete(c);
  setDefault = (c: string) => this.repo.setDefault(c);
  syncRates = () => this.repo.sync();
}

export const adminCurrencyService = new AdminCurrencyService();
