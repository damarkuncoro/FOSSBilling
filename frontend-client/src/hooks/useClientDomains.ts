import { useState, useEffect } from 'react';
import type { DomainRecord } from '../types/clientModules';
import { domainService } from '../services/domain.service';

export function useClientDomains(initial: DomainRecord[] = []) {
  const [domains, setDomains] = useState<DomainRecord[]>(initial);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [checkQuery, setCheckQuery] = useState('');
  const [checkResult, setCheckResult] = useState<{ domain: string; available: boolean; price: number; currency?: string } | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const [editingDomain, setEditingDomain] = useState<DomainRecord | null>(null);

  // Fetch registered domains via DomainService
  const fetchDomains = async () => {
    try {
      setLoading(true);
      const res = await domainService.listClientDomains();
      if (res && Array.isArray(res)) {
        setDomains(res);
      }
    } catch {
      // Retain state on error / offline
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDomains();
  }, []);

  // Live WHOIS availability check via DomainService
  const checkAvailability = async () => {
    if (!checkQuery.trim()) return;
    setIsSearching(true);
    try {
      const res = await domainService.checkAvailability(checkQuery);
      setCheckResult({
        domain: res.domain,
        available: res.available,
        price: res.price,
        currency: res.currency,
      });
    } catch (err) {
      console.error('Failed to check domain availability:', err);
      setCheckResult(null);
    } finally {
      setIsSearching(false);
    }
  };

  const updateNameservers = async (id: number, ns: string[]) => {
    try {
      await domainService.updateNameservers(id, ns);
    } catch {
      // Continue with optimistic UI update
    }
    setDomains((prev) =>
      prev.map((d) => (d.id === id ? { ...d, nameservers: ns } : d))
    );
    setEditingDomain(null);
  };

  const toggleAutoRenew = async (id: number) => {
    try {
      await domainService.toggleAutoRenew(id);
    } catch {
      // Continue with optimistic UI update
    }
    setDomains((prev) =>
      prev.map((d) => (d.id === id ? { ...d, auto_renew: !d.auto_renew } : d))
    );
  };

  const filteredDomains = domains.filter((d) =>
    d.domain_name.toLowerCase().includes(search.toLowerCase())
  );

  return {
    domains: filteredDomains,
    loading,
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
    refreshDomains: fetchDomains,
  };
}
