import { AdminSeoRepositoryInterface } from '../repositories/admin_seo.repository';

export class AdminSeoService {
  constructor(private readonly repo: AdminSeoRepositoryInterface) {}

  async getSeoConfig(): Promise<any> {
    return this.repo.getSettings();
  }

  async updateSeoConfig(config: any): Promise<any> {
    if (!config) {
      throw new Error('SEO configuration payload is required');
    }
    return this.repo.updateSettings(config);
  }
}
