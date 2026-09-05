import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { SecuritySettings } from '@/types/api';

export function useAntiSpam() {
  const [settings, setSettings] = useState<SecuritySettings>({
    recaptcha_enabled: true,
    recaptcha_provider: 'cloudflare_turnstile',
    site_key: '0x4AAAAAAAxMockSiteKey',
    ip_blacklist: ['198.51.100.4', '203.0.113.88'],
    max_login_attempts: 5,
    lockout_time_minutes: 15,
    force_ssl: true,
  });
  const [newIp, setNewIp] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  const fetchSettings = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getSecuritySettings().catch(() => null);
      if (data) setSettings(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setSaveMessage(null);
    try {
      await api.updateSecuritySettings(settings);
      setSaveMessage('Anti-Spam & Security policies successfully updated!');
    } catch (err: any) {
      alert(`Save failed: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  const handleAddIp = () => {
    if (!newIp.trim()) return;
    if (settings.ip_blacklist.includes(newIp.trim())) return;
    setSettings((prev) => ({
      ...prev,
      ip_blacklist: [...prev.ip_blacklist, newIp.trim()],
    }));
    setNewIp('');
  };

  const handleRemoveIp = (ip: string) => {
    setSettings((prev) => ({
      ...prev,
      ip_blacklist: prev.ip_blacklist.filter((item) => item !== ip),
    }));
  };

  return {
    settings,
    setSettings,
    newIp,
    setNewIp,
    loading,
    saving,
    saveMessage,
    handleSave,
    handleAddIp,
    handleRemoveIp,
  };
}
