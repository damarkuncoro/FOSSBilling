import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { domainService } from '../services/domain.service';

export function useClientDomains() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');
  const [checkQuery, setCheckQuery] = useState('');
  const [editingDomain, setEditingDomain] = useState<any>(null);

  const { data: domains = [], isLoading: loading } = useQuery({
    queryKey: ['client', 'domains'],
    queryFn: () => domainService.listClientDomains(),
  });

  const checkMutation = useMutation({
    mutationFn: (q: string) => domainService.checkAvailability(q),
  });

  const nsMutation = useMutation({
    mutationFn: ({ id, ns }: { id: number; ns: string[] }) => domainService.updateNameservers(id, ns),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'domains'] }); setEditingDomain(null); },
  });

  const renewMutation = useMutation({
    mutationFn: (id: number) => domainService.toggleAutoRenew(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'domains'] }),
  });

  const filtered = useMemo(() => domains.filter((d: any) => d.domain_name.toLowerCase().includes(search.toLowerCase())), [domains, search]);

  return {
    domains: filtered, loading, search, setSearch, checkQuery, setCheckQuery, editingDomain, setEditingDomain,
    checkResult: checkMutation.data, isSearching: checkMutation.isPending,
    checkAvailability: () => checkQuery && checkMutation.mutate(checkQuery),
    updateNameservers: (id: number, ns: string[]) => nsMutation.mutate({ id, ns }),
    toggleAutoRenew: (id: number) => renewMutation.mutate(id),
    refreshDomains: () => queryClient.invalidateQueries({ queryKey: ['client', 'domains'] }),
  };
}
