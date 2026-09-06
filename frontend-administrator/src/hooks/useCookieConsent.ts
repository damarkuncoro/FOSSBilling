import { useState } from 'react';
import type { CookieConsentSettings, ConsentLog } from '../types/cookieConsent';

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

const initialLogs: ConsentLog[] = [
  { id: 1, ip_address: '103.24.112.5', country: 'ID', decision: 'accepted', created_at: '2026-09-05T14:15:00Z' },
  { id: 2, ip_address: '198.51.100.44', country: 'US', decision: 'accepted', created_at: '2026-09-05T13:40:00Z' },
  { id: 3, ip_address: '213.180.204.8', country: 'DE', decision: 'declined', created_at: '2026-09-05T12:10:00Z' },
];

const STORAGE_KEY = 'fossbilling_cookie_consent_settings';

export function useCookieConsent() {
  const [settings, setSettings] = useState<CookieConsentSettings>(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) return JSON.parse(saved);
    } catch (e) {
      console.error('Failed to parse cookie consent settings', e);
    }
    return initialSettings;
  });
  const [logs] = useState<ConsentLog[]>(initialLogs);
  const [isSaved, setIsSaved] = useState(false);

  const updateSetting = (key: keyof CookieConsentSettings, value: any) => {
    setSettings((prev) => {
      const updated = { ...prev, [key]: value };
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
      } catch (e) {
        console.error('Failed to persist cookie consent settings', e);
      }
      return updated;
    });
  };

  const handleSave = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
    } catch (e) {
      console.error('Failed to persist cookie consent settings', e);
    }
    setIsSaved(true);
    setTimeout(() => setIsSaved(false), 2500);
  };

  return {
    settings,
    logs,
    isSaved,
    updateSetting,
    handleSave,
  };
}
