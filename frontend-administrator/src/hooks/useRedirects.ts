import { useState, useEffect, useCallback } from 'react';
import type { UrlRedirect } from '../types/redirects';
import { systemApi } from '../lib/api/system';

export function useRedirects() {
  const [redirects, setRedirects] = useState<UrlRedirect[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [isAddOpen, setIsAddOpen] = useState(false);

  const fetchRedirects = useCallback(async () => {
    setLoading(true);
    try {
      const data = await systemApi.getRedirects();
      // Map backend fields to frontend types
      const mapped = (data || []).map((r: any) => ({
        id: r.id,
        source_path: r.path,
        target_url: r.target,
        status_code: r.status_code,
        is_active: r.is_enabled,
        hits_count: r.hit_count,
        created_at: r.created_at,
        updated_at: r.updated_at,
      }));
      setRedirects(mapped);
    } catch (err) {
      console.error('Failed to fetch redirects');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchRedirects();
  }, [fetchRedirects]);

  const createRedirect = async (data: { source_path: string; target_url: string; status_code: number }) => {
    try {
      await systemApi.createRedirect({
        path: data.source_path,
        target: data.target_url,
        status_code: data.status_code,
        is_enabled: true,
      });
      fetchRedirects();
      setIsAddOpen(false);
    } catch (err: any) {
      alert(`Failed to create redirect: ${err.message}`);
    }
  };

  const toggleRedirect = async (id: number) => {
    const item = redirects.find(r => r.id === id);
    if (!item) return;

    try {
      await systemApi.updateRedirect(id, { is_enabled: !item.is_active });
      fetchRedirects();
    } catch (err) {
      console.error('Failed to toggle redirect');
    }
  };

  const deleteRedirect = async (id: number) => {
    if (!confirm('Are you sure you want to delete this redirect?')) return;
    try {
      await systemApi.deleteRedirect(id);
      fetchRedirects();
    } catch (err) {
      console.error('Failed to delete redirect');
    }
  };

  const filteredRedirects = redirects.filter(
    (r) =>
      r.source_path.toLowerCase().includes(search.toLowerCase()) ||
      r.target_url.toLowerCase().includes(search.toLowerCase())
  );

  return {
    redirects: filteredRedirects,
    loading,
    search,
    setSearch,
    isAddOpen,
    setIsAddOpen,
    createRedirect,
    toggleRedirect,
    deleteRedirect,
    refresh: fetchRedirects,
  };
}
