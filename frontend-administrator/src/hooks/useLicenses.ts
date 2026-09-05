import { useState, useEffect } from 'react';
import type { SoftwareLicense, LicenseValidationLog } from '../types/licenses';

const initialLicenses: SoftwareLicense[] = [
  {
    id: 1,
    client_id: 1,
    client_name: 'Ahmad Dhani',
    product_id: 10,
    product_title: 'SaaS Billing Enterprise Edition',
    license_key: 'FOSS-ENT-8842-991A-BC72-9104',
    status: 'active',
    licensed_domain: 'billing.example.com',
    licensed_ip: '198.51.100.24',
    version: '2.5.0',
    max_instances: 5,
    instances_count: 1,
    expires_at: '2027-12-31T23:59:59Z',
    created_at: '2026-01-15T10:00:00Z',
    updated_at: '2026-09-01T08:00:00Z',
  },
  {
    id: 2,
    client_id: 2,
    client_name: 'Budi Santoso',
    product_id: 11,
    product_title: 'Cloud VPN Pro Client',
    license_key: 'FOSS-VPN-3112-78AE-992F-4410',
    status: 'active',
    licensed_domain: 'vpn.santoso.io',
    licensed_ip: '203.0.113.88',
    version: '1.2.0',
    max_instances: 1,
    instances_count: 1,
    expires_at: '2026-11-30T23:59:59Z',
    created_at: '2026-03-10T12:00:00Z',
    updated_at: '2026-08-20T14:30:00Z',
  },
];

const initialLogs: LicenseValidationLog[] = [
  {
    id: 1,
    license_key: 'FOSS-ENT-8842-991A-BC72-9104',
    ip_address: '198.51.100.24',
    domain: 'billing.example.com',
    result: 'valid',
    created_at: '2026-09-05T14:10:00Z',
  },
  {
    id: 2,
    license_key: 'FOSS-VPN-3112-78AE-992F-4410',
    ip_address: '203.0.113.88',
    domain: 'vpn.santoso.io',
    result: 'valid',
    created_at: '2026-09-05T13:50:00Z',
  },
];

export function useLicenses() {
  const [licenses, setLicenses] = useState<SoftwareLicense[]>(initialLicenses);
  const [logs, setLogs] = useState<LicenseValidationLog[]>(initialLogs);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [selectedLicense, setSelectedLicense] = useState<SoftwareLicense | null>(null);

  const generateKey = (prefix = 'FOSS') => {
    const part = () => Math.random().toString(36).substring(2, 6).toUpperCase();
    return `${prefix}-${part()}-${part()}-${part()}-${part()}`;
  };

  const createLicense = (data: Partial<SoftwareLicense>) => {
    const newLic: SoftwareLicense = {
      id: Date.now(),
      client_id: data.client_id || 1,
      client_name: data.client_name || 'Client',
      product_id: data.product_id || 1,
      product_title: data.product_title || 'Software License',
      license_key: data.license_key || generateKey(),
      status: 'active',
      licensed_domain: data.licensed_domain || '',
      licensed_ip: data.licensed_ip || '',
      version: data.version || '1.0.0',
      max_instances: Number(data.max_instances) || 1,
      instances_count: 0,
      expires_at: data.expires_at || '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    setLicenses((prev) => [newLic, ...prev]);
    setIsAddOpen(false);
  };

  const toggleStatus = (id: number) => {
    setLicenses((prev) =>
      prev.map((l) => {
        if (l.id !== id) return l;
        const nextStatus = l.status === 'active' ? 'suspended' : 'active';
        return { ...l, status: nextStatus, updated_at: new Date().toISOString() };
      })
    );
  };

  const reIssueKey = (id: number) => {
    const newKey = generateKey();
    setLicenses((prev) =>
      prev.map((l) => (l.id === id ? { ...l, license_key: newKey, updated_at: new Date().toISOString() } : l))
    );
  };

  const resetLock = (id: number) => {
    setLicenses((prev) =>
      prev.map((l) => (l.id === id ? { ...l, licensed_domain: '', licensed_ip: '', updated_at: new Date().toISOString() } : l))
    );
  };

  const filteredLicenses = licenses.filter((l) => {
    const matchesSearch =
      l.license_key.toLowerCase().includes(search.toLowerCase()) ||
      l.client_name.toLowerCase().includes(search.toLowerCase()) ||
      (l.licensed_domain && l.licensed_domain.toLowerCase().includes(search.toLowerCase()));
    const matchesStatus = statusFilter === 'all' || l.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  return {
    licenses: filteredLicenses,
    logs,
    search,
    setSearch,
    statusFilter,
    setStatusFilter,
    isAddOpen,
    setIsAddOpen,
    selectedLicense,
    setSelectedLicense,
    createLicense,
    toggleStatus,
    reIssueKey,
    resetLock,
    generateKey,
  };
}
