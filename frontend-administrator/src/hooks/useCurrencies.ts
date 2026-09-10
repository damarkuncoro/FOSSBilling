import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminCurrencyService } from '@/services/admin_currency.service';

export function useCurrencies() {
  const queryClient = useQueryClient();
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ code: '', title: '', conversion_rate: 1.0, format: '$ {{price}}' });

  const { data: currencies = [], isLoading: loading, refetch: fetchCurrencies } = useQuery({
    queryKey: ['admin', 'currencies'],
    queryFn: () => adminCurrencyService.listCurrencies(),
  });

  const createMutation = useMutation({
    mutationFn: (d: any) => adminCurrencyService.createCurrency(d),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'currencies'] });
      setOpenModal(false);
      setForm({ code: '', title: '', conversion_rate: 1.0, format: '$ {{price}}' });
    },
  });

  const syncMutation = useMutation({
    mutationFn: () => adminCurrencyService.syncRates(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'currencies'] });
      alert('Rates updated.');
    },
  });

  const defaultMutation = useMutation({
    mutationFn: (c: string) => adminCurrencyService.setDefault(c),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'currencies'] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (c: string) => adminCurrencyService.deleteCurrency(c),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'currencies'] }),
  });

  return {
    currencies,
    loading,
    fetchCurrencies,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving: createMutation.isPending,
    syncing: syncMutation.isPending,
    handleCreate: (e: React.FormEvent) => { e.preventDefault(); createMutation.mutate(form); },
    handleSync: () => syncMutation.mutate(),
    handleSetDefault: (c: string) => defaultMutation.mutate(c),
    handleDelete: (c: string) => { if (confirm('Delete?')) deleteMutation.mutate(c); },
  };
}
