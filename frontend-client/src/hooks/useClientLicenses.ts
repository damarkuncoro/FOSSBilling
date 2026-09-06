import { useState, useEffect } from 'react';
import type { ClientLicense } from '../types/clientModules';
import { licenseService } from '../services/license.service';

export function useClientLicenses(initial: ClientLicense[] = []) {
  const [licenses, setLicenses] = useState<ClientLicense[]>(initial);
  const [loading, setLoading] = useState(false);
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  const fetchLicenses = async () => {
    try {
      setLoading(true);
      const res = await licenseService.listClientLicenses();
      if (res && Array.isArray(res)) {
        setLicenses(res);
      } else {
        setLicenses([]);
      }
    } catch {
      // Retain current state if offline or mocked
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLicenses();
  }, []);

  const copyKey = (key: string) => {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(key);
    }
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  const resetLock = async (id: number) => {
    try {
      await licenseService.resetLicenseLock(id);
    } catch {
      // Optimistic fallback
    }
    setLicenses((prev) =>
      prev.map((l) => (l.id === id ? { ...l, licensed_domain: '', licensed_ip: '' } : l))
    );
  };

  return {
    licenses,
    loading,
    copiedKey,
    copyKey,
    resetLock,
    refreshLicenses: fetchLicenses,
    setLicenses,
  };
}

