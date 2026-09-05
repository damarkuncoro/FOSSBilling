import { useState, useEffect } from 'react';
import { api, CompanySettings } from '@/lib/api';

export const initialCompanyState: CompanySettings = {
  name: '',
  email: '',
  phone: '',
  address_1: '',
  address_2: '',
  city: '',
  state: '',
  postcode: '',
  country: '',
  vat_number: '',
  logo_url: '',
  logo_dark_url: '',
  favicon_url: '',
  terms_url: '',
  email_signature: '',
};

export function useCompany() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [form, setForm] = useState<CompanySettings>(initialCompanyState);

  const fetchCompanySettings = async () => {
    setLoading(true);
    setErrorMessage(null);
    try {
      const data = await api.getCompany();
      if (data) {
        setForm({
          name: data.name || '',
          email: data.email || '',
          phone: data.phone || '',
          address_1: data.address_1 || '',
          address_2: data.address_2 || '',
          city: data.city || '',
          state: data.state || '',
          postcode: data.postcode || '',
          country: data.country || '',
          vat_number: data.vat_number || '',
          logo_url: data.logo_url || '',
          logo_dark_url: data.logo_dark_url || '',
          favicon_url: data.favicon_url || '',
          terms_url: data.terms_url || '',
          email_signature: data.email_signature || '',
          updated_at: data.updated_at,
        });
      }
    } catch (err: any) {
      setErrorMessage(err.message || 'Failed to load company settings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCompanySettings();
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setSuccessMessage(null);
    setErrorMessage(null);

    try {
      const res = await api.updateCompany(form);
      if (res) {
        setForm((prev) => ({ ...prev, updated_at: res.updated_at }));
      }
      setSuccessMessage('Company profile and branding settings updated successfully!');
      setTimeout(() => setSuccessMessage(null), 4000);
    } catch (err: any) {
      setErrorMessage(err.message || 'Failed to update company settings');
    } finally {
      setSaving(false);
    }
  };

  return {
    form,
    setForm,
    loading,
    saving,
    successMessage,
    errorMessage,
    fetchCompanySettings,
    handleChange,
    handleSave,
  };
}
