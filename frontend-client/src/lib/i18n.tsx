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
    storefront: 'Storefront',
    cart: 'Shopping Cart',
    services: 'My Services',
    domains: 'Domains',
    invoices: 'Invoices & Billing',
    support: 'Support Tickets',
    downloads: 'Downloads',
    knowledgebase: 'Knowledgebase',
    login: 'Log In',
    register: 'Register',
    checkout: 'Proceed to Checkout',
    empty_cart: 'Your shopping cart is empty',
    order_now: 'Order Now',
  },
  id_ID: {
    storefront: 'Katalog Layanan',
    cart: 'Keranjang Belanja',
    services: 'Layanan Saya',
    domains: 'Domain',
    invoices: 'Faktur & Tagihan',
    support: 'Bantuan & Tiket',
    downloads: 'Berkas Unduhan',
    knowledgebase: 'Pusat Pengetahuan',
    login: 'Masuk',
    register: 'Daftar Akun',
    checkout: 'Lanjut ke Pembayaran',
    empty_cart: 'Keranjang belanja Anda kosong',
    order_now: 'Pesan Sekarang',
  },
  de_DE: {
    storefront: 'Storefront',
    cart: 'Warenkorb',
    services: 'Meine Dienste',
    domains: 'Domains',
    invoices: 'Rechnungen',
    support: 'Support-Tickets',
    downloads: 'Downloads',
    knowledgebase: 'Wissensdatenbank',
    login: 'Anmelden',
    register: 'Registrieren',
    checkout: 'Zur Kasse',
    empty_cart: 'Ihr Warenkorb ist leer',
    order_now: 'Jetzt Bestellen',
  },
  fr_FR: {
    storefront: 'Boutique',
    cart: 'Panier',
    services: 'Mes services',
    domains: 'Domaines',
    invoices: 'Facturation',
    support: 'Support client',
    downloads: 'Téléchargements',
    knowledgebase: 'Base de connaissances',
    login: 'Connexion',
    register: 'Inscription',
    checkout: 'Passer à la caisse',
    empty_cart: 'Votre panier est vide',
    order_now: 'Commander maintenant',
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
