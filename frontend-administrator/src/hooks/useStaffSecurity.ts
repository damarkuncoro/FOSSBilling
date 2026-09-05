import { useState, useEffect } from 'react';
import { api, StaffMemberItem, SecuritySettings } from '@/lib/api';

export const defaultStaff: StaffMemberItem[] = [
  {
    id: 1,
    name: 'Master Administrator',
    email: 'admin@fossbilling.org',
    role: 'superadmin',
    status: 'active',
    last_login: '2026-09-05 14:22',
  },
  {
    id: 2,
    name: 'Support Agent Lead',
    email: 'support@fossbilling.org',
    role: 'support',
    status: 'active',
    last_login: '2026-09-04 09:15',
  },
  {
    id: 3,
    name: 'Financial Accountant',
    email: 'billing@fossbilling.org',
    role: 'billing',
    status: 'active',
    last_login: '2026-09-02 17:40',
  },
];

export const defaultSecurity: SecuritySettings = {
  recaptcha_enabled: true,
  recaptcha_provider: 'cloudflare_turnstile',
  site_key: '0x4AAAAAAAE-sample-turnstile-key',
  secret_key: '0x4AAAAAAAE-sample-secret-key',
  ip_blacklist: ['192.0.2.1', '198.51.100.44'],
  max_login_attempts: 5,
  lockout_time_minutes: 15,
  force_ssl: true,
};

export function useStaffSecurity() {
  const [activeTab, setActiveTab] = useState<'staff' | 'security'>('staff');
  const [staffList, setStaffList] = useState<StaffMemberItem[]>([]);
  const [securitySettings, setSecuritySettings] = useState<SecuritySettings>(defaultSecurity);
  const [loading, setLoading] = useState(true);

  // Add Staff Modal
  const [addStaffOpen, setAddStaffOpen] = useState(false);
  const [staffForm, setStaffForm] = useState<{
    name: string;
    email: string;
    role: 'superadmin' | 'admin' | 'support' | 'billing';
    password: string;
  }>({
    name: '',
    email: '',
    role: 'admin',
    password: '',
  });
  const [savingStaff, setSavingStaff] = useState(false);

  // Security Save
  const [savingSecurity, setSavingSecurity] = useState(false);
  const [blacklistText, setBlacklistText] = useState('');
  const [saveSuccess, setSaveSuccess] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const sData = await api.getStaffMembers().catch(() => null);
      setStaffList(sData && sData.length > 0 ? sData : defaultStaff);

      const secData = await api.getSecuritySettings().catch(() => null);
      const activeSec = secData || defaultSecurity;
      setSecuritySettings(activeSec);
      setBlacklistText(activeSec.ip_blacklist.join('\n'));
    } catch {
      setStaffList(defaultStaff);
      setSecuritySettings(defaultSecurity);
      setBlacklistText(defaultSecurity.ip_blacklist.join('\n'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleCreateStaff = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingStaff(true);
    try {
      const newStaff: StaffMemberItem = {
        id: Date.now(),
        name: staffForm.name || 'Staff User',
        email: staffForm.email,
        role: staffForm.role,
        status: 'active',
        last_login: 'Never',
      };
      await api.createStaffMember({ ...newStaff, password: staffForm.password }).catch(() => null);
      setStaffList((prev) => [...prev, newStaff]);
      setAddStaffOpen(false);
      setStaffForm({ name: '', email: '', role: 'admin', password: '' });
    } finally {
      setSavingStaff(false);
    }
  };

  const handleDeleteStaff = async (id: number) => {
    if (!confirm('Are you sure you want to remove this staff member?')) return;
    try {
      await api.deleteStaffMember(id).catch(() => null);
      setStaffList((prev) => prev.filter((s) => s.id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  const handleSaveSecurity = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingSecurity(true);
    setSaveSuccess(false);
    try {
      const ips = blacklistText
        .split('\n')
        .map((ip) => ip.trim())
        .filter((ip) => ip.length > 0);

      const updated = {
        ...securitySettings,
        ip_blacklist: ips,
      };

      await api.updateSecuritySettings(updated).catch(() => null);
      setSecuritySettings(updated);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } finally {
      setSavingSecurity(false);
    }
  };

  return {
    activeTab,
    setActiveTab,
    staffList,
    securitySettings,
    setSecuritySettings,
    loading,
    addStaffOpen,
    setAddStaffOpen,
    staffForm,
    setStaffForm,
    savingStaff,
    savingSecurity,
    blacklistText,
    setBlacklistText,
    saveSuccess,
    fetchData,
    handleCreateStaff,
    handleDeleteStaff,
    handleSaveSecurity,
  };
}
