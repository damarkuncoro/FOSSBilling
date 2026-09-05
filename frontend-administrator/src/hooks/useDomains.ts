import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { TldPricingItem, RegistrarConfig } from '@/types/modules';

export function useDomains() {
  const [tlds, setTlds] = useState<TldPricingItem[]>([]);
  const [registrars, setRegistrars] = useState<RegistrarConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [openAddModal, setOpenAddModal] = useState(false);
  const [saving, setSaving] = useState(false);
  const [tldForm, setTldForm] = useState({
    tld: '.com',
    registrar: 'namecheap',
    price_registration: 12.99,
    price_renewal: 14.99,
    price_transfer: 12.99,
    min_years: 1,
    is_active: true,
  });

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [tldData, regData] = await Promise.all([
        api.getTlds().catch(() => [
          { id: 1, tld: '.com', registrar: 'namecheap', price_registration: 12.99, price_renewal: 14.99, price_transfer: 12.99, min_years: 1, is_active: true },
          { id: 2, tld: '.id', registrar: 'custom', price_registration: 18.00, price_renewal: 18.00, price_transfer: 18.00, min_years: 1, is_active: true },
          { id: 3, tld: '.net', registrar: 'namecheap', price_registration: 13.50, price_renewal: 15.50, price_transfer: 13.50, min_years: 1, is_active: true },
          { id: 4, tld: '.org', registrar: 'enom', price_registration: 14.00, price_renewal: 16.00, price_transfer: 14.00, min_years: 1, is_active: true },
        ]),
        api.getRegistrars().catch(() => [
          { id: 'namecheap', name: 'Namecheap API', enabled: true, api_user: 'api_fossbilling', test_mode: false },
          { id: 'enom', name: 'eNom Reseller', enabled: true, api_user: 'reseller_demo', test_mode: true },
          { id: 'resellerclub', name: 'ResellerClub / LogicBoxes', enabled: false, test_mode: true },
          { id: 'custom', name: 'DigitalRegistrar (Pandi .ID)', enabled: true, test_mode: false },
        ]),
      ]);
      setTlds(tldData || []);
      setRegistrars(regData || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCreateTld = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.createTld(tldForm);
      setOpenAddModal(false);
      await fetchData();
    } catch (err: any) {
      alert(`Failed to save TLD: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteTld = async (id: number) => {
    if (!confirm('Are you sure you want to delete this TLD pricing rule?')) return;
    try {
      await api.deleteTld(id);
      await fetchData();
    } catch (err: any) {
      alert(`Failed to delete TLD: ${err.message}`);
    }
  };

  return {
    tlds,
    registrars,
    loading,
    saving,
    openAddModal,
    setOpenAddModal,
    tldForm,
    setTldForm,
    fetchData,
    handleCreateTld,
    handleDeleteTld,
  };
}
