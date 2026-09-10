import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { invoiceService } from '@/services/invoice.service';
import { useClientAuth } from '@/lib/auth';
import type { Invoice } from '@/types/api';

export function useClientInvoices() {
  const queryClient = useQueryClient();
  const { user, balance, refreshProfile } = useClientAuth();
  const [payModal, setPayModal] = useState<Invoice | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const { data: invoices = [], isLoading: loading, refetch } = useQuery({
    queryKey: ['client', 'invoices'],
    queryFn: () => invoiceService.listClientInvoices(),
  });

  const pBM = useMutation({
    mutationFn: (id: number) => invoiceService.payWithBalance(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'invoices'] }); refreshProfile(); setPayModal(null); setMessage('Invoice paid successfully!'); setTimeout(() => setMessage(null), 5000); },
    onError: (err: any) => alert(err.message),
  });

  const pGM = useMutation({
    mutationFn: ({ id, gateway }: { id: number; gateway: string }) => invoiceService.payWithGateway(id, gateway),
    onSuccess: (res) => { if (res.redirect_url) window.location.href = res.redirect_url; },
    onError: (err: any) => alert(err.message),
  });

  return {
    user, balance, invoices, loading, payModal, setPayModal, message,
    paying: pBM.isPending || pGM.isPending,
    fetchInvoices: () => refetch(),
    handlePayBalance: (id: number) => pBM.mutateAsync(id),
    handlePayGateway: (id: number, gateway = 'midtrans') => pGM.mutateAsync({ id, gateway }),
  };
}
