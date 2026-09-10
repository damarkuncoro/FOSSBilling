import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { SecuritySettings } from '@/types/api';

export function useAntiSpam() {
  const queryClient = useQueryClient();
  const [newIp, setNewIp] = useState('');
  const [localSettings, setLocalSettings] = useState<SecuritySettings | null>(null);

  const { data: serverSettings, isLoading: loading } = useQuery({
    queryKey: ['admin', 'antispam', 'settings'],
    queryFn: () => api.getSecuritySettings(),
  });

  useEffect(() => {
    if (serverSettings) {
      setLocalSettings(serverSettings);
    }
  }, [serverSettings]);

  const saveMutation = useMutation({
    mutationFn: (d: any) => api.updateSecuritySettings(d),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'antispam', 'settings'] });
    },
  });

  const settings = localSettings || {
    recaptcha_enabled: false,
    recaptcha_provider: 'cloudflare_turnstile',
    site_key: '',
    ip_blacklist: [],
    max_login_attempts: 5,
    lockout_time_minutes: 15,
    force_ssl: true,
  } as SecuritySettings;

  return {
    settings,
    setSettings: (s: SecuritySettings) => setLocalSettings(s),
    loading,
    saving: saveMutation.isPending,
    saveMessage: saveMutation.isSuccess ? 'Updated' : null,
    newIp,
    setNewIp,
    handleSave: (e: any) => {
      if (e && e.preventDefault) e.preventDefault();
      saveMutation.mutate(settings);
    },
    handleAddIp: () => {
      if (newIp && localSettings) {
        const updated = { ...localSettings, ip_blacklist: [...localSettings.ip_blacklist, newIp.trim()] };
        setLocalSettings(updated);
        saveMutation.mutate(updated);
      }
      setNewIp('');
    },
    handleRemoveIp: (ip: string) => {
      if (localSettings) {
        const updated = { ...localSettings, ip_blacklist: localSettings.ip_blacklist.filter((i: string) => i !== ip) };
        setLocalSettings(updated);
        saveMutation.mutate(updated);
      }
    },
  };
}
