import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { systemApi } from '../lib/api/system';

export function useRedirects() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');
  const [isAddOpen, setIsAddOpen] = useState(false);

  const { data: redirects = [], isLoading: loading, refetch: refresh } = useQuery({
    queryKey: ['admin', 'redirects'],
    queryFn: async () => {
      const d = await systemApi.getRedirects();
      return (d || []).map((r: any) => ({
        id: r.id,
        source_path: r.path,
        target_url: r.target,
        status_code: r.status_code,
        is_active: r.is_enabled,
        hits_count: r.hit_count,
        created_at: r.created_at,
        updated_at: r.updated_at || r.created_at
      }));
    },
  });

  const createMutation = useMutation({
    mutationFn: (d: any) => systemApi.createRedirect({ path: d.source_path, target: d.target_url, status_code: d.status_code, is_enabled: true }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'redirects'] }); setIsAddOpen(false); },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, active }: any) => systemApi.updateRedirect(id, { is_enabled: active }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'redirects'] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => systemApi.deleteRedirect(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'redirects'] }),
  });

  const filtered = useMemo(() => (redirects || []).filter(r => r.source_path.includes(search) || r.target_url.includes(search)), [redirects, search]);

  return {
    redirects: filtered, loading, search, setSearch, isAddOpen, setIsAddOpen,
    refresh,
    createRedirect: createMutation.mutate,
    toggleRedirect: (id: number) => { const r = redirects.find(x => x.id === id); if (r) updateMutation.mutate({ id, active: !r.is_active }); },
    deleteRedirect: (id: number) => { if (confirm('Delete?')) deleteMutation.mutate(id); },
  };
}
