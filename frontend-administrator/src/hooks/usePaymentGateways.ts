import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminGatewayService } from '@/services/admin_gateway.service';
import type { PaymentGatewayItem, TaxRuleItem } from '@/types/api';

export function usePaymentGateways() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'gateways' | 'tax'>('gateways');
  const [selectedGw, setSelectedGw] = useState<PaymentGatewayItem | null>(null);
  const [editGwOpen, setEditGwOpen] = useState(false);
  const [gwForm, setGwForm] = useState<Partial<PaymentGatewayItem>>({});
  const [taxModalOpen, setTaxModalOpen] = useState(false);
  const [taxForm, setTaxForm] = useState<Partial<TaxRuleItem>>({ name: '', country: 'US', state: '', rate: 10, is_active: true });

  const { data: gateways = [], isLoading: gLoading, refetch: refetchGateways } = useQuery({ queryKey: ['admin', 'gateways'], queryFn: () => adminGatewayService.listPaymentGateways().catch(() => []) });
  const { data: taxRules = [], isLoading: tLoading, refetch: refetchTax } = useQuery({ queryKey: ['admin', 'tax'], queryFn: () => adminGatewayService.listTaxRules().catch(() => []) });

  const gwMutation = useMutation({
    mutationFn: (d: any) => adminGatewayService.updatePaymentGateway(selectedGw!.id, d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'gateways'] }); setEditGwOpen(false); },
  });

  const taxMutation = useMutation({
    mutationFn: (d: any) => adminGatewayService.createTaxRule(d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'tax'] }); setTaxModalOpen(false); },
  });

  const deleteTaxMutation = useMutation({
    mutationFn: (id: number) => adminGatewayService.deleteTaxRule(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'tax'] }),
  });

  const fetchData = async () => {
    await Promise.all([refetchGateways(), refetchTax()]);
  };

  return {
    activeTab, setActiveTab, gateways, taxRules, loading: gLoading || tLoading,
    fetchData,
    selectedGw, editGwOpen, setEditGwOpen, gwForm, setGwForm, savingGw: gwMutation.isPending,
    taxModalOpen, setTaxModalOpen, taxForm, setTaxForm, savingTax: taxMutation.isPending,
    openEditGateway: (gw: any) => { setSelectedGw(gw); setGwForm(gw); setEditGwOpen(true); },
    handleSaveGateway: (e: any) => { e.preventDefault(); gwMutation.mutate(gwForm); },
    toggleGatewayEnabled: (id: string, st: boolean) => { gwMutation.mutate({ enabled: !st }); },
    handleAddTaxRule: (e: any) => { e.preventDefault(); taxMutation.mutate(taxForm); },
    handleDeleteTaxRule: (id: number) => { if (confirm('Delete?')) deleteTaxMutation.mutate(id); },
  };
}
