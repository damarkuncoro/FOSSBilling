import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminMassMailService } from '@/services/admin_massmail.service';
import { api } from '@/lib/api';

export function useMassMail() {
  const queryClient = useQueryClient();
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ subject: '', content: '' });

  const { data: campaigns = [], isLoading: loading, refetch: fetchCampaigns } = useQuery({
    queryKey: ['admin', 'mass-mail'],
    queryFn: () => api.getMassMailCampaigns(),
  });

  const createMutation = useMutation({
    mutationFn: (d: any) => api.createMassMailCampaign(d),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'mass-mail'] });
      setOpenModal(false);
      setForm({ subject: '', content: '' });
    },
  });

  const sendMutation = useMutation({
    mutationFn: (id: number) => api.sendMassMailCampaign(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'mass-mail'] }),
  });

  return {
    campaigns,
    loading,
    fetchCampaigns,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving: createMutation.isPending,
    sendingId: sendMutation.isPending ? -1 : null,
    handleCreate: (e: React.FormEvent) => { e.preventDefault(); createMutation.mutate(form); },
    handleSend: (id: number) => sendMutation.mutate(id),
  };
}
