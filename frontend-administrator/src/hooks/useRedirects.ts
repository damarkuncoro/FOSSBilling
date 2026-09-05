import { useState } from 'react';
import type { UrlRedirect } from '../types/redirects';

const initialRedirects: UrlRedirect[] = [
  {
    id: 1,
    source_path: '/promo-merdeka',
    target_url: '/cart?promo=MERDEKA20',
    status_code: 301,
    is_active: true,
    hits_count: 342,
    created_at: '2026-08-01T10:00:00Z',
    updated_at: '2026-08-17T12:00:00Z',
  },
  {
    id: 2,
    source_path: '/discord',
    target_url: 'https://discord.gg/fossbilling',
    status_code: 302,
    is_active: true,
    hits_count: 1250,
    created_at: '2026-03-15T09:00:00Z',
    updated_at: '2026-09-01T15:00:00Z',
  },
];

export function useRedirects() {
  const [redirects, setRedirects] = useState<UrlRedirect[]>(initialRedirects);
  const [search, setSearch] = useState('');
  const [isAddOpen, setIsAddOpen] = useState(false);

  const createRedirect = (data: { source_path: string; target_url: string; status_code: 301 | 302 }) => {
    let source = data.source_path.trim();
    if (!source.startsWith('/')) source = `/${source}`;

    const newRed: UrlRedirect = {
      id: Date.now(),
      source_path: source,
      target_url: data.target_url.trim(),
      status_code: data.status_code,
      is_active: true,
      hits_count: 0,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    setRedirects((prev) => [newRed, ...prev]);
    setIsAddOpen(false);
  };

  const toggleRedirect = (id: number) => {
    setRedirects((prev) =>
      prev.map((r) => (r.id === id ? { ...r, is_active: !r.is_active, updated_at: new Date().toISOString() } : r))
    );
  };

  const deleteRedirect = (id: number) => {
    setRedirects((prev) => prev.filter((r) => r.id !== id));
  };

  const filteredRedirects = redirects.filter(
    (r) =>
      r.source_path.toLowerCase().includes(search.toLowerCase()) ||
      r.target_url.toLowerCase().includes(search.toLowerCase())
  );

  return {
    redirects: filteredRedirects,
    search,
    setSearch,
    isAddOpen,
    setIsAddOpen,
    createRedirect,
    toggleRedirect,
    deleteRedirect,
  };
}
