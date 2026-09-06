import { AdminLicenseRepository, adminLicenseRepository, IAdminLicenseRepository } from '../repositories/admin_license.repository';
import type { SoftwareLicense, LicenseValidationLog } from '../types/licenses';

export class AdminLicenseService {
  constructor(private repo: IAdminLicenseRepository = adminLicenseRepository) {}

  async listLicenses(): Promise<SoftwareLicense[]> {
    return this.repo.getLicenses();
  }

  async createLicense(dto: Partial<SoftwareLicense>): Promise<SoftwareLicense> {
    const key = dto.license_key || this.generateKey();
    return this.repo.createLicense({
      ...dto,
      license_key: key,
      status: dto.status || 'active',
      max_instances: Number(dto.max_instances) || 1,
    });
  }

  async updateLicense(id: number, dto: Partial<SoftwareLicense>): Promise<SoftwareLicense> {
    if (!id || id <= 0) {
      throw new Error('Valid license ID is required');
    }
    return this.repo.updateLicense(id, dto);
  }

  async resetLock(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid license ID is required');
    }
    return this.repo.resetLock(id);
  }

  async listValidationLogs(): Promise<LicenseValidationLog[]> {
    return this.repo.getValidationLogs();
  }

  generateKey(prefix = 'FOSS'): string {
    const part = () => Math.random().toString(36).substring(2, 6).toUpperCase();
    return `${prefix}-${part()}-${part()}-${part()}-${part()}`;
  }

  filterLicenses(licenses: SoftwareLicense[], query: string, statusFilter: string): SoftwareLicense[] {
    const q = query.toLowerCase().trim();
    return licenses.filter((l) => {
      const matchesSearch =
        !q ||
        l.license_key.toLowerCase().includes(q) ||
        l.client_name.toLowerCase().includes(q) ||
        (l.licensed_domain && l.licensed_domain.toLowerCase().includes(q));
      const matchesStatus = !statusFilter || statusFilter === 'all' || l.status === statusFilter;
      return matchesSearch && matchesStatus;
    });
  }
}

export const adminLicenseService = new AdminLicenseService();
