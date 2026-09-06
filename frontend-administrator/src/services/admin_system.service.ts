import { AdminSystemRepository, adminSystemRepository, IAdminSystemRepository } from '../repositories/admin_system.repository';
import type { SystemInfo, SystemStatusInfo } from '@/types/api';

export class AdminSystemService {
  constructor(private repo: IAdminSystemRepository = adminSystemRepository) {}

  async getSystemInfo(): Promise<SystemInfo> {
    return this.repo.getSystemInfo();
  }

  async getSystemStatus(): Promise<SystemStatusInfo> {
    return this.repo.getSystemStatus();
  }

  async triggerCron(): Promise<{ success: boolean; message: string }> {
    return this.repo.triggerCron();
  }

  async clearCache(): Promise<{ success: boolean; message: string }> {
    return this.repo.clearCache();
  }

  async listCurrencies(): Promise<any[]> {
    return this.repo.listCurrencies();
  }

  async updateCurrencyRate(code: string, rate: number): Promise<any> {
    if (!code) {
      throw new Error('Currency code is required');
    }
    if (rate <= 0) {
      throw new Error('Currency rate must be greater than zero');
    }
    return this.repo.updateCurrency(code.toUpperCase(), rate);
  }

  async listNewsArticles(): Promise<any[]> {
    return this.repo.listNews();
  }

  async createNewsArticle(title: string, content: string, slug?: string): Promise<any> {
    if (!title.trim()) {
      throw new Error('Article title is required');
    }
    if (!content.trim()) {
      throw new Error('Article content is required');
    }
    return this.repo.createNews({
      title: title.trim(),
      content: content.trim(),
      slug: slug?.trim() || undefined,
    });
  }

  async deleteNewsArticle(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid article ID is required');
    }
    return this.repo.deleteNews(id);
  }
}

export const adminSystemService = new AdminSystemService();
