import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { ApiKey } from '@/types/api';

export function useClientSettings() {
  const { user, refreshProfile } = useClientAuth();
  const [profileForm, setProfileForm] = useState({
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    company: user?.company || '',
    country: user?.country || 'ID',
  });
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
  const [keyName, setKeyName] = useState('');
  const [savingProfile, setSavingProfile] = useState(false);
  const [generatingKey, setGeneratingKey] = useState(false);
  const [profileMessage, setProfileMessage] = useState<string | null>(null);

  const fetchKeys = useCallback(async () => {
    try {
      const data = await api.getApiKeys();
      setApiKeys(data || []);
    } catch {
      // ignore
    }
  }, []);

  useEffect(() => {
    fetchKeys();
  }, [fetchKeys]);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingProfile(true);
    setProfileMessage(null);
    try {
      await api.updateProfile(profileForm);
      await refreshProfile();
      setProfileMessage('Your profile has been saved successfully!');
    } catch (err: any) {
      alert(`Update failed: ${err.message}`);
    } finally {
      setSavingProfile(false);
    }
  };

  const handleGenerateKey = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!keyName.trim()) return;

    setGeneratingKey(true);
    try {
      await api.generateApiKey(keyName.trim());
      setKeyName('');
      await fetchKeys();
    } catch (err: any) {
      alert(`Failed to generate key: ${err.message}`);
    } finally {
      setGeneratingKey(false);
    }
  };

  const handleRevokeKey = async (id: number) => {
    if (!confirm('Are you sure you want to revoke this API key?')) return;
    try {
      await api.revokeApiKey(id);
      await fetchKeys();
    } catch (err: any) {
      alert(`Failed to revoke key: ${err.message}`);
    }
  };

  return {
    user,
    profileForm,
    setProfileForm,
    apiKeys,
    keyName,
    setKeyName,
    savingProfile,
    generatingKey,
    profileMessage,
    handleUpdateProfile,
    handleGenerateKey,
    handleRevokeKey,
  };
}
