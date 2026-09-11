import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useExtensions() {
  const queryClient = useQueryClient();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');

  const { data: extensions = [], isLoading: loading, refetch: fetchExtensions } = useQuery({
    queryKey: ['admin', 'extensions'],
    queryFn: () => api.getExtensions().catch(() => []),
  });

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) => api.toggleExtension(id, enabled),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'extensions'] }),
  });

  const filtered = useMemo(() => {
    return (extensions || []).filter(e => {
      const mS = (e.name || '').toLowerCase().includes(searchQuery.toLowerCase()) || (e.description || '').toLowerCase().includes(searchQuery.toLowerCase());
      const mT = selectedType === 'all' || e.type === selectedType;
      return mS && mT;
    });
  }, [extensions, searchQuery, selectedType]);

  return {
    extensions: filtered || [],
    loading,
    fetchExtensions,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    handleToggle: (id: string, st: boolean) => toggleMutation.mutate({ id, enabled: !st }),
  };
}
