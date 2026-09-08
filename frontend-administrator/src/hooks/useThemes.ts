import { useState, useEffect, useCallback } from 'react';
import { request } from '../lib/api/client';

export function useThemes() {
  const [clientThemes, setClientThemes] = useState<any[]>([]);
  const [adminThemes, setAdminThemes] = useState<any[]>([]);
  const [currentClientTheme, setCurrentClientTheme] = useState<any>(null);
  const [currentAdminTheme, setCurrentAdminTheme] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const fetchThemes = useCallback(async () => {
    setLoading(true);
    try {
      const [cList, aList, cCurr, aCurr] = await Promise.all([
        request<any>('/admin/themes?target=client'),
        request<any>('/admin/themes?target=admin'),
        request<any>('/admin/themes/current?target=client'),
        request<any>('/admin/themes/current?target=admin'),
      ]);

      setClientThemes(cList.list || []);
      setAdminThemes(aList.list || []);
      setCurrentClientTheme(cCurr);
      setCurrentAdminTheme(aCurr);
    } catch (err) {
      console.error('Failed to fetch themes');
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

  return {
    clientThemes,
    adminThemes,
    currentClientTheme,
    currentAdminTheme,
    loading,
    fetchThemes,
    handleSelectTheme,
  };
}
