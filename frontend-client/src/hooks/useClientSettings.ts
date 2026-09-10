import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';

export function useClientSettings() {
  const queryClient = useQueryClient();
  const { user, refreshProfile } = useClientAuth();
  const [profileForm, setProfileForm] = useState({ first_name: user?.first_name || '', last_name: user?.last_name || '', company: user?.company || '', country: user?.country || 'ID' });
  const [keyName, setKeyName] = useState('');
  const [profileMessage, setProfileMessage] = useState<string | null>(null);

  const { data: apiKeys = [] } = useQuery({ queryKey: ['client', 'api-keys'], queryFn: () => api.getApiKeys().catch(() => []) });

  const profileMutation = useMutation({
    mutationFn: (d: any) => api.updateProfile(d),
    onSuccess: () => { refreshProfile(); setProfileMessage('Saved'); },
  });

  const keyMutation = useMutation({
    mutationFn: (n: string) => api.generateApiKey(n),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'api-keys'] }); setKeyName(''); },
  });

  const revokeMutation = useMutation({
    mutationFn: (id: number) => api.revokeApiKey(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'api-keys'] }),
  });

  return {
    user, profileForm, setProfileForm, apiKeys, keyName, setKeyName, profileMessage,
    savingProfile: profileMutation.isPending, generatingKey: keyMutation.isPending,
    handleUpdateProfile: (e: any) => { e.preventDefault(); setProfileMessage(null); profileMutation.mutate(profileForm); },
    handleGenerateKey: (e: any) => { e.preventDefault(); if (keyName) keyMutation.mutate(keyName); },
    handleRevokeKey: (id: number) => { if (confirm('Revoke?')) revokeMutation.mutate(id); },
  };
}
