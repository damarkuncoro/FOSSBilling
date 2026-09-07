import { AdminThemeRepositoryInterface } from '../repositories/admin_theme.repository';

export class AdminThemeService {
  constructor(private readonly repo: AdminThemeRepositoryInterface) {}

  async listThemes(): Promise<any[]> {
    return this.repo.listThemes();
  }

  async getActiveTheme(): Promise<any> {
    return this.repo.getActiveTheme();
  }

  async selectTheme(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Theme ID is required');
    }
    return this.repo.setActiveTheme(id);
  }

  async saveSettings(id: string, settings: any): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Theme ID is required');
    }
    return this.repo.updateThemeSettings(id, settings);
  }
}
