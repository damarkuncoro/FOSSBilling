import { request } from '../lib/api/client';
import type { SoftwareLicense, LicenseValidationLog } from '../types/licenses';

export interface IAdminLicenseRepository {
  getLicenses(): Promise<SoftwareLicense[]>;
  createLicense(dto: Partial<SoftwareLicense>): Promise<SoftwareLicense>;
  updateLicense(id: number, dto: Partial<SoftwareLicense>): Promise<SoftwareLicense>;
  getValidationLogs(): Promise<LicenseValidationLog[]>;
  resetLock(id: number): Promise<any>;
}

export class AdminLicenseRepository implements IAdminLicenseRepository {
  async getLicenses(): Promise<SoftwareLicense[]> {
    return request<SoftwareLicense[]>('/admin/licenses');
  }

  async createLicense(dto: Partial<SoftwareLicense>): Promise<SoftwareLicense> {
    return request<SoftwareLicense>('/admin/licenses', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateLicense(id: number, dto: Partial<SoftwareLicense>): Promise<SoftwareLicense> {
    return request<SoftwareLicense>(`/admin/licenses/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async getValidationLogs(): Promise<LicenseValidationLog[]> {
    return request<LicenseValidationLog[]>('/admin/licenses/logs');
  }

  async resetLock(id: number): Promise<any> {
    return request<any>(`/admin/licenses/${id}/reset`, {
      method: 'POST',
    });
  }
}

export const adminLicenseRepository = new AdminLicenseRepository();
