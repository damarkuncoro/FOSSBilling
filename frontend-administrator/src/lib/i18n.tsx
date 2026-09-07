import React, { createContext, useContext, useState, useEffect } from 'react';

export interface LocaleItem {
  code: string;
  name: string;
  native_name: string;
  flag: string;
  direction: 'ltr' | 'rtl';
}

export const SUPPORTED_LOCALES: LocaleItem[] = [
  { code: 'en_US', name: 'English (United States)', native_name: 'English', flag: '🇺🇸', direction: 'ltr' },
  { code: 'id_ID', name: 'Indonesian', native_name: 'Bahasa Indonesia', flag: '🇮🇩', direction: 'ltr' },
  { code: 'de_DE', name: 'German', native_name: 'Deutsch', flag: '🇩🇪', direction: 'ltr' },
  { code: 'fr_FR', name: 'French', native_name: 'Français', flag: '🇫🇷', direction: 'ltr' },
  { code: 'es_ES', name: 'Spanish', native_name: 'Español', flag: '🇪🇸', direction: 'ltr' },
  { code: 'ja_JP', name: 'Japanese', native_name: '日本語', flag: '🇯🇵', direction: 'ltr' },
  { code: 'ar_SA', name: 'Arabic', native_name: 'العربية', flag: '🇸🇦', direction: 'rtl' },
];

const TRANSLATIONS: Record<string, Record<string, string>> = {
  en_US: {
    dashboard: 'Dashboard',
    clients: 'Clients',
    orders: 'Orders',
    invoices: 'Invoices',
    support: 'Support Tickets',
    settings: 'Settings',
    extensions: 'Extensions',
    system_health: 'System Health',
    logout: 'Logout',
    welcome: 'Welcome back, :name',
    total_revenue: 'Total Revenue',
    active_clients: 'Active Clients',
    pending_orders: 'Pending Orders',
    open_tickets: 'Open Tickets',
  },
  id_ID: {
    dashboard: 'Dasbor',
    clients: 'Klien',
    orders: 'Pesanan',
    invoices: 'Faktur & Tagihan',
    support: 'Tiket Bantuan',
    settings: 'Pengaturan',
    extensions: 'Ekstensi & Modul',
    system_health: 'Kesehatan Sistem',
    logout: 'Keluar',
    welcome: 'Selamat datang kembali, :name',
    total_revenue: 'Total Pendapatan',
    active_clients: 'Klien Aktif',
    pending_orders: 'Pesanan Tertunda',
    open_tickets: 'Tiket Terbuka',
  },
  de_DE: {
    dashboard: 'Übersicht',
    clients: 'Kunden',
    orders: 'Bestellungen',
    invoices: 'Rechnungen',
    support: 'Support-Tickets',
    settings: 'Einstellungen',
    extensions: 'Erweiterungen',
    system_health: 'Systemstatus',
    logout: 'Abmelden',
    welcome: 'Willkommen zurück, :name',
    total_revenue: 'Gesamteinnahmen',
    active_clients: 'Aktive Kunden',
    pending_orders: 'Ausstehende Bestellungen',
    open_tickets: 'Offene Tickets',
  },
  fr_FR: {
    dashboard: 'Tableau de bord',
    clients: 'Clients',
    orders: 'Commandes',
    invoices: 'Factures',
    support: 'Tickets de support',
    settings: 'Paramètres',
    extensions: 'Extensions',
    system_health: 'Santé du système',
    logout: 'Se déconnecter',
    welcome: 'Bon retour, :name',
    total_revenue: 'Revenu total',
    active_clients: 'Clients actifs',
    pending_orders: 'Commandes en attente',
    open_tickets: 'Tickets ouverts',
  },
};

interface I18nContextType {
  locale: string;
  setLocale: (code: string) => void;
  t: (key: string, params?: Record<string, string | number>) => string;
  currentLocaleInfo: LocaleItem;
  locales: LocaleItem[];
}

const I18nContext = createContext<I18nContextType | null>(null);

export const I18nProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [locale, setLocaleState] = useState<string>(() => {
    return localStorage.getItem('fossbilling_locale') || 'en_US';
  });

  const setLocale = (code: string) => {
    setLocaleState(code);
    localStorage.setItem('fossbilling_locale', code);
    document.cookie = `BBLANG=${code};path=/;max-age=2592000`;
  };

  const currentLocaleInfo = SUPPORTED_LOCALES.find((l) => l.code === locale) || SUPPORTED_LOCALES[0];

  useEffect(() => {
    document.documentElement.lang = locale.replace('_', '-');
    document.documentElement.dir = currentLocaleInfo.direction;
  }, [locale, currentLocaleInfo]);

  const t = (key: string, params?: Record<string, string | number>): string => {
    let text = TRANSLATIONS[locale]?.[key] || TRANSLATIONS.en_US?.[key] || key;
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        text = text.replace(new RegExp(`:${k}`, 'g'), String(v));
      });
    }
    return text;
  };

  return (
    <I18nContext.Provider value={{ locale, setLocale, t, currentLocaleInfo, locales: SUPPORTED_LOCALES }}>
      {children}
    </I18nContext.Provider>
  );
};

export const useTranslation = () => {
  const context = useContext(I18nContext);
  if (!context) {
    throw new Error('useTranslation must be used within an I18nProvider');
  }
  return context;
};
