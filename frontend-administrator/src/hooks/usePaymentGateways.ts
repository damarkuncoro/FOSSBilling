import { useState, useEffect } from 'react';
import { api, PaymentGatewayItem, TaxRuleItem } from '@/lib/api';

export const defaultGateways: PaymentGatewayItem[] = [
  {
    id: 'stripe',
    name: 'Stripe Payments',
    description: 'Credit Cards, Apple Pay, Google Pay, and localized European methods',
    type: 'card',
    enabled: true,
    test_mode: true,
    public_key: 'pk_test_51Mz00000000000000000',
    secret_key: 'sk_test_51Mz00000000000000000',
  },
  {
    id: 'paypal',
    name: 'PayPal Commerce Platform',
    description: 'Standard PayPal wallet, Pay Later, and Debit/Credit Card gateway',
    type: 'wallet',
    enabled: true,
    test_mode: true,
    public_key: 'sb-client-id-sample-000',
  },
  {
    id: 'midtrans',
    name: 'Midtrans SNAP Gateway',
    description: 'Indonesian payment methods: QRIS, GoPay, BCA/Mandiri/BRI Virtual Accounts',
    type: 'wallet',
    enabled: true,
    test_mode: true,
    public_key: 'SB-Mid-client-000000',
  },
  {
    id: 'bank_transfer',
    name: 'Manual Bank Transfer / Wire',
    description: 'Offline payment option with manual receipt approval in Invoices menu',
    type: 'bank_transfer',
    enabled: true,
    test_mode: false,
    instructions: 'Transfer to BCA: 1234567890 a/n PT FOSSBilling. Confirm via ticket.',
  },
];

export const defaultTaxRules: TaxRuleItem[] = [
  { id: 1, name: 'Indonesia PPN', country: 'ID', rate: 11.0, is_active: true, apply_to_all_clients: true },
  { id: 2, name: 'UK Standard VAT', country: 'GB', rate: 20.0, is_active: true, apply_to_all_clients: false },
  { id: 3, name: 'EU Standard VAT', country: 'EU', rate: 19.0, is_active: false, apply_to_all_clients: false },
];

export function usePaymentGateways() {
  const [activeTab, setActiveTab] = useState<'gateways' | 'tax'>('gateways');
  const [gateways, setGateways] = useState<PaymentGatewayItem[]>([]);
  const [taxRules, setTaxRules] = useState<TaxRuleItem[]>([]);
  const [loading, setLoading] = useState(true);

  const [selectedGw, setSelectedGw] = useState<PaymentGatewayItem | null>(null);
  const [editGwOpen, setEditGwOpen] = useState(false);
  const [gwForm, setGwForm] = useState<Partial<PaymentGatewayItem>>({});
  const [savingGw, setSavingGw] = useState(false);

  const [taxModalOpen, setTaxModalOpen] = useState(false);
  const [taxForm, setTaxForm] = useState<Partial<TaxRuleItem>>({
    name: '',
    country: 'US',
    state: '',
    rate: 10,
    is_active: true,
    apply_to_all_clients: false,
  });
  const [savingTax, setSavingTax] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const gwData = await api.getPaymentGateways().catch(() => null);
      setGateways(gwData && gwData.length > 0 ? gwData : defaultGateways);

      const taxData = await api.getTaxRules().catch(() => null);
      setTaxRules(taxData && taxData.length > 0 ? taxData : defaultTaxRules);
    } catch {
      setGateways(defaultGateways);
      setTaxRules(defaultTaxRules);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openEditGateway = (gw: PaymentGatewayItem) => {
    setSelectedGw(gw);
    setGwForm({ ...gw });
    setEditGwOpen(true);
  };

  const handleSaveGateway = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedGw) return;
    setSavingGw(true);
    try {
      await api.updatePaymentGateway(selectedGw.id, gwForm).catch(() => null);
      setGateways((prev) =>
        prev.map((g) => (g.id === selectedGw.id ? ({ ...g, ...gwForm } as PaymentGatewayItem) : g))
      );
      setEditGwOpen(false);
    } finally {
      setSavingGw(false);
    }
  };

  const toggleGatewayEnabled = (id: string) => {
    setGateways((prev) =>
      prev.map((g) => (g.id === id ? { ...g, enabled: !g.enabled } : g))
    );
  };

  const handleAddTaxRule = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingTax(true);
    try {
      const newRule: TaxRuleItem = {
        id: Date.now(),
        name: taxForm.name || 'New Tax Rule',
        country: taxForm.country || 'GLOBAL',
        state: taxForm.state || '',
        rate: Number(taxForm.rate) || 0,
        is_active: taxForm.is_active ?? true,
        apply_to_all_clients: taxForm.apply_to_all_clients ?? false,
      };
      await api.createTaxRule(newRule).catch(() => null);
      setTaxRules((prev) => [newRule, ...prev]);
      setTaxModalOpen(false);
      setTaxForm({ name: '', country: 'US', state: '', rate: 10, is_active: true, apply_to_all_clients: false });
    } finally {
      setSavingTax(false);
    }
  };

  const handleDeleteTaxRule = async (id: number) => {
    if (!confirm('Are you sure you want to delete this tax rule?')) return;
    try {
      await api.deleteTaxRule(id).catch(() => null);
      setTaxRules((prev) => prev.filter((r) => r.id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  return {
    activeTab,
    setActiveTab,
    gateways,
    taxRules,
    loading,
    selectedGw,
    editGwOpen,
    setEditGwOpen,
    gwForm,
    setGwForm,
    savingGw,
    taxModalOpen,
    setTaxModalOpen,
    taxForm,
    setTaxForm,
    savingTax,
    fetchData,
    openEditGateway,
    handleSaveGateway,
    toggleGatewayEnabled,
    handleAddTaxRule,
    handleDeleteTaxRule,
  };
}
