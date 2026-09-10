import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { licenseService } from '../services/license.service';

export function useClientLicenses() {
  const queryClient = useQueryClient();
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  const { data: licenses = [], isLoading: loading } = useQuery({
    queryKey: ['client', 'licenses'],
    queryFn: () => licenseService.listClientLicenses(),
  });

  const resetMutation = useMutation({
    mutationFn: (id: number) => licenseService.resetLicenseLock(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'licenses'] }),
  });

  return {
    licenses, loading, copiedKey,
    copyKey: (k: string) => { navigator.clipboard?.writeText(k); setCopiedKey(k); setTimeout(() => setCopiedKey(null), 2000); },
    resetLock: resetMutation.mutate,
    refreshLicenses: () => queryClient.invalidateQueries({ queryKey: ['client', 'licenses'] }),
  };
}
