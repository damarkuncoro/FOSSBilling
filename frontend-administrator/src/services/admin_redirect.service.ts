import { AdminRedirectRepositoryInterface } from '../repositories/admin_redirect.repository';

export class AdminRedirectService {
  constructor(private readonly repo: AdminRedirectRepositoryInterface) {}

  async listRedirects(): Promise<any[]> {
    return this.repo.listRedirects();
  }

  async createRedirectRule(path: string, target: string, code = 301): Promise<any> {
    if (!path || !path.trim()) {
      throw new Error('Source path is required');
    }
    if (!target || !target.trim()) {
      throw new Error('Target URL is required');
    }
    return this.repo.createRedirect({ path, target, status_code: code, is_active: true });
  }

  async removeRedirect(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid redirect ID is required');
    }
    return this.repo.deleteRedirect(id);
  }
}
