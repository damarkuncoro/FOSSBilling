import { useState, useEffect, useMemo } from 'react';
import { api } from '@/lib/api';

export interface CreateClientInput {
  first_name: string;
  last_name: string;
  email: string;
  password?: string;
  company?: string;
  country?: string;
  currency?: string;
  status?: string;
}

export function useClients() {
  const [clients, setClients] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [search, setSearch] = useState('');

  const fetchClients = async () => {
    setLoading(true);
    try {
      const data = await api.getClients();
      setClients(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const createClient = async (input: CreateClientInput): Promise<boolean> => {
    setSaving(true);
    try {
      const newClient = await api.createClient(input);
      setClients((prev) => [newClient, ...prev]);
      return true;
    } catch (err) {
      console.error('Failed to create client:', err);
      return false;
    } finally {
      setSaving(false);
    }
  };

  const updateClient = async (id: number, input: Partial<CreateClientInput>): Promise<boolean> => {
    setSaving(true);
    try {
      const updated = await api.updateClient(id, input);
      setClients((prev) => prev.map((c) => (c.id === id ? { ...c, ...updated } : c)));
      return true;
    } catch (err) {
      console.error('Failed to update client:', err);
      return false;
    } finally {
      setSaving(false);
    }
  };

  const deleteClient = async (id: number): Promise<boolean> => {
    setSaving(true);
    try {
      await api.deleteClient(id);
      setClients((prev) => prev.filter((c) => c.id !== id));
      return true;
    } catch (err) {
      console.error('Failed to delete client:', err);
      return false;
    } finally {
      setSaving(false);
    }
  };

  useEffect(() => {
    fetchClients();
  }, []);

  const filtered = useMemo(() => {
    return clients.filter(
      (c) =>
        c.email?.toLowerCase().includes(search.toLowerCase()) ||
        c.first_name?.toLowerCase().includes(search.toLowerCase()) ||
        c.last_name?.toLowerCase().includes(search.toLowerCase()) ||
        c.company?.toLowerCase().includes(search.toLowerCase())
    );
  }, [clients, search]);

  return {
    clients,
    loading,
    saving,
    search,
    setSearch,
    filtered,
    fetchClients,
    createClient,
    updateClient,
    deleteClient,
  };
}
