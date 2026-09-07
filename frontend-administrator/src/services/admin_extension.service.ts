import { AdminExtensionRepositoryInterface } from '../repositories/admin_extension.repository';

export class AdminExtensionService {
  constructor(private readonly repo: AdminExtensionRepositoryInterface) {}

  async listInstalledExtensions(): Promise<any[]> {
    return this.repo.listExtensions();
  }

  async listMarketplaceExtensions(): Promise<any[]> {
    return this.repo.listMarketplace();
  }

  async activate(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.activateExtension(id);
  }

  async deactivate(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.deactivateExtension(id);
  }

  async install(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.installExtension(id);
  }

  async uninstall(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.uninstallExtension(id);
  }

  async getConfiguration(id: string): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.getConfig(id);
  }

  async saveConfiguration(id: string, config: any): Promise<any> {
    if (!id || !id.trim()) {
      throw new Error('Extension ID is required');
    }
    return this.repo.updateConfig(id, config);
  }

  filterByType(extensions: any[], type: string): any[] {
    if (!type || type === 'all') return extensions;
    return extensions.filter((ext) => ext.type === type);
  }
}
