import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminServerService } from '@/services/admin_server.service';
import type { ServerItem } from '@/types/api';

export const defaultServers: ServerItem[] = [
  { id: 1, name: 'US Node 01', hostname: 'srv1.example.net', ip: '198.51.100.24', manager: 'cpanel', status: 'online', active_accounts: 84, max_accounts: 150, nameserver_1: 'ns1.example.net', nameserver_2: 'ns2.example.net', is_default: true },
];

export const initialServerForm: Partial<ServerItem> & { api_token?: string } = {
  name: '', hostname: '', ip: '', manager: 'cpanel', max_accounts: 100, nameserver_1: 'ns1.example.com', nameserver_2: 'ns2.example.com', is_default: false, api_token: '',
};

export function useServers() {
  const queryClient = useQueryClient();
  const [openModal, setOpenModal] = useState(false);
  const [testResult, setTestResult] = useState<{ id: number; message: string; success: boolean } | null>(null);
  const [form, setForm] = useState(initialServerForm);

  const { data: servers = [], isLoading: loading, refetch: fetchServers } = useQuery({
    queryKey: ['admin', 'servers'],
    queryFn: async () => {
      const data = await adminServerService.listServers();
      return data && data.length > 0 ? data : defaultServers;
    },
  });

  const createMutation = useMutation({
    mutationFn: (input: Partial<ServerItem>) => adminServerService.createServer(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'servers'] });
      setOpenModal(false);
      setForm(initialServerForm);
    },
  });

  const testMutation = useMutation({
    mutationFn: (id: number) => adminServerService.testConnection(id),
    onSuccess: (res, id) => setTestResult({ id, success: res.success, message: res.message }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminServerService.deleteServer(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'servers'] }),
  });

  return {
    servers,
    loading,
    fetchServers,
    openModal,
    setOpenModal,
    testingId: testMutation.isPending ? -1 : null,
    testResult,
    setTestResult,
    saving: createMutation.isPending,
    form,
    setForm,
    handleCreate: (e: React.FormEvent) => { e.preventDefault(); createMutation.mutate(form); },
    handleTestConnection: (id: number) => testMutation.mutate(id),
    handleDelete: (id: number) => { if (confirm('Delete server?')) deleteMutation.mutate(id); },
  };
}
