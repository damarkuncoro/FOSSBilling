import { useQuery, useMutation } from '@tanstack/react-query';
import { downloadService } from '../services/download.service';

export function useDownloads() {
  const { data: downloads = [], isLoading: loading } = useQuery({
    queryKey: ['client', 'downloads'],
    queryFn: () => downloadService.listDownloads(),
  });

  const downloadMutation = useMutation({
    mutationFn: (id: number) => downloadService.getSecureDownloadUrl(id),
    onSuccess: (url) => window.open(url, '_blank'),
  });

  return {
    downloads, loading,
    downloadingId: downloadMutation.isPending ? -1 : null,
    triggerDownload: (item: any) => item.download_url ? window.open(item.download_url, '_blank') : downloadMutation.mutate(item.id),
  };
}
