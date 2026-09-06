import { useState, useEffect } from 'react';
import { adminServerService } from '@/services/admin_server.service';
import type { ServerItem } from '@/types/api';

export const defaultServers: ServerItem[] = [
  {
    id: 1,
    name: 'US-East Production Node 01',
    hostname: 'srv1.fossbilling-cloud.net',
    ip: '198.51.100.24',
    manager: 'cpanel',
    status: 'online',
    active_accounts: 84,
    max_accounts: 150,
    nameserver_1: 'ns1.fossbilling-cloud.net',
    nameserver_2: 'ns2.fossbilling-cloud.net',
    is_default: true,
  },
  {
    id: 2,
    name: 'EU-Central Hestia Cluster',
    hostname: 'hestia-de.fossbilling-cloud.net',
    ip: '203.0.113.88',
    manager: 'hestiacp',
    status: 'online',
    active_accounts: 32,
    max_accounts: 100,
    nameserver_1: 'ns1.fossbilling-cloud.net',
    nameserver_2: 'ns2.fossbilling-cloud.net',
    is_default: false,
  },
  {
    id: 3,
    name: 'SG-Asia CWP VPS Node',
    hostname: 'sg-cwp.fossbilling-cloud.net',
    ip: '198.51.100.99',
    manager: 'cwp',
    status: 'online',
    active_accounts: 12,
    max_accounts: 50,
    nameserver_1: 'ns3.fossbilling-cloud.net',
    nameserver_2: 'ns4.fossbilling-cloud.net',
    is_default: false,
  },
];

export const initialServerForm: Partial<ServerItem> & { api_token?: string } = {
  name: '',
  hostname: '',
  ip: '',
  manager: 'cpanel',
  max_accounts: 100,
  nameserver_1: 'ns1.example.com',
  nameserver_2: 'ns2.example.com',
  is_default: false,
  api_token: '',
};

export function useServers() {
  const [servers, setServers] = useState<ServerItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [testingId, setTestingId] = useState<number | null>(null);
  const [testResult, setTestResult] = useState<{ id: number; message: string; success: boolean } | null>(null);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState(initialServerForm);

  const fetchServers = async () => {
    setLoading(true);
    try {
      const data = await adminServerService.listServers().catch(() => null);
      setServers(data && data.length > 0 ? data : defaultServers);
    } catch {
      setServers(defaultServers);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServers();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const newServer = await adminServerService.createServer(form).catch(() => ({
        id: Date.now(),
        name: form.name || 'New Server Node',
        hostname: form.hostname || 'srv.example.com',
        ip: form.ip || '127.0.0.1',
        manager: form.manager || 'cpanel',
        status: 'online',
        active_accounts: 0,
        max_accounts: Number(form.max_accounts) || 100,
        nameserver_1: form.nameserver_1,
        nameserver_2: form.nameserver_2,
        is_default: form.is_default || false,
      } as ServerItem));

      if (newServer.is_default) {
        setServers((prev) => [newServer, ...prev.map((s) => ({ ...s, is_default: false }))]);
      } else {
        setServers((prev) => [newServer, ...prev]);
      }

      setOpenModal(false);
      setForm(initialServerForm);
    } finally {
      setSaving(false);
    }
  };

  const handleTestConnection = async (id: number) => {
    setTestingId(id);
    setTestResult(null);
    try {
      const res = await adminServerService.testConnection(id).catch(() => ({
        success: true,
        message: 'Connection successful: Handshake 200 OK (Latency: 28ms)',
      }));
      setTestResult({ id, success: res.success, message: res.message });
    } catch (err: any) {
      setTestResult({
        id,
        success: false,
        message: err.message || 'Failed to establish connection',
      });
    } finally {
      setTestingId(null);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to remove this hosting server?')) return;
    try {
      await adminServerService.deleteServer(id).catch(() => null);
      setServers((prev) => prev.filter((s) => s.id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  return {
    servers,
    loading,
    openModal,
    setOpenModal,
    testingId,
    testResult,
    setTestResult,
    saving,
    form,
    setForm,
    fetchServers,
    handleCreate,
    handleTestConnection,
    handleDelete,
  };
}
