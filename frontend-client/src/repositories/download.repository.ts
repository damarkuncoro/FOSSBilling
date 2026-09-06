import { request } from '../lib/api/client';
import type { DownloadItem } from '@/types/clientModules';

export interface IDownloadRepository {
  listDownloads(): Promise<DownloadItem[]>;
  generateToken(orderId: number): Promise<{ token: string; expires_at: string }>;
}

export class DownloadRepository implements IDownloadRepository {
  async listDownloads(): Promise<DownloadItem[]> {
    return request<DownloadItem[]>('/client/downloads');
  }

  async generateToken(orderId: number): Promise<{ token: string; expires_at: string }> {
    return request<{ token: string; expires_at: string }>(`/client/downloads/${orderId}/token`, {
      method: 'POST',
    });
  }
}

export const downloadRepository = new DownloadRepository();
