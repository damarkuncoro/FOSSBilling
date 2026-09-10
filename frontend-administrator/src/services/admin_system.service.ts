import { adminSystemRepository, AdminSystemRepository } from '../repositories/admin_system.repository';

export class AdminSystemService {
  constructor(private repo: AdminSystemRepository = adminSystemRepository) {}

  getSystemInfo = () => this.repo.getSystemInfo();
  getSystemStatus = () => this.repo.getSystemStatus();
  triggerCron = () => this.repo.triggerCron();
  clearCache = () => this.repo.clearCache();
  listNewsArticles = () => this.repo.listNews();
  createNewsArticle = (title: string, content: string, slug?: string) => {
    if (!title) throw new Error('Article title is required');
    if (!content) throw new Error('Article content is required');
    return this.repo.createNews({ title, content, slug });
  };
  deleteNewsArticle = (id: number) => this.repo.deleteNews(id);
  exportBackup = () => this.repo.exportBackup();
  listAuditLogs = (l?: number, o?: number) => this.repo.getAuditLogs(l, o);
  getActivityTrend = (days: number) => this.repo.getActivityTrend(days);
}

export const adminSystemService = new AdminSystemService();
