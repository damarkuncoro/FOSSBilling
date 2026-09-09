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
      if (data) {
        setSettings({
          ...data,
          ip_blacklist: Array.isArray(data.ip_blacklist) ? data.ip_blacklist : [],
        });
      }
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
      console.error(err);
      setSaveMessage(null);
      // Fallback if no alert/toast system
      alert(`Save failed: ${err.message || 'Unknown error'}`);
    } finally {
      setSaving(false);
    }
  };

  const handleAddIp = () => {
    if (!newIp.trim()) return;
    const blacklist = Array.isArray(settings.ip_blacklist) ? settings.ip_blacklist : [];
    if (blacklist.includes(newIp.trim())) return;
    setSettings((prev) => ({
      ...prev,
      ip_blacklist: [...blacklist, newIp.trim()],
    }));
    setNewIp('');
  };

  const handleRemoveIp = (ip: string) => {
    const blacklist = Array.isArray(settings.ip_blacklist) ? settings.ip_blacklist : [];
    setSettings((prev) => ({
      ...prev,
      ip_blacklist: blacklist.filter((item) => item !== ip),
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
