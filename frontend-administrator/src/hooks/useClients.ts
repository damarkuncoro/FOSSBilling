import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminClientService } from '@/services/admin_client.service';
import type { ClientProfile } from '@/types/api';

export interface CreateClientInput {
  first_name: string;
  last_name: string;
  email: string;
  password?: string;
  company?: string;
  country?: string;
  currency?: string;
  status?: string;
}

export function useClients() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');

  const { data: clients = [], isLoading: loading, refetch: fetchClients } = useQuery({
    queryKey: ['admin', 'clients'],
    queryFn: () => adminClientService.listClients(),
  });

  const createMutation = useMutation({
    mutationFn: (input: CreateClientInput) => adminClientService.createClient(input as any),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] }),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<CreateClientInput> }) =>
      adminClientService.updateClient(id, input as any),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminClientService.deleteClient(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] }),
  });

  const impersonateClient = async (id: number) => {
    try {
      const token = await adminClientService.impersonateClient(id);
      window.open(`/login?token=${token}`, '_blank');
    } catch (err) {
      console.error('Impersonation failed:', err);
      alert('Could not impersonate client.');
    }
  };

  const filtered = useMemo(() => {
    return adminClientService.filterClients(clients, search);
  }, [clients, search]);

  return {
    clients,
    loading,
    saving: createMutation.isPending || updateMutation.isPending || deleteMutation.isPending,
    search,
    setSearch,
    filtered,
    fetchClients,
    createClient: async (input: CreateClientInput) => {
      try {
        await createMutation.mutateAsync(input as any);
        return true;
      } catch {
        return false;
      }
    },
    updateClient: async (id: number, input: Partial<CreateClientInput>) => {
      try {
        await updateMutation.mutateAsync({ id, input: input as any });
        return true;
      } catch {
        return false;
      }
    },
    deleteClient: deleteMutation.mutateAsync,
    impersonateClient,
  };
}
