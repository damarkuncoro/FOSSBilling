import { useState, useEffect, useMemo } from 'react';
import { adminClientService } from '@/services/admin_client.service';
import type { ClientProfile } from '@/types/api';

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
  const [clients, setClients] = useState<ClientProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [search, setSearch] = useState('');

  const fetchClients = async () => {
    setLoading(true);
    try {
      const data = await adminClientService.listClients();
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
      const newClient = await adminClientService.createClient(input as any);
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
      const updated = await adminClientService.updateClient(id, input as any);
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
      await adminClientService.deleteClient(id);
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
    return adminClientService.filterClients(clients, search);
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
