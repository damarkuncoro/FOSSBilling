import { useState } from 'react';
import type { ClientLicense } from '../types/clientModules';

const initialLicenses: ClientLicense[] = [
  {
    id: 1,
    product_title: 'SaaS Billing Enterprise Edition',
    license_key: 'FOSS-ENT-8842-991A-BC72-9104',
    status: 'active',
    licensed_domain: 'billing.example.com',
    licensed_ip: '198.51.100.24',
    max_instances: 5,
    expires_at: '2027-12-31T23:59:59Z',
  },
];

export function useClientLicenses() {
  const [licenses, setLicenses] = useState<ClientLicense[]>(initialLicenses);
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  const copyKey = (key: string) => {
    navigator.clipboard.writeText(key);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  const resetLock = (id: number) => {
    setLicenses((prev) =>
      prev.map((l) => (l.id === id ? { ...l, licensed_domain: '', licensed_ip: '' } : l))
    );
  };

  return {
    licenses,
    copiedKey,
    copyKey,
    resetLock,
  };
}
