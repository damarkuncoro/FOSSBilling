import React from 'react';
import { useTranslation } from '@/lib/i18n';
import { Globe } from 'lucide-react';

export const LanguageSwitcher: React.FC<{ className?: string }> = ({ className = '' }) => {
  const { locale, setLocale, locales } = useTranslation();

  return (
    <div className={`relative inline-flex items-center gap-1.5 ${className}`}>
      <Globe className="h-4 w-4 text-muted-foreground shrink-0" />
      <select
        value={locale}
        onChange={(e) => setLocale(e.target.value)}
        aria-label="Select Language"
        className="bg-secondary/60 hover:bg-secondary border text-xs font-medium rounded-lg px-2.5 py-1.5 focus:outline-none focus:ring-2 focus:ring-primary cursor-pointer transition-colors"
      >
        {locales.map((loc) => (
          <option key={loc.code} value={loc.code} className="bg-popover text-foreground">
            {loc.flag} {loc.native_name} ({loc.code})
          </option>
        ))}
      </select>
    </div>
  );
};

export default LanguageSwitcher;
