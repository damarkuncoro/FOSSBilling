import { useState, useEffect } from 'react';
import { api } from '@/lib/api';

export function useCurrencies() {
  const [currencies, setCurrencies] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({
    code: '',
    title: '',
    conversion_rate: 1.0,
    format: '$ {{price}}',
  });
  const [saving, setSaving] = useState(false);

  const fetchCurrencies = async () => {
    setLoading(true);
    try {
      const data = await api.getCurrencies();
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
      await api.createCurrency({
        code: form.code.toUpperCase(),
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
      await api.setDefaultCurrency(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    }
  };

  const handleDelete = async (code: string) => {
    if (!confirm(`Are you sure you want to delete currency ${code}?`)) return;
    try {
      await api.deleteCurrency(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
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
    fetchCurrencies,
    handleCreate,
    handleSetDefault,
    handleDelete,
  };
}
