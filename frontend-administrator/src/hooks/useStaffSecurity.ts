import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useStaffSecurity() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'staff' | 'security'>('staff');
  const [addStaffOpen, setAddStaffOpen] = useState(false);
  const [staffForm, setStaffForm] = useState<any>({ name: '', email: '', role: 'admin', password: '' });
  const [blacklistText, setBlacklistText] = useState('');

  const { data: staffList = [], isLoading: sLoading, refetch: refetchStaff } = useQuery({ queryKey: ['admin', 'staff'], queryFn: () => api.getStaffMembers().catch(() => []) });
  const { data: securitySettings = {
    recaptcha_enabled: false,
    recaptcha_provider: 'cloudflare_turnstile',
    site_key: '',
    ip_blacklist: [],
    max_login_attempts: 5,
    lockout_time_minutes: 15,
    force_ssl: true
  }, isLoading: secLoading, refetch: refetchSec } = useQuery({ queryKey: ['admin', 'security', 'settings'], queryFn: () => api.getSecuritySettings() });

  const staffMutation = useMutation({
    mutationFn: (d: any) => api.createStaffMember(d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'staff'] }); setAddStaffOpen(false); setStaffForm({ name: '', email: '', role: 'admin', password: '' }); },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.deleteStaffMember(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'staff'] }),
  });

  const secMutation = useMutation({
    mutationFn: (d: any) => api.updateSecuritySettings(d),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'security', 'settings'] });
    },
  });

  const fetchData = async () => {
    await Promise.all([refetchStaff(), refetchSec()]);
  };

  return {
    activeTab, setActiveTab, staffList, securitySettings, loading: sLoading || secLoading,
    addStaffOpen, setAddStaffOpen, staffForm, setStaffForm,
    blacklistText, setBlacklistText,
    saveSuccess: secMutation.isSuccess,
    fetchData,
    savingStaff: staffMutation.isPending, savingSecurity: secMutation.isPending,
    handleCreateStaff: (e: any) => { e.preventDefault(); staffMutation.mutate(staffForm); },
    handleDeleteStaff: (id: number) => { if (confirm('Delete?')) deleteMutation.mutate(id); },
    handleSaveSecurity: (d: any) => secMutation.mutate(d),
    setSecuritySettings: (s: any) => { /* Mock for state binding if needed by component */ },
  };
}
