import { useState, useEffect } from 'react';
import type { DownloadItem } from '../types/clientModules';
import { downloadService } from '../services/download.service';

export function useDownloads(initial: DownloadItem[] = []) {
  const [downloads, setDownloads] = useState<DownloadItem[]>(initial);
  const [loading, setLoading] = useState(false);
  const [downloadingId, setDownloadingId] = useState<number | null>(null);

  const fetchDownloads = async () => {
    try {
      setLoading(true);
      const res = await downloadService.listDownloads();
      if (res && Array.isArray(res)) {
        setDownloads(res);
      } else {
        setDownloads([]);
      }
    } catch {
      // Retain current state if offline or mocked
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDownloads();
  }, []);

  const triggerDownload = async (item: DownloadItem) => {
    setDownloadingId(item.id);
    try {
      let url = item.download_url;
      if (!url) {
        url = await downloadService.getSecureDownloadUrl(item.id);
      }
      setTimeout(() => {
        setDownloadingId(null);
        window.open(url, '_blank');
      }, 600);
    } catch {
      setTimeout(() => {
        setDownloadingId(null);
        if (item.download_url) {
          window.open(item.download_url, '_blank');
        }
      }, 600);
    }
  };

  return {
    downloads,
    loading,
    downloadingId,
    triggerDownload,
    refreshDownloads: fetchDownloads,
    setDownloads,
  };
}

