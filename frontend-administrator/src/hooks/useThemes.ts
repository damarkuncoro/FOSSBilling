import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { request } from '../lib/api/client';

export function useThemes() {
  const queryClient = useQueryClient();

  const { data: themes = { client: [], admin: [] }, isLoading: loading, refetch: fetchThemes } = useQuery({
    queryKey: ['admin', 'themes'],
    queryFn: async () => {
      const [c, a, cc, ac] = await Promise.all([
        request<any>('/admin/themes?target=client').catch(() => ({ list: [] })),
        request<any>('/admin/themes?target=admin').catch(() => ({ list: [] })),
        request<any>('/admin/themes/current?target=client').catch(() => null),
        request<any>('/admin/themes/current?target=admin').catch(() => null),
      ]);
      return { client: c.list || [], admin: a.list || [], currentClient: cc, currentAdmin: ac };
    },
  });

  const { data: branding = { company_name: 'FOSSBilling', primary_color: '#4f46e5' } } = useQuery({
    queryKey: ['admin', 'settings', 'branding'],
    queryFn: () => request<any>('/admin/settings/branding'),
  });

  const selectMutation = useMutation({
    mutationFn: (args: { code: string; target: string }) =>
      request('/admin/themes/select', { method: 'POST', body: JSON.stringify(args) }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'themes'] }),
  });

  const brandingMutation = useMutation({
    mutationFn: (data: any) => request('/admin/settings/branding', { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'settings', 'branding'] }),
  });

  return {
    clientThemes: (themes as any).client,
    adminThemes: (themes as any).admin,
    currentClientTheme: (themes as any).currentClient,
    currentAdminTheme: (themes as any).currentAdmin,
    branding,
    loading,
    fetchThemes,
    savingBranding: brandingMutation.isPending,
    handleSelectTheme: (code: string, target: string) => selectMutation.mutate({ code, target }),
    handleUpdateBranding: (data: any) => brandingMutation.mutate(data),
    setBranding: (d: any) => { /* Stub */ },
  };
}
