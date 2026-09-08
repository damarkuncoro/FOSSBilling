import { useState, useEffect, useCallback } from 'react';
import type { CookieConsentSettings, ConsentLog } from '../types/cookieConsent';
import { systemApi } from '../lib/api/system';

const initialSettings: CookieConsentSettings = {
  is_enabled: true,
  banner_position: 'bottom',
  theme: 'dark',
  message: 'We use cookies to enhance your browsing experience, serve personalized ads or content, and analyze our traffic.',
  accept_button_text: 'Accept All Cookies',
  decline_button_text: 'Essential Only',
  show_decline_button: true,
  privacy_policy_url: '/privacy-policy',
  cookie_expiration_days: 180,
};

export function useCookieConsent() {
  const [settings, setSettings] = useState<CookieConsentSettings>(initialSettings);
  const [logs, setLogs] = useState<ConsentLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [isSaved, setIsSaved] = useState(false);

  const fetchSettings = useCallback(async () => {
    setLoading(true);
    try {
      const data = await systemApi.getCookieConsent();
      if (data && data.config) {
        setSettings(data.config);
      }
      setLogs(data?.logs || []);
    } catch (err) {
      console.error('Failed to fetch cookie consent settings');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const updateSetting = (key: keyof CookieConsentSettings, value: any) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = async () => {
    try {
      await systemApi.updateCookieConsent(settings);
      setIsSaved(true);
      setTimeout(() => setIsSaved(false), 2500);
    } catch (err: any) {
      alert(`Failed to save: ${err.message}`);
    }
  };

  return {
    settings,
    logs,
    loading,
    isSaved,
    updateSetting,
    handleSave,
  };
}
