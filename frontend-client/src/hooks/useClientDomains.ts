import { useState } from 'react';
import type { DomainRecord } from '../types/clientModules';

const initialDomains: DomainRecord[] = [
  {
    id: 1,
    domain_name: 'mycompanycloud.com',
    tld: '.com',
    status: 'active',
    nameservers: ['ns1.fossbilling.org', 'ns2.fossbilling.org'],
    epp_code: 'EPP-8891-9921',
    auto_renew: true,
    expires_at: '2027-08-15T00:00:00Z',
  },
  {
    id: 2,
    domain_name: 'startupapps.io',
    tld: '.io',
    status: 'active',
    nameservers: ['ns1.cloudflare.com', 'ns2.cloudflare.com'],
    epp_code: 'EPP-4412-1189',
    auto_renew: false,
    expires_at: '2026-12-01T00:00:00Z',
  },
];

export function useClientDomains() {
  const [domains, setDomains] = useState<DomainRecord[]>(initialDomains);
  const [search, setSearch] = useState('');
  const [checkQuery, setCheckQuery] = useState('');
  const [checkResult, setCheckResult] = useState<{ domain: string; available: boolean; price: number } | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const [editingDomain, setEditingDomain] = useState<DomainRecord | null>(null);

  const checkAvailability = () => {
    if (!checkQuery.trim()) return;
    setIsSearching(true);
    setTimeout(() => {
      const isAvailable = !checkQuery.toLowerCase().includes('google') && !checkQuery.toLowerCase().includes('foss');
      setCheckResult({
        domain: checkQuery.trim().toLowerCase(),
        available: isAvailable,
        price: 12.99,
      });
      setIsSearching(false);
    }, 500);
  };

  const updateNameservers = (id: number, ns: string[]) => {
    setDomains((prev) =>
      prev.map((d) => (d.id === id ? { ...d, nameservers: ns } : d))
    );
    setEditingDomain(null);
  };

  const toggleAutoRenew = (id: number) => {
    setDomains((prev) =>
      prev.map((d) => (d.id === id ? { ...d, auto_renew: !d.auto_renew } : d))
    );
  };

  const filteredDomains = domains.filter((d) =>
    d.domain_name.toLowerCase().includes(search.toLowerCase())
  );

  return {
    domains: filteredDomains,
    search,
    setSearch,
    checkQuery,
    setCheckQuery,
    checkResult,
    isSearching,
    editingDomain,
    setEditingDomain,
    checkAvailability,
    updateNameservers,
    toggleAutoRenew,
  };
}
