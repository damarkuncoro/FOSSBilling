import { useState } from 'react';
import type { DownloadItem } from '../types/clientModules';

const initialDownloads: DownloadItem[] = [
  {
    id: 1,
    title: 'FOSSBilling Client Windows Desktop Tool',
    category: 'Client Software',
    version: '1.4.2',
    file_size: '24.8 MB',
    description: 'Manage instances, check live server metrics, and SSH into VPS with one click.',
    requires_active_service: false,
    download_url: 'https://cdn.fossbilling.org/downloads/fossbilling-desktop-v1.4.2.exe',
    updated_at: '2026-08-15T10:00:00Z',
  },
  {
    id: 2,
    title: 'Enterprise Billing Integration Module for WHMCS / cPanel',
    category: 'Addon Modules',
    version: '2.1.0',
    file_size: '4.2 MB',
    description: 'Drop-in extension for automated billing synchronization.',
    requires_active_service: true,
    download_url: 'https://cdn.fossbilling.org/downloads/enterprise-addon-2.1.0.zip',
    updated_at: '2026-09-01T12:00:00Z',
  },
];

export function useDownloads() {
  const [downloads] = useState<DownloadItem[]>(initialDownloads);
  const [downloadingId, setDownloadingId] = useState<number | null>(null);

  const triggerDownload = (item: DownloadItem) => {
    setDownloadingId(item.id);
    setTimeout(() => {
      setDownloadingId(null);
      window.open(item.download_url, '_blank');
    }, 600);
  };

  return {
    downloads,
    downloadingId,
    triggerDownload,
  };
}
