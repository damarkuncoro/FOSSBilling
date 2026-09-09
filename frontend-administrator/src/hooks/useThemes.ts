import { useState, useEffect, useCallback } from 'react';
import { request } from '../lib/api/client';

export function useThemes() {
  const [clientThemes, setClientThemes] = useState<any[]>([]);
  const [adminThemes, setAdminThemes] = useState<any[]>([]);
  const [currentClientTheme, setCurrentClientTheme] = useState<any>(null);
  const [currentAdminTheme, setCurrentAdminTheme] = useState<any>(null);
  const [branding, setBranding] = useState<any>({
    company_name: 'FOSSBilling',
    primary_color: '#4f46e5',
  });
  const [loading, setLoading] = useState(true);
  const [savingBranding, setSavingBranding] = useState(false);

  const fetchThemes = useCallback(async () => {
    setLoading(true);
    try {
      const [cList, aList, cCurr, aCurr, brandData] = await Promise.all([
        request<any>('/admin/themes?target=client'),
        request<any>('/admin/themes?target=admin'),
        request<any>('/admin/themes/current?target=client'),
        request<any>('/admin/themes/current?target=admin'),
        request<any>('/admin/settings/branding').catch(() => ({})),
      ]);

      setClientThemes(cList.list || []);
      setAdminThemes(aList.list || []);
      setCurrentClientTheme(cCurr);
      setCurrentAdminTheme(aCurr);
      if (brandData) setBranding(brandData);
    } catch (err) {
      console.error('Failed to fetch themes/branding');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchThemes();
  }, [fetchThemes]);

  const handleSelectTheme = async (code: string, target: 'client' | 'admin') => {
    try {
      await request('/admin/themes/select', {
        method: 'POST',
        body: JSON.stringify({ code, target }),
      });
      fetchThemes();
    } catch (err: any) {
      alert(`Failed to activate theme: ${err.message}`);
    }
  };

  const handleUpdateBranding = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingBranding(true);
    try {
      await request('/admin/settings/branding', {
        method: 'PUT',
        body: JSON.stringify(branding),
      });
      alert('Branding settings updated successfully!');
    } catch (err: any) {
      alert(`Failed to update branding: ${err.message}`);
    } finally {
      setSavingBranding(false);
    }
  };

  return {
    clientThemes,
    adminThemes,
    currentClientTheme,
    currentAdminTheme,
    branding,
    setBranding,
    loading,
    savingBranding,
    fetchThemes,
    handleSelectTheme,
    handleUpdateBranding,
  };
}
