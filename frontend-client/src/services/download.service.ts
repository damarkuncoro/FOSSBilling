import { DownloadRepository, downloadRepository, IDownloadRepository } from '../repositories/download.repository';
import { API_BASE } from '../lib/api/client';
import type { DownloadItem } from '@/types/clientModules';

export class DownloadService {
  constructor(private repo: IDownloadRepository = downloadRepository) {}

  async listDownloads(): Promise<DownloadItem[]> {
    return this.repo.listDownloads();
  }

  async getSecureDownloadUrl(orderId: number): Promise<string> {
    const res = await this.repo.generateToken(orderId);
    if (!res?.token) {
      throw new Error('Failed to generate secure download token');
    }
    return `${API_BASE}/client/downloads/${orderId}/file?token=${encodeURIComponent(res.token)}`;
  }
}

export const downloadService = new DownloadService();
