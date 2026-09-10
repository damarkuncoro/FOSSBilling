import i18n from 'i18next';
import { initReactI18next, I18nextProvider } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import React from 'react';

import enCommon from '../locales/en/common.json';
import idCommon from '../locales/id/common.json';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: { common: enCommon },
      id: { common: idCommon },
    },
    fallbackLng: 'en',
    ns: ['common'],
    defaultNS: 'common',
    interpolation: {
      escapeValue: false,
    },
    detection: {
      order: ['localStorage', 'cookie', 'navigator'],
      caches: ['localStorage', 'cookie'],
    },
  });

export const I18nProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return React.createElement(I18nextProvider, { i18n }, children);
};

export default i18n;
