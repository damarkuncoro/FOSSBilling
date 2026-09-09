import { useState, useEffect } from 'react';
import { adminCurrencyService } from '@/services/admin_currency.service';
import type { CurrencyItem } from '@/repositories/admin_currency.repository';

export function useCurrencies() {
  const [currencies, setCurrencies] = useState<CurrencyItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({
    code: '',
    title: '',
    conversion_rate: 1.0,
    format: '$ {{price}}',
  });
  const [saving, setSaving] = useState(false);
  const [syncing, setSyncing] = useState(false);

  const fetchCurrencies = async () => {
    setLoading(true);
    try {
      const data = await adminCurrencyService.listCurrencies();
      setCurrencies(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCurrencies();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await adminCurrencyService.createCurrency({
        code: form.code,
        title: form.title,
        conversion_rate: Number(form.conversion_rate),
        format: form.format,
      });
      setOpenModal(false);
      setForm({ code: '', title: '', conversion_rate: 1.0, format: '$ {{price}}' });
      await fetchCurrencies();
    } finally {
      setSaving(false);
    }
  };

  const handleSetDefault = async (code: string) => {
    try {
      await adminCurrencyService.setDefault(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    }
  };

  const handleDelete = async (code: string) => {
    if (!confirm(`Are you sure you want to delete currency ${code}?`)) return;
    try {
      await adminCurrencyService.deleteCurrency(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    }
  };

  const handleSync = async () => {
    setSyncing(true);
    try {
      await adminCurrencyService.syncRates();
      await fetchCurrencies();
      alert('Currency exchange rates have been updated to latest global market values.');
    } catch (err: any) {
      alert(`Sync failed: ${err.message}`);
    } finally {
      setSyncing(false);
    }
  };

  return {
    currencies,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    syncing,
    fetchCurrencies,
    handleCreate,
    handleSetDefault,
    handleDelete,
    handleSync,
  };
}
