import { LicenseRepository, licenseRepository, ILicenseRepository } from '../repositories/license.repository';
import type { ClientLicense } from '@/types/clientModules';

export class LicenseService {
  constructor(private repo: ILicenseRepository = licenseRepository) {}

  async listClientLicenses(): Promise<ClientLicense[]> {
    return this.repo.listLicenses();
  }

  async resetLicenseLock(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid license ID is required');
    }
    return this.repo.resetLock(id);
  }
}

export const licenseService = new LicenseService();
