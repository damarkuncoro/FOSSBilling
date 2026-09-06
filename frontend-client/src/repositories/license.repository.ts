import { request } from '../lib/api/client';
import type { ClientLicense } from '@/types/clientModules';

export interface ILicenseRepository {
  listLicenses(): Promise<ClientLicense[]>;
  resetLock(id: number): Promise<any>;
}

export class LicenseRepository implements ILicenseRepository {
  async listLicenses(): Promise<ClientLicense[]> {
    return request<ClientLicense[]>('/client/licenses');
  }

  async resetLock(id: number): Promise<any> {
    return request(`/client/licenses/${id}/reset`, {
      method: 'POST',
    });
  }
}

export const licenseRepository = new LicenseRepository();
